// Package key converts uint64 identifiers into short, reversible string keys
// using a custom alphabet, and decodes them back. It is handy for shortening
// URLs, obfuscating sequential database ids, and minting human-readable codes
// such as tickets or coupons.
//
// Under the hood a Locksmith is a generalized base-N positional codec: with an
// alphabet of N characters the value is written in base N. The mapping is
// deterministic and bidirectional - each id has exactly one canonical key, and
// each valid key decodes back to its id. All arithmetic is integer-only
// (uint64 with math/bits overflow checks), so the codec is exact across the
// whole uint64 range, not just for small values.
//
// # Two modes
//
// A dynamic Locksmith yields variable-length keys, as short as the value
// allows:
//
//	ls, _ := key.NewDynamic(key.Base62)
//	k, _ := ls.Marshal(1234567) // "5BAH"
//	id, _ := ls.Unmarshal(k)    // 1234567
//
// A fixed Locksmith always yields keys of a given length, left-padded with the
// first alphabet character:
//
//	ls, _ := key.NewFixed(key.Base62, 8)
//	k, _ := ls.Marshal(1234567) // "00005BAH"
//
// # Strict decoding
//
// Unmarshal enforces a bijection. It rejects the empty string, keys of the
// wrong length (fixed mode), non-canonical keys with redundant leading
// lead-characters (dynamic mode), unknown characters, and values that overflow
// uint64. Use the sentinel errors with errors.Is to react to each case, or
// Valid to test a key without decoding it.
//
// # Key space
//
// For a fixed Locksmith the space is alphabet_length**size; for a dynamic one
// it spans the whole uint64 range. When the space reaches or exceeds 2**64 it
// is reported as saturated: Total returns MaxUint64, Saturated returns true,
// and every uint64 id is encodable.
//
// # Common alphabets
//
// The package ships ready-made alphabets (Base16, Base32, Base36, Base62,
// Crockford, Unambiguous); any string of unique characters works too. Prefer
// an unambiguous alphabet when keys are typed by hand, and a fixed size when a
// constant length matters.
//
// A Locksmith is immutable after construction and safe for concurrent use.
package key
