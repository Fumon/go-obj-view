package glmatrix

import (
	"testing"
	"testing/quick"
	"math"
)

/*
func Test4x4Transpose(t *testing.T) {
	f := func (a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16 float32) bool {
		m := &Mat4{a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16}
		mt := &Mat4{a1, a5, a9, a13, a2, a6, a10, a14, a3, a7, a11, a15, a4, a8, a12, a16}
		m = m.Transpose()

		for i, me := range m {
			if me != mt[i] {
				return false
			}
		}
		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
*/

func TestUpper3x3(t *testing.T) {
	m := &Mat4{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}
	mu := &Mat3{1.0, 2.0, 3.0, 5.0, 6.0, 7.0, 9.0, 10.0, 11.0}
	mcu := m.Upper3by3()
	for i, me := range mcu {
		if me != mu[i] {
			t.FailNow()
		}
	}
}

func TestDeterminant3x3(t *testing.T) {
	m := &Mat3{-2.0, 5.0, 4.0, 3.0, -1.0, -8.0, -1.0, 4.0, 2.0}
	d := m.Determinant()
	if d != -6.0 {
		t.Fatalf("Determinant = %v instead of -6.0", d)
	}
}

func TestInverse3x3(t *testing.T) {
	m := &Mat3{1.0, 1.0, 1.0, 3.0, 4.0, 3.0, 3.0, 3.0, 4.0}
	mI := m.Inverse()
	//ma := &Mat3{7.0, -1.0, -1.0, -3.0, 1.0, 0.0, -3.0, 0.0, 1.0}
	mt := mI.Product(m)
	ma := IdMat3()
	for i, me := range mt {
		if me != ma[i] {
			t.Fatalf("Fails on %v\nm\n%v\ninv\n%v\nmt\n%v\nmId\n%v", i, m, mI, mt, ma)
		}
	}
}

func TestInverse3x3Quick(t *testing.T) {
	i := 0
	f := func (a1 float32) bool {
		mat := AxisAngleRotation([]float32{0.0, 1.0, 0.0}, (a1/math.MaxFloat32)*math.Pi*3.0).Upper3by3()
		//mat := &Mat3{a1, a2, a3, a4, a5, a6, a7, a8, a9}
		det := mat.Determinant()
		if det==0 || math.IsNaN(float64(det)){
			i++
			return true
		}

		mI := mat.Inverse()
		mProd := mat.Product(mI)
		mIdent := IdMat3()
		for i, e := range mProd {
				if math.Abs(float64(mIdent[i] - e)) > 0.000001 {
				t.Logf("Fails\ndet: %v a1: %v\nm\n%v\nmI\n%v\nmProd\n%v\nmIdent\n%v\n", det, a1, mat, mI, mProd, mIdent)
				return false
			}
		}
		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
	t.Logf("Didn't check %v NaNs and 0's", i)
}
