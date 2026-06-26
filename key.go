package key

import (
	"fmt"
	"iter"
	"math"
	"math/bits"
	"unicode/utf8"
)

// stackDigits is the size of the on-stack digit buffer used by Marshal and
// MarshalAppend. A uint64 needs at most 64 digits (base 2), so any dynamic
// key and any fixed key up to this length is encoded without a heap
// allocation. Larger fixed sizes fall back to a single make.
const stackDigits = 64

// Locksmith encodes uint64 identifiers into string keys and back, using a
// custom alphabet. It is a generalized base-N positional codec: with an
// alphabet of length N the value is written in base N, optionally padded to a
// fixed width with the first alphabet character.
//
// A Locksmith is immutable after construction and therefore safe for
// concurrent use by multiple goroutines. Build one only through NewDynamic,
// NewFixed or their Must* counterparts.
type Locksmith struct {
	alphabet []rune       // characters of the alphabet, index == digit value
	indexOf  map[rune]int // reverse lookup: character -> digit value
	base     uint64       // len(alphabet); the numeric base
	size     uint64       // fixed key length; 0 means dynamic length
	total    uint64       // number of representable ids (saturated to MaxUint64)
	full     bool         // true when every uint64 id fits the key space
}

// newLocksmith validates the alphabet and builds a Locksmith. size == 0 with
// dynamic == true selects a variable-length codec; otherwise size is the fixed
// key width.
func newLocksmith(alphabet string, size uint64, dynamic bool) (*Locksmith, error) {
	runes := []rune(alphabet)
	switch {
	case len(runes) == 0:
		return nil, ErrBlankAlphabet
	case len(runes) < 2:
		return nil, ErrShortAlphabet
	}

	// Build the reverse lookup and reject duplicates in one pass. A duplicate
	// would make two digit values share a character, breaking the bijection.
	indexOf := make(map[rune]int, len(runes))
	for i, ch := range runes {
		if _, dup := indexOf[ch]; dup {
			return nil, fmt.Errorf("%w: %q", ErrDuplicateChar, ch)
		}
		indexOf[ch] = i
	}

	ls := &Locksmith{
		alphabet: runes,
		indexOf:  indexOf,
		base:     uint64(len(runes)),
		size:     size,
	}

	// Determine the size of the key space.
	//
	// Dynamic keys can grow as long as needed, so every uint64 is
	// representable: the space is saturated. For fixed keys the space is
	// base**size; if that overflows uint64 it is likewise saturated (every
	// uint64 id fits into size digits). Only a space strictly smaller than
	// 2**64 has a meaningful, enforced upper bound.
	if dynamic {
		ls.full = true
		ls.total = math.MaxUint64
		return ls, nil
	}

	total, overflow := powU64(ls.base, size)
	if overflow {
		ls.full = true
		ls.total = math.MaxUint64
	} else {
		ls.total = total
	}

	return ls, nil
}

// NewDynamic returns a Locksmith that produces variable-length keys: the key
// is as short as the value allows, with no padding. The canonical key for
// zero is the first alphabet character.
//
// The alphabet must contain at least two unique characters and no duplicates.
func NewDynamic(alphabet string) (*Locksmith, error) {
	return newLocksmith(alphabet, 0, true)
}

// NewFixed returns a Locksmith that always produces keys of exactly size
// characters, left-padded with the first alphabet character when the value is
// short. size must be positive.
//
// The alphabet must contain at least two unique characters and no duplicates.
func NewFixed(alphabet string, size int) (*Locksmith, error) {
	if size < 1 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidSize, size)
	}
	return newLocksmith(alphabet, uint64(size), false)
}

// MustNewDynamic is like NewDynamic but panics on error. It is meant for
// package-level codecs with a constant alphabet, where a bad alphabet is a
// programming error that should fail at startup.
func MustNewDynamic(alphabet string) *Locksmith {
	ls, err := NewDynamic(alphabet)
	if err != nil {
		panic(err)
	}
	return ls
}

// MustNewFixed is like NewFixed but panics on error. It is meant for
// package-level codecs with constant parameters.
func MustNewFixed(alphabet string, size int) *Locksmith {
	ls, err := NewFixed(alphabet, size)
	if err != nil {
		panic(err)
	}
	return ls
}

// Alphabet returns the alphabet the Locksmith was built with.
func (ls *Locksmith) Alphabet() string {
	return string(ls.alphabet)
}

// Size returns the fixed key length, or 0 for a dynamic-length Locksmith.
func (ls *Locksmith) Size() uint64 {
	return ls.size
}

// Total returns the number of representable ids: valid ids run from 0 up to
// Total()-1.
//
// When the key space is at least 2**64 the value is saturated to MaxUint64
// and every uint64 id is valid (including MaxUint64 itself); use Saturated to
// tell that case apart from a space whose true size happens to be MaxUint64.
func (ls *Locksmith) Total() uint64 {
	return ls.total
}

// Saturated reports whether the key space holds every uint64 id. It is true
// for any dynamic Locksmith and for a fixed one whose base**size reaches or
// exceeds 2**64. When true, Marshal never rejects an id.
func (ls *Locksmith) Saturated() bool {
	return ls.full
}

