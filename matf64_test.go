package matrix

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewf64(t *testing.T) {
	t.Helper()
	rows := 13
	cols := 7
	m := Newf64()
	assert.Equal(t, 0, m.r, "should be zero")
	assert.Equal(t, 0, m.c, "should be zero")
	assert.NotNil(t, m.vals, "should not be nil")
	assert.Equal(t, 0, len(m.vals), "should be zero")
	assert.Equal(t, 0, cap(m.vals), "should be zero")

	m = Newf64(rows)
	assert.Equal(t, rows, m.r, "should be equal")
	assert.Equal(t, rows, m.c, "should be equal")
	assert.NotNil(t, m.vals, "should not be nil")
	assert.Equal(t, rows*rows, len(m.vals), "should be equal")
	assert.Equal(t, 2*rows*rows, cap(m.vals), "should have twice the capacity")

	m = Newf64(rows, cols)
	assert.Equal(t, rows, m.r, "should be equal")
	assert.Equal(t, cols, m.c, "should be equal")
	assert.NotNil(t, m.vals, "should not be nil")
	assert.Equal(t, rows*cols, len(m.vals), "should be equal")
	assert.Equal(t, 2*rows*cols, cap(m.vals), "should have twice the capacity")
}

func TestMatf64FromData(t *testing.T) {
	t.Helper()
	rows := 50
	cols := 2

	v := make([]float64, rows*cols)
	for i := range v {
		v[i] = float64(i * i)
	}

	m := Matf64FromData(v)
	assert.Equal(t, 1, m.r, "should have one row")
	assert.Equal(t, len(v), len(m.vals), "should have the same # of elements")
	for i := range v {
		assert.Equal(t, v[i], m.vals[i], "should be equal")
	}
	v[0] = 1321.0
	assert.NotEqual(t, v[0], m.vals[0], "changing data should not effect mat")
	m.vals[0] = 1201.0
	assert.NotEqual(t, m.vals[0], v[0], "changing mat should not effect data")

	v[0] = 0.0
	m = Matf64FromData(v, rows*cols)
	assert.Equal(t, rows*cols, m.r, "should be equal")
	assert.Equal(t, 1, m.c, "should have one col")
	assert.Equal(t, len(v), len(m.vals), "should have the same # of elements")
	for i := range v {
		assert.Equal(t, v[i], m.vals[i], "should be equal")
	}
	v[0] = 1321.0
	assert.NotEqual(t, v[0], m.vals[0], "changing data should not effect mat")
	m.vals[0] = 1201.0
	assert.NotEqual(t, m.vals[0], v[0], "changing mat should not effect data")

	v[0] = 0.0
	m = Matf64FromData(v, rows, cols)
	assert.Equal(t, rows, m.r, "should be equal")
	assert.Equal(t, cols, m.c, "should be equal")
	assert.Equal(t, len(v), len(m.vals), "should have the same # of elements")
	for i := range v {
		assert.Equal(t, v[i], m.vals[i], "should be equal")
	}
	v[0] = 1321.0
	assert.NotEqual(t, v[0], m.vals[0], "changing data should not effect mat")
	m.vals[0] = 1201.0
	assert.NotEqual(t, m.vals[0], v[0], "changing mat should not effect data")

	s := make([][]float64, rows)
	for i := range s {
		s[i] = make([]float64, cols)
	}
	for i := range s {
		for j := range s[i] {
			s[i][j] = float64(i + j)
		}
	}
	m = Matf64FromData(s)
	assert.Equal(t, rows*cols, len(m.vals), "should be equal")
	assert.Equal(t, 2*rows*cols, cap(m.vals), "should be equal")
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

	s[0][0] = 0.0
	m = Matf64FromData(s, 10)
	assert.Equal(t, 10, m.r, "should be equal")
	assert.Equal(t, 10, m.c, "should be equal")
	assert.Equal(t, 100, len(m.vals), "should be equal")
	assert.Equal(t, 200, cap(m.vals), "should be equal")
	idx = 0
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

	s[0][0] = 0.0
	m = Matf64FromData(s, rows, cols)
	assert.Equal(t, rows, m.r, "should be equal")
	assert.Equal(t, cols, m.c, "should be equal")
	assert.Equal(t, rows*cols, len(m.vals), "should be equal")
	assert.Equal(t, 2*rows*cols, cap(m.vals), "should be equal")
	idx = 0
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

func TestMatf64FromCSV(t *testing.T) {
	t.Helper()
	rows := 3
	cols := 4

	filename := "non-exitant-file"

	filename = "test.csv"
	str := "1.0,1.0,2.0,3.0\n5.0,8.0,13.0,21.0\n34.0,55.0,89.0,144.0"
	if _, err := os.Stat(filename); err == nil {
		err = os.Remove(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write([]byte(str))
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := Matf64FromCSV(filename)
	assert.Equal(t, rows*cols, len(m.vals), "should be equal")
	assert.Equal(t, 1.0, m.vals[0], "should be equal")
	assert.Equal(t, 1.0, m.vals[1], "should be equal")
	for i := 2; i < m.r*m.c; i++ {
		assert.Equal(t, (m.vals[i-1] + m.vals[i-2]), m.vals[i], "should be equal")
	}
	err = os.Remove(filename)
	if err != nil {
		log.Fatal(err)
	}
}

func TestRandf64(t *testing.T) {
	t.Helper()
	rows := 31
	cols := 42

	m := RandMatf64(rows, cols)
	for i := 0; i < rows*cols; i++ {
		if m.vals[i] < 0.0 || m.vals[i] >= 1.0 {
			t.Errorf("at index %d, expected [0, 1.0), got %f", i, m.vals[i])
		}
	}
	m = RandMatf64(rows, cols, 100.0)
	for i := 0; i < rows*cols; i++ {
		if m.vals[i] < 0.0 || m.vals[i] >= 100.0 {
			t.Errorf("at index %d, expected [0, 100.0), got %f", i, m.vals[i])
		}
	}
	m = RandMatf64(rows, cols, -12.0, 2.0)
	for i := 0; i < rows*cols; i++ {
		if m.vals[i] < -12.0 || m.vals[i] >= 2.0 {
			t.Errorf("at index %d, expected [-12.0, 2.0), got %f", i, m.vals[i])
		}
	}
}

func TestReshapef64(t *testing.T) {
	t.Helper()
	rows, cols := 10, 12
	s := make([]float64, 120)
	for i := 0; i < len(s); i++ {
		s[i] = float64(i * 3)
	}
	m := Matf64FromData(s).Reshape(rows, cols)
	assert.Equal(t, rows, m.r, "should be equal")
	assert.Equal(t, cols, m.c, "should be equal")
	for i := 0; i < len(s); i++ {
		assert.Equal(t, s[i], m.vals[i], "should be equal")
	}
}

func TestShapef64(t *testing.T) {
	t.Helper()
	m := Newf64(11, 10)
	r, c := m.Shape()
	assert.Equal(t, r, m.r, "should be equal")
	assert.Equal(t, c, m.c, "should be equal")
}

func TestValsf64(t *testing.T) {
	t.Helper()
	rows, cols := 22, 22
	m := Newf64(rows, cols)
	m.SetAll(1.0)
	assert.Equal(t, rows*cols, len(m.vals), "should be equal")
	for i := range m.vals {
		assert.Equal(t, 1.0, m.vals[i], "should be equal")
	}
}

func TestToSlicef64(t *testing.T) {
	t.Helper()
	rows := 13
	cols := 21
	m := Newf64(rows, cols)
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] = float64(i)
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
}

func TestToCSVf64(t *testing.T) {
	t.Helper()
	m := Newf64(23, 17)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	filename := "tocsv_test.csv"
	m.ToCSV(filename)
	n := Matf64FromCSV(filename)
	if !n.Equals(m) {
		t.Errorf("m and n are not equal")
	}
	os.Remove(filename)
}

func TestGetf64(t *testing.T) {
	t.Helper()
	rows := 17
	cols := 13
	m := Newf64(rows, cols)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	idx := 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			assert.Equal(t, m.vals[idx], m.Get(i, j), "should be equal")
			idx++
		}
	}
}

