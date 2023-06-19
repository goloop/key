[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/key)](https://goreportcard.com/report/github.com/goloop/key) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/key/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/key) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)

# Key

Package key allows to find the sequence that will be formed from the permutation of some characters defined in the alphabet by the specified iteration number.

The library can be used to create unique string identifiers based on a numeric identifier, or to create short URLs in redirect systems, or to mask certain indexes in data reports, etc.

## Theory

Use the arbitrary alphabet and size of key can be created a sequence of unique combinations, where each new combination has its unique numeric index (from 0 to N - where the N is maximum number of possible combinations).

If specify the iteration index (for example, it can be ID field from the some table of database) - will be returned the combination (key) for this index. And if the key is specified - decoding allows to determine the iteration index.

For example, the length of the key is 3 characters and there are several elements to iterate (alphabet): `a`, `b` and `c`. The package allows to answer the following questions:

  - How many maximum possible combinations of permutation of
    characters from the alphabet for a given key size?
  - What is the combination of permutation for N iteration?
  - What is the iteration index for some combination, for example "abc"?

For "abc" alphabet and 3 key size can be created the next iterations:

```
    0. aaa    1. aab    2. aac    3. aba    4. abb    5. abc
    6. aca    7. acb    8. acc    9. baa   10. bab   11. bac
   12. bba   13. bbb   14. bbc   15. bca   16. bcb   17. bcc
   18. caa   19. cab   20. cac   21. cba   22. cbb   23. cbc
   24. cca   25. ccb   26. ccc
```

So, the maximum number of iterations is 27. For the 10 iteration will correspond to the "bab" sequence and for example the "aba" combination is the 3 iteration.

## Install

Install key:

```shell
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
    ls, _ := key.New("abcde", 3)
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
ls, _ = key.New("abc") // size not specified
ls.Marshal(1)        // "b", <nil>
ls.Marshal(10)       // "bab", <nil>
ls.Marshal(100)      // "bacab", <nil>
ls.Marshal(1000)     // "bbabaab", <nil>
ls.Marshal(10000)    // "bbbcabbab", <nil>
ls.Marshal(100000)   // "bcaacabbcab", <nil>
ls.Marshal(1000000)  // "bcbccbacacaab", <nil>
ls.Marshal(10000000) // "caacbbaabbacbab", <nil>

ls.Unmarshal("b")               // 1, <nil>
ls.Unmarshal("bab")             // 10, <nil>
ls.Unmarshal("bacab")           // 100, <nil>
ls.Unmarshal("bbabaab")         // 1000, <nil>
ls.Unmarshal("bbbcabbab")       // 10000, <nil>
ls.Unmarshal("bcaacabbcab")     // 100000, <nil>
ls.Unmarshal("bcbccbacacaab")   // 1000000, <nil>
ls.Unmarshal("caacbbaabbacbab") // 10000000, <nil>
```

## Functions

- **New**(alphabet string, size ...uint) (*Locksmith, error)

  New returns a pointer to a Locksmith object as the first value and an error if something went wrong, or nil as the second value.

  As the first argument, the function takes the size of the key. If the size of the key is set to zero, the key size will be dynamic, i.e. the minimum key size will be one character, and the maximum will depend on the length of the alphabet and the possible maximum iteration index.

  The second value is the sequence elements for permutation (alphabet). The alphabet mustn't contain duplicate chars or be empty.


## Locksmith Methods

The Locksmith struct represents a key generation object. It provides methods for working with keys, including generating keys from IDs and retrieving IDs from keys.

### Alphabet

- **Alphabet**() string

  The `Alphabet` method returns the current alphabet value used by the Locksmith object. The alphabet is a string of unique characters from which the keys will be generated.

- **Size**() uint64

  The Size method returns the size of the key set for the Locksmith object. If the size is set to zero, the key size will be dynamic, meaning it can vary depending on the ID.

- **Total**() uint64

  The Total method returns the highest possible iteration number for the Locksmith object. It represents the total number of possible keys that can be generated.

- **Marshal**(id uint64) (string, error)

  The Marshal method converts an ID into a key. It takes an ID as input and generates a corresponding key based on the Locksmith's alphabet and size.

  The id is the ID to be converted into a key.

  The method returns a string representing the generated key and an error if something went wrong. If the function is successful, the error will be nil.

- **Unmarshal**(key string) (uint64, error)

  The Unmarshal method decodes a key and returns its corresponding ID. It converts a key back into its ID. The key should be a string composed of characters from the Locksmith's alphabet.

  The key is the key to be decoded into an ID.

  The method returns an integer representing the decoded ID and an error if something went wrong. If the function is successful, the error will be nil.
