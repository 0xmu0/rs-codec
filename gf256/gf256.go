// Package gf256 implements arithmetic over GF(2^8) with
// irreducible polynomial x⁸ + x⁴ + x³ + x + 1 (0x11b).
//
// Generator g = 0x03 (x + 1).
// Every non-zero element is a power of g.
// g²⁵⁵ ≡ 1 (mod 0x11b).
package gf256

func Mul(a, b byte) byte {
	if a == 0 || b == 0 {
		return 0
	}
	logA := logTable[a]
	logB := logTable[b]
	logResult := (int(logA) + int(logB)) % 255
	return expTable[logResult]
}

func Div(a, b byte) byte {
	if b == 0 {
		panic("division by zero")
	}
	if a == 0 {
		return 0
	}
	logA := logTable[a]
	logB := logTable[b]
	logResult := (int(logA) - int(logB) + 255) % 255 // Add 255 to ensure non-negative before mod
	return expTable[logResult]
}

func Add(a, b byte) byte {
	return a ^ b
}

func Sub(a, b byte) byte {
	return a ^ b
}

func Inv(a byte) byte {
	if a == 0 {
		panic("zero has no multiplicative inverse")
	}
	logA := logTable[a]
	logResult := (255 - int(logA)) % 255
	return expTable[logResult]
}

// mulByX multiplies a single element by x (0x02) in GF(2^8).
//
// This is a left shift by 1. If the high bit was set (meaning x⁸ appeared),
// XOR with the full modulus 0x11b cancels the x⁸ term.
//
// In 9-bit form:
//
//	result_9bit = a << 1          (may produce x⁸ term)
//	result_9bit ⊕ 0x11b          (0x11b has x⁸ too, so x⁸ cancels, leaving 0x1b correction)
//
// Shortcut: if bit 7 was set, drop it (it becomes x⁸) and XOR with 0x1b (the modulus minus x⁸).
func mulByX(a byte) byte {
	shifted := uint16(a) << 1 // 9-bit result, no information lost
	if shifted&0x100 != 0 {   // x⁸ appeared (9th bit set)
		shifted ^= 0x11b // XOR full modulus: x⁸ ⊕ x⁸ = 0, rest reduces
	}
	return byte(shifted)
}

// mulByG multiplies a single element by the generator g = 0x03 = (x + 1).
//
// a × (x + 1) = a×x + a×1 = mulByX(a) ⊕ a
//
// Distributive law: multiplication over addition (XOR).
func mulByG(a byte) byte {
	return mulByX(a) ^ a
}

// expTable[i] = g^i for i in [0, 255].
// expTable[255] = g^255 ≡ 1 (mod 0x11b), same as expTable[0].
// This duplicate entry at index 255 simplifies mod 255 lookups.
var expTable [256]byte

var logTable [256]byte

func init() {
	var val byte = 1 // g⁰ = 1
	for i := 0; i < 256; i++ {
		expTable[i] = val
		val = mulByG(val)
	}

	val = 1
	for i := 0; i < 255; i++ {
		logTable[val] = byte(i)
		val = mulByG(val)
	}

}
