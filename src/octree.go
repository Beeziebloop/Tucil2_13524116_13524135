package main

import (
	"sync"
)

type OctreeNode struct {
	isLeaf           bool
	children         [8]*OctreeNode
	minX, minY, minZ float64
	size             float64
	depth            int
}

func newOctNode(minX, minY, minZ, size float64, depth int) *OctreeNode {
	return &OctreeNode{
		minX:  minX,
		minY:  minY,
		minZ:  minZ,
		size:  size,
		depth: depth,
	}
}

func calculateBoundingBox(triangles []Triangle) (minB, maxB Vector3) {
	minB = triangles[0].a
	maxB = triangles[0].a
	for _, tri := range triangles {
		verts := [3]Vector3{tri.a, tri.b, tri.c}
		for _, v := range verts {
			if v.x < minB.x {
				minB.x = v.x
			}
			if v.y < minB.y {
				minB.y = v.y
			}
			if v.z < minB.z {
				minB.z = v.z
			}
			if v.x > maxB.x {
				maxB.x = v.x
			}
			if v.y > maxB.y {
				maxB.y = v.y
			}
			if v.z > maxB.z {
				maxB.z = v.z
			}
		}
	}
	return
}

func nextPowerOfTwo(x float64) float64 {
	p := 1.0
	for p < x {
		p *= 2
	}
	return p
}

func createRootNode(triangles []Triangle) *OctreeNode {
	minB, maxB := calculateBoundingBox(triangles)

	//cari sisi terpanjang
	dx := maxB.x - minB.x
	dy := maxB.y - minB.y
	dz := maxB.z - minB.z

	size := nextPowerOfTwo(dx)
	if dy > size {
		size = dy
	}
	if dz > size {
		size = dz
	}

	//padding biar objek tidak tepat di sisi kubus
	size *= 1.001

	//center bounding box
	cx := (minB.x + maxB.x) / 2
	cy := (minB.y + maxB.y) / 2
	cz := (minB.z + maxB.z) / 2

	return newOctNode(
		cx-size/2,
		cy-size/2,
		cz-size/2,
		size,
		0,
	)
}

func buildOctree(node *OctreeNode, triangles []Triangle, maxDepth int, stats *Stats) {
	stats.mu.Lock()
	//tandai node ini di stats
	stats.nodeCount[node.depth]++
	stats.mu.Unlock()
	//kalau sudah di max depth, ini adalah voxel
	if node.depth == maxDepth {
		node.isLeaf = true
		return
	}

	childSize := node.size / 2
	var wg sync.WaitGroup

	//8 anak kubus, index 0-7 berdasarkan posisi x, y, z (0 atau 1)
	for i := 0; i < 8; i++ {
		//tentukan offset tiap anak berdasarkan index binary
		//bit 0 = x, bit 1 = y, bit 2 = z
		ox := float64((i >> 0) & 1)
		oy := float64((i >> 1) & 1)
		oz := float64((i >> 2) & 1)

		childMin := Vector3{
			x: node.minX + ox*childSize,
			y: node.minY + oy*childSize,
			z: node.minZ + oz*childSize,
		}
		childCenter := Vector3{
			x: childMin.x + childSize/2,
			y: childMin.y + childSize/2,
			z: childMin.z + childSize/2,
		}
		childHalf := Vector3{childSize / 2, childSize / 2, childSize / 2}

		//cek segitiga mana yang berpotongan dengan anak ini
		var childTriangles []Triangle
		for _, tri := range triangles {
			if isTriBoxOverlap(childCenter, childHalf, tri) {
				childTriangles = append(childTriangles, tri)
			}
		}

		//kalau tidak ada segitiga, skip
		if len(childTriangles) == 0 {
			stats.mu.Lock()
			stats.skippedCount[node.depth+1]++
			stats.mu.Unlock()
			continue
		}

		//kalau ada segitiga, buat node anak dan rekursi
		child := newOctNode(childMin.x, childMin.y, childMin.z, childSize, node.depth+1)
		node.children[i] = child
		//hanya spawn goroutine di depth awal, kalau terlalu dalam nanti malah overhead
		if node.depth < 3 {
			wg.Add(1)
			go func(c *OctreeNode, tris []Triangle) {
				defer wg.Done()
				buildOctree(c, tris, maxDepth, stats)
			}(child, childTriangles)
		} else {
			buildOctree(child, childTriangles, maxDepth, stats)
		}
	}
	wg.Wait()
}

//kumpulkan semua leaf nodes (voxel)
func collectLeaves(node *OctreeNode) []*OctreeNode {
	if node == nil {
		return nil
	}
	if node.isLeaf {
		return []*OctreeNode{node}
	}
	var leaves []*OctreeNode
	for _, child := range node.children {
		leaves = append(leaves, collectLeaves(child)...)
	}
	return leaves
}