func TestMapf64(t *testing.T) {
	t.Helper()
	rows := 132
	cols := 24
	f := func(i *float64) {
		*i = 1.0
		return
	}
	m := Newf64(rows, cols).Map(f)
	for i := 0; i < rows*cols; i++ {
		assert.Equal(t, 1.0, m.vals[i], "should be equal")
	}
}

func BenchmarkMapf64(b *testing.B) {
	m := Newf64(17, 31)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	f := func(i *float64) {
		*i = 1.0
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Map(f)
	}
}

func TestSetAllf64(t *testing.T) {
	t.Helper()
	row := 3
	col := 4
	val := 11.0
	m := Newf64(row, col).SetAll(val)
	for i := 0; i < row*col; i++ {
		assert.Equal(t, val, m.vals[i], "should be equal")
	}
}

func TestSetf64(t *testing.T) {
	t.Helper()
	m := Newf64(5)
	m.Set(2, 3, 10.0)
	assert.Equal(t, 10.0, m.vals[13], "should be equal")
}

func TestSetColf64(t *testing.T) {
	t.Helper()
	m := Newf64(3, 4)
	m.SetCol(-1, 3.0)
	n := m.Col(-1)
	for i := range n.vals {
		assert.Equal(t, 3.0, n.vals[i], "should be equal")
	}
	m.SetCol(-1, []float64{0.0, 0.0, 0.0})
	n = m.Col(-1)
	for i := range n.vals {
		assert.Equal(t, 0.0, n.vals[i], "should be equal")
	}
	m.SetCol(1, 3.0)
	n = m.Col(1)
	for i := range n.vals {
		assert.Equal(t, 3.0, n.vals[i], "should be equal")
	}
	m.SetCol(1, []float64{0.0, 0.0, 0.0})
	n = m.Col(1)
	for i := range n.vals {
		assert.Equal(t, 0.0, n.vals[i], "should be equal")
	}
}

