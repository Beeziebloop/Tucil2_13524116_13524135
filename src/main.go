package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	//validasi argument
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . <nama_file.obj> <max_depth>")
		fmt.Println("Contoh: go run . pumpkin.obj 5")
		os.Exit(1)
	}

	fileName := os.Args[1]
	maxDepth, err := strconv.Atoi(os.Args[2])
	if err != nil || maxDepth < 1 {
		fmt.Println("Error: max_depth harus berupa integer positif")
		os.Exit(1)
	}

	//validasi ekstensi file
	if strings.ToLower(filepath.Ext(fileName)) != ".obj" {
		fmt.Println("Error: file harus berekstensi .obj")
		os.Exit(1)
	}

	//path ke folder test (relatif dari src/)
	testDir := filepath.Join("..", "test")

	//input path
	inputPath := filepath.Join(testDir, fileName)

	//validasi file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Printf("Error: file '%s' tidak ditemukan di folder test/\n", fileName)
		os.Exit(1)
	}

	//output path: [nama]_voxed.obj
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	outputPath := filepath.Join(testDir, baseName+"_voxed.obj")

	fmt.Printf("Input file : %s\n", inputPath)
	fmt.Printf("Output file: %s\n", outputPath)
	fmt.Printf("Max depth  : %d\n", maxDepth)
	fmt.Println("Memproses...")

	//step 1: parse
	triangles, err := parseObjFile(inputPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("Parse selesai: %d triangles\n", len(triangles))

	//step 2: buat root node
	root := createRootNode(triangles)
	fmt.Printf("Bounding box: min = (%.3f, %.3f, %.3f) size = %.3f\n",
		root.minX, root.minY, root.minZ, root.size)

	//step 3: build octree
	stats := newStats(maxDepth)
	fmt.Println("Membangun octree...")
	buildOctree(root, triangles, maxDepth, stats)

	//step 4: kumpulkan leaves
	leaves := collectLeaves(root)
	fmt.Printf("Octree selesai: %d voxel ditemukan\n", len(leaves))

	//step 5: tulis output
	fmt.Println("Menulis output...")
	err = writeOBJ(leaves, outputPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	//step 6: print stats
	stats.printStats(leaves, maxDepth, outputPath)
}