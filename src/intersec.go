package main

func isBoxOverlap(normal, vert, maxbox Vector3)bool{
	var vmin, vmax Vector3
	if normal.x > 0 {
		vmin.x = -maxbox.x - vert.x
		vmax.x = maxbox.x - vert.x
	} else {
		vmin.x = maxbox.x - vert.x
		vmax.x = -maxbox.x - vert.x
	}
	if normal.y > 0 {
		vmin.y = -maxbox.y - vert.y
		vmax.y = maxbox.y - vert.y
	} else {
		vmin.y = maxbox.y - vert.y
		vmax.y = -maxbox.y - vert.y
	}
	if normal.z > 0 {
		vmin.z = -maxbox.z - vert.z
		vmax.z = maxbox.z - vert.z
	} else {
		vmin.z = maxbox.z - vert.z
		vmax.z = -maxbox.z - vert.z
	}

	if dot(normal, vmin) > 0 {
		return false
	}
	if dot(normal, vmax) >= 0 {
		return true
	}

	return false
}

func isTriBoxOverlap(boxcenter, boxhalfsize Vector3, tri Triangle) bool {
	//translate dulu ke origin biar lebih uniform
	v1 := sub(tri.a, boxcenter)
	v2 := sub(tri.b, boxcenter)
	v3 := sub(tri.c, boxcenter)

	//hitung edges dari segitiga
	e1 := sub(v2, v1)
	e2 := sub(v3, v2)
	e3 := sub(v1, v3)

	//9 axis tests
	//pake bantuan helper buat ngecek salah satu axis, akan return false kalau ada separating axis
	axTest := func(p0, p2, rad float64) bool {
		minimum, maximum := minMax(p0, p0, p2)
		return !(minimum > rad || maximum < -rad)
	}

	var p0, p2, rad float64

	//e1 tests
	fex := abs(e1.x)
	fey := abs(e1.y)
	fez := abs(e1.z)

	p0 = e1.z * v1.y - e1.y * v1.z
	p2 = e1.z * v3.y - e1.y * v3.z
	rad = fez * boxhalfsize.y + fey * boxhalfsize.z
	if !axTest(p0, p2, rad) {
		return false
	}

	p0 = -e1.z * v1.x + e1.x * v1.z
	p2 = -e1.z * v3.x + e1.x * v3.z
	rad = fez * boxhalfsize.x + fex * boxhalfsize.z
	if !axTest(p0, p2, rad) {
		return false
	}

	p0 = e1.y * v2.x - e1.x * v2.y
	p2 = e1.y * v3.x - e1.x * v3.y
	rad = fey * boxhalfsize.x + fex * boxhalfsize.y
	if !axTest(p0, p2, rad) {
		return false
	}

	//e2 tests
	fex = abs(e2.x)
	fey = abs(e2.y)
	fez = abs(e2.z)

	p0 = e2.z * v1.y - e2.y * v1.z
	p2 = e2.z * v3.y - e2.y * v3.z
	rad = fez * boxhalfsize.y + fey * boxhalfsize.z
	if !axTest(p0, p2, rad) {
		return false
	}

	p0 = -e2.z * v1.x + e2.x * v1.z
	p2 = -e2.z * v3.x + e2.x * v3.z
	rad = fez * boxhalfsize.x + fex * boxhalfsize.z
	if !axTest(p0, p2, rad) {
		return false
	}

	p0 = e2.y * v1.x - e2.x * v1.y
	p2 = e2.y * v2.x - e2.x * v2.y
	rad = fey * boxhalfsize.x + fex * boxhalfsize.y
	if !axTest(p0, p2, rad) {
		return false
	}

	//e3 tests
	fex = abs(e3.x)
	fey = abs(e3.y)
	fez = abs(e3.z)

	p0 = e3.z * v1.y - e3.y * v1.z
	p2 = e3.z * v2.y - e3.y * v2.z
	rad = fez * boxhalfsize.y + fey * boxhalfsize.z
	if !axTest(p0, p2, rad) {
		return false
	}

	p0 = -e3.z * v1.x + e3.x * v1.z
	p2 = -e3.z * v2.x + e3.x * v2.z
	rad = fez * boxhalfsize.x + fex * boxhalfsize.z
	if !axTest(p0, p2, rad) {
		return false
	}

	p0 = e3.y * v1.x - e3.x * v1.y
	p2 = e3.y * v2.x - e3.x * v2.y
	rad = fey * boxhalfsize.x + fex * boxhalfsize.y
	if !axTest(p0, p2, rad) {
		return false
	}

	//AABB box intersection test di tiap sumbu
	minimum, maximum := minMax(v1.x, v2.x, v3.x)
	if minimum > boxhalfsize.x || maximum < -boxhalfsize.x {
		return false
	}

	minimum, maximum = minMax(v1.y, v2.y, v3.y)
	if minimum > boxhalfsize.y || maximum < -boxhalfsize.y {
		return false
	}

	minimum, maximum = minMax(v1.z, v2.z, v3.z)
	if minimum > boxhalfsize.z || maximum < -boxhalfsize.z {
		return false
	}

	//plane-box overlap test
	normal := cross(e1, e2)
	if !isBoxOverlap(normal, v1, boxhalfsize) {
		return false
	}

	return true
}