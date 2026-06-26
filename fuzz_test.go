package key

import (
	"math"
	"testing"
)

// fuzzCodecs is the matrix of codecs exercised by the fuzz targets: a
// saturated dynamic one, a saturated fixed one, a bounded fixed one and a
// unicode one. Each stresses a different arithmetic and padding path.
func fuzzCodecs() []*Locksmith {
	return []*Locksmith{
		MustNewDynamic("ab"),          // base 2, longest keys
		MustNewDynamic(Base62),        // wide base
		MustNewFixed(Base16, 16),      // saturated fixed (2**64)
		MustNewFixed("0123456789", 9), // bounded fixed (1e9 keys)
		MustNewDynamic("αβγδε"),       // multi-byte runes
	}
}

// FuzzRoundTrip asserts the core contract for every codec: whenever Marshal
// accepts an id, Unmarshal must return that exact id, the key must be Valid,
// and re-encoding must reproduce the same key (canonical stability). This is
// the property that the old float64/int arithmetic silently violated.
func FuzzRoundTrip(f *testing.F) {
	for _, id := range []uint64{0, 1, 2, 255, 1000000, 1 << 53, 1 << 63, math.MaxUint64} {
		f.Add(id)
	}

	codecs := fuzzCodecs()
	f.Fuzz(func(t *testing.T, id uint64) {
		for _, ls := range codecs {
			k, err := ls.Marshal(id)
			if err != nil {
				// Only a bounded space may reject, and only for id >= Total.
				if !ls.Saturated() && id < ls.Total() {
					t.Fatalf("Marshal(%d) rejected within bounds (Total %d): %v",
						id, ls.Total(), err)
				}
				continue
			}

			back, err := ls.Unmarshal(k)
			if err != nil {
				t.Fatalf("Unmarshal(%q) after Marshal(%d): %v", k, id, err)
			}
			if back != id {
				t.Fatalf("round trip: %d -> %q -> %d", id, k, back)
			}
			if !ls.Valid(k) {
				t.Fatalf("Valid(%q) = false for a freshly marshaled key", k)
			}
			if re, _ := ls.Marshal(back); re != k {
				t.Fatalf("non-canonical: %q -> %d -> %q", k, back, re)
			}
		}
	})
}

// FuzzUnmarshal feeds arbitrary strings to Unmarshal. The invariants: it must
// never panic, and every accepted key must be canonical, i.e. round-trip back
// to itself through Marshal. This pins the strict-bijection guarantee against
// any input, not just keys we generated.
func FuzzUnmarshal(f *testing.F) {
	for _, s := range []string{"", "a", "b", "ab", "aab", "ffff", "zzzzzzzzzzzzzzzzz", "αβ"} {
		f.Add(s)
	}

	codecs := fuzzCodecs()
	f.Fuzz(func(t *testing.T, s string) {
		for _, ls := range codecs {
			id, err := ls.Unmarshal(s)
			if err != nil {
				continue // rejected inputs are fine
			}
			// Accepted => must be the canonical encoding of its id.
			re, err := ls.Marshal(id)
			if err != nil {
				t.Fatalf("Marshal(%d) failed for accepted key %q: %v", id, s, err)
			}
			if re != s {
				t.Fatalf("accepted non-canonical key %q (id %d re-encodes to %q)",
					s, id, re)
			}
		}
	})
}
