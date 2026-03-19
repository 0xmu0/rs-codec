package poly

import gf256 "github.com/0xmu0/rs-codec/gf256"

//Horner's Method
//Example : p(x) = 0x53·x² + 0xCA·x + 0x01
//p := []{0x53,0xCA,0x01}, x = a
//result := 0x53

//Until last element
//result = Mul(result,a) // Mul(0x53,a)
//result = Add(result,0xCA) // Add(Mul(0x53,a),0xCA)
//result = Mul(result,a) // Mul(Add(Mul(0x53,a),0xCA),a)
//result = Add(result,0x01) // Add(Mul(Add(Mul(0x53,a),0x01),a),a)

func Eval(p []byte, x byte) byte {

	if len(p) == 0 {
		panic("Empty data")
	}

	l := len(p)
	result := p[0]

	for i := 1; i < l; i++ {
		result = gf256.Mul(result, x)
		result = gf256.Add(result, p[i])
	}

	return result

}

func Add(p, q []byte) []byte {
	if len(p) == 0 {
		return q
	}
	if len(q) == 0 {
		return p
	}

	lp := len(p)
	lq := len(q)

	larger := p
	smaller := q
	if lq > lp {
		larger = q
		smaller = p
	}
	result := make([]byte, len(larger))
	copy(result, larger)

	j := len(smaller) - 1
	for i := len(result) - 1; i >= 0 && j >= 0; i-- {
		result[i] = gf256.Add(result[i], smaller[j])
		j--
	}

	return result
}

func Mul(p, q []byte) []byte {
	if len(p) == 0 || len(q) == 0 {
		return []byte{}
	}

	result := make([]byte, len(p)+len(q)-1)

	for i := 0; i < len(p); i++ {
		for j := 0; j < len(q); j++ {
			result[i+j] = gf256.Add(result[i+j], gf256.Mul(p[i], q[j]))
		}
	}
	return result
}
