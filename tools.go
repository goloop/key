package key

import "math/bits"

// powU64 computes base**exp using binary exponentiation, entirely in uint64.
//
// Unlike a plain integer pow, every multiplication is checked for overflow
// with bits.Mul64: if the true result does not fit into uint64 the function
// returns (0, true). This lets the caller treat an over-large key space as
// "saturated" instead of silently wrapping around to a wrong (or zero) value.
//
// Examples:
//
//	powU64(3, 3)   -> (27, false)
//	powU64(16, 16) -> (0, true)   // 16**16 == 2**64 does not fit uint64
func powU64(base, exp uint64) (uint64, bool) {
	result := uint64(1)
	for exp > 0 {
		if exp&1 == 1 {
			hi, lo := bits.Mul64(result, base)
			if hi != 0 {
				return 0, true
			}
			result = lo
		}

		exp >>= 1
		if exp == 0 {
			break // avoid a final, unused (and possibly overflowing) square
		}

		hi, lo := bits.Mul64(base, base)
		if hi != 0 {
			return 0, true
		}
		base = lo
	}

	return result, false
}
