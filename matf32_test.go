package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewf32(t *testing.T) {
	t.Helper()
	rows := 13
	cols := 7
	m := Newf32()
	assert.Equal(t, 0, m.r, "should be zero")
	assert.Equal(t, 0, m.c, "should be zero")
	assert.NotNil(t, m.vals, "should not be nil")
	assert.Equal(t, 0, len(m.vals), "should be zero")

	m = Newf32(rows)
	assert.Equal(t, rows, m.r, "should be equal")
	assert.Equal(t, rows, m.c, "should be equal")
	assert.NotNil(t, m.vals, "should not be nil")
	assert.Equal(t, rows*rows, len(m.vals), "should be equal")

	m = Newf32(rows, cols)
	assert.Equal(t, rows, m.r, "should be equal")
	assert.Equal(t, cols, m.c, "should be equal")
	assert.NotNil(t, m.vals, "should not be nil")
	assert.Equal(t, rows*cols, len(m.vals), "should be equal")

	// assert.Panics(t, func() { Newf32(1, 2, 3, 4) }, "should panic with 3+ args")
}

func TestMatf32FromData(t *testing.T) {
	t.Helper()
	rows := 50
	cols := 2

	// assert.Panics(t, func() { Matf32FromData(1.0) }, "should panic with wrong arg")

	v := make([]float32, rows*cols)
	for i := range v {
		v[i] = float32(i * i)
	}

	m := Matf32FromData(v)
	assert.Equal(t, 1, m.r, "should have one row")
	assert.Equal(t, len(v), len(m.vals), "should have the same # of elements")
	for i := range v {
		assert.Equal(t, v[i], m.vals[i], "should be equal")
	}
	v[0] = 1321.0
	assert.NotEqual(t, v[0], m.vals[0], "changing data should not effect mat")
	m.vals[0] = 1201.0
	assert.NotEqual(t, m.vals[0], v[0], "changing mat should not effect data")

	s := make([][]float32, rows)
	for i := range s {
		s[i] = make([]float32, cols)
	}
	for i := range s {
		for j := range s[i] {
			s[i][j] = float32(i + j)
		}
	}
	m = Matf32FromData(s)
	assert.Equal(t, rows*cols, len(m.vals), "should be equal")
	idx := 0
	for i := range s {
		for j := range s[i] {
			assert.Equal(t, s[i][j], m.vals[idx], "should be equal")
			idx++
		}
	}
	s[0][0] = 1021.0
	assert.NotEqual(t, s[0][0], m.vals[0], "changing data should not effect mat")
	m.vals[0] = 1201.0
	assert.NotEqual(t, m.vals[0], s[0][0], "changing mat should not effect data")
}

func TestRandf32(t *testing.T) {
	t.Helper()
	rows := 31
	cols := 42

	m := RandMatf32(rows, cols)
	for i := 0; i < rows*cols; i++ {
		if m.vals[i] < 0.0 || m.vals[i] >= 1.0 {
			t.Errorf("at index %d, expected [0, 1.0), got %f", i, m.vals[i])
		}
	}
}

func TestReshapef32(t *testing.T) {
	t.Helper()
	rows, cols := 10, 12
	s := make([]float32, 120)
	for i := 0; i < len(s); i++ {
		s[i] = float32(i * 3)
	}
	m := Matf32FromData(s).Reshape(rows, cols)
	assert.Equal(t, rows, m.r, "should be equal")
	assert.Equal(t, cols, m.c, "should be equal")
	for i := 0; i < len(s); i++ {
		assert.Equal(t, s[i], m.vals[i], "should be equal")
	}

	// assert.Panics(t, func() { m.Reshape(rows, rows) }, "should panic")
}

func TestShapef32(t *testing.T) {
	t.Helper()
	m := Newf32(11, 10)
	r, c := m.Shape()
	assert.Equal(t, r, m.r, "should be equal")
	assert.Equal(t, c, m.c, "should be equal")
}

