package main
import(
	"os"
	"fmt"
	"bufio"
	"strings"
	"strconv"
)

func parseObjFile(filepath string)([]Triangle, error){
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("file %s tidak ditemukan", filepath)
	}
	defer file.Close()

	var vertices []Vector3
    var triangles []Triangle

    scanner := bufio.NewScanner(file)
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        line := strings.TrimSpace(scanner.Text())

        //strip inline comment
        if idx := strings.Index(line, "#"); idx != -1 {
            line = strings.TrimSpace(line[:idx])
        }
        //skip baris kosong
        if line == "" {
            continue
        }

        parts := strings.Fields(line)

        switch parts[0] {
        case "v":
            //harus ada 3 koordinat
            if len(parts) != 4 {
                return nil, fmt.Errorf("baris %d: format vertex tidak valid", lineNum)
            }
            x, errX := strconv.ParseFloat(parts[1], 64)
            y, errY := strconv.ParseFloat(parts[2], 64)
            z, errZ := strconv.ParseFloat(parts[3], 64)
            if errX != nil || errY != nil || errZ != nil {
                return nil, fmt.Errorf("baris %d: koordinat vertex bukan angka", lineNum)
            }
            vertices = append(vertices, Vector3{x, y, z})
        case "f":
            //harus ada 3 indeks untuk segitiga
            if len(parts) != 4 {
                return nil, fmt.Errorf("baris %d: format face tidak valid", lineNum)
            }
            //asumsinya isi .obj filenya tidak akan mengandung format face yang kompleks pada tugas ini
			parseIndex := func(s string) (int, error) {
				return strconv.Atoi(strings.Split(s, "/")[0])
			}
            i, errI := parseIndex(parts[1])
            j, errJ := parseIndex(parts[2])
            k, errK := parseIndex(parts[3])
            if errI != nil || errJ != nil || errK != nil {
                return nil, fmt.Errorf("baris %d: indeks face bukan integer", lineNum)
            }
            //indeks tidak boleh out of range, dengan idx .obj dimulai dari 1
            if i < 1 || j < 1 || k < 1 || i > len(vertices) || j > len(vertices) || k > len(vertices) {
                return nil, fmt.Errorf("baris %d: indeks face out of range", lineNum)
            }
            triangles = append(triangles, Triangle{vertices[i-1], vertices[j-1], vertices[k-1]})

        default:
            //skip baris yang bukan v atau f (vt, vn, dll karena asumsinya hanya v dan f yang diproses dan isi .obj itu simpel)
            continue
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("error membaca file %s: %v", filepath, err)
    }
    if len(vertices) == 0 {
        return nil, fmt.Errorf("file tidak mengandung vertex")
    }
    if len(triangles) == 0 {
        return nil, fmt.Errorf("file tidak mengandung faces")
    }

    return triangles, nil
}