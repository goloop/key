package key

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// New returns a pointer to a Locksmith object as the first value and
// an error if something went wrong, or nil as the second value.
//
// As the first argument, the function takes the size of the key.
// If the size of the key is set to zero, the key size will be dynamic, i.e.
// the minimum key size will be one character, and the maximum will depend
// on the length of the alphabet and the possible maximum iteration index.
//
// The second value is the sequence elements for permutation (alphabet).
// The alphabet mustn't contain duplicate chars or be empty.
func New(size uint, alphabet string) (*Locksmith, error) {
	var locksmith *Locksmith

	// The alphabet must contain at least one character.
	if len(alphabet) == 0 {
		return &Locksmith{}, errors.New("blank alphabet string")
	}

	// Create a pointer to the Locksmith object.
	locksmith = &Locksmith{
		size:     size,
		alphabet: []rune(alphabet),
		indexOf:  make(map[rune]int),
		total:    uint64(math.MaxUint64), // recalculate below if size != 0
	}

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
	for i, char := range locksmith.alphabet {
		// Check the presence of duplicates in the alphabet.
		// The alphabet shouldn't contain duplicates.
		if _, ok := locksmith.indexOf[char]; ok {
			return locksmith, fmt.Errorf(
				"the %c item is repeated in the alphabet",
				char,
			)
		}

		locksmith.indexOf[char] = i
	}

	// If the size is set to zero - the key size will be dynamic.
	// The dynamic index iteration is limited to MaxUint64 size.
	//
	// Otherwise the value of the last iteration index is calculated
	// according to the formula L to the power of S, where L is the
	// size of the alphabet, and S is the size of the key. But and
	// this value is limited to MaxUint64 too.
	if locksmith.size != 0 {
		l, s := float64(len(locksmith.alphabet)), float64(locksmith.size)
		if total := uint64(math.Pow(l, s)); total < math.MaxUint64 {
			locksmith.total = total
		}
	}

	locksmith.isValid = true
	return locksmith, nil
}

// Locksmith is a key generation object.
// It can be created correctly through the New function only.
type Locksmith struct {
	size     uint         // length of the generated key
	total    uint64       // maximum allowable key value
	alphabet []rune       // list of characters to generate the key
	indexOf  map[rune]int // the map of matching characters of alphabet
	isValid  bool         // true if the object was created correctly
}

// IsValid returns true if Locksmith object is valid.
func (ls *Locksmith) IsValid() bool {
	return ls.isValid // len(ls.alphabet) > 0 && len(ls.alphabet) == len(ls.indexOf)
}

// Alphabet returns current alphabet value.
func (ls *Locksmith) Alphabet() string {
	return string(ls.alphabet)
}

// Size return size of the key.
func (ls *Locksmith) Size() uint {
	return ls.size
}

// Total returns the highest possible iteration number.
//
// For example, for "abc" alphabet and key size as 3 - can be
// created the 27 iterations: aaa, aab, aac, ..., cca, ccb, ccc.
// So can be used indexs as 0 <= ID < 27 to generate a key.
func (ls *Locksmith) Total() uint64 {
	return ls.total
}

// Marshal returns the key (sequence element) by ID.
func (ls *Locksmith) Marshal(id uint64) (string, error) {
	var result string

	if id > ls.Total() {
		return "", fmt.Errorf("%d is large ID for key generation", id)
	}

	// Create key.
	al := len(ls.alphabet)
	l, r := int(id/uint64(al)), int(math.Mod(float64(id), float64(al)))
	result = string(ls.alphabet[r])
	for l >= al {
		l, r = l/al, int(math.Mod(float64(l), float64(al)))
		result = string(ls.alphabet[r]) + result
	}

	// If there is a balance only.
	if l != 0 {
		result = string(ls.alphabet[l]) + result
	}

	// Create the right size wrench.
	if repeat := int(ls.size) - len(result); ls.size > 0 && repeat > 0 {
		result = strings.Repeat(string(ls.alphabet[0]), repeat) + result
	}

	return result, nil
}

// Unmarshal returns ID of the specified sequence.
func (ls *Locksmith) Unmarshal(key string) (uint64, error) {
	var value = []rune(key)

	// The key is the wrong size.
	if l := uint(len(value)); ls.size > 0 && l != ls.size {
		return 0, fmt.Errorf("invalid key length, "+
			"must be %d char(s) but %d char(s)", ls.size, l)
	}

	id, value := uint64(0), reverse(unlead(ls.alphabet[0], value))
	for i, char := range value {
		index, ok := ls.indexOf[char]
		if !ok {
			return 0, fmt.Errorf("key contains a char that isn't "+
				"set in the alphabet: %c", char)
		}

		if i == 0 {
			id += uint64(index)
			continue
		}

		iter := int(math.Pow(float64(len(ls.alphabet)), float64(i)))
		id += uint64(index * iter)
	}

	return id, nil
}
