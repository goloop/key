package key

import "errors"

// Sentinel errors returned by the package. They are exported so callers can
// match them with errors.Is and react programmatically instead of parsing
// error strings. Constructors and methods wrap these with additional detail
// (the offending character, length or id) using fmt.Errorf("%w: ...").
var (
	// ErrBlankAlphabet is returned when the alphabet is an empty string.
	ErrBlankAlphabet = errors.New("alphabet is blank")

	// ErrShortAlphabet is returned when the alphabet has fewer than two
	// unique characters. A base-1 alphabet cannot encode anything beyond
	// zero and would loop forever, so it is rejected up front.
	ErrShortAlphabet = errors.New("alphabet must contain at least 2 unique characters")

	// ErrDuplicateChar is returned when the alphabet contains the same
	// character more than once.
	ErrDuplicateChar = errors.New("alphabet contains a duplicate character")

	// ErrInvalidSize is returned by NewFixed when size is not positive.
	ErrInvalidSize = errors.New("fixed size must be a positive number")

	// ErrIDTooLarge is returned by Marshal when the id does not fit into the
	// key space (id >= Total, for non-saturated spaces).
	ErrIDTooLarge = errors.New("id is out of the key space")

	// ErrEmptyKey is returned by Unmarshal when the key is an empty string.
	// The empty string is never produced by Marshal: the canonical key for
	// zero is the first alphabet character.
	ErrEmptyKey = errors.New("key is empty")

	// ErrInvalidLength is returned by Unmarshal in fixed mode when the key
	// length differs from the configured size.
	ErrInvalidLength = errors.New("key has invalid length")

	// ErrUnknownChar is returned by Unmarshal when the key contains a
	// character that is not part of the alphabet.
	ErrUnknownChar = errors.New("key contains a character not in the alphabet")

	// ErrNonCanonical is returned by Unmarshal in dynamic mode when the key
	// carries redundant leading lead-characters (e.g. "aab" instead of "b").
	// Strict decoding keeps the mapping a bijection: one id, one key.
	ErrNonCanonical = errors.New("key is not in canonical form")

	// ErrKeyOutOfRange is returned by Unmarshal when the decoded value would
	// overflow uint64.
	ErrKeyOutOfRange = errors.New("key is out of the uint64 range")
)