// writeDigits appends the base-N digits of id to dst, least-significant first,
// then pads with the lead character up to the fixed size. The caller reverses
// the written segment to obtain the most-significant-first key.
func (ls *Locksmith) writeDigits(dst []rune, id uint64) []rune {
	start := len(dst)
	for {
		dst = append(dst, ls.alphabet[id%ls.base])
		id /= ls.base
		if id == 0 {
			break
		}
	}

	for ls.size > 0 && uint64(len(dst)-start) < ls.size {
		dst = append(dst, ls.alphabet[0])
	}

	return dst
}

// Marshal converts an id into its key.
//
// For a fixed-length Locksmith the key is left-padded with the first alphabet
// character to the configured size. For a dynamic one the key is as short as
// the value allows. An id outside the key space yields ErrIDTooLarge.
func (ls *Locksmith) Marshal(id uint64) (string, error) {
	if !ls.full && id >= ls.total {
		return "", fmt.Errorf("%w: %d", ErrIDTooLarge, id)
	}

	var stack [stackDigits]rune
	var buf []rune
	if ls.size > stackDigits {
		buf = make([]rune, 0, ls.size)
	} else {
		buf = stack[:0]
	}

	buf = ls.writeDigits(buf, id)
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf), nil
}

// MarshalAppend encodes id and appends the resulting key (UTF-8 encoded) to
// dst, returning the extended slice. It mirrors the strconv.Append* family and
// avoids allocating a fresh string on each call, which is useful when emitting
// many keys into a shared buffer.
func (ls *Locksmith) MarshalAppend(dst []byte, id uint64) ([]byte, error) {
	if !ls.full && id >= ls.total {
		return dst, fmt.Errorf("%w: %d", ErrIDTooLarge, id)
	}

	var stack [stackDigits]rune
	var digits []rune
	if ls.size > stackDigits {
		digits = make([]rune, 0, ls.size)
	} else {
		digits = stack[:0]
	}

	digits = ls.writeDigits(digits, id)
	for i := len(digits) - 1; i >= 0; i-- {
		dst = utf8.AppendRune(dst, digits[i])
	}

	return dst, nil
}

// Unmarshal decodes a key back into its id.
//
// Decoding is strict: it enforces a bijection so that each valid key maps to
// exactly one id and round-trips back to the same key. Concretely it rejects
//   - the empty string (ErrEmptyKey);
//   - in fixed mode, a key whose length differs from size (ErrInvalidLength);
//   - in dynamic mode, redundant leading lead-characters, e.g. "aab"
//     (ErrNonCanonical);
//   - a character outside the alphabet (ErrUnknownChar);
//   - a value that overflows uint64 (ErrKeyOutOfRange).
func (ls *Locksmith) Unmarshal(key string) (uint64, error) {
	runes := []rune(key)
	if len(runes) == 0 {
		return 0, ErrEmptyKey
	}

	if ls.size > 0 {
		if l := uint64(len(runes)); l != ls.size {
			return 0, fmt.Errorf("%w: want %d char(s), got %d",
				ErrInvalidLength, ls.size, l)
		}
	} else if len(runes) > 1 && runes[0] == ls.alphabet[0] {
		// Dynamic keys are not padded, so a leading lead-character (other than
		// the single-character key for zero) is a non-canonical encoding.
		return 0, fmt.Errorf("%w: redundant leading %q",
			ErrNonCanonical, ls.alphabet[0])
	}

	var id uint64
	for _, ch := range runes {
		idx, ok := ls.indexOf[ch]
		if !ok {
			return 0, fmt.Errorf("%w: %q", ErrUnknownChar, ch)
		}

		// Horner's scheme with overflow control: id = id*base + idx.
		hi, lo := bits.Mul64(id, ls.base)
		if hi != 0 {
			return 0, fmt.Errorf("%w: %q", ErrKeyOutOfRange, key)
		}
		sum, carry := bits.Add64(lo, uint64(idx), 0)
		if carry != 0 {
			return 0, fmt.Errorf("%w: %q", ErrKeyOutOfRange, key)
		}
		id = sum
	}

	return id, nil
}

// Valid reports whether key is a well-formed, canonical key for this
// Locksmith. It is true exactly when Unmarshal would succeed.
func (ls *Locksmith) Valid(key string) bool {
	_, err := ls.Unmarshal(key)
	return err == nil
}

// Iter returns a range-over-func sequence of (id, key) pairs for ids in the
// inclusive range [from, to]. For a bounded Locksmith to is clamped to the
// last valid id (Total-1); if from is past the end the sequence is empty.
//
//	for id, k := range ls.Iter(0, 99) {
//	    fmt.Println(id, k)
//	}
//
// The inclusive upper bound lets the sequence reach the final id (including
// MaxUint64 on a saturated space), which a half-open range could not express.
func (ls *Locksmith) Iter(from, to uint64) iter.Seq2[uint64, string] {
	return func(yield func(uint64, string) bool) {
		hi := to
		if !ls.full && hi > ls.total-1 {
			hi = ls.total - 1
		}
		if from > hi {
			return
		}

		for id := from; ; id++ {
			k, _ := ls.Marshal(id) // cannot fail: id is within the space
			if !yield(id, k) || id == hi {
				return
			}
		}
	}
}