func TestSetRowf64(t *testing.T) {
	t.Helper()
	m := Newf64(3, 4)
	m.SetRow(-1, 3.0)
	n := m.Row(-1)
	for i := range n.vals {
		assert.Equal(t, 3.0, n.vals[i], "should be equal")
	}
	m.SetRow(-1, []float64{0.0, 0.0, 0.0, 0.0})
	n = m.Row(-1)
	for i := range n.vals {
		assert.Equal(t, 0.0, n.vals[i], "should be equal")
	}
	m.SetRow(1, 3.0)
	n = m.Row(1)
	for i := range n.vals {
		assert.Equal(t, 3.0, n.vals[i], "should be equal")
	}
	m.SetRow(1, []float64{0.0, 0.0, 0.0, 0.0})
	n = m.Row(1)
	for i := range n.vals {
		assert.Equal(t, 0.0, n.vals[i], "should be equal")
	}
}

func TestColf64(t *testing.T) {
	t.Helper()
	row := 3
	col := 4
	m := Newf64(row, col)
	for i := range m.vals {
		m.vals[i] = float64(i)
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

func BenchmarkColf64(b *testing.B) {
	m := Newf64(17, 31)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Col(21)
	}
}

func TestRowf64(t *testing.T) {
	t.Helper()
	row := 3
	col := 4
	m := Newf64(row, col)
	for i := range m.vals {
		m.vals[i] = float64(i)
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

func BenchmarkRowf64(b *testing.B) {
	m := Newf64(17, 31)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Row(11)
	}
}

func TestMinf64(t *testing.T) {
	t.Helper()
	m := Newf64(3, 4)
	m.Set(2, 1, -100.0)
	_, minVal := m.Min()
	assert.Equal(t, -100.0, minVal, "should be equal")
	idx, minVal := m.Min(0, 2)
	assert.Equal(t, -100.0, minVal, "should be equal")
	assert.Equal(t, 1, idx, "should be equal")
	idx, minVal = m.Min(1, 1)
	assert.Equal(t, -100.0, minVal, "should be equal")
	assert.Equal(t, 2, idx, "should be equal")
}

func TestMaxf64(t *testing.T) {
	t.Helper()
	m := Newf64(3, 4)
	m.Set(2, 1, 100.0)
	_, maxVal := m.Max()
	assert.Equal(t, 100.0, maxVal, "should be equal")
	idx, maxVal := m.Max(0, 2)
	assert.Equal(t, 100.0, maxVal, "should be equal")
	assert.Equal(t, 1, idx, "should be equal")
	idx, maxVal = m.Max(1, 1)
	assert.Equal(t, 100.0, maxVal, "should be equal")
	assert.Equal(t, 2, idx, "should be equal")
}

func TestEqualsf64(t *testing.T) {
	t.Helper()
	m := Newf64(13, 12)
	if !m.Equals(m) {
		t.Errorf("m is not equal itself")
	}
}

func TestCopyf64(t *testing.T) {
	t.Helper()
	rows, cols := 17, 13
	m := Newf64(rows, cols)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	n := m.Copy()
	for i := 0; i < rows*cols; i++ {
		assert.Equal(t, m.vals[i], n.vals[i], "should be equal")
	}
}

func TestTf64(t *testing.T) {
	t.Helper()
	m := Newf64(12, 3)
	for i := range m.vals {
		m.vals[i] = float64(i)
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
	res := m.Copy()
	res.T().T()
	assert.True(t, m.Equals(res), "should be equal")
}

func BenchmarkTf64(b *testing.B) {
	m := Newf64(11, 21)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.T()
	}
}

func BenchmarkTf64Vanilla(b *testing.B) {
	m := make([][]float64, 11)
	for i := range m {
		m[i] = make([]float64, 21)
	}
	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		n := make([][]float64, len(m[0]))
		for i := range n {
			n[i] = make([]float64, len(m))
		}
		for i := range m {
			for j := range m[i] {
				n[j][i] = m[i][j]
			}
		}
	}
}

func TestAllf64(t *testing.T) {
	t.Helper()
	m := Newf64(100, 21)
	for i := range m.vals {
		m.vals[i] = float64(i + 1)
	}
	positive := func(i *float64) bool {
		return *i >= 0.0
	}
	assert.True(t, m.All(positive), "All should be > 0")
	isOne := func(i *float64) bool {
		return *i == 1.0
	}
	m.SetAll(1.0)
	assert.True(t, m.All(isOne), "All should be 1.0s")
}

func TestAnyf64(t *testing.T) {
	t.Helper()
	m := Newf64(100, 21)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	positive := func(i *float64) bool {
		return *i >= 0.0
	}
	negative := func(i *float64) bool {
		return *i < 0.0
	}
	assert.False(t, m.Any(negative), "should have no negatives")
	assert.True(t, m.Any(positive), "should have positives")
}

func TestMulf64(t *testing.T) {
	t.Helper()
	rows, cols := 13, 90
	m := Newf64(rows, cols)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	n := m.Copy()
	m.Mul(m)
	for i := 0; i < rows*cols; i++ {
		assert.Equal(t, n.vals[i]*n.vals[i], m.vals[i], "should be equal")
	}
}

func BenchmarkMulf64(b *testing.B) {
	n := Newf64(10)
	for i := range n.vals {
		n.vals[i] = float64(i)
	}
	m := Newf64(10)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Mul(n)
	}
}

