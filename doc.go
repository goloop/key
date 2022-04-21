/*
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
correspond to the "baa" sequence and the "abc" combination - this is the
fifth iteration.


## Theory

Use the arbitrary alphabet (slice of runes to create the key) and size of key
can be created a sequence of unique combinations, where each new combination
has its unique numeric index (from 0 to N - where the N is maximum number of
possible combinations).

If specify the iteration index (for example, it can be ID field from some
database table) - will be returned the combination (key) for this index.
And if the key is specified, decoding allows to determine the iteration index.

Examples:

```go

package main

import (
	"fmt"
	"github.com/goloop/key"
)

func main() {
    k, _ := key.New(3, 'a', 'b', 'c', 'd', 'e')
    v, _ := k.Marshal(122)     // v == eec
    i, _ := k.Unmarshal("eec") // i == 122

    // ...
}

```

If you don't specify a custom alphabet for the New method, a random alphabet
consisting of uppercase Latin characters a-z and numbers 0-9 will be used.
Where k1 := key.New(3); k2 := key.New(3); k1 == k2 is false.

If you specify the key size as 0 - the key length will be from one character
to unknown length (depends on the size of the alphabet and is generated on
the k.LastID() iteration).

Example:

```go

	// Size not specified.
	k, _ := key.New(0, 'a', 'b', 'c')
	k.Marshal(1) // "b", <nil>
	k.Marshal(10) // "bab", <nil>
	k.Marshal(100) // "bacab", <nil>
	k.Marshal(1000) // "bbabaab", <nil>
	k.Marshal(10000) // "bbbcabbab", <nil>
	k.Marshal(100000) // "bcaacabbcab", <nil>
	k.Marshal(1000000) // "bcbccbacacaab", <nil>
	k.Marshal(10000000) // "caacbbaabbacbab", <nil>

```

*/
package key

const version = "0.0.3"

// Version returns the version of the module.
func Version() string {
	return "v" + version
}