func TestToSlice1Df32(t *testing.T) {
	t.Helper()
	rows, cols := 22, 22
	m := Newf32(rows, cols)
	m.SetAll(1.0)
	v := m.ToSlice1D()
	assert.Equal(t, rows*cols, len(v), "should be equal")
	for i := range v {
		assert.Equal(t, float32(1.0), v[i], "should be equal")
	}
	w := m.ToSlice1Df64()
	assert.Equal(t, rows*cols, len(w), "should be equal")
	for i := range v {
		assert.Equal(t, 1.0, w[i], "should be equal")
	}
}

func TestToSlice2Df32(t *testing.T) {
	t.Helper()
	rows := 13
	cols := 21
	m := Newf32(rows, cols)
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] = float32(i)
	}

	s := m.ToSlice2D()
	assert.Equal(t, m.r, len(s), "should be equal")
	assert.Equal(t, m.c, len(s[0]), "should be equal")
	idx := 0
	for i := range s {
		for j := range s[i] {
			assert.Equal(t, s[i][j], m.vals[idx], "should be equal")
			idx++
		}
	}
	s[0][0] = 1021.0
	assert.NotEqual(t, s[0][0], m.vals[0], "changing data should not effect mat")
	m.vals[0] = 1201.0
	assert.NotEqual(t, m.vals[0], s[0][0], "changing mat should not effect data")

	u := m.ToSlice2Df64()
	assert.Equal(t, m.r, len(u), "should be equal")
	assert.Equal(t, m.c, len(u[0]), "should be equal")
	idx = 0
	for i := range u {
		for j := range u[i] {
			assert.Equal(t, u[i][j], float64(m.vals[idx]), "should be equal")
			idx++
		}
	}
	u[0][0] = 1021.0
	assert.NotEqual(t, u[0][0], m.vals[0], "changing data should not effect mat")
	m.vals[0] = 1201.0
	assert.NotEqual(t, m.vals[0], u[0][0], "changing mat should not effect data")
}

func TestGetf32(t *testing.T) {
	t.Helper()
	rows := 17
	cols := 13
	m := Newf32(rows, cols)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	idx := 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			assert.Equal(t, m.vals[idx], m.Get(i, j), "should be equal")
			idx++
		}
	}
}

func TestMapf32(t *testing.T) {
	t.Helper()
	rows := 132
	cols := 24
	f := func(i *float32) {
		*i = 1.0
		return
	}
	m := Newf32(rows, cols).Map(f)
	for i := 0; i < rows*cols; i++ {
		assert.Equal(t, float32(1.0), m.vals[i], "should be equal")
	}
}

func BenchmarkMapf32(b *testing.B) {
	m := Newf32(17, 31)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	f := func(i *float32) {
		*i = 1.0
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Map(f)
	}
}

func TestSetAllf32(t *testing.T) {
	t.Helper()
	row := 3
	col := 4
	val := 11.0
	m := Newf32(row, col).SetAll(val)
	for i := 0; i < row*col; i++ {
		assert.Equal(t, float32(val), m.vals[i], "should be equal")
	}
}

func TestSetf32(t *testing.T) {
	t.Helper()
	m := Newf32(5)
	m.Set(2, 3, 10.0)
	assert.Equal(t, float32(10.0), m.vals[13], "should be equal")
}

