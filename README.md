[//]: # (!!!Don't modify the README.md, use `make readme` to change/generate one!!!)


[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/key)](https://goreportcard.com/report/github.com/goloop/key) [![License](https://img.shields.io/badge/license-BSD-blue)](https://github.com/goloop/scs/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/key)

*Version: v1.1.0*

# Key

Package key allows to find the sequence that will be formed from the permutation
of some characters defined in the alphabet by the specified iteration number.

For example, the length of the key is 3 characters and there are several
elements to iterate (alphabet): `a`, `b` and `c`. The package allows to
answer the following questions:

	- How many maximum possible combinations of permutation of
	  characters from the alphabet for a given key size?
	- What is the combination of permutation for N iteration?
	- What is the iteration index for some combination, such as "abc"?

For "abc" alphabet and 3 key size can be created the next iterations:

```
    0. aaa    1. aab    2. aac    3. aba    4. abb    5. abc
    6. aca    7. acb    8. acc    9. baa   10. bab   11. bac
   12. bba   13. bbb   14. bbc   15. bca   16. bcb   17. bcc
   18. caa   19. cab   20. cac   21. cba   22. cbb   23. cbc
   24. cca   25. ccb   26. ccc
```

So, the maximum number of iterations is 27. For the 10 iteration
will correspond to the "bab" sequence and for example the "aba"
combination is the 3 iteration.

## Theory

Use the arbitrary alphabet and size of key can be created a sequence of
unique combinations, where each new combination has its unique numeric index
(from 0 to N - where the N is maximum number of possible combinations).

If specify the iteration index (for example, it can be ID field from
the some table of database) - will be returned the combination (key)
for this index. And if the key is specified - decoding allows to
determine the iteration index.

Install:

```
$ go get github.com/goloop/key
```

Examples:

```go
package main

import (
	"fmt"
	"github.com/goloop/key"
)

func main() {
    ls, _ := key.New(3, "abcde")
    v, _ := ls.Marshal(122)     // eec
    i, _ := ls.Unmarshal("eec") // 122
    
    fmt.Println(v, i)
    // Output: eec 122
}
```

If you specify the key size as 0 - the key length will be from one character
to unknown length (depends on the size of the alphabet).

Example:

```go
    ls, _ = key.New(0, "abc") // size not specified
    ls.Marshal(1) // "b", <nil>
    ls.Marshal(10) // "bab", <nil>
    ls.Marshal(100) // "bacab", <nil>
    ls.Marshal(1000) // "bbabaab", <nil>
    ls.Marshal(10000) // "bbbcabbab", <nil>
    ls.Marshal(100000) // "bcaacabbcab", <nil>
    ls.Marshal(1000000) // "bcbccbacacaab", <nil>
    ls.Marshal(10000000) // "caacbbaabbacbab", <nil>
```

## Usage

#### func  Version

    func Version() string

Version returns the version of the module.

#### type Locksmith

    type Locksmith struct {}

Locksmith is a key generation object.

#### func  New

    func New(size uint, alphabet string) (*Locksmith, error)

New returns a pointer to a Locksmith object as the first value and an error if
something went wrong, or nil as the second value.

As the first argument, the function takes the size of the key. If the size of
the key is set to zero, the key size will be dynamic, i.e. the minimum key size
will be one character, and the maximum will depend on the length of the alphabet
and the possible maximum iteration index.

The second value is the sequence elements for permutation (alphabet). The
alphabet mustn't contain duplicate chars or be empty.

#### func (*Locksmith) Alphabet

    func (ls *Locksmith) Alphabet() string

Alphabet returns current alphabet value.

#### func (*Locksmith) IsValid

    func (ls *Locksmith) IsValid() bool

IsValid returns true if Locksmith object is valid.

#### func (*Locksmith) Marshal

    func (ls *Locksmith) Marshal(id uint64) (string, error)

Marshal returns the key (sequence element) by ID.

#### func (*Locksmith) Size

    func (ls *Locksmith) Size() uint

Size return size of the key.

#### func (*Locksmith) Total

    func (ls *Locksmith) Total() uint64

Total returns the highest possible iteration number.

For example, for "abc" alphabet and key size as 3 - can be created the 27
iterations: aaa, aab, aac, ..., cca, ccb, ccc. So can be used indexs as 0 <= ID
< 27 to generate a key.

#### func (*Locksmith) Unmarshal

    func (ls *Locksmith) Unmarshal(key string) (uint64, error)

Unmarshal returns ID of the specified sequence.
