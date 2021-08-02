package key

import (
	"fmt"
	"math"
	"strings"
)

// The defaultAlphabetCharacters the string of characters from which
// the alphabet will be generated if the user doesn't set a custom one.
const defaultAlphabetCharacters = "abcdefghijklmnopqrstuvwxyz0123456789"

// New returns a pointer to the Key object as the first value.
// The second value can contains an error if something went wrong.
//
// The function takes size of key as first argument. If any positive value
// greater than zero is specified the key length will match this value.
// Otherwise, if the size is set to zero - the key size will be dynamic, i.e.
// the minimum key size will be one character and the maximum will depend on
// the length of the alphabet and the possible maximum index of the iteration.
//
// Second, third, etc. is optional arguments of function.
// These are the elements of the sequence for permutation (alphabet).
// If the custom alphabet is missing, will be used an alphabet that
// randomly generated from the characters a-z and 0-9.
// The alphabet should not contain duplicate values.
func New(size uint, alphabet ...rune) (*Key, error) {
	var key *Key

	// Generate default alphabet if custom one is missing.
	if len(alphabet) == 0 {
		// Use a standard sequence of characters to generate the alphabet.
		// The alphabet should to be shuffled. Two different objects must
		// have a different alphabet.
		alphabet = shuffle([]rune(defaultAlphabetCharacters))
	}

	// Create a pointer to the Key object.
	key = &Key{size: size, alphabet: alphabet, indexOf: make(map[rune]int)}

	// In the operation of the algorithm, it is necessary to determine the
	// character in the sequence by the specified index (just slice[index]),
	// and the character index in the sequence by the character
	// (like the indexOf method).
	//
	// The classical indexOf method iterates over the characters in
	// the sequence, which slows down the algorithm. Instead, we use
	// map to store the matched character and index in the sequence.
	// This requires quite a bit more memory to store a copy of the
	// alphabet with the matches but increases the speed of the
	// algorithm as a whole.
	for i, char := range alphabet {
		// Check the presence of duplicates in the alphabet.
		// The alphabet should not contain duplicates.
		if _, ok := key.indexOf[char]; ok {
			return key, fmt.Errorf("some elements of the alphabet "+
				"are repeated: %c", char)
		}

		key.indexOf[char] = i
	}

	return key, nil
}

// Key is a key object.
type Key struct {
	// size is the length of the generated key
	size uint

	// alphabet is a list of characters to generate the key
	alphabet []rune

	// indexOf it the map of matching alphabet characters as
	// character and index (as well as an alphabet duplicate checker tool)
	indexOf map[rune]int
}

// IsValid returns true if Key object is valid.
// True only when the New method was executed without error.
func (k *Key) IsValid() bool {
	return len(k.alphabet) > 0 && len(k.alphabet) == len(k.indexOf)
}

// Alphabet returns current alphabet as rune slice.
func (k *Key) Alphabet() []rune {
	return k.alphabet
}

// Size return size of the key.
func (k *Key) Size() uint {
	return k.size
}

// Total returns the highest possible iteration number.
// For example, for "abc" alphabet and 3 key size can be created
// the 27 iterations: aaa, aab, aac, ..., cca, ccb, ccc.
// So indexs as 0 <= ID < Totla() can be used to generate a key.
func (k *Key) Total() uint64 {
	// If the size is set to zero - the key size is dynamic.
	// Dynamic index iteration is limited to MaxUint64 size.
	if k.size == 0 {
		return math.MaxUint64
	}

	// The value of the last iteration index is calculated according to the
	// formula A to the power of S, where A is the size of the alphabet,
	// and S is the size of the key.
	//
	// But this value is limited to MaxUint64 size too.
	tmp := math.Pow(float64(len(k.alphabet)), float64(k.size))
	if tmp > math.MaxUint64 {
		return math.MaxUint64
	}

	return uint64(tmp)
}

// Marshal returns the key (sequence element) by ID.
func (k *Key) Marshal(id uint64) (string, error) {
	var result string

	if id > k.Total() {
		return "", fmt.Errorf("large ID for key generation: %d", id)
	}

	// Create key.
	al := len(k.alphabet)
	l, r := int(id/uint64(al)), int(math.Mod(float64(id), float64(al)))
	result = string(k.alphabet[r])
	for l >= al {
		l, r = l/al, int(math.Mod(float64(l), float64(al)))
		result = string(k.alphabet[r]) + result
	}

	// Only if there is a balance.
	if l != 0 {
		result = string(k.alphabet[l]) + result
	}

	// Create the right size wrench.
	if repeat := int(k.size) - len(result); k.size > 0 && repeat > 0 {
		result = strings.Repeat(string(k.alphabet[0]), repeat) + result
	}

	return result, nil
}

// Unmarshal returns ID of the specified sequence.
func (k *Key) Unmarshal(key string) (uint64, error) {
	var value = []rune(key)

	if k.size > 0 && uint(len(value)) != k.size {
		return 0, fmt.Errorf("invalid key length, "+
			"max %d but %d", k.size, len(value))
	}

	id, value := uint64(0), reverse(unlead(k.alphabet[0], value))
	for i, char := range value {
		index, ok := k.indexOf[char]
		if !ok {
			return 0, fmt.Errorf("key contains a char that isn't "+
				"set in the alphabet: %c", char)
		}

		if i == 0 {
			id += uint64(index)
			continue
		}

		iter := int(math.Pow(float64(len(k.alphabet)), float64(i)))
		id += uint64(index * iter)
	}

	return id, nil
}
