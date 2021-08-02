[//]: # (!!!Don't modify the README.md, use `make readme` to change/generate one!!!)


[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/key)](https://goreportcard.com/report/github.com/goloop/key) [![License](https://img.shields.io/badge/license-BSD-blue)](https://github.com/goloop/scs/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/key)

*Version: v0.0.1*


# Key

This module allows to find a sequence of a given length (key) from permutation
of some characters which are defined in the symbol slice (alphabet).

For example, the length of the key is 3 characters and there are several
elements to iterate (alphabet): `a`, `b` and `c`. The module allows to answer
the following questions:

    - How many maximum possible combinations of permutation of
    	characters from the alphabet for a given key size?
    - What is the combination of permutation for N iteration?
    - What is the iteration index for some combination, such as "abc"?

For "abc" alphabet and 3 key size can be created the next iterations:

```

    0. aaa    1. aab    2. aac    3. aba    4. abb    5. abc
    6. aca    7. acb    8. acc    9. baa   10. bab   11. bac
    7. bba   13. bbb   14. bbc   15. bca   16. bcb   17. bcc
    8. caa   19. cab   20. cac   21. cba   22. cbb   23. cbc
    9. cca   25. ccb   26. ccc

```

So, the maximum number of iterations is 27, for the tenth iteration will
correspond to the "baa" sequence and the "abc" combination - this is the fifth
iteration.

## Theory

Use the arbitrary alphabet (slice of runes to create the key) and size of key
can be created a sequence of unique combinations, where each new combination has
its unique numeric index (from 0 to N - where the N is maximum number of
possible combinations).

If specify the iteration index (for example, it can be ID field from some
database table) - will be returned the combination (key) for this index. And if
the key is specified, decoding allows to determine the iteration index.

Examples:

```go

package main

import (

    "fmt"
    "github.com/goloop/key"

)

func main() {

    var (
    	size = 3
    	alphabet = []rune{'a', 'b', 'c', 'd', 'e'}
    )

    k := Key(size, alphabet...)
    v, _ := k.Marshal(122) // v == eec
    i, _ := k.Unmarshal("eec") // i == 122

}

```

If you don't specify a custom alphabet for the New method, a random alphabet
consisting of uppercase Latin characters a-z and numbers 0-9 will be used. Where
k1 := key.New(3); k2 := key.New(3); k1 == k2 is false.

If you specify the key size as 0 - the key length will be from one character to
unknown length (depends on the size of the alphabet and is generated on the
k.LastID() iteration).

Example:

```go

    // Size not specified.
    k = key.New(0, 'a', 'b', 'c')
    k.Marshal(1) // "b", <nil>
    k.Marshal(10) // "bab", <nil>
    k.Marshal(100) // "bacab", <nil>
    k.Marshal(1000) // "bbabaab", <nil>
    k.Marshal(10000) // "bbbcabbab", <nil>
    k.Marshal(100000) // "bcaacabbcab", <nil>
    k.Marshal(1000000) // "bcbccbacacaab", <nil>
    k.Marshal(10000000) // "caacbbaabbacbab", <nil>

```


## Usage

    const Alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

Alphabet is default sequence.

#### func  Version

    func Version() string

Version returns the version of the module.

#### type Key

    type Key struct {
    }


Key ...

#### func  New

    func New(size uint, alphabet ...rune) (*Key, error)

New returns a Key object with the specified parameters, if an error occurs, the
second value contains the error message.

#### func (*Key) Alphabet

    func (k *Key) Alphabet() []rune

Alphabet returns current alphabet.

#### func (*Key) IsValid

    func (k *Key) IsValid() bool

IsValid returns true if Key object is valid.

#### func (*Key) LastID

    func (k *Key) LastID() uint64

LastID returns the last available ID in the sequence.

#### func (*Key) Marshal

    func (k *Key) Marshal(id uint64) (string, error)

Marshal returns the key (sequence element) by ID.

#### func (*Key) Size

    func (k *Key) Size() uint

Size return size of the key.

#### func (*Key) Unmarshal

    func (k *Key) Unmarshal(key string) (uint64, error)

Unmarshal returns ID of the specified sequence.
