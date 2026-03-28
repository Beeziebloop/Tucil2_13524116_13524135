# Tucil2 IF2211 - Voxelization Objek 3D menggunakan Octree

##Deskripsi Program
Program ini mengkonversi model 3D dalam format `.obj` menjadi model voxelized (tersusun dari kubus-kubus kecil seperti Minecraft) menggunakan struktur data Octree dengan algoritma Divide and Conquer. Program juga memanfaatkan concurrency (goroutines) untuk mempercepat proses konversi.

## Requirements
- Go 1.21 atau lebih baru
- Tidak ada dependency eksternal (hanya menggunakan standard library Go)

## Instalasi Go
Download dan install Go dari: https://go.dev/dl/
Pastikan Go sudah terdaftar di PATH dengan menjalankan:
```bash
go version
```

## Struktur Repository
```
Tucil2_13524116_13524135/
├── doc/
│   └── laporan.pdf
├── src/
│   ├── go.mod
│   ├── main.go
│   ├── parser.go
│   ├── geometry.go
│   ├── intersection.go
│   ├── octree.go
│   ├── voxel.go
│   └── stats.go
└── test/
    ├── (input .obj files)
    └── (output _voxed.obj files)
```

## Cara Kompilasi
```bash
cd src
go build -o ../bin/tucil2
```

## Cara Menjalankan Program

### Menggunakan go run
```bash
cd src
go run . <nama_file.obj> <max_depth>
```

### Menggunakan executable (setelah kompilasi)
```bash
cd bin
./tucil2 <nama_file.obj> <max_depth>
```

### Contoh
```bash
cd src
go run . pumpkin.obj 7
go run . teapot.obj 5
go run . cube.obj 6
```

### Catatan
- File `.obj` input harus diletakkan di folder `test/`
- Output akan tersimpan di folder `test/` dengan nama `[nama_file]_voxed.obj`
- `max_depth` harus berupa integer positif (rekomendasi: 5-7)
- Semakin tinggi `max_depth`, semakin detail hasil voxelization tapi semakin lama waktu proses

## Format Input
Program hanya menerima file `.obj` dengan format:
```
v x y z        (vertex)
f i j k        (face, hanya segitiga saja, tidak support quad ataupun poligon)
```
Format face yang didukung:
- `f 1 2 3`           (vertex only)
- `f 1/2 3/4 5/6`     (vertex/texture)
- `f 1//3 4//6 7//9`  (vertex//normal)
- `f 1/2/3 4/5/6 7/8/9` (vertex/texture/normal)
Program hanya membaca vertex index dari tiap face dan mengabaikan texture coordinate dan normal.
Contoh isi file .obj valid:
```
v 2.229345 -0.992723 -0.862826
v 2.292449 -0.871852 -0.882400
v 2.410367 -0.777999 -0.841105
f 1 2 3
```

## Output
Program menghasilkan:
- File `.obj` hasil voxelization di folder `test/`
- Informasi statistik di CLI:
  - Jumlah voxel, vertex, dan faces
  - Statistik node octree per depth
  - Statistik node yang di-skip per depth
  - Kedalaman octree
  - Waktu eksekusi
  - Path file output

## Performa
Program menggunakan goroutines untuk memparalelkan proses build octree.
Contoh perbandingan waktu pada `pumpkin.obj` depth 7:
- Tanpa concurrency  6.97 detik 
- Dengan concurrency  3.04 detik 
Speedup = ~2.3x 

## Author
Eliana Natalie Widjojo 13524116
Varistha Devi 13524135
