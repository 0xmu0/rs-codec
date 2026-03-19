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

func TestAdd(t *testing.T) {
	cases := []struct {
		name string
		p    []byte
		q    []byte
		want []byte
	}{
		{
			name: "same length",
			p:    []byte{0x53, 0xCA, 0x01},
			q:    []byte{0xFF, 0x02, 0x03},
			want: []byte{0xAC, 0xC8, 0x02},
		},
		{
			name: "different length, p longer",
			p:    []byte{0x53, 0xCA, 0x01},
			q:    []byte{0xFF, 0x02},
			want: []byte{0x53, 0x35, 0x03},
		},
		{
			name: "different length, q longer",
			p:    []byte{0xFF, 0x02},
			q:    []byte{0x53, 0xCA, 0x01},
			want: []byte{0x53, 0x35, 0x03},
		},
		{
			name: "self addition gives zero",
			p:    []byte{0x53, 0xCA, 0x01},
			q:    []byte{0x53, 0xCA, 0x01},
			want: []byte{0x00, 0x00, 0x00},
		},
		{
			name: "empty p",
			p:    []byte{},
			q:    []byte{0x53, 0xCA},
			want: []byte{0x53, 0xCA},
		},
		{
			name: "empty q",
			p:    []byte{0x53, 0xCA},
			q:    []byte{},
			want: []byte{0x53, 0xCA},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Add(c.p, c.q)
			if len(got) != len(c.want) {
				t.Fatalf("Add(%v, %v) len = %d, want %d", c.p, c.q, len(got), len(c.want))
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Fatalf("Add(%v, %v)[%d] = 0x%02x, want 0x%02x", c.p, c.q, i, got[i], c.want[i])
				}
			}
		})
	}
}

func TestMulIdentity(t *testing.T) {
	testPoly := [][]byte{
		{0x05},
		{0x53, 0xCA, 0x01},
		{0xFF, 0x00, 0xAB, 0x02},
	}
	identityElement := []byte{0x01}
	for _, c := range testPoly {
		got := Mul(c, identityElement)
		if len(got) != len(c) {
			t.Fatalf("Mul(%v, [0x01]) len = %d, want %d", c, len(got), len(c))
		}
		for i := range got {
			if got[i] != c[i] {
				t.Fatalf("Mul(%v, [0x01])[%d] = 0x%02x, want 0x%02x", c, i, got[i], c[i])
			}
		}
	}
}

func TestMulZero(t *testing.T) {
	testPoly := [][]byte{
		{0x05},
		{0x53, 0xCA, 0x01},
		{0xFF, 0x00, 0xAB, 0x02},
	}
	zeroElement := []byte{0x00}
	for _, c := range testPoly {
		got := Mul(c, zeroElement)
		for i := range got {
			if got[i] != 0x00 {
				t.Fatalf("Mul(%v, [0x00])[%d] = 0x%02x, want 0x%02x", c, i, got[i], zeroElement[i])
			}
		}
	}

}

func TestMulCommutativity(t *testing.T) {
	p := []byte{0x53, 0xCA, 0x01}
	q := []byte{0xFF, 0x02}
	pq := Mul(p, q)
	qp := Mul(q, p)
	if len(pq) != len(qp) {
		t.Fatalf("p*q len %d != q*p len %d", len(pq), len(qp))
	}
	for i := range pq {
		if pq[i] != qp[i] {
			t.Fatalf("p*q[%d] = 0x%02x, q*p[%d] = 0x%02x", i, pq[i], i, qp[i])
		}
	}
}

func TestMulDistributivity(t *testing.T) {
	p := []byte{0x53, 0x01}
	q := []byte{0xCA, 0x02}
	r := []byte{0xFF, 0xAB}
	lhs := Mul(p, Add(q, r))
	rhs := Add(Mul(p, q), Mul(p, r))
	if len(lhs) != len(rhs) {
		t.Fatalf("distributivity: len %d != %d", len(lhs), len(rhs))
	}
	for i := range lhs {
		if lhs[i] != rhs[i] {
			t.Fatalf("distributivity[%d] = 0x%02x, want 0x%02x", i, lhs[i], rhs[i])
		}
	}
}