func BenchmarkMulVanillaf64(b *testing.B) {
	n := make([][]float64, 10)
	m := make([][]float64, 10)
	for i := range n {
		n[i] = make([]float64, 10)
		m[i] = make([]float64, 10)
		for j := range n[i] {
			n[i][j] = float64(i*10 + j)
			m[i][j] = float64(i*10 + j)
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

func TestAddf64(t *testing.T) {
	t.Helper()
	rows, cols := 13, 90
	m := Newf64(rows, cols)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	n := m.Copy()
	m.Add(m)
	for i := 0; i < rows*cols; i++ {
		assert.Equal(t, n.vals[i]+n.vals[i], m.vals[i], "should be equal")
	}
}

func BenchmarkAddf64(b *testing.B) {
	n := Newf64(10)
	for i := range n.vals {
		n.vals[i] = float64(i)
	}
	m := Newf64(10)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Add(n)
	}
}

func TestSubf64(t *testing.T) {
	t.Helper()
	rows, cols := 13, 90
	m := Newf64(rows, cols)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	m.Sub(m)
	for i := 0; i < rows*cols; i++ {
		assert.Equal(t, 0.0, m.vals[i], "should be equal")
	}
}

func BenchmarkSubf64(b *testing.B) {
	n := Newf64(10)
	for i := range n.vals {
		n.vals[i] = float64(i)
	}
	m := Newf64(10)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Sub(n)
	}
}

func TestDivf64(t *testing.T) {
	t.Helper()
	rows, cols := 13, 90
	m := Newf64(rows, cols)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	m.vals[0] = 1.0
	m.Div(m)
	for i := 0; i < rows*cols; i++ {
		assert.Equal(t, 1.0, m.vals[i], "should be equal")
	}
}

func BenchmarkDivf64(b *testing.B) {
	n := Newf64(10)
	for i := range n.vals {
		n.vals[i] = float64(i)
	}
	m := Newf64(10)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Div(n)
	}
}

func TestSumf64(t *testing.T) {
	t.Helper()
	row := 12
	col := 17
	m := Newf64(row, col).SetAll(1.0)
	for i := 0; i < row; i++ {
		assert.Equal(t, float64(col), m.Sum(0, i), "should be equal")
	}
	for i := 0; i < col; i++ {
		assert.Equal(t, float64(row), m.Sum(1, i), "should be equal")
	}
}

func TestAvgf64(t *testing.T) {
	t.Helper()
	row := 12
	col := 17
	m := Newf64(row, col).SetAll(1.0)
	for i := 0; i < row; i++ {
		assert.Equal(t, 1.0, m.Avg(0, i), "should be equal")
	}
	for i := 0; i < col; i++ {
		assert.Equal(t, 1.0, m.Avg(1, i), "should be equal")
	}
}

