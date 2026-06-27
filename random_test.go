package key

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"testing"
)

// u64Reader builds an io.Reader that yields the given uint64 values as 8-byte
// big-endian blocks, in order. It feeds deterministic "randomness" to Random.
func u64Reader(vals ...uint64) io.Reader {
	buf := make([]byte, 0, len(vals)*8)
	for _, v := range vals {
		var b [8]byte
		binary.BigEndian.PutUint64(b[:], v)
		buf = append(buf, b[:]...)
	}
	return bytes.NewReader(buf)
}

// errReader fails after returning n good bytes, to test error propagation.
type errReader struct {
	left int
}

func (e *errReader) Read(p []byte) (int, error) {
	if e.left <= 0 {
		return 0, errors.New("boom")
	}
	n := len(p)
	if n > e.left {
		n = e.left
	}
	e.left -= n
	return n, nil
}

// TestRandomBounded checks that Random maps the raw uint64 through the bounded
// uniform reduction and then Marshals it.
func TestRandomBounded(t *testing.T) {
	ls, _ := NewFixed("0123456789", 1) // Total 10, valid ids 0..9

	// 42 % 10 == 2 -> key "2".
	k, err := ls.Random(u64Reader(42))
	if err != nil {
		t.Fatal(err)
	}
	if k != "2" {
		t.Fatalf("Random = %q, want %q", k, "2")
	}
}

// TestRandomSaturated checks that on a saturated space the raw uint64 is used
// directly, including MaxUint64 which a bounded reduction could not produce.
func TestRandomSaturated(t *testing.T) {
	ls, _ := NewFixed(Base16, 16) // saturated (2**64)
	k, err := ls.Random(u64Reader(math.MaxUint64))
	if err != nil {
		t.Fatal(err)
	}
	if k != "ffffffffffffffff" {
		t.Fatalf("Random = %q, want all-f", k)
	}
	if id, _ := ls.Unmarshal(k); id != math.MaxUint64 {
		t.Fatalf("round trip id = %d, want MaxUint64", id)
	}
}

// TestRandomRejection forces the rejection branch: the first draw lands in the
// discarded tail [limit, MaxUint64] and must be skipped, the second is used.
func TestRandomRejection(t *testing.T) {
	ls, _ := NewFixed("0123456789", 1) // Total 10
	// MaxUint64 % 10 == 5, so limit == MaxUint64-5; MaxUint64 is rejected.
	k, err := ls.Random(u64Reader(math.MaxUint64, 7))
	if err != nil {
		t.Fatal(err)
	}
	if k != "7" { // 7 < limit, 7 % 10 == 7
		t.Fatalf("Random after rejection = %q, want %q", k, "7")
	}
}

// TestRandomError propagates a reader error from both code paths.
func TestRandomError(t *testing.T) {
	bounded, _ := NewFixed("0123456789", 1)
	if _, err := bounded.Random(&errReader{left: 0}); err == nil {
		t.Error("bounded Random: expected error from reader")
	}
	saturated, _ := NewFixed(Base16, 16)
	if _, err := saturated.Random(&errReader{left: 3}); err == nil {
		t.Error("saturated Random: expected error from short reader")
	}
}

// TestRandomDistribution is a light statistical check: every draw must be a
// valid key, decode within range, and the sampler must cover the whole small
// space without ever escaping it.
func TestRandomDistribution(t *testing.T) {
	ls, _ := NewFixed("abcd", 2) // Total 16
	counts := make(map[uint64]int)

	// A simple counter-based reader gives reproducible, well-spread draws.
	var seed uint64
	reader := readerFunc(func(p []byte) (int, error) {
		binary.BigEndian.PutUint64(p, seed*2654435761+12345)
		seed++
		return 8, nil
	})

	for i := 0; i < 20000; i++ {
		k, err := ls.Random(reader)
		if err != nil {
			t.Fatal(err)
		}
		if !ls.Valid(k) {
			t.Fatalf("Random produced invalid key %q", k)
		}
		id, _ := ls.Unmarshal(k)
		if id >= ls.Total() {
			t.Fatalf("id %d out of range (Total %d)", id, ls.Total())
		}
		counts[id]++
	}
	if len(counts) != int(ls.Total()) {
		t.Fatalf("covered %d of %d ids", len(counts), ls.Total())
	}
}

// TestRandomCrypto exercises the secure wrapper end-to-end: every draw is a
// valid, in-range key, and the source produces good variety.
func TestRandomCrypto(t *testing.T) {
	ls, _ := NewFixed("abcd", 4) // Total 256
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		k, err := ls.RandomCrypto()
		if err != nil {
			t.Fatal(err)
		}
		if !ls.Valid(k) {
			t.Fatalf("RandomCrypto produced invalid key %q", k)
		}
		if id, _ := ls.Unmarshal(k); id >= ls.Total() {
			t.Fatalf("id %d out of range (Total %d)", id, ls.Total())
		}
		seen[k] = true
	}
	if len(seen) < 50 {
		t.Fatalf("only %d distinct keys from 1000 draws; entropy source suspect", len(seen))
	}
}

// TestUniformUint64Single covers the n == 1 short-circuit directly.
func TestUniformUint64Single(t *testing.T) {
	got, err := uniformUint64(u64Reader(), 1)
	if err != nil || got != 0 {
		t.Fatalf("uniformUint64(_, 1) = %d, %v; want 0, nil", got, err)
	}
}

// readerFunc adapts a function to io.Reader.
type readerFunc func([]byte) (int, error)

func (f readerFunc) Read(p []byte) (int, error) { return f(p) }
