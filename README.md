[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/key/v2)](https://goreportcard.com/report/github.com/goloop/key/v2) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/key/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://pkg.go.dev/github.com/goloop/key/v2) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)

# Key

Package `key` converts `uint64` identifiers into short, reversible string keys
using a custom alphabet, and decodes them back. It is handy for shortening
URLs, obfuscating sequential database IDs, and minting human-readable codes
such as tickets or coupons.

Under the hood a `Locksmith` is a generalized base-N positional codec: with an
alphabet of N characters the value is written in base N. The mapping is
deterministic and bidirectional — each ID has exactly one canonical key, and
each valid key decodes back to its ID. All arithmetic is integer-only
(`uint64` with `math/bits` overflow checks), so the codec is exact across the
**whole** `uint64` range, not just for small values.

## Theory

With an arbitrary alphabet and key size you get a sequence of unique
combinations, where each combination has a unique numeric index (from `0` to
`Total-1`). Given an index (e.g. a database ID) you get its key; given a key
you recover the index.

For the `"abc"` alphabet and a fixed key size of 3:

```
    0. aaa    1. aab    2. aac    3. aba    4. abb    5. abc
    6. aca    7. acb    8. acc    9. baa   10. bab   11. bac
   12. bba   13. bbb   14. bbc   15. bca   16. bcb   17. bcc
   18. caa   19. cab   20. cac   21. cba   22. cbb   23. cbc
   24. cca   25. ccb   26. ccc
```

So the maximum number of iterations is 27. Iteration 10 maps to `"bab"`, and
the `"aba"` combination is iteration 3.

## Install

```shell
$ go get github.com/goloop/key/v2
```

## Usage

```go
package main

import (
	"fmt"

	key "github.com/goloop/key/v2"
)

func main() {
	ls, _ := key.NewFixed("abcde", 3)
	v, _ := ls.Marshal(122)     // "eec"
	i, _ := ls.Unmarshal("eec") // 122

	fmt.Println(v, i)
	// Output: eec 122
}
```

### Fixed vs dynamic length

A **fixed** `Locksmith` always produces keys of a given length, left-padded
with the first alphabet character:

```go
ls, _ := key.NewFixed("abc", 8)
ls.Marshal(10) // "aaaaabab"
```

A **dynamic** `Locksmith` produces keys as short as the value allows:

```go
ls, _ := key.NewDynamic("abc")
ls.Marshal(1)        // "b"
ls.Marshal(10)       // "bab"
ls.Marshal(10000000) // "caacbbaabbacbab"
```

### Strict decoding

`Unmarshal` enforces a bijection, so a key always round-trips to itself. It
returns a typed sentinel error you can match with `errors.Is`:

```go
ls, _ := key.NewDynamic("abc")

_, err := ls.Unmarshal("")     // key.ErrEmptyKey
_, err = ls.Unmarshal("aab")   // key.ErrNonCanonical (canonical key for 1 is "b")
_, err = ls.Unmarshal("abX")   // key.ErrUnknownChar
_ = err

ok := ls.Valid("bab") // true; Valid is true iff Unmarshal succeeds
```

### Ready-made alphabets

The package ships common alphabets — `Base16`, `Base32`, `Base36`, `Base62`,
`Crockford` and `Unambiguous` (drops look-alike characters for hand-typed
codes). Any string of unique characters works too.

```go
coupon := key.MustNewFixed(key.Unambiguous, 6)
```

## API

### Constructors

- **NewDynamic**(alphabet string) (\*Locksmith, error) — variable-length keys.
- **NewFixed**(alphabet string, size int) (\*Locksmith, error) — fixed-length keys.
- **MustNewDynamic** / **MustNewFixed** — like the above but panic on a bad
  alphabet; meant for package-level codecs with constant parameters.

The alphabet must contain at least two unique characters and no duplicates. On
error the returned `*Locksmith` is `nil`.

### Methods

- **Alphabet**() string — the alphabet the Locksmith was built with.
- **Size**() uint64 — the fixed key length, or `0` for a dynamic Locksmith.
- **Total**() uint64 — number of representable IDs; valid IDs run `0..Total-1`.
  Saturated to `MaxUint64` for spaces of at least 2⁶⁴.
- **Saturated**() bool — whether every `uint64` ID fits the key space.
- **Marshal**(id uint64) (string, error) — encode an ID into its key.
- **MarshalAppend**(dst []byte, id uint64) ([]byte, error) — encode and append
  the key to `dst`, avoiding a string allocation on hot paths.
- **Unmarshal**(key string) (uint64, error) — decode a key back into its ID
  (strict).
- **Valid**(key string) bool — whether `key` is a well-formed, canonical key.
- **Iter**(from, to uint64) iter.Seq2[uint64, string] — range-over-func over the
  inclusive ID range `[from, to]`, yielding `(id, key)` pairs.
- **IterN**(from, n uint64) iter.Seq2[uint64, string] — at most `n` keys from
  `from`; count-bounded, safer than `Iter` for bulk generation.
- **Random**(r io.Reader) (string, error) — a key for a uniformly random ID
  (rejection sampling, no modulo bias). Pass `crypto/rand.Reader` for secure keys.
- **RandomCrypto**() (string, error) — secure-by-default wrapper over
  `Random(crypto/rand.Reader)`.

### Errors

`ErrBlankAlphabet`, `ErrShortAlphabet`, `ErrDuplicateChar`, `ErrInvalidSize`,
`ErrIDTooLarge`, `ErrEmptyKey`, `ErrInvalidLength`, `ErrUnknownChar`,
`ErrNonCanonical`, `ErrKeyOutOfRange` — all matchable with `errors.Is`.
