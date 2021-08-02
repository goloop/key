package key

import (
	"testing"
	"unicode/utf8"
)

// TestNewWithoutAlphabet tests New method without custom alphabet.
func TestNewWithoutAlphabet(t *testing.T) {
	key, err := New(3)
	if err != nil {
		t.Error(err)
	}

	defLen := utf8.RuneCountInString(defaultAlphabetCharacters)
	if len(key.Alphabet()) != defLen {
		t.Error("the default alphabet should be used")
	}
}

// TestMarshal ...
func TestMarshal(t *testing.T) {
	var tests = []struct {
		size   int
		value  uint64
		expect string
	}{
		{3, 0, "aaa"},
		{3, 1, "aab"},
		{3, 10, "aak"},
		{5, 1024, "aaa2q"},
		{7, 1024, "aaaaa2q"},
	}

	for _, test := range tests {
		key, err := New(uint(test.size), []rune(defaultAlphabetCharacters)...)
		if err != nil {
			t.Error(err)
		}

		abc, err := key.Marshal(test.value)
		if err != nil {
			t.Error(err)
		}

		if abc != test.expect {
			t.Errorf("expected result %v but %v", test.expect, abc)
		}
	}
}

// TestUnmarshal ...
func TestUnmarshal(t *testing.T) {
	var tests = []struct {
		size   int
		value  string
		expect uint64
	}{
		{3, "aaa", 0},
		{3, "aab", 1},
		{3, "aak", 10},
		{5, "aaa2q", 1024},
		{7, "aaaaa2q", 1024},
	}

	for _, test := range tests {
		key, err := New(uint(test.size), []rune(defaultAlphabetCharacters)...)
		if err != nil {
			t.Error(err)
		}

		abc, err := key.Unmarshal(test.value)
		if err != nil {
			t.Error(err)
		}

		if abc != test.expect {
			t.Errorf("expected result %v but %v", test.expect, abc)
		}
	}
}
