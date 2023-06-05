package key

import (
	"math"
	"testing"
)

// TestNew tests New function.
func TestNew(t *testing.T) {
	// Test when the alphabet is empty
	_, err := New("")
	if err == nil {
		t.Error("Expected an error when the alphabet is empty")
	}

	// Test when the alphabet has duplicate characters
	_, err = New("abcabc")
	if err == nil {
		t.Error("Expected an error when the alphabet has duplicates")
	}

	// Test when the size is less than zero
	_, err = New("abc", -1)
	if err == nil {
		t.Error("Expected an error when the size is less than zero")
	}

	// Test when the alphabet and size are valid
	ls, err := New("abc", 3)
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if ls.Size() != 3 {
		t.Errorf("Expected size to be 3, got %d", ls.Size())
	}

	if ls.Alphabet() != "abc" {
		t.Errorf("Expected alphabet to be 'abc', got %s", ls.Alphabet())
	}

	if ls.Total() != 27 {
		t.Errorf("Expected total to be 27, got %d", ls.Total())
	}
}

// TestTotalMax tests Total method with 0 size of key.
func TestTotalMax(t *testing.T) {
	ls, err := New("abc")
	if err != nil {
		t.Error(err)
	}

	if ls.Total() != math.MaxUint64 {
		t.Error("for a dynamic key size, the total must be as MaxUint64")
	}
}

// TestMarshalOverflow overflow testing Marshal method.
func TestMarshalOverflow(t *testing.T) {
	ls, err := New("abc", 3)
	if err != nil {
		t.Error(err)
	}

	_, err = ls.Marshal(33) // max 26
	if err == nil {
		t.Error("the ID is more than total value")
	}
}

// TestMarshalLogic tests Marshal method.
func TestMarshalLogic(t *testing.T) {
	tests := []struct {
		size   int
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
		ls, err := New("abcdefghijklmnopqrstuvwxyz0123456789", test.size)
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

func TestMarshal(t *testing.T) {
	// Test when the id is larger than the total number of keys.
	ls, _ := New("abc", 3)
	_, err := ls.Marshal(1000) // this will be more than total keys for size 3
	if err == nil {
		t.Error("Expected an error when the id is larger than total keys")
	}

	// Test when the id is valid and size is fixed
	key, err := ls.Marshal(13) // bbb
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	// Check key.
	if key != "bbb" {
		t.Errorf("Expected %s, got %s", "bbb", key)
	}

	if len(key) != 3 {
		t.Errorf("Expected key to be of length 3, got %d", len(key))
	}

	// Test when the size is zero (i.e., dynamic size).
	ls, _ = New("abc")
	key, err = ls.Marshal(10) // bab
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	// Check key.
	if key != "bab" {
		t.Errorf("Expected %s, got %s", "bab", key)
	}

	// The key length should be less than or equal to 3.
	if len(key) > 3 {
		t.Errorf("Expected key to be of length <= 3, got %d", len(key))
	}

	key, _ = ls.Marshal(10000000) // caacbbaabbacbab
	if key != "caacbbaabbacbab" {
		t.Errorf("Expected %s, got %s", "caacbbaabbacbab", key)
	}
}

// TestUnmarshalWrongSize tests Unmarshal method with wrong size of key.
func TestUnmarshalWrongSize(t *testing.T) {
	ls, err := New("abc", 3)
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
	ls, err := New("abc", 3)
	if err != nil {
		t.Error(err)
	}

	_, err = ls.Unmarshal("aXc") // X is incorrect char
	if err == nil {
		t.Error("key has wrong char")
	}
}

// TestUnmarshalLogic tests Unmarshal method.
func TestUnmarshalLogic(t *testing.T) {
	tests := []struct {
		size   int
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
		ls, err := New("abcdefghijklmnopqrstuvwxyz0123456789", test.size)
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

func TestUnmarshal(t *testing.T) {
	// Test when the key is not of the correct length.
	ls, _ := New("abc", 3)
	_, err := ls.Unmarshal("ab") // this is not of length 3
	if err == nil {
		t.Error("Expected an error when the key is not of the correct length")
	}

	// Test when the key contains a character not in the alphabet.
	_, err = ls.Unmarshal("abd") // 'd' is not in the alphabet
	if err == nil {
		t.Error("Expected an error when the key contains " +
			"a character not in the alphabet")
	}

	// Test when the key is valid and size is fixed.
	id, err := ls.Unmarshal("aab")
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	// The ID should match the ID we used in the Marshal test.
	expectedId := uint64(1)
	if id != expectedId {
		t.Errorf("Expected id to be %d, got %d", expectedId, id)
	}

	// Test when the size is zero (i.e., dynamic size).
	ls, _ = New("abc")
	id, err = ls.Unmarshal("bab") // id 10
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	// The ID should match the ID we used in the Marshal test.
	expectedId = uint64(10)
	if id != expectedId {
		t.Errorf("Expected id to be %d, got %d", expectedId, id)
	}

	id, err = ls.Unmarshal("bbb") // id 13
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	// The ID should match the ID we used in the Marshal test.
	expectedId = uint64(13)
	if id != expectedId {
		t.Errorf("Expected id to be %d, got %d", expectedId, id)
	}

	id, err = ls.Unmarshal("caacbbaabbacbab") // 10000000
	expectedId = uint64(10000000)
	if id != expectedId {
		t.Errorf("Expected id to be %d, got %d", expectedId, id)
	}
}
