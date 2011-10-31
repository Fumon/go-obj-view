package glmatrix

import (
	"fmt"
	"math"
)

//A column major ordered matrix
type mat4 [16]float32

func MakeMat4(a float32) (m *mat4) {
	m = new(mat4)
	for i := range (*m) {
		(*m)[i] = a
	}
	return
}

func IdMat4() (m *mat4) {
	m = new(mat4)
	for i := range m {
		switch i {
		case 0, 5, 10, 15: m[i] = 1.0
		default: m[i] = 0.0
		}
	}
	return
}

func ScaleMat4(s float32) (m *mat4) {
	m = new(mat4)
	for i := range m {
		switch i {
		case 0, 5, 10, 15: m[i] = s
		default: m[i] = 0.0
		}
	}
	return
}

func (m mat4) String() string {
	var result string
	for row := 0; row < 4; row++ {
		result +=
		  fmt.Sprintln(m[row], m[row + 4], m[row + 8], m[row + 12])
	}
	return result
}

//Remember, column major and zero indexed
func (m *mat4) At(row, col int) float32 {
	return m[4*col + row]
}

//a * b in written order
func (a *mat4) Product(b *mat4) (mv *mat4) {
	mv = new(mat4)
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			var sum float32
			for i := 0; i < 4; i++ {
				sum += a.At(row, i) * b.At(i, col)
			}
			mv[col*4 + row] = sum
		}
	}

	return
}

func TranslateMat4(t []float32) (m *mat4) {
	m = IdMat4()
	m[12] = t[0]
	m[13] = t[1]
	m[14] = t[2]
	return
}

func AxisAngleRotation(axis []float32, angle float32) (mv *mat4) {
	mv = IdMat4()
	c := math.Cos(float64(angle))
	s := math.Sin(float64(angle))
	t := 1 - c
	x := float64(axis[0])
	y := float64(axis[1])
	z := float64(axis[2])
	//Exploded
	mv[0] = float32(t * x * x + c)
	mv[1] = float32(t * x * y + z * s)
	mv[2] = float32(t * x * z - y * s)

	mv[4] = float32(t * x * y - z * s)
	mv[5] = float32(t * y * y + c)
	mv[6] = float32(t * y * z + x * s)

	mv[8] = float32(t * x * z + y * s)
	mv[9] = float32(t * y * z - x * s)
	mv[10] = float32(t * z * z + c)

	return
}
