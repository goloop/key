package key

import (
	"reflect"
	"testing"
)

// TestReverse tests the private reverse method.
func TestReverse(t *testing.T) {
	var tests = []struct {
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
	var tests = []struct {
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

func TestShuffle(t *testing.T) {
	var abc = []rune{'a', 'b', 'c'}

	origin := make([]rune, len(abc))
	copy(origin, abc)

	shuffle(abc)
	if reflect.DeepEqual(abc, origin) {
		t.Error("unmixed slice", abc, origin)
	}

	if len(abc) != len(origin) {
		t.Error("slice length broken")
	}
}
