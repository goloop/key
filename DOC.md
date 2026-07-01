# key — reference

The full reference for the `key` package: the mental model, the base-N theory,
every constructor and method, the ready-made alphabets, the error set and
practical recipes.

Ukrainian version: **[DOC.UK.md](DOC.UK.md)**.

## Contents

- [Mental model](#mental-model)
- [Theory: a base-N positional codec](#theory-a-base-n-positional-codec)
- [Constructors](#constructors)
- [Fixed vs dynamic length](#fixed-vs-dynamic-length)
- [Encoding and decoding](#encoding-and-decoding)
- [Strict decoding and canonical keys](#strict-decoding-and-canonical-keys)
- [Inspecting a Locksmith](#inspecting-a-locksmith)
- [Iteration](#iteration)
- [Random keys](#random-keys)
- [Ready-made alphabets](#ready-made-alphabets)
- [Errors](#errors)
- [Recipes and tips](#recipes-and-tips)

## Mental model

`key` converts a `uint64` identifier into a short, reversible string key over a
custom alphabet, and decodes it back. It is useful for shortening URLs,
obfuscating sequential database IDs, and minting human-readable codes such as
tickets or coupons.

The one type you work with is a `*Locksmith`. It is:

- **Deterministic and bidirectional** — each ID has exactly one canonical key,
  and each valid key decodes to exactly one ID.
- **Exact across the whole `uint64` range** — all arithmetic is integer-only
  (`uint64` with `math/bits` overflow checks), so there is no floating-point
  drift and no silent truncation for large values.
- **Immutable** — a `Locksmith` is configured once at construction; there is no
  mutable global state to corrupt.

```go
import key "github.com/goloop/key/v2"
```

## Theory: a base-N positional codec

With an alphabet of N characters, a value is written in base N — each character
is a digit. With a fixed key size you get a finite sequence of unique
combinations, each with a unique numeric index from `0` to `Total-1`. Given an
index (say a database ID) you get its key; given a key you recover the index.

For the `"abc"` alphabet and a fixed key size of 3:

```
    0. aaa    1. aab    2. aac    3. aba    4. abb    5. abc
    6. aca    7. acb    8. acc    9. baa   10. bab   11. bac
   12. bba   13. bbb   14. bbc   15. bca   16. bcb   17. bcc
   18. caa   19. cab   20. cac   21. cba   22. cbb   23. cbc
   24. cca   25. ccb   26. ccc
```

There are 3³ = 27 combinations. Iteration 10 maps to `"bab"`; the `"aba"`
combination is iteration 3.

## Constructors

```go
func NewDynamic(alphabet string) (*Locksmith, error)
func NewFixed(alphabet string, size int) (*Locksmith, error)
func MustNewDynamic(alphabet string) *Locksmith
func MustNewFixed(alphabet string, size int) *Locksmith
```

`NewDynamic` builds a variable-length codec; `NewFixed` a fixed-length one. The
`MustNew…` variants panic on a bad alphabet and are meant for package-level
codecs with constant parameters.

The alphabet must contain **at least two unique characters and no duplicates**.
On error the returned `*Locksmith` is `nil` and the error is one of the
[alphabet/size sentinels](#errors).

```go
ls, err := key.NewFixed("abcde", 3)
coupon := key.MustNewFixed(key.Unambiguous, 6) // package-level, constant params
```

## Fixed vs dynamic length

A **fixed** `Locksmith` always produces keys of the configured length,
left-padded with the first alphabet character:

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

Choose fixed when the key length must be constant (columns, printed codes,
sortable identifiers); choose dynamic when you want the shortest possible key.

## Encoding and decoding

```go
func (l *Locksmith) Marshal(id uint64) (string, error)
func (l *Locksmith) MarshalAppend(dst []byte, id uint64) ([]byte, error)
func (l *Locksmith) Unmarshal(key string) (uint64, error)
```

`Marshal` encodes an ID into its key; it returns [`ErrIDTooLarge`](#errors) when
the ID does not fit a fixed key space. `MarshalAppend` encodes and appends the
key bytes to `dst`, avoiding a string allocation on hot paths. `Unmarshal`
decodes a key back into its ID (strictly — see below).

```go
ls, _ := key.NewFixed("abcde", 3)
v, _ := ls.Marshal(122)     // "eec"
i, _ := ls.Unmarshal("eec") // 122

buf := ls.Alphabet()[:0]
buf, _ = ls.MarshalAppend(buf[:0], 122) // append into a reusable buffer
```

## Strict decoding and canonical keys

`Unmarshal` enforces a bijection: a key always round-trips to itself, so there
is exactly one canonical spelling for each ID. Non-canonical, empty, wrongly
sized or out-of-range keys are rejected with a typed sentinel:

```go
ls, _ := key.NewDynamic("abc")

_, err := ls.Unmarshal("")   // key.ErrEmptyKey
_, err = ls.Unmarshal("aab") // key.ErrNonCanonical (canonical key for 1 is "b")
_, err = ls.Unmarshal("abX") // key.ErrUnknownChar

ok := ls.Valid("bab") // true; Valid is true iff Unmarshal succeeds
```

```go
func (l *Locksmith) Valid(key string) bool
```

`Valid` is a convenience predicate — it is true exactly when `Unmarshal` would
succeed, without allocating the decoded value.

## Inspecting a Locksmith

```go
func (l *Locksmith) Alphabet() string
func (l *Locksmith) Size() uint64
func (l *Locksmith) Total() uint64
func (l *Locksmith) Saturated() bool
```

- `Alphabet` — the alphabet the Locksmith was built with.
- `Size` — the fixed key length, or `0` for a dynamic Locksmith.
- `Total` — the number of representable IDs; valid IDs run `0..Total-1`. It
  saturates to `MaxUint64` for spaces of at least 2⁶⁴.
- `Saturated` — whether every `uint64` ID fits the key space.

```go
ls, _ := key.NewFixed("abc", 8)
ls.Size()  // 8
ls.Total() // 6561 (3^8)
```

## Iteration

```go
func (l *Locksmith) Iter(from, to uint64) iter.Seq2[uint64, string]
func (l *Locksmith) IterN(from, n uint64) iter.Seq2[uint64, string]
```

Both are range-over-func sequences yielding `(id, key)` pairs. `Iter` walks the
inclusive ID range `[from, to]`; `IterN` yields at most `n` keys starting at
`from` and is count-bounded, which is safer for bulk generation near the top of
the ID space.

```go
for id, k := range ls.Iter(0, 99) {
    fmt.Println(id, k)
}
for _, k := range ls.IterN(1_000_000, 50) { // 50 keys, no overflow risk
    _ = k
}
```

## Random keys

```go
func (l *Locksmith) Random(r io.Reader) (string, error)
func (l *Locksmith) RandomCrypto() (string, error)
```

`Random` draws a uniformly distributed ID from the supplied reader using
rejection sampling, so there is **no modulo bias**, and returns its key. Pass
`crypto/rand.Reader` for unpredictable keys. `RandomCrypto` is the
secure-by-default wrapper over `Random(crypto/rand.Reader)`.

```go
k, _ := ls.RandomCrypto()          // secure by default
k, _ = ls.Random(mathrand.Reader)  // fast, non-secure — for tests/sampling
```

## Ready-made alphabets

The package ships common alphabets; any string of unique characters also works:

| Constant | Alphabet |
|----------|----------|
| `Base16`      | `0123456789abcdef` |
| `Base32`      | `ABCDEFGHIJKLMNOPQRSTUVWXYZ234567` |
| `Base36`      | `0123456789abcdefghijklmnopqrstuvwxyz` |
| `Base62`      | `0-9 A-Z a-z` |
| `Crockford`   | Crockford's Base32 (no `I L O U`) |
| `Unambiguous` | drops look-alike characters for hand-typed codes |

```go
coupon := key.MustNewFixed(key.Unambiguous, 6) // human-friendly codes
```

## Errors

All are sentinels, matchable with `errors.Is`:

| Error | Raised when |
|-------|-------------|
| `ErrBlankAlphabet`  | the alphabet is empty |
| `ErrShortAlphabet`  | fewer than two characters |
| `ErrDuplicateChar`  | the alphabet has a repeated character |
| `ErrInvalidSize`    | a fixed size is out of range |
| `ErrIDTooLarge`     | the ID does not fit a fixed key space |
| `ErrEmptyKey`       | an empty key was decoded |
| `ErrInvalidLength`  | a fixed key has the wrong length |
| `ErrUnknownChar`    | the key contains a non-alphabet character |
| `ErrNonCanonical`   | the key is not the canonical spelling of its ID |
| `ErrKeyOutOfRange`  | the decoded ID is outside the representable range |

```go
if _, err := ls.Unmarshal(s); errors.Is(err, key.ErrNonCanonical) {
    // reject a non-canonical key
}
```

## Recipes and tips

**Obfuscate sequential IDs.** Use a dynamic Locksmith over `Base62` to turn a
database primary key into a short opaque token, and `Unmarshal` on the way back:

```go
ls := key.MustNewDynamic(key.Base62)
token, _ := ls.Marshal(userID)
```

**Human-typed codes.** Prefer `Unambiguous` (or `Crockford`) with `NewFixed` so
coupons and tickets have a constant width and avoid `0/O`, `1/l` confusion.

**Secure tokens.** Reach for `RandomCrypto` — it draws from `crypto/rand` with
no modulo bias, so keys are both unpredictable and uniform.

**Bulk generation.** Prefer `IterN(from, n)` over `Iter(from, to)` when
generating many keys near the top of the ID space; the count bound removes any
overflow footgun.

**Avoid allocations on hot paths.** `MarshalAppend` writes into a reusable
`[]byte`, which matters when minting keys in a tight loop.