func TestSetColf32(t *testing.T) {
	t.Helper()
	m := Newf32(3, 4)
	m.SetCol(-1, 3.0)
	n := m.Col(-1)
	for i := range n.vals {
		assert.Equal(t, float32(3.0), n.vals[i], "should be equal")
	}
	m.SetCol(-1, []float32{0.0, 0.0, 0.0})
	n = m.Col(-1)
	for i := range n.vals {
		assert.Equal(t, float32(0.0), n.vals[i], "should be equal")
	}
	m.SetCol(1, 3.0)
	n = m.Col(1)
	for i := range n.vals {
		assert.Equal(t, float32(3.0), n.vals[i], "should be equal")
	}
	m.SetCol(1, []float32{0.0, 0.0, 0.0})
	n = m.Col(1)
	for i := range n.vals {
		assert.Equal(t, float32(0.0), n.vals[i], "should be equal")
	}

	// assert.Panics(t, func() { m.SetCol(-5, 2.0) }, "should panic")
	// assert.Panics(t, func() { m.SetCol(5, 2.0) }, "should panic")
	// assert.Panics(t, func() { m.SetCol(-1, []float32{0.0}) }, "should panic")
	// assert.Panics(t, func() { m.SetCol(1, []float32{0.0}) }, "should panic")
	// assert.Panics(t, func() { m.SetCol(-1, 1) }, "should panic")
	// assert.Panics(t, func() { m.SetCol(1, 1) }, "should panic")
}

func TestSetRowf32(t *testing.T) {
	t.Helper()
	m := Newf32(3, 4)
	m.SetRow(-1, 3.0)
	n := m.Row(-1)
	for i := range n.vals {
		assert.Equal(t, float32(3.0), n.vals[i], "should be equal")
	}
	m.SetRow(-1, []float32{0.0, 0.0, 0.0, 0.0})
	n = m.Row(-1)
	for i := range n.vals {
		assert.Equal(t, float32(0.0), n.vals[i], "should be equal")
	}
	m.SetRow(1, 3.0)
	n = m.Row(1)
	for i := range n.vals {
		assert.Equal(t, float32(3.0), n.vals[i], "should be equal")
	}
	m.SetRow(1, []float32{0.0, 0.0, 0.0, 0.0})
	n = m.Row(1)
	for i := range n.vals {
		assert.Equal(t, float32(0.0), n.vals[i], "should be equal")
	}

	// assert.Panics(t, func() { m.SetRow(-5, 2.0) }, "should panic")
	// assert.Panics(t, func() { m.SetRow(5, 2.0) }, "should panic")
	// assert.Panics(t, func() { m.SetRow(-1, []float32{0.0}) }, "should panic")
	// assert.Panics(t, func() { m.SetRow(1, []float32{0.0}) }, "should panic")
	// assert.Panics(t, func() { m.SetRow(-1, 1) }, "should panic")
	// assert.Panics(t, func() { m.SetRow(1, 1) }, "should panic")
}

func TestColf32(t *testing.T) {
	t.Helper()
	row := 3
	col := 4
	m := Newf32(row, col)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	for i := 0; i < col; i++ {
		got := m.Col(i)
		for j := 0; j < row; j++ {
			assert.Equal(t, m.vals[j*m.c+i], got.vals[j], "should be equal")
		}
	}
	for i := col; i < 0; i-- {
		got := m.Col(-i)
		for j := 0; j < row; j++ {
			assert.Equal(t, m.vals[j*m.c+(row-i)], got.vals[j], "should be equal")
		}
	}
}

func BenchmarkColf32(b *testing.B) {
	m := Newf32(17, 31)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Col(21)
	}
}

func TestRowf32(t *testing.T) {
	t.Helper()
	row := 3
	col := 4
	m := Newf32(row, col)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	idx := 0
	for i := 0; i < row; i++ {
		got := m.Row(i)
		for j := 0; j < col; j++ {
			assert.Equal(t, m.vals[idx], got.vals[j], "should be equal")
			idx++
		}
	}
}

func BenchmarkRowf32(b *testing.B) {
	m := Newf32(17, 31)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Row(11)
	}
}

