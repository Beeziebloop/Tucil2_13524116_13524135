package main

import (
	"fmt"
	"time"
	"sync"
)

type Stats struct {
	mu sync.Mutex
	nodeCount map[int]int
	skippedCount map[int]int
	startTime time.Time
}

func newStats(maxDepth int) *Stats {
	nodeCount := make(map[int]int)
	skippedCount := make(map[int]int)
	//inisialisasi semua depth-0
	for i := 0; i <= maxDepth; i++ {
		nodeCount[i] = 0
		skippedCount[i] = 0
	}
	return &Stats{
		nodeCount: nodeCount,
		skippedCount: skippedCount,
		startTime: time.Now(),
	}
}

func (s *Stats) printStats(leaves []*OctreeNode, maxDepth int, outputPath string) {
	elapsed := time.Since(s.startTime)
	//hitung voxel, vertex, dan faces
	totVox := len(leaves)
	totVertices := totVox * 8 //karena tiap kubus itu ada 8 vertex
	totFaces := totVox * 12 //karena tiap kubus itu ada 6 sisi * 2 segitiga diagonal = 12 faces

	fmt.Println("================ Stats ================")
	fmt.Printf("Jumlah voxel : %d\n", totVox)
	fmt.Printf("Jumlah vertices : %d\n", totVertices)
	fmt.Printf("Jumlah faces : %d\n", totFaces)
	fmt.Printf("Kedalaman octree: %d\n", maxDepth)
	fmt.Println("Statistik node per depth:")
	for i := 1; i <= maxDepth; i++ {
		fmt.Printf("  %d : %d\n", i, s.nodeCount[i])
	}
	fmt.Println("Statistik node yang dilewati per depth:")
	for i := 1; i <= maxDepth; i++ {
		fmt.Printf("  %d : %d\n", i, s.skippedCount[i])
	}
	fmt.Printf("Waktu eksekusi: %s\n", elapsed)
	fmt.Printf("Output .obj disimpan: %s\n", outputPath)
	fmt.Println("=====================================")
}