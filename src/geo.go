package main

//structs
type Vector3 struct{
	x,y,z float64
}

type Triangle struct{
	a,b,c Vector3
}

//vector operations helpers
func sub(a, b Vector3) Vector3{
	return Vector3{a.x - b.x, a.y - b.y, a.z - b.z}
}

func cross(a, b Vector3) Vector3{
	return Vector3{a.y * b.z - a.z * b.y, a.z * b.x - a.x * b.z, a.x * b.y - a.y * b.x}
}

func dot(a, b Vector3) float64{
	return a.x * b.x + a.y * b.y + a.z * b.z
}

func abs(x float64) float64{
	if x < 0 {
		return -x
	}
	return x
}

func minMax(a, b, c float64) (minimum, maximum float64){
	minimum = min(a, min(b, c))
	maximum = max(a, max(b, c))
	return
}