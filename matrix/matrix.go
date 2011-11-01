package glmatrix

import (
	"fmt"
	"math"
)

//A column major ordered matrix
type Mat4 [16]float32

func MakeMat4(a float32) (m *Mat4) {
	m = new(Mat4)
	for i := range (*m) {
		(*m)[i] = a
	}
	return
}

func IdMat4() (m *Mat4) {
	m = new(Mat4)
	for i := range m {
		switch i {
		case 0, 5, 10, 15: m[i] = 1.0
		default: m[i] = 0.0
		}
	}
	return
}

func ScaleMat4(s float32) (m *Mat4) {
	m = new(Mat4)
	for i := range m {
		switch i {
		case 0, 5, 10, 15: m[i] = s
		default: m[i] = 0.0
		}
	}
	return
}

func (m Mat4) String() string {
	var result string
	for row := 0; row < 4; row++ {
		result +=
		  fmt.Sprintln(m[row], m[row + 4], m[row + 8], m[row + 12])
	}
	return result
}

//Remember, column major and zero indexed
func (m *Mat4) At(row, col int) float32 {
	return m[4*col + row]
}



func (m *Mat4) Transpose() (mp *Mat4) {
	mp = new(Mat4)
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			mp[col * 4 + row] = m[row * 4 + col]
		}
	}
	return
}

//a * b in written order
func (a *Mat4) Product(b *Mat4) (mv *Mat4) {
	mv = new(Mat4)
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

//a * b where by is a 4x1
func (a *Mat4) ProductV(b []float32) (bp []float32) {
	bp = make([]float32, 4)
	for row := 0; row < 4; row++ {
		var sum float32
		for col := 0; col < 4; col++ {
			sum += a.At(row, col) * b[col]
		}
		bp[row] = sum
	}
	return
}

//Return a translation in homogeneous coordinates
func TranslateMat4(t []float32) (m *Mat4) {
	m = IdMat4()
	m[12] = t[0]
	m[13] = t[1]
	m[14] = t[2]
	return
}

func mag(vec []float32) (mag float32) {
	mag = float32(math.Sqrt(math.Pow(float64(vec[0]), 2) + math.Pow(float64(vec[1]), 2) + math.Pow(float64(vec[2]), 2)))
	return
}

func cross(a, b []float32) (*[3]float32) {
	//cross = (a2*b3 - a3*b2, a3*b1 - a1*b3, a1*b2 - a2*b1)
	cu := new([3]float32)
	cu[0] = a[1]*b[2] - a[2]*b[1]
	cu[1] = a[2]*b[0] - a[0]*b[2]
	cu[2] = a[0]*b[1] - a[1]*b[0]
	return cu
}

func AxisAngleRotation(axis []float32, angle float32) (mv *Mat4) {
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

//Return a view matrix for a camera of position pos, view direction direction and up vector up.
func ViewLookAt(pos, direction, up []float32) (mv *Mat4) {
	x := new([3]float32)
	y := new([3]float32)
	z := new([3]float32)
	//z = negative normed direction
	m := mag(direction)
	z[0] = -direction[0]/m
	z[1] = -direction[1]/m
	z[2] = -direction[2]/m

	//x = normed direction cross up
	x = cross(direction, up)
	m = mag(x[:])
	x[0] = x[0] / m
	x[1] = x[1] / m
	x[2] = x[2] / m

	//y = z cross x
	y = cross(z[:], x[:])

	//Make Rtranspose
	RT := Mat4{x[0], y[0], z[0], 0, x[1], y[1], z[1], 0, x[2], y[2], z[2], 0, 0, 0, 0, 1}
	t := RT.ProductV([]float32{pos[0], pos[1], pos[2], 1})
	RT[12] = -t[0]
	RT[13] = -t[1]
	RT[14] = -t[2]
	mv = &RT
	return
}

func StdProjection(fovy, near, far, ar float32) (m *Mat4) {
	m = new(Mat4)
	d := float32(1.0/math.Tan(float64(fovy)/2.0))

	m[0] = d/ar
	m[5] = d
	m[10] = (near + far) / (near - far)
	m[11] = -1
	m[14] = (2 * near * far) / (near - far)

	return m
}