func TestMinf32(t *testing.T) {
	t.Helper()
	m := Newf32(3, 4)
	m.Set(2, 1, -100.0)
	_, minVal := m.Min()
	assert.Equal(t, -float32(100.0), minVal, "should be equal")
	idx, minVal := m.Min(0, 2)
	assert.Equal(t, -float32(100.0), minVal, "should be equal")
	assert.Equal(t, 1, idx, "should be equal")
	idx, minVal = m.Min(1, 1)
	assert.Equal(t, -float32(100.0), minVal, "should be equal")
	assert.Equal(t, 2, idx, "should be equal")
}

func TestMaxf32(t *testing.T) {
	t.Helper()
	m := Newf32(3, 4)
	m.Set(2, 1, 100.0)
	_, maxVal := m.Max()
	assert.Equal(t, float32(100.0), maxVal, "should be equal")
	idx, maxVal := m.Max(0, 2)
	assert.Equal(t, float32(100.0), maxVal, "should be equal")
	assert.Equal(t, 1, idx, "should be equal")
	idx, maxVal = m.Max(1, 1)
	assert.Equal(t, float32(100.0), maxVal, "should be equal")
	assert.Equal(t, 2, idx, "should be equal")
}

func TestEqualsf32(t *testing.T) {
	t.Helper()
	m := Newf32(13, 12)
	if !m.Equals(m) {
		t.Errorf("m is not equal itself")
	}
}

func TestCopyf32(t *testing.T) {
	t.Helper()
	rows, cols := 17, 13
	m := Newf32(rows, cols)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	n := m.Copy()
	for i := 0; i < rows*cols; i++ {
		assert.Equal(t, m.vals[i], n.vals[i], "should be equal")
	}
}

func TestTf32(t *testing.T) {
	t.Helper()
	m := Newf32(12, 3)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	n := m.Copy()
	m.T()
	p := m.ToSlice2D()
	q := n.ToSlice2D()
	for i := 0; i < m.r; i++ {
		for j := 0; j < m.c; j++ {
			assert.Equal(t, p[i][j], q[j][i], "should be equal")
		}
	}
	res := m.Dot(n)
	res2 := res.Copy()
	res.T()
	assert.True(t, res2.Equals(res), "should be equal")
}

func BenchmarkTf32(b *testing.B) {
	m := Newf32(11, 21)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.T()
	}
}

func BenchmarkTf32Vanilla(b *testing.B) {
	m := make([][]float32, 11)
	for i := range m {
		m[i] = make([]float32, 21)
	}
	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		n := make([][]float32, len(m[0]))
		for i := range n {
			n[i] = make([]float32, len(m))
		}
		for i := range m {
			for j := range m[i] {
				n[j][i] = m[i][j]
			}
		}
	}
}

func TestAllf32(t *testing.T) {
	t.Helper()
	m := Newf32(100, 21)
	for i := range m.vals {
		m.vals[i] = float32(i + 1)
	}
	positive := func(i *float32) bool {
		return *i >= 0
	}
	assert.True(t, m.All(positive), "All should be > 0")
	isOne := func(i *float32) bool {
		return *i == 1.0
	}
	m.SetAll(1.0)
	assert.True(t, m.All(isOne), "All should be 1.0s")
}

func TestAnyf32(t *testing.T) {
	t.Helper()
	m := Newf32(100, 21)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	positive := func(i *float32) bool {
		return *i >= 0
	}
	negative := func(i *float32) bool {
		return *i < 0
	}
	assert.False(t, m.Any(negative), "should have no negatives")
	assert.True(t, m.Any(positive), "should have positives")
}

func TestMulf32(t *testing.T) {
	t.Helper()
	rows, cols := 13, 90
	m := Newf32(rows, cols)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	n := m.Copy()
	m.Mul(m)
	for i := 0; i < rows*cols; i++ {
		assert.Equal(t, n.vals[i]*n.vals[i], m.vals[i], "should be equal")
	}
}

func BenchmarkMulf32(b *testing.B) {
	n := Newf32(10)
	for i := range n.vals {
		n.vals[i] = float32(i)
	}
	m := Newf32(10)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Mul(n)
	}
}

