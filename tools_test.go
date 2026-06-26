package key

import (
	"math"
	"testing"
)

// TestPowU64 checks binary exponentiation and, crucially, its overflow flag:
// the whole point of powU64 over a plain pow is that it reports when the true
// result does not fit into uint64 instead of silently wrapping.
func TestPowU64(t *testing.T) {
	tests := []struct {
		base, exp uint64
		want      uint64
		overflow  bool
	}{
		{2, 0, 1, false},       // anything**0 == 1
		{0, 5, 0, false},       // 0**5 == 0
		{1, 1 << 20, 1, false}, // 1**n == 1, must terminate (no base-1 trap)
		{2, 5, 32, false},
		{3, 3, 27, false},
		{10, 19, 10000000000000000000, false}, // 1e19 < MaxUint64
		{2, 63, 1 << 63, false},
		{62, 10, 839299365868340224, false}, // 62**10 fits
		// Boundaries that wrap a naive int/uint64 multiply:
		{2, 64, 0, true},  // 2**64 == exactly MaxUint64+1
		{16, 16, 0, true}, // 16**16 == 2**64 (the BUG-02 trigger)
		{36, 16, 0, true}, // 36**16 is astronomically larger
		{62, 11, 0, true}, // 62**11 overflows
		{10, 20, 0, true}, // 1e20 > MaxUint64
	}

	for _, tt := range tests {
		got, overflow := powU64(tt.base, tt.exp)
		if overflow != tt.overflow {
			t.Errorf("powU64(%d, %d) overflow = %v, want %v",
				tt.base, tt.exp, overflow, tt.overflow)
		}
		if !tt.overflow && got != tt.want {
			t.Errorf("powU64(%d, %d) = %d, want %d",
				tt.base, tt.exp, got, tt.want)
		}
	}
}

// TestPowU64ExactBoundary pins the exact uint64 ceiling: 2**63 must fit and
// equal the known constant, while one more doubling must report overflow.
func TestPowU64ExactBoundary(t *testing.T) {
	if got, of := powU64(2, 63); of || got != 1<<63 {
		t.Fatalf("powU64(2,63) = %d, overflow %v; want %d, false", got, of, uint64(1)<<63)
	}
	if _, of := powU64(2, 64); !of {
		t.Fatalf("powU64(2,64) must overflow (true value is MaxUint64+1)")
	}
	// Sanity: the largest power of two that fits is one below the overflow.
	if 1<<63 >= uint64(math.MaxUint64) {
		t.Fatalf("sanity check broken")
	}
}