func TestPrdf64(t *testing.T) {
	t.Helper()
	row := 12
	col := 17
	m := Newf64(row, col).SetAll(1.0)
	for i := 0; i < row; i++ {
		assert.Equal(t, 1.0, m.Prd(0, i), "should be equal")
	}
	for i := 0; i < col; i++ {
		assert.Equal(t, 1.0, m.Prd(1, i), "should be equal")
	}
}

func TestStdf64(t *testing.T) {
	t.Helper()
	row := 12
	col := 17
	m := Newf64(row, col).SetAll(1.0)
	for i := 0; i < row; i++ {
		assert.Equal(t, 0.0, m.Std(0, i), "should be equal")
	}
	for i := 0; i < col; i++ {
		assert.Equal(t, 0.0, m.Std(1, i), "should be equal")
	}
}

func TestDotf64(t *testing.T) {
	t.Helper()
	var (
		row = 10
		col = 4
	)
	m := Newf64(row, col)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	n := Newf64(col, row)
	for i := range n.vals {
		n.vals[i] = float64(i)
	}
	o := m.Dot(n)
	assert.Equal(t, row, o.r, "should be equal")
	assert.Equal(t, row, o.c, "should be equal")
	p := Newf64(row, row)
	q := o.Dot(p)
	for i := 0; i < row*row; i++ {
		assert.Equal(t, 0.0, q.vals[i], "should be zero")
	}
	bra := Newf64(3, 1).SetAll(2.0)
	ket := Newf64(1, 3).SetAll(3.0)
	bracket := bra.Dot(ket)
	v := bracket.ToSlice1D()
	assert.Equal(t, 9, len(v))
	for i := range v {
		assert.Equal(t, 6.0, v[i])
	}
	bra = Newf64(3, 1).SetAll(2.0)
	ket = Newf64(1, 3).SetAll(3.0)
	bracket = ket.Dot(bra)
	v = bracket.ToSlice1D()
	assert.Equal(t, 1, len(v))
	for i := range v {
		assert.Equal(t, 18.0, v[i])
	}
	x := Newf64(13)
	y := If64(13)
	x1 := x.Dot(y)
	assert.True(t, x1.Equals(x), "A times I should equal A")
}

func BenchmarkDotf64(b *testing.B) {
	m := Newf64(10)
	n := Newf64(10)
	for i := range m.vals {
		m.vals[i] = float64(i + i)
	}
	for i := range n.vals {
		n.vals[i] = float64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Dot(n)
	}
}

func BenchmarkDotf64Vanilla(b *testing.B) {
	m := make([][]float64, 10)
	n := make([][]float64, 10)
	for i := range m {
		m[i] = make([]float64, 10)
		n[i] = make([]float64, 10)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o := make([][]float64, 10)
		for j := range o {
			o[j] = make([]float64, 10)
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

func TestAppendColf64(t *testing.T) {
	t.Helper()
	var (
		row = 10
		col = 4
	)
	m := Newf64(row, col)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	v := make([]float64, row)
	m.AppendCol(v)
	assert.Equal(t, col+1, m.c, "should have one more column")
	m.AppendCol(v)
	assert.Equal(t, col+2, m.c, "should have two more columns")
	m.AppendCol(v)
	assert.Equal(t, col+3, m.c, "should have three more columns")
}

func TestAppendRowf64(t *testing.T) {
	t.Helper()
	var (
		row = 3
		col = 4
	)
	m := Newf64(row, col)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	v := make([]float64, col)
	for i := range v {
		v[i] = float64(i * i * i)
	}
	m.AppendRow(v)
	assert.Equal(t, row+1, m.r, "should have one more row")
	m.AppendRow(v)
	assert.Equal(t, row+2, m.r, "should have two more rows")
	m.AppendRow(v)
	assert.Equal(t, row+3, m.r, "should have three more rows")
}

func TestConcatf64(t *testing.T) {
	t.Helper()
	var (
		row = 10
		col = 4
	)
	m := Newf64(row, col)
	for i := range m.vals {
		m.vals[i] = float64(i)
	}
	n := Newf64(row, row)
	for i := range n.vals {
		n.vals[i] = float64(i)
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
				assert.Equal(t, float64(idx1), m.vals[i*m.c+j], "should be equal")
				idx1++
				continue
			}
			assert.Equal(t, float64(idx2), m.vals[i*m.c+j], "should be equal")
			idx2++
		}
	}
}