func BenchmarkMulVanillaf32(b *testing.B) {
	n := make([][]float32, 10)
	m := make([][]float32, 10)
	for i := range n {
		n[i] = make([]float32, 10)
		m[i] = make([]float32, 10)
		for j := range n[i] {
			n[i][j] = float32(i*10 + j)
			m[i][j] = float32(i*10 + j)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range n {
			for k := range n[j] {
				m[j][k] = m[j][k] * n[j][k]
			}
		}
	}
}

func TestAddf32(t *testing.T) {
	t.Helper()
	rows, cols := 13, 90
	m := Newf32(rows, cols)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	n := m.Copy()
	m.Add(m)
	for i := 0; i < rows*cols; i++ {
		assert.Equal(t, n.vals[i]+n.vals[i], m.vals[i], "should be equal")
	}
}

func BenchmarkAddf32(b *testing.B) {
	n := Newf32(10)
	for i := range n.vals {
		n.vals[i] = float32(i)
	}
	m := Newf32(10)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Add(n)
	}
}

func TestSubf32(t *testing.T) {
	t.Helper()
	rows, cols := 13, 90
	m := Newf32(rows, cols)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	m.Sub(m)
	for i := 0; i < rows*cols; i++ {
		assert.Equal(t, float32(0.0), m.vals[i], "should be equal")
	}
}

func BenchmarkSubf32(b *testing.B) {
	n := Newf32(10)
	for i := range n.vals {
		n.vals[i] = float32(i)
	}
	m := Newf32(10)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Sub(n)
	}
}

func TestDivf32(t *testing.T) {
	t.Helper()
	rows, cols := 13, 90
	m := Newf32(rows, cols)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	m.vals[0] = float32(1.0)
	m.Div(m)
	for i := 0; i < rows*cols; i++ {
		assert.Equal(t, float32(1.0), m.vals[i], "should be equal")
	}
}

func BenchmarkDivf32(b *testing.B) {
	n := Newf32(10)
	for i := range n.vals {
		n.vals[i] = float32(i)
	}
	m := Newf32(10)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Div(n)
	}
}

func TestSumf32(t *testing.T) {
	t.Helper()
	row := 12
	col := 17
	m := Newf32(row, col).SetAll(1.0)
	for i := 0; i < row; i++ {
		assert.Equal(t, float32(col), m.Sum(0, i), "should be equal")
	}
	for i := 0; i < col; i++ {
		assert.Equal(t, float32(row), m.Sum(1, i), "should be equal")
	}
}

func TestAvgf32(t *testing.T) {
	t.Helper()
	row := 12
	col := 17
	m := Newf32(row, col).SetAll(1.0)
	for i := 0; i < row; i++ {
		assert.Equal(t, float32(1.0), m.Avg(0, i), "should be equal")
	}
	for i := 0; i < col; i++ {
		assert.Equal(t, float32(1.0), m.Avg(1, i), "should be equal")
	}
}

func TestPrdf32(t *testing.T) {
	t.Helper()
	row := 12
	col := 17
	m := Newf32(row, col).SetAll(1.0)
	for i := 0; i < row; i++ {
		assert.Equal(t, float32(1.0), m.Prd(0, i), "should be equal")
	}
	for i := 0; i < col; i++ {
		assert.Equal(t, float32(1.0), m.Prd(1, i), "should be equal")
	}
}

func TestStdf32(t *testing.T) {
	t.Helper()
	row := 12
	col := 17
	m := Newf32(row, col).SetAll(1.0)
	for i := 0; i < row; i++ {
		assert.Equal(t, float32(0.0), m.Std(0, i), "should be equal")
	}
	for i := 0; i < col; i++ {
		assert.Equal(t, float32(0.0), m.Std(1, i), "should be equal")
	}
}

