package poly

import "testing"

func TestEvalAtZero(t *testing.T) {
	cases := []struct {
		p    []byte
		want byte
	}{
		{[]byte{0x05}, 0x05},
		{[]byte{0x53, 0xCA, 0x01}, 0x01},
		{[]byte{0xFF, 0x00, 0xAB}, 0xAB},
	}
	for _, c := range cases {
		got := Eval(c.p, 0x00)
		if got != c.want {
			t.Errorf("Eval(%v, 0x00) = 0x%02x, want 0x%02x", c.p, got, c.want)
		}
	}
}

func TestEvalConstant(t *testing.T) {
	got := Eval([]byte{0x05}, 0x02)
	want := byte(0x05)
	if got != want {
		t.Errorf("Eval([0x05], 0x02) = 0x%02x, want 0x%02x", got, want)
	}
}

func TestEvalAtOne(t *testing.T) {
	p := []byte{0x53, 0xCA, 0x01}
	got := Eval(p, 0x01)
	want := byte(0x53 ^ 0xCA ^ 0x01)
	if got != want {
		t.Errorf("Eval at 1: got 0x%02x, want 0x%02x", got, want)
	}
}

func TestEvalIdentity(t *testing.T) {
	for a := 0; a < 256; a++ {
		got := Eval([]byte{0x01, 0x00}, byte(a))
		if got != byte(a) {
			t.Fatalf("Eval([0x01,0x00], 0x%02x) = 0x%02x, want 0x%02x", a, got, a)
		}
	}
}

func TestEvalHandComputed(t *testing.T) {
	got := Eval([]byte{0x53, 0xCA, 0x01}, 0x02)
	want := byte(0xD9)
	if got != want {
		t.Errorf("Eval([0x53,0xCA,0x01], 0x02) = 0x%02x, want 0x%02x", got, want)
	}
}
