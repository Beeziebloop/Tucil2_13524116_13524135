package main

import (
	"fmt"
	"os"
)

func writeOBJ(leaves []*OctreeNode, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("gagal membuat file output: %v", err)
	}
	defer file.Close()

	writer := fmt.Fprintf
	vertexOffset := 0

	for _, leaf := range leaves {
		//8 sudut kubus berdasarkan minX, minY, minZ dan size
		x0 := leaf.minX
		y0 := leaf.minY
		z0 := leaf.minZ
		x1 := leaf.minX + leaf.size
		y1 := leaf.minY + leaf.size
		z1 := leaf.minZ + leaf.size

		//8 vertex kubus
		//urutannya bottom face dulu baru top face
		writer(file, "v %f %f %f\n", x0, y0, z0) // 1
		writer(file, "v %f %f %f\n", x1, y0, z0) // 2
		writer(file, "v %f %f %f\n", x1, y1, z0) // 3
		writer(file, "v %f %f %f\n", x0, y1, z0) // 4
		writer(file, "v %f %f %f\n", x0, y0, z1) // 5
		writer(file, "v %f %f %f\n", x1, y0, z1) // 6
		writer(file, "v %f %f %f\n", x1, y1, z1) // 7
		writer(file, "v %f %f %f\n", x0, y1, z1) // 8
		//offset untuk indeks face, tiap kubus nambah 8
		o := vertexOffset
		//12 faces (6 sisi * 2 segitiga per sisi)
		//bottom (z0)
		writer(file, "f %d %d %d\n", o+1, o+2, o+3)
		writer(file, "f %d %d %d\n", o+1, o+3, o+4)
		//top (z1)
		writer(file, "f %d %d %d\n", o+5, o+7, o+6)
		writer(file, "f %d %d %d\n", o+5, o+8, o+7)
		//front (y0)
		writer(file, "f %d %d %d\n", o+1, o+6, o+2)
		writer(file, "f %d %d %d\n", o+1, o+5, o+6)
		//back (y1)
		writer(file, "f %d %d %d\n", o+4, o+3, o+7)
		writer(file, "f %d %d %d\n", o+4, o+7, o+8)
		//left (x0)
		writer(file, "f %d %d %d\n", o+1, o+4, o+8)
		writer(file, "f %d %d %d\n", o+1, o+8, o+5)
		//right (x1)
		writer(file, "f %d %d %d\n", o+2, o+6, o+7)
		writer(file, "f %d %d %d\n", o+2, o+7, o+3)
		vertexOffset += 8
	}
	
	return nil
}