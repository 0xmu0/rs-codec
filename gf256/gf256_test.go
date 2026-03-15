package gf256

import "testing"

// g²⁵⁵ ≡ 1 (mod 0x11b)
// The generator cycles through all 255 non-zero elements and wraps back.
func TestExpTableWrapsToOne(t *testing.T) {
	if expTable[255] != 0x01 {
		t.Errorf("g^255 = 0x%02x, want 0x01", expTable[255])
	}
}

// g⁰ = 1 by definition.
func TestExpTableStartsAtOne(t *testing.T) {
	if expTable[0] != 0x01 {
		t.Errorf("g^0 = 0x%02x, want 0x01", expTable[0])
	}
}

// First few known powers of g = 0x03 with modulus 0x11b.
// We computed these by hand:
//
//	g⁰ = 0x01
//	g¹ = 0x03
//	g² = 0x05
//	g³ = 0x0f
//	g⁴ = 0x11
//	g⁵ = 0x33
//	g⁶ = 0x55
//	g⁷ = 0xff
//	g⁸ = 0x1a  (first overflow + reduction)
func TestExpTableKnownValues(t *testing.T) {
	known := []struct {
		power int
		value byte
	}{
		{0, 0x01},
		{1, 0x03},
		{2, 0x05},
		{3, 0x0f},
		{4, 0x11},
		{5, 0x33},
		{6, 0x55},
		{7, 0xff},
		{8, 0x1a},
	}
	for _, k := range known {
		if expTable[k.power] != k.value {
			t.Errorf("g^%d = 0x%02x, want 0x%02x", k.power, expTable[k.power], k.value)
		}
	}
}

// Every non-zero element must appear exactly once in expTable[0..254].
// This proves g = 0x03 is a generator of GF(2^8).
func TestExpTableAllDistinct(t *testing.T) {
	seen := [256]bool{}
	for i := 0; i < 255; i++ {
		v := expTable[i]
		if v == 0x00 {
			t.Errorf("g^%d = 0x00, but zero should never appear", i)
		}
		if seen[v] {
			t.Errorf("g^%d = 0x%02x is a duplicate", i, v)
		}
		seen[v] = true
	}
}

// mulByX matches manual computation.
// 0x80 × 0x02: x⁷ × x = x⁸ ≡ x⁴ + x³ + x + 1 = 0x1b (mod 0x11b)
func TestMulByX(t *testing.T) {
	cases := []struct {
		in   byte
		want byte
	}{
		{0x01, 0x02}, // 1 × x = x
		{0x02, 0x04}, // x × x = x²
		{0x80, 0x1b}, // x⁷ × x = x⁸ ≡ 0x1b (overflow + reduction)
		{0x00, 0x00}, // 0 × anything = 0
	}
	for _, c := range cases {
		got := mulByX(c.in)
		if got != c.want {
			t.Errorf("mulByX(0x%02x) = 0x%02x, want 0x%02x", c.in, got, c.want)
		}
	}
}

// mulByG matches manual computation.
// a × 0x03 = a×x ⊕ a = mulByX(a) ⊕ a
func TestMulByG(t *testing.T) {
	cases := []struct {
		in   byte
		want byte
	}{
		{0x01, 0x03}, // 1 × (x+1) = x+1
		{0x03, 0x05}, // (x+1) × (x+1) = x²+1
		{0x05, 0x0f}, // (x²+1) × (x+1) = x³+x²+x+1
		{0xff, 0x1a}, // g⁷ × g = g⁸ = 0x1a
	}
	for _, c := range cases {
		got := mulByG(c.in)
		if got != c.want {
			t.Errorf("mulByG(0x%02x) = 0x%02x, want 0x%02x", c.in, got, c.want)
		}
	}
}

// logTable is the reverse of expTable.
// logTable[g^i] = i for all i in [0, 254].
func TestLogTableKnownValues(t *testing.T) {
	known := []struct {
		element byte
		power   byte
	}{
		{0x01, 0},
		{0x03, 1},
		{0x05, 2},
		{0x0f, 3},
		{0x11, 4},
		{0x33, 5},
		{0x55, 6},
		{0xff, 7},
		{0x1a, 8},
	}
	for _, k := range known {
		if logTable[k.element] != k.power {
			t.Errorf("log(0x%02x) = %d, want %d", k.element, logTable[k.element], k.power)
		}
	}
}

// expTable and logTable are inverses of each other.
// For every non-zero element: expTable[logTable[a]] = a
func TestLogExpRoundTrip(t *testing.T) {
	for a := 1; a < 256; a++ {
		got := expTable[logTable[byte(a)]]
		if got != byte(a) {
			t.Errorf("expTable[logTable[0x%02x]] = 0x%02x, want 0x%02x", a, got, a)
		}
	}
}

// Mul(a, b) = Mul(b, a) for all elements — commutativity.
func TestMulCommutative(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			if Mul(byte(a), byte(b)) != Mul(byte(b), byte(a)) {
				t.Fatalf("Mul(0x%02x, 0x%02x) != Mul(0x%02x, 0x%02x)", a, b, b, a)
			}
		}
	}
}

// Mul(a, 0) = 0 — zero absorbs.
func TestMulZero(t *testing.T) {
	for a := 0; a < 256; a++ {
		if Mul(byte(a), 0) != 0 {
			t.Fatalf("Mul(0x%02x, 0) = 0x%02x, want 0", a, Mul(byte(a), 0))
		}
	}
}

// Mul(a, 1) = a — multiplicative identity.
func TestMulIdentity(t *testing.T) {
	for a := 0; a < 256; a++ {
		if Mul(byte(a), 1) != byte(a) {
			t.Fatalf("Mul(0x%02x, 1) = 0x%02x, want 0x%02x", a, Mul(byte(a), 1), a)
		}
	}
}

// Mul(a, Inv(a)) = 1 — inverse property.
func TestMulInverse(t *testing.T) {
	for a := 1; a < 256; a++ {
		got := Mul(byte(a), Inv(byte(a)))
		if got != 1 {
			t.Fatalf("Mul(0x%02x, Inv(0x%02x)) = 0x%02x, want 0x01", a, a, got)
		}
	}
}

// Div(a, b) = Mul(a, Inv(b)) — division is multiplication by inverse.
func TestDivIsMulByInverse(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 1; b < 256; b++ {
			got := Div(byte(a), byte(b))
			want := Mul(byte(a), Inv(byte(b)))
			if got != want {
				t.Fatalf("Div(0x%02x, 0x%02x) = 0x%02x, want 0x%02x", a, b, got, want)
			}
		}
	}
}

// Div(a, a) = 1 — self-division.
func TestDivSelf(t *testing.T) {
	for a := 1; a < 256; a++ {
		got := Div(byte(a), byte(a))
		if got != 1 {
			t.Fatalf("Div(0x%02x, 0x%02x) = 0x%02x, want 0x01", a, a, got)
		}
	}
}

// Add(a, a) = 0 — every element is its own additive inverse in GF(2).
func TestAddSelfInverse(t *testing.T) {
	for a := 0; a < 256; a++ {
		if Add(byte(a), byte(a)) != 0 {
			t.Fatalf("Add(0x%02x, 0x%02x) = 0x%02x, want 0", a, a, Add(byte(a), byte(a)))
		}
	}
}

// Sub = Add in GF(2). Subtraction is the same as addition.
func TestSubEqualsAdd(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			if Sub(byte(a), byte(b)) != Add(byte(a), byte(b)) {
				t.Fatalf("Sub(0x%02x, 0x%02x) != Add(0x%02x, 0x%02x)", a, b, a, b)
			}
		}
	}
}
