package key

import (
	"fmt"
	"math"
	"strings"
)

// Alphabet is default sequence.
const Alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

// New returns a Key object with the specified parameters,
// if an error occurs, the second value contains the error message.
func New(size uint, alphabet ...rune) (*Key, error) {
	var key *Key

	// // The key cannot consist of an empty string only.
	// if size < 1 {
	// 	return key, errors.New("the key size cannot be less than one")
	// }

	// Set the default alphabet if the custom alphabet not set.
	if len(alphabet) == 0 {
		alphabet = shuffle([]rune(Alphabet))
	}

	key = &Key{size: size, alphabet: alphabet, model: make(map[rune]int)}
	for i, char := range alphabet {
		if _, ok := key.model[char]; ok {
			return key, fmt.Errorf("some elements of the alphabet "+
				"are repeated: %c", char)
		}

		key.model[char] = i
	}

	return key, nil
}

// Key ...
type Key struct {
	// size is the size of the generated key
	size uint

	// alphabet is a list of words to generate the key
	alphabet []rune

	// model is a character identifiers in the alphabet
	model map[rune]int
}

// IsValid returns true if Key object is valid.
func (k *Key) IsValid() bool {
	return len(k.alphabet) > 0
}

// Alphabet returns current alphabet.
func (k *Key) Alphabet() []rune {
	return k.alphabet
}

// Size return size of the key.
func (k *Key) Size() uint {
	return k.size
}

// LastID returns the last available ID in the sequence.
func (k *Key) LastID() uint64 {
	if k.size == 0 {
		return math.MaxUint64
	}

	tmp := math.Pow(float64(len(k.alphabet)), float64(k.size))
	if tmp > math.MaxUint64 {
		return math.MaxUint64
	}

	return uint64(tmp)
}

// Marshal returns the key (sequence element) by ID.
func (k *Key) Marshal(id uint64) (string, error) {
	var result string

	if id > k.LastID() {
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
		index, ok := k.model[char]
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