func TestDotf32(t *testing.T) {
	t.Helper()
	var (
		row = 10
		col = 4
	)
	m := Newf32(row, col)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	n := Newf32(col, row)
	for i := range n.vals {
		n.vals[i] = float32(i)
	}
	o := m.Dot(n)
	assert.Equal(t, row, o.r, "should be equal")
	assert.Equal(t, row, o.c, "should be equal")
	p := Newf32(row, row)
	q := o.Dot(p)
	for i := 0; i < row*row; i++ {
		assert.Equal(t, float32(0.0), q.vals[i], "should be zero")
	}
	bra := Newf32(3, 1).SetAll(2.0)
	bra.SetAll(2.0)
	ket := Newf32(1, 3).SetAll(3.0)
	bracket := bra.Dot(ket)
	v := bracket.ToSlice1D()
	assert.Equal(t, 9, len(v))
	for i := range v {
		assert.Equal(t, float32(6.0), v[i])
	}
	bracket = ket.Dot(bra)
	v = bracket.ToSlice1D()
	assert.Equal(t, 1, len(v))
	for i := range v {
		assert.Equal(t, float32(18.0), v[i])
	}
	x := Newf32(13)
	y := Eyef32(13)
	z := x.Dot(y)
	assert.True(t, x.Equals(z), "A times I should equal A")
}

func BenchmarkDotf32(b *testing.B) {
	m := Newf32(10)
	n := Newf32(10)
	for i := range m.vals {
		m.vals[i] = float32(i + i)
	}
	for i := range n.vals {
		n.vals[i] = float32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Dot(n)
	}
}

func BenchmarkDotf32Vanilla(b *testing.B) {
	m := make([][]float32, 10)
	n := make([][]float32, 10)
	for i := range m {
		m[i] = make([]float32, 10)
		n[i] = make([]float32, 10)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o := make([][]float32, 10)
		for j := range o {
			o[j] = make([]float32, 10)
		}
		for j := range m {
			for k := range n[j] {
				for p := range m[j] {
					o[j][k] += m[j][p] * n[p][k]
				}
			}
		}
	}
}

func TestAppendColf32(t *testing.T) {
	t.Helper()
	var (
		row = 10
		col = 4
	)
	m := Newf32(row, col)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	v := make([]float32, row)
	m.AppendCol(v)
	assert.Equal(t, col+1, m.c, "should have one more column")
	m.AppendCol(v)
	assert.Equal(t, col+2, m.c, "should have two more columns")
	m.AppendCol(v)
	assert.Equal(t, col+3, m.c, "should have three more columns")
}

func TestAppendRowf32(t *testing.T) {
	t.Helper()
	var (
		row = 3
		col = 4
	)
	m := Newf32(row, col)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	v := make([]float32, col)
	for i := range v {
		v[i] = float32(i * i * i)
	}
	m.AppendRow(v)
	assert.Equal(t, row+1, m.r, "should have one more row")
	m.AppendRow(v)
	assert.Equal(t, row+2, m.r, "should have two more rows")
	m.AppendRow(v)
	assert.Equal(t, row+3, m.r, "should have three more rows")
}

func TestConcatf32(t *testing.T) {
	t.Helper()
	var (
		row = 10
		col = 4
	)
	m := Newf32(row, col)
	for i := range m.vals {
		m.vals[i] = float32(i)
	}
	n := Newf32(row, row)
	for i := range n.vals {
		n.vals[i] = float32(i)
	}
	m.Concat(n)
	if m.c != row+col {
		t.Errorf("Expected number of cols to be %d, but got %d", row+col, m.c)
	}
	idx1 := 0
	idx2 := 0
	for i := 0; i < row; i++ {
		for j := 0; j < col+row; j++ {
			if j < col {
				assert.Equal(t, float32(idx1), m.vals[i*m.c+j], "should be equal")
				idx1++
				continue
			}
			assert.Equal(t, float32(idx2), m.vals[i*m.c+j], "should be equal")
			idx2++
		}
	}
}
