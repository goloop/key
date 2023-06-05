package key

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// New returns a new Locksmith object. It takes in three arguments:
// alphabet (string) and size (int).
//
//   - alphabet is a string of unique characters from which the keys will
//     be generated. It must contain at least 2 unique characters and no
//     duplicates.
//
//   - size' is the fixed length of the keys to be generated. If size is
//     set to zero, the key size will be dynamic, with a minimum size of 1.
//     In this case, the maximum key size will depend on the length of
//     the alphabet and the maximum iteration index.
//
//     The value can be specified as a sequence of numbers, in this case,
//     size will be the sum of these numbers New("abc", 1, 2, 3) // size == 6
//
//     The size value can be omitted for dynamic keys (or set size as 0).
//
// The function returns a pointer to the created Locksmith object and
// an error if something went wrong, or nil if it was successful.
//
// Example usage:
//
//	ls, err := New("abc")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	key, err := ls.Marshal(10)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(key) // Output: "bab"
func New(alphabet string, args ...int) (*Locksmith, error) {
	var size int64

	// The alphabet must contain at least one character.
	if len(alphabet) == 0 {
		return &Locksmith{}, errors.New("blank alphabet string")
	}

	// Size is a sum of all arguments.
	// It must be zero or positive value.
	for _, v := range args {
		size += int64(v)
	}

	if size < 0 {
		return &Locksmith{}, errors.New("incorrect size")
	}

	// Create a pointer to the Locksmith object.
	locksmith := &Locksmith{
		size:     uint64(size),
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
			return &Locksmith{}, fmt.Errorf(
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

	return locksmith, nil
}

// Locksmith is a key generation object.
// It can be created correctly through the New function only.
type Locksmith struct {
	size     uint64       // length of the generated key
	total    uint64       // maximum allowable key value
	alphabet []rune       // list of characters to generate the key
	indexOf  map[rune]int // the map of matching characters of alphabet
}

// Alphabet returns current alphabet value.
func (ls *Locksmith) Alphabet() string {
	return string(ls.alphabet)
}

// Size return size of the key.
func (ls *Locksmith) Size() uint64 {
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

// Marshal converts an ID into a key. This function takes an ID as
// input and generates a corresponding key based on the Locksmith's
// alphabet and size.
//
// The ID should be less than the total number of possible keys,
// otherwise, an error will be returned. The total number of possible
// keys can be obtained by calling the 'Total' method.
//
// If the Locksmith's size is set to a fixed length, the generated key
// will be padded with the first character of the alphabet to reach the
// required length. If the size is set to zero (i.e., dynamic size),
// the key will not be padded and its length will vary depending on the
// ID.
//
// This function returns a string representing the key and an error if
// something went wrong. If the function is successful, the error will
// be nil.
//
// Example usage:
//
//	ls, _ := New("abc")
//	key, err := ls.Marshal(10)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(key) // Output: "bab"
func (ls *Locksmith) Marshal(id uint64) (string, error) {
	var result string

	if id >= ls.Total() {
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

// Unmarshal decodes a key and returns its corresponding ID.
// This function converts a key back into its ID. The key should be a
// string composed of characters from the Locksmith's alphabet.
//
// If the Locksmith's size is set to a fixed length, the key must have
// the same length. If the size is set to zero (i.e., dynamic size), the
// key can have any length.
//
// This function returns an integer representing the ID and an error if
// something went wrong. If the function is successful, the error will
// be nil.
//
// Example usage:
//
//	ls, _ := New("abc")
//	id, err := ls.Unmarshal("bab")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(id) // Output: 10
func (ls *Locksmith) Unmarshal(key string) (uint64, error) {
	value := []rune(key)

	// The key is the wrong size.
	if l := uint64(len(value)); ls.size > 0 && l != ls.size {
		return 0, fmt.Errorf("invalid key length, "+
			"must be %d char(s) but %d char(s)", ls.size, l)
	}

	id, value := uint64(0), reverse(unlead(ls.alphabet[0], value))
	alphabetLength := len(ls.alphabet)
	for i, char := range value {
		index, ok := ls.indexOf[char]
		if !ok {
			return 0, fmt.Errorf("key contains a char that isn't "+
				"set in the alphabet: %c", char)
		}

		// Replacing loop with math.Pow.
		power := math.Pow(float64(alphabetLength), float64(i))
		id += uint64(float64(index) * power)
	}

	return id, nil
}
