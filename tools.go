package key

// The reverse returns a slice of rune in reverse order.
func reverse(v []rune) []rune {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v[i], v[j] = v[j], v[i]
	}

	return v
}

// The unlead removes the leading characters from the slice of rune.
//
// The leading character is the first character in the alphabet,
// for example alphabet contains {'a', 'b', 'c', ..., 'z'} i.e.
// first char is 'a' - it is lead char.
//
// So, if the rune slice is {'a', 'a', 'a', 'c', 'a', 'b'} -
// function removes the first duplicates and returns a slice
// without lead chars {'c', 'a', 'b'}.
func unlead(lead rune, v []rune) []rune {
	var seek int

	for seek < len(v)-1 && v[seek] == lead {
		seek++
	}

	return v[seek:]
}

// Pow calculates the exponentiation of a base to an exponent using
// binary exponentiation. It returns the result of base raised to
// the power of exponent.
func pow(base, exponent int) int {
	if exponent < 0 {
		return 0
	}

	result := 1
	for exponent > 0 {
		if exponent&1 == 1 {
			result *= base
		}
		base *= base
		exponent >>= 1
	}

	return result
}
