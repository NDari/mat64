package mat64

import "testing"

func TestFunctions(t *testing.T) {
	m := NewMat(10, 12)
	for i := range m.vals {
		m.vals[i] = float64(i * 2)
	}
	m.vals[0] = 2.0
	if m.Any(Negative) {
		t.Errorf("found negatives")
	}
	if !m.All(Positive) {
		t.Errorf("Some are not positive")
	}
	if m.Any(Odd) {
		t.Errorf("Some are odd")
	}
	if !m.All(Even) {
		t.Errorf("Some are not even")
	}
	m.vals[0] = 0.0
	m.Foreach(Square)
	for i := range m.vals {
		if m.vals[i] != float64(i*i*4) {
			t.Errorf("At %d, expected %f, got %f", i, float64(i*i*4), m.vals[i])
		}
	}
}
