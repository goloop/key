package key

import (
	"math/rand"
	"time"
)

// The reverse writes in reverse order a slice of runes.
//
// Usage:
// 		// Reverse by the link.
// 		var abc = []rune{'a', 'b', 'c'}
// 		reverse(abc) // abc == []rune{'c', 'b', 'a'}
//
//		// As the returned result.
//		def := reverse([]rune{'d', 'e', 'f'}) // def == []rune{'f', 'e', 'd'}
func reverse(v []rune) []rune {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v[i], v[j] = v[j], v[i]
	}
	return v
}

// The unlead removes the leading characters from rune slice.
//
// The leading character is the first character in the alphabet,
// for example alphabet contains {'a', 'b', 'c', ..., 'z'} i.e.
// first char is 'a'.
//
// So if the rune slice is {'a', 'a', 'a', 'c', 'a', 'b'} -
// function removes the first duplicates and returns a slice
// without lead chars {'c', 'a', 'b'}.
//
// Usage:
//		// Remove leading chars.
//		abc = unlead('a', []rune{'a', 'a', 'a', 'b'}) // abc == []rune{'b'}
func unlead(lead rune, v []rune) []rune {
	var seek int

	for seek < len(v)-1 && v[seek] == lead {
		seek++
	}

	return v[seek:]
}

// The shuffle method takes a sequence of runes and
// reorganize the order of the items.
func shuffle(v []rune) []rune {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(v), func(i, j int) { v[i], v[j] = v[j], v[i] })
	return v
}
