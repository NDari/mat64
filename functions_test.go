package matrix

import "testing"

func TestFunctions(t *testing.T) {
	m := Newf64(10, 12)
	for i := range m.vals {
		m.vals[i] = float64(i * 2)
	}
	m.vals[0] = 2.0
	if m.Any(Negativef64) {
		t.Errorf("found negatives")
	}
	if !m.All(Positivef64) {
		t.Errorf("Some are not positive")
	}
	if m.Any(Oddf64) {
		t.Errorf("Some are odd")
	}
	if !m.All(Evenf64) {
		t.Errorf("Some are not even")
	}
}
