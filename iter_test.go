package key

import (
	"math"
	"testing"
)

// TestIterFullSpace enumerates the entire space of a small fixed codec and
// checks the pairs match Marshal exactly and in order.
func TestIterFullSpace(t *testing.T) {
	ls, _ := NewFixed("abc", 3) // Total 27

	var got int
	for id, k := range ls.Iter(0, ls.Total()-1) {
		if id != uint64(got) {
			t.Fatalf("out-of-order id: got %d, want %d", id, got)
		}
		want, _ := ls.Marshal(id)
		if k != want {
			t.Fatalf("Iter id %d key = %q, want %q", id, k, want)
		}
		got++
	}
	if got != 27 {
		t.Fatalf("iterated %d ids, want 27", got)
	}
}

// TestIterSubRange checks the inclusive [from, to] semantics on a small range.
func TestIterSubRange(t *testing.T) {
	ls, _ := NewFixed("abc", 3)

	var ids []uint64
	for id := range ls.Iter(5, 10) {
		ids = append(ids, id)
	}
	want := []uint64{5, 6, 7, 8, 9, 10} // inclusive of both ends
	if len(ids) != len(want) {
		t.Fatalf("got %v, want %v", ids, want)
	}
	for i := range want {
		if ids[i] != want[i] {
			t.Fatalf("got %v, want %v", ids, want)
		}
	}
}

// TestIterClampsToTotal verifies that an over-large to is clamped to the last
// valid id of a bounded codec.
func TestIterClampsToTotal(t *testing.T) {
	ls, _ := NewFixed("abc", 3) // Total 27, last id 26

	var last uint64
	var n int
	for id := range ls.Iter(20, 1_000_000) {
		last = id
		n++
	}
	if last != 26 {
		t.Fatalf("last id = %d, want 26", last)
	}
	if n != 7 { // 20..26 inclusive
		t.Fatalf("count = %d, want 7", n)
	}
}

// TestIterEmpty covers the from > to (after clamping) early return.
func TestIterEmpty(t *testing.T) {
	ls, _ := NewFixed("abc", 3)
	for range ls.Iter(10, 5) {
		t.Fatal("expected no iterations for from > to")
	}
	// from past the end of a bounded space.
	for range ls.Iter(100, 200) {
		t.Fatal("expected no iterations for from past Total")
	}
}

// TestIterEarlyStop checks that breaking out of the range stops the producer
// (yield returning false is honoured).
func TestIterEarlyStop(t *testing.T) {
	ls, _ := NewDynamic("abc") // saturated: exercises the !full == false branch

	var seen []uint64
	for id := range ls.Iter(0, math.MaxUint64) {
		seen = append(seen, id)
		if id == 3 {
			break
		}
	}
	if len(seen) != 4 { // 0,1,2,3
		t.Fatalf("collected %v, want 0..3", seen)
	}
}

// TestIterReachesMaxUint64 confirms the inclusive bound can yield the very last
// id on a saturated space, which a half-open range could never reach.
func TestIterReachesMaxUint64(t *testing.T) {
	ls, _ := NewDynamic("ab")

	var got bool
	for id := range ls.Iter(math.MaxUint64-2, math.MaxUint64) {
		if id == math.MaxUint64 {
			got = true
		}
	}
	if !got {
		t.Fatal("Iter never yielded MaxUint64")
	}
}
