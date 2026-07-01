[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/key/v2)](https://goreportcard.com/report/github.com/goloop/key/v2) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/key/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://pkg.go.dev/github.com/goloop/key/v2) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)

# key

`key` converts `uint64` identifiers into short, reversible string keys over a
custom alphabet, and decodes them back. It is handy for shortening URLs,
obfuscating sequential database IDs, and minting human-readable codes such as
tickets or coupons.

Under the hood a `Locksmith` is a base-N positional codec: with an alphabet of
N characters, the value is written in base N. The mapping is deterministic and
bidirectional — each ID has exactly one canonical key, and each valid key
decodes back to its ID. All arithmetic is integer-only (`uint64` with
`math/bits` overflow checks), so the codec is exact across the **whole**
`uint64` range, not just for small values.

## Features

- **Fixed or dynamic length** — constant-width keys, or the shortest key the
  value allows.
- **Strict, canonical decoding** — one canonical spelling per ID; anything else
  is rejected with a typed error.
- **Whole-range accuracy** — integer-only math with overflow checks across all
  of `uint64`.
- **Iteration** — `Iter`/`IterN` as range-over-func sequences of `(id, key)`.
- **Random keys** — uniform, no modulo bias; `RandomCrypto` is secure by default.
- **Ready-made alphabets** — `Base16/32/36/62`, `Crockford`, `Unambiguous`.

## Installation

```shell
go get github.com/goloop/key/v2
```

```go
import key "github.com/goloop/key/v2"
```

Requires Go 1.24 or newer. The package has no third-party dependencies.

## Quick start

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

Fixed keys are constant width; dynamic keys are as short as possible:

```go
fixed, _ := key.NewFixed("abc", 8)
fixed.Marshal(10) // "aaaaabab"

dyn, _ := key.NewDynamic("abc")
dyn.Marshal(10)   // "bab"

// Human-friendly codes: constant width, no look-alike characters.
coupon := key.MustNewFixed(key.Unambiguous, 6)
```

Decoding is strict — a key always round-trips to itself:

```go
ls, _ := key.NewDynamic("abc")
_, err := ls.Unmarshal("aab") // key.ErrNonCanonical ("b" is canonical for 1)
ok := ls.Valid("bab")         // true
```

## Documentation

- Full reference and recipes: [DOC.md](DOC.md) · [DOC.UK.md](DOC.UK.md)
- Package API: [pkg.go.dev/github.com/goloop/key/v2](https://pkg.go.dev/github.com/goloop/key/v2)
- Changes between versions: [CHANGELOG.md](CHANGELOG.md)

## Contributing

Contributions are welcome. Please run `go test ./...`, `go vet ./...` and
`gofmt -l .` before submitting a pull request.

## License

`key` is released under the MIT License. See [LICENSE](LICENSE).
