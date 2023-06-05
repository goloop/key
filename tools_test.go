package key

import (
	"reflect"
	"testing"
)

// TestReverse tests the private reverse method.
func TestReverse(t *testing.T) {
	tests := []struct {
		value  []rune
		expect []rune
	}{
		{[]rune{'c'}, []rune{'c'}},
		{[]rune{'a', 'b', 'a'}, []rune{'a', 'b', 'a'}},
		{[]rune{'a', 'b', 'c'}, []rune{'c', 'b', 'a'}},
	}

	for _, test := range tests {
		// Reverse by link.
		abc := make([]rune, len(test.value))
		copy(abc, test.value)
		reverse(abc)
		if !reflect.DeepEqual(abc, test.expect) {
			t.Errorf("abc: expected result %v but %v", test.expect, abc)
		}

		// Returned result.
		def := make([]rune, len(test.value))
		copy(def, test.value)
		def = reverse(def)
		if !reflect.DeepEqual(def, test.expect) {
			t.Errorf("def: expected result %v but %v", test.expect, def)
		}
	}
}

// TestUnlead tests the private unlead method.
func TestUnlead(t *testing.T) {
	tests := []struct {
		value  []rune
		expect []rune
	}{
		{[]rune{'b', 'c'}, []rune{'b', 'c'}},
		{[]rune{'a', 'a', 'b', 'c'}, []rune{'b', 'c'}},
		{[]rune{'a', 'a'}, []rune{'a'}},
	}

	for _, test := range tests {
		abc := unlead('a', test.value)
		if !reflect.DeepEqual(abc, test.expect) {
			t.Errorf("expected result %v but %v\n", test.expect, abc)
		}
	}
}

// TestPow tests pow functions.
func TestPow(t *testing.T) {
	// Test when the exponent is 0
	result := pow(2, 0)
	if result != 1 {
		t.Errorf("Expected result to be 1, got %d", result)
	}

	// Test when the exponent is positive
	result = pow(2, 5)
	if result != 32 {
		t.Errorf("Expected result to be 32, got %d", result)
	}

	// Test when the exponent is negative
	result = pow(2, -3)
	if result != 0 {
		t.Errorf("Expected result to be 0, got %d", result)
	}

	// Test when the base is negative
	result = pow(-2, 4)
	if result != 16 {
		t.Errorf("Expected result to be 16, got %d", result)
	}

	// Test when the base is 0
	result = pow(0, 5)
	if result != 0 {
		t.Errorf("Expected result to be 0, got %d", result)
	}
}
