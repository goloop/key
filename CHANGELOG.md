# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.1.0]

Minor release: closes the last bijection gap and fixes the documentation.

### Fixed
- The constructors reject an alphabet that is not valid UTF-8 or contains the
  replacement character U+FFFD (new `ErrInvalidAlphabet` sentinel). Such an
  alphabet let a malformed byte in a key decode to U+FFFD and be accepted, so
  several distinct byte strings aliased one id, breaking the strict bijection.
- Corrected the Base62 vectors in the package documentation (`5BAN` /
  `00005BAN`, previously `5BAH` / `00005BAH`) and pinned them with an example
  test so they cannot drift again.
- Fixed the non-compiling `MarshalAppend` snippet and the reference to a
  non-existent reader in `DOC.md`.

## [2.0.0]

A correctness-focused rewrite. The codec is now exact across the whole `uint64`
range and decoding is a strict bijection. The module path is now
`github.com/goloop/key/v2`.

### Changed (breaking)
- **Constructors split** into `NewDynamic(alphabet)` and
  `NewFixed(alphabet, size)`, replacing `New(alphabet, ...int)` and its
  "size is the sum of the arguments" semantics.
- **Strict decoding.** `Unmarshal` now rejects the empty string, keys of the
  wrong length (fixed mode), non-canonical keys with redundant leading
  lead-characters (dynamic mode), unknown characters, and values that overflow
  `uint64`. Each case returns a typed sentinel error (`errors.Is`).
- **Constructors return `nil`** (not a half-built `*Locksmith`) on error.
- **`Total()` saturates** to `MaxUint64` for spaces of at least 2⁶⁴ instead of
  returning `0`; new `Saturated()` reports that case.
- Minimum Go version raised to 1.24.

### Added
- `MarshalAppend(dst, id)`, `Valid(key)`, `Saturated()`.
- `Iter(from, to) iter.Seq2[uint64, string]` — range-over-func over an
  inclusive ID range, and `IterN(from, n)` — a count-bounded variant.
- `Random(r io.Reader)` — a key for a uniformly random ID via rejection
  sampling — and `RandomCrypto()`, its secure-by-default `crypto/rand` wrapper.
- `MustNewDynamic` / `MustNewFixed`.
- Ready-made alphabets: `Base16`, `Base32`, `Base36`, `Base62`, `Crockford`,
  `Unambiguous`.
- Sentinel errors: `ErrBlankAlphabet`, `ErrShortAlphabet`, `ErrDuplicateChar`,
  `ErrInvalidSize`, `ErrIDTooLarge`, `ErrEmptyKey`, `ErrInvalidLength`,
  `ErrUnknownChar`, `ErrNonCanonical`, `ErrKeyOutOfRange`.

### Fixed
- Lost precision for IDs above 2⁵³: `Marshal` used `float64`/`math.Mod`, which
  collided distinct IDs (e.g. 2⁵³ and 2⁵³+1). Arithmetic is now integer-only.
- Silent `int` overflow above 2⁶³ in `Marshal`, producing wrong keys.
- `Total()` returning `0` (bricking the instance) or a plausible-but-wrong
  value when `alphabet^size` reached or exceeded 2⁶⁴.
- `Unmarshal` silently overflowing on over-long keys instead of reporting an
  error.
- A single-character alphabet being accepted and then hanging `Marshal` in an
  infinite loop; such an alphabet is now rejected at construction.

### Performance
- `Marshal` rebuilt to a single allocation and linear time (no per-character
  string concatenation).
- `Unmarshal` uses Horner's scheme (no per-character `pow`) and iterates the
  key as a string, so it is allocation-free even for long keys (no `[]rune`
  materialization).
