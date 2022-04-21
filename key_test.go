package key

import (
	"math"
	"strings"
	"testing"
)

// TestVersion tests the package version.
// Note: each time you change the major version, you need to fix the tests.
func TestVersion(t *testing.T) {
	var expected = "v1." // change it for major version

	version := Version()
	if strings.HasPrefix(version, expected) != true {
		t.Error("incorrect version")
	}

	if len(strings.Split(version, ".")) != 3 {
		t.Error("version format should be as " +
			"v{major_version}.{minor_version}.{patch_version}")
	}
}

// TestNewWithEmptyAlphabet tests New method with empty alphabet.
func TestNewWithEmptyAlphabet(t *testing.T) {
	if _, err := New(3, ""); err == nil {
		t.Error("the alphabet cannot be empty")
	}
}

// TestNewWithDuplicateCharacters tests New method with duplicate
// characters in alphabet.
func TestNewWithDuplicateCharacters(t *testing.T) {
	if _, err := New(3, "abcade"); err == nil {
		t.Error("duplicates in the alphabet are not allowed")
	}
}

// TestAlphabet tests Alphabet method.
func TestAlphabet(t *testing.T) {
	ls, err := New(3, "abc")
	if err != nil {
		t.Error(err)
	}

	if ls.Alphabet() != "abc" {
		t.Error("the alphabet doesn't match")
	}
}

// TestSize tests Size method.
func TestSize(t *testing.T) {
	ls, err := New(3, "abc")
	if err != nil {
		t.Error(err)
	}

	if ls.Size() != 3 {
		t.Error("the size doesn't match")
	}
}

// TestTotal tests Total method.
func TestTotal(t *testing.T) {
	ls, err := New(3, "abc")
	if err != nil {
		t.Error(err)
	}

	if ls.Total() != 27 {
		t.Error("for key size as 3 and alphabet from 3 chars, " +
			"the total must be as 27")
	}
}

// TestTotalMax tests Total method with 0 size of key.
func TestTotalMax(t *testing.T) {
	ls, err := New(0, "abc")
	if err != nil {
		t.Error(err)
	}

	if ls.Total() != math.MaxUint64 {
		t.Error("for a dynamic key size, the total must be as MaxUint64")
	}
}

// TestMarshalOverflow overflow testing Marshal method.
func TestMarshalOverflow(t *testing.T) {
	ls, err := New(3, "abc")
	if err != nil {
		t.Error(err)
	}

	_, err = ls.Marshal(33) // max 26
	if err == nil {
		t.Error("the ID is more than total value")
	}
}

// TestMarshal tests Marshal method.
func TestMarshal(t *testing.T) {
	var tests = []struct {
		size   uint
		value  uint64
		expect string
	}{
		{0, 3333333, "b9qav"}, // dynamic size
		{3, 0, "aaa"},
		{3, 1, "aab"},
		{3, 10, "aak"},
		{5, 1024, "aaa2q"},
		{7, 1024, "aaaaa2q"},
	}

	for _, test := range tests {
		ls, err := New(test.size, "abcdefghijklmnopqrstuvwxyz0123456789")
		if err != nil {
			t.Error(err)
		}

		abc, err := ls.Marshal(test.value)
		if err != nil {
			t.Error(err)
		}

		if abc != test.expect {
			t.Errorf("expected result %v but %v", test.expect, abc)
		}
	}
}

// TestUnmarshalWrongSize tests Unmarshal method with wrong size of key.
func TestUnmarshalWrongSize(t *testing.T) {
	ls, err := New(3, "abc")
	if err != nil {
		t.Error(err)
	}

	_, err = ls.Unmarshal("aabbcc") // key must contain only 3 chars
	if err == nil {
		t.Error("key must contain only 3 chars")
	}
}

// TestUnmarshalWrongChar tests Unmarshal method with wrong char in the key.
func TestUnmarshalWrongChar(t *testing.T) {
	ls, err := New(3, "abc")
	if err != nil {
		t.Error(err)
	}

	_, err = ls.Unmarshal("aXc") // X is incorrect char
	if err == nil {
		t.Error("key has wrong char")
	}
}

// TestUnmarshal tests Unmarshal method.
func TestUnmarshal(t *testing.T) {
	var tests = []struct {
		size   uint
		value  string
		expect uint64
	}{
		{0, "b9qav", 3333333}, // dynamic size
		{3, "aaa", 0},
		{3, "aab", 1},
		{3, "aak", 10},
		{5, "aaa2q", 1024},
		{7, "aaaaa2q", 1024},
	}

	for _, test := range tests {
		ls, err := New(test.size, "abcdefghijklmnopqrstuvwxyz0123456789")
		if err != nil {
			t.Error(err)
		}

		abc, err := ls.Unmarshal(test.value)
		if err != nil {
			t.Error(err)
		}

		if abc != test.expect {
			t.Errorf("expected result %v but %v", test.expect, abc)
		}
	}
}

// TestIsValis tests IsValid mthod.
func TestIsValid(t *testing.T) {
	ls := &Locksmith{}
	if ls.IsValid() {
		t.Error("an empty Locksmith object cannot be valid")
	}
}
