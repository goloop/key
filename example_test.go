package key_test

import (
	"errors"
	"fmt"

	key "github.com/goloop/key/v2"
)

func ExampleNewDynamic() {
	ls, _ := key.NewDynamic("abc")

	k, _ := ls.Marshal(10)
	fmt.Println(k)
	// Output: bab
}

func ExampleNewFixed() {
	ls, _ := key.NewFixed("abc", 4)

	k, _ := ls.Marshal(10) // left-padded to 4 characters
	fmt.Println(k)
	// Output: abab
}

func ExampleLocksmith_Unmarshal() {
	ls, _ := key.NewDynamic("abc")

	id, _ := ls.Unmarshal("bab")
	fmt.Println(id)
	// Output: 10
}

func ExampleLocksmith_Unmarshal_strict() {
	ls, _ := key.NewDynamic("abc")

	// "aab" is not canonical: the canonical key for 1 is "b".
	_, err := ls.Unmarshal("aab")
	fmt.Println(errors.Is(err, key.ErrNonCanonical))
	// Output: true
}

func ExampleLocksmith_Valid() {
	ls, _ := key.NewDynamic("abc")

	fmt.Println(ls.Valid("bab")) // well-formed
	fmt.Println(ls.Valid("aab")) // non-canonical
	fmt.Println(ls.Valid(""))    // empty
	// Output:
	// true
	// false
	// false
}

func ExampleLocksmith_MarshalAppend() {
	ls := key.MustNewFixed(key.Base62, 6)

	buf := []byte("id=")
	buf, _ = ls.MarshalAppend(buf, 1)
	fmt.Println(string(buf))
	// Output: id=000001
}
