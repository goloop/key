package key

import (
	"errors"
	"math"
	"sync"
	"testing"
)

const base36 = "abcdefghijklmnopqrstuvwxyz0123456789"

// TestConstructorValidation covers every rejection path of both constructors
// and asserts the sentinel error (errors.Is) plus the BUG-07 contract: on
// error the returned *Locksmith must be nil, never a half-built object.
func TestConstructorValidation(t *testing.T) {
	tests := []struct {
		name string
		make func() (*Locksmith, error)
		want error
	}{
		{"dynamic blank", func() (*Locksmith, error) { return NewDynamic("") }, ErrBlankAlphabet},
		{"dynamic single char", func() (*Locksmith, error) { return NewDynamic("a") }, ErrShortAlphabet},
		{"dynamic duplicate", func() (*Locksmith, error) { return NewDynamic("aab") }, ErrDuplicateChar},
		{"dynamic duplicate rune", func() (*Locksmith, error) { return NewDynamic("abca") }, ErrDuplicateChar},
		{"fixed blank", func() (*Locksmith, error) { return NewFixed("", 3) }, ErrBlankAlphabet},
		{"fixed single char", func() (*Locksmith, error) { return NewFixed("a", 5) }, ErrShortAlphabet},
		{"fixed duplicate", func() (*Locksmith, error) { return NewFixed("abcc", 3) }, ErrDuplicateChar},
		{"fixed zero size", func() (*Locksmith, error) { return NewFixed("abc", 0) }, ErrInvalidSize},
		{"fixed negative size", func() (*Locksmith, error) { return NewFixed("abc", -1) }, ErrInvalidSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls, err := tt.make()
			if !errors.Is(err, tt.want) {
				t.Errorf("error = %v, want errors.Is(_, %v)", err, tt.want)
			}
			if ls != nil {
				t.Errorf("Locksmith must be nil on error, got %#v", ls)
			}
		})
	}
}

// TestConstructorValidInputs verifies that valid inputs build a usable codec
// and that unicode (multi-byte) alphabets are honoured by rune, not byte.
func TestConstructorValidInputs(t *testing.T) {
	if ls, err := NewDynamic("αβγ"); err != nil || ls == nil {
		t.Fatalf("NewDynamic(unicode) = %v, %v", ls, err)
	}
	if ls, err := NewFixed(base36, 4); err != nil || ls == nil {
		t.Fatalf("NewFixed = %v, %v", ls, err)
	}
}

// TestAccessors locks the metadata reported by the small-space fixed codec and
// the saturated dynamic one.
func TestAccessors(t *testing.T) {
	fixed, _ := NewFixed("abc", 3)
	if got := fixed.Alphabet(); got != "abc" {
		t.Errorf("Alphabet() = %q, want abc", got)
	}
	if got := fixed.Size(); got != 3 {
		t.Errorf("Size() = %d, want 3", got)
	}
	if got := fixed.Total(); got != 27 {
		t.Errorf("Total() = %d, want 27", got)
	}
	if fixed.Saturated() {
		t.Error("Saturated() = true for a 27-key space, want false")
	}

	dyn, _ := NewDynamic("abc")
	if got := dyn.Size(); got != 0 {
		t.Errorf("dynamic Size() = %d, want 0", got)
	}
	if got := dyn.Total(); got != math.MaxUint64 {
		t.Errorf("dynamic Total() = %d, want MaxUint64", got)
	}
	if !dyn.Saturated() {
		t.Error("dynamic Saturated() = false, want true")
	}
}

// TestKnownVectors pins the encoding so the bijection's actual byte output
// cannot drift between versions. Values are unchanged base-N encodings.
func TestKnownVectors(t *testing.T) {
	tests := []struct {
		size int
		id   uint64
		key  string
	}{
		{0, 3333333, "b9qav"}, // dynamic
		{0, 10, "k"},
		{3, 0, "aaa"},
		{3, 1, "aab"},
		{3, 10, "aak"},
		{5, 1024, "aaa2q"},
		{7, 1024, "aaaaa2q"},
	}

	for _, tt := range tests {
		var ls *Locksmith
		var err error
		if tt.size == 0 {
			ls, err = NewDynamic(base36)
		} else {
			ls, err = NewFixed(base36, tt.size)
		}
		if err != nil {
			t.Fatal(err)
		}

		got, err := ls.Marshal(tt.id)
		if err != nil {
			t.Fatalf("Marshal(%d): %v", tt.id, err)
		}
		if got != tt.key {
			t.Errorf("Marshal(%d) = %q, want %q", tt.id, got, tt.key)
		}

		back, err := ls.Unmarshal(tt.key)
		if err != nil {
			t.Fatalf("Unmarshal(%q): %v", tt.key, err)
		}
		if back != tt.id {
			t.Errorf("Unmarshal(%q) = %d, want %d", tt.key, back, tt.id)
		}
	}
}

// TestRoundTripExhaustiveFixed walks the entire key space of a small fixed
// codec and asserts three invariants at once: round-trip identity, fixed
// length, and global injectivity (no two ids share a key).
func TestRoundTripExhaustiveFixed(t *testing.T) {
	ls, _ := NewFixed("abc", 4) // 3**4 == 81 keys
	if ls.Total() != 81 {
		t.Fatalf("Total() = %d, want 81", ls.Total())
	}

	seen := make(map[string]uint64, 81)
	for id := uint64(0); id < ls.Total(); id++ {
		k, err := ls.Marshal(id)
		if err != nil {
			t.Fatalf("Marshal(%d): %v", id, err)
		}
		if len([]rune(k)) != 4 {
			t.Fatalf("Marshal(%d) = %q, length %d, want 4", id, k, len([]rune(k)))
		}
		if prev, dup := seen[k]; dup {
			t.Fatalf("collision: ids %d and %d both map to %q", prev, id, k)
		}
		seen[k] = id

		back, err := ls.Unmarshal(k)
		if err != nil || back != id {
			t.Fatalf("round trip %d -> %q -> %d (err %v)", id, k, back, err)
		}
	}
	if len(seen) != 81 {
		t.Fatalf("expected 81 distinct keys, got %d", len(seen))
	}
}

// TestRoundTripExhaustiveDynamic checks the dynamic codec over a wide prefix of
// ids for both round-trip identity and canonical stability
// (Marshal(Unmarshal(k)) == k).
func TestRoundTripExhaustiveDynamic(t *testing.T) {
	ls, _ := NewDynamic("abcde")
	for id := uint64(0); id < 200000; id++ {
		k, err := ls.Marshal(id)
		if err != nil {
			t.Fatalf("Marshal(%d): %v", id, err)
		}
		back, err := ls.Unmarshal(k)
		if err != nil || back != id {
			t.Fatalf("round trip %d -> %q -> %d (err %v)", id, k, back, err)
		}
		// Canonical: re-encoding the decoded id must reproduce the key.
		if re, _ := ls.Marshal(back); re != k {
			t.Fatalf("non-canonical: %q -> %d -> %q", k, back, re)
		}
	}
}

// TestLargeIDRoundTrip is the BUG-01/BUG-05 regression: with base 2 the codec
// must stay exact past 2**53 (float64 mantissa) and past 2**63 (signed int),
// all the way to MaxUint64, with no collisions among the boundary values.
func TestLargeIDRoundTrip(t *testing.T) {
	ls, _ := NewDynamic("ab") // base 2: maximally stresses the arithmetic

	ids := []uint64{
		(1 << 53) - 1, 1 << 53, (1 << 53) + 1, // float64 boundary
		(1 << 63) - 1, 1 << 63, (1 << 63) + 1, // signed-int boundary
		math.MaxInt64 + 1000,
		math.MaxUint64 - 1, math.MaxUint64,
	}

	seen := make(map[string]uint64)
	for _, id := range ids {
		k, err := ls.Marshal(id)
		if err != nil {
			t.Fatalf("Marshal(%d): %v", id, err)
		}
		if prev, dup := seen[k]; dup {
			t.Fatalf("collision: %d and %d both -> %q", prev, id, k)
		}
		seen[k] = id

		back, err := ls.Unmarshal(k)
		if err != nil {
			t.Fatalf("Unmarshal(%q): %v", k, err)
		}
		if back != id {
			t.Fatalf("round trip broke at %d: -> %q -> %d", id, k, back)
		}
	}
}

// TestNoCollisionAcrossFloatBoundary is the focused BUG-01 reproduction: the
// old float64 path mapped 2**53 and 2**53+1 to the same key.
func TestNoCollisionAcrossFloatBoundary(t *testing.T) {
	ls, _ := NewDynamic("ab")
	a, _ := ls.Marshal(1 << 53)
	b, _ := ls.Marshal((1 << 53) + 1)
	if a == b {
		t.Fatalf("2**53 and 2**53+1 collide: both %q", a)
	}
}

// TestTotalSaturation is the BUG-02/BUG-03 regression. Spaces of exactly or
// above 2**64 must report MaxUint64 and stay fully usable, instead of the old
// Total()==0 (which bricked the instance) or a plausible-but-wrong value.
func TestTotalSaturation(t *testing.T) {
	tests := []struct {
		name     string
		alphabet string
		size     int
	}{
		{"16^16 == 2^64", Base16, 16},
		{"2^64 via base2", "01", 64},
		{"base20 ^16 (far above)", "0123456789abcdefghij", 16},
		{"base36 ^16", base36, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls, err := NewFixed(tt.alphabet, tt.size)
			if err != nil {
				t.Fatal(err)
			}
			if !ls.Saturated() {
				t.Fatalf("Saturated() = false, want true")
			}
			if ls.Total() != math.MaxUint64 {
				t.Fatalf("Total() = %d, want MaxUint64", ls.Total())
			}
			// Every uint64 id must encode, including MaxUint64 itself.
			for _, id := range []uint64{0, 1, math.MaxInt64, math.MaxUint64} {
				k, err := ls.Marshal(id)
				if err != nil {
					t.Fatalf("Marshal(%d): %v", id, err)
				}
				if l := uint64(len([]rune(k))); l != ls.Size() {
					t.Fatalf("Marshal(%d) length %d, want %d", id, l, ls.Size())
				}
				if back, err := ls.Unmarshal(k); err != nil || back != id {
					t.Fatalf("round trip %d -> %q -> %d (err %v)", id, k, back, err)
				}
			}
		})
	}
}

// TestTotalExactBelowBoundary checks that a space just under 2**64 keeps an
// enforced, exact bound and rejects the first out-of-range id.
func TestTotalExactBelowBoundary(t *testing.T) {
	ls, _ := NewFixed("01", 63) // 2**63 keys exactly
	if ls.Saturated() {
		t.Fatal("Saturated() = true, want false for a 2**63 space")
	}
	if ls.Total() != 1<<63 {
		t.Fatalf("Total() = %d, want 2**63", ls.Total())
	}
	if _, err := ls.Marshal((1 << 63) - 1); err != nil {
		t.Fatalf("Marshal(2**63-1) must succeed: %v", err)
	}
	if _, err := ls.Marshal(1 << 63); !errors.Is(err, ErrIDTooLarge) {
		t.Fatalf("Marshal(2**63) error = %v, want ErrIDTooLarge", err)
	}
}

// TestMarshalIDTooLarge checks the ordinary out-of-range rejection.
func TestMarshalIDTooLarge(t *testing.T) {
	ls, _ := NewFixed("abc", 3) // 27 keys: valid ids 0..26
	if _, err := ls.Marshal(26); err != nil {
		t.Fatalf("Marshal(26) must succeed: %v", err)
	}
	if _, err := ls.Marshal(27); !errors.Is(err, ErrIDTooLarge) {
		t.Fatalf("Marshal(27) error = %v, want ErrIDTooLarge", err)
	}
}

// TestUnmarshalStrict is the BUG-06 regression: decoding is a strict bijection.
func TestUnmarshalStrict(t *testing.T) {
	dyn, _ := NewDynamic("abc")

	// Empty string is never a valid key.
	if _, err := dyn.Unmarshal(""); !errors.Is(err, ErrEmptyKey) {
		t.Errorf("Unmarshal(\"\") error = %v, want ErrEmptyKey", err)
	}

	// Canonical forms decode.
	for k, want := range map[string]uint64{"a": 0, "b": 1, "c": 2, "ba": 3} {
		if got, err := dyn.Unmarshal(k); err != nil || got != want {
			t.Errorf("Unmarshal(%q) = %d, %v; want %d, nil", k, got, err, want)
		}
	}

	// Non-canonical: redundant leading lead-character must be rejected, even
	// though the old code silently normalized "aab" and "ab" down to "b".
	for _, k := range []string{"aab", "ab", "aac", "aba"} {
		if _, err := dyn.Unmarshal(k); !errors.Is(err, ErrNonCanonical) {
			t.Errorf("Unmarshal(%q) error = %v, want ErrNonCanonical", k, err)
		}
	}

	// Unknown character (leading char is non-lead, so the only fault is 'X').
	if _, err := dyn.Unmarshal("bXc"); !errors.Is(err, ErrUnknownChar) {
		t.Errorf("Unmarshal(\"bXc\") error = %v, want ErrUnknownChar", err)
	}

	// Fixed mode: leading lead-character is legitimate padding, not redundant.
	fixed, _ := NewFixed("abc", 3)
	if got, err := fixed.Unmarshal("aab"); err != nil || got != 1 {
		t.Errorf("fixed Unmarshal(\"aab\") = %d, %v; want 1, nil", got, err)
	}
	// Wrong length is rejected in fixed mode.
	for _, k := range []string{"ab", "aabb", ""} {
		if _, err := fixed.Unmarshal(k); err == nil {
			t.Errorf("fixed Unmarshal(%q) = nil error, want non-nil", k)
		}
	}
	if _, err := fixed.Unmarshal("ab"); !errors.Is(err, ErrInvalidLength) {
		t.Errorf("fixed Unmarshal(\"ab\") error = %v, want ErrInvalidLength", err)
	}
}

// TestUnmarshalOverflow is the BUG-04 regression: an over-long key must be
// rejected with ErrKeyOutOfRange instead of returning silent garbage.
func TestUnmarshalOverflow(t *testing.T) {
	ls, _ := NewDynamic(Base16) // base 16

	// The longest canonical key for MaxUint64 round-trips exactly.
	maxKey, _ := ls.Marshal(math.MaxUint64)
	if back, err := ls.Unmarshal(maxKey); err != nil || back != math.MaxUint64 {
		t.Fatalf("Unmarshal(MaxUint64 key) = %d, %v", back, err)
	}

	// "f" repeated 17 times is 17 hex digits == far beyond uint64. This trips
	// the multiply check (id*base overflows).
	over := "fffffffffffffffff" // 17 chars
	if _, err := ls.Unmarshal(over); !errors.Is(err, ErrKeyOutOfRange) {
		t.Fatalf("Unmarshal(%q) error = %v, want ErrKeyOutOfRange", over, err)
	}

	// The add-carry overflow path: with base 10, MaxUint64+4 == 18446744073709551619
	// has id*10 fitting (no multiply overflow) but id*10+digit carrying past
	// uint64. A power-of-two base can never reach this branch.
	dec, _ := NewDynamic("0123456789")
	if _, err := dec.Unmarshal("18446744073709551619"); !errors.Is(err, ErrKeyOutOfRange) {
		t.Fatalf("Unmarshal(MaxUint64+4) error = %v, want ErrKeyOutOfRange", err)
	}
	// And the exact maximum still decodes.
	if back, err := dec.Unmarshal("18446744073709551615"); err != nil || back != math.MaxUint64 {
		t.Fatalf("Unmarshal(MaxUint64) = %d, %v; want MaxUint64, nil", back, err)
	}
}

// TestSingleCharRejected is the BUG-09 regression. A base-1 alphabet used to be
// accepted and then hang Marshal in an infinite loop; now it cannot even be
// constructed, so the dangerous input is unreachable.
func TestSingleCharRejected(t *testing.T) {
	if _, err := NewDynamic("a"); !errors.Is(err, ErrShortAlphabet) {
		t.Errorf("NewDynamic(\"a\") error = %v, want ErrShortAlphabet", err)
	}
	if _, err := NewFixed("a", 5); !errors.Is(err, ErrShortAlphabet) {
		t.Errorf("NewFixed(\"a\", 5) error = %v, want ErrShortAlphabet", err)
	}
}

// TestMarshalAppend asserts MarshalAppend is byte-for-byte consistent with
// Marshal, preserves any existing prefix, and exercises the large-size
// (heap fallback) path.
func TestMarshalAppend(t *testing.T) {
	cases := []struct {
		ls *Locksmith
		id uint64
	}{
		{mustFixed("abc", 4), 10},
		{mustDyn(Base62), 1234567},
		{mustFixed(Base16, 16), math.MaxUint64}, // saturated, stack path
		{mustFixed("01", 70), 0xDEADBEEF},       // size > stackDigits, heap path
	}

	for _, c := range cases {
		want, err := c.ls.Marshal(c.id)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}

		// Fresh buffer.
		got, err := c.ls.MarshalAppend(nil, c.id)
		if err != nil {
			t.Fatalf("MarshalAppend: %v", err)
		}
		if string(got) != want {
			t.Errorf("MarshalAppend(nil) = %q, want %q", got, want)
		}

		// Existing prefix must be preserved.
		got, _ = c.ls.MarshalAppend([]byte("id="), c.id)
		if string(got) != "id="+want {
			t.Errorf("MarshalAppend(prefix) = %q, want %q", got, "id="+want)
		}
	}

	// Out-of-range id is reported and leaves dst untouched.
	ls := mustFixed("abc", 3)
	dst := []byte("x")
	got, err := ls.MarshalAppend(dst, 27)
	if !errors.Is(err, ErrIDTooLarge) {
		t.Errorf("MarshalAppend out-of-range error = %v, want ErrIDTooLarge", err)
	}
	if string(got) != "x" {
		t.Errorf("MarshalAppend kept %q, want unchanged %q", got, "x")
	}
}

// TestValid mirrors Unmarshal's accept/reject decisions.
func TestValid(t *testing.T) {
	dyn, _ := NewDynamic("abc")
	checks := map[string]bool{
		"a":   true,
		"bab": true,
		"":    false,
		"aab": false, // non-canonical
		"abX": false, // unknown char
	}
	for k, want := range checks {
		if got := dyn.Valid(k); got != want {
			t.Errorf("Valid(%q) = %v, want %v", k, got, want)
		}
	}
}

// TestMust verifies the Must* helpers panic on bad input and not on good input.
func TestMust(t *testing.T) {
	if got := MustNewDynamic("abc").Alphabet(); got != "abc" {
		t.Errorf("MustNewDynamic alphabet = %q", got)
	}
	if MustNewFixed("abc", 3).Size() != 3 {
		t.Error("MustNewFixed size mismatch")
	}
	assertPanics(t, "MustNewDynamic(\"a\")", func() { MustNewDynamic("a") })
	assertPanics(t, "MustNewFixed(\"abc\", 0)", func() { MustNewFixed("abc", 0) })
}

// TestConcurrentReads stresses a shared Locksmith from many goroutines. It is
// meaningful under `go test -race`: it proves the codec is safe for concurrent
// use, backing the documentation's immutability claim.
func TestConcurrentReads(t *testing.T) {
	ls, _ := NewFixed(base36, 8)

	var wg sync.WaitGroup
	for g := 0; g < 16; g++ {
		wg.Add(1)
		go func(off uint64) {
			defer wg.Done()
			for i := uint64(0); i < 5000; i++ {
				id := off*5000 + i
				k, err := ls.Marshal(id)
				if err != nil {
					t.Errorf("Marshal(%d): %v", id, err)
					return
				}
				if back, err := ls.Unmarshal(k); err != nil || back != id {
					t.Errorf("round trip %d -> %q -> %d (%v)", id, k, back, err)
					return
				}
				_ = ls.Valid(k)
			}
		}(uint64(g))
	}
	wg.Wait()
}

// TestUnmarshalLongKeyNoAlloc guards the BUG-01 fix: decoding even a maximal
// 64-character key must not allocate (no []rune materialization).
func TestUnmarshalLongKeyNoAlloc(t *testing.T) {
	ls, _ := NewDynamic("ab") // base 2 -> 64-char key for MaxUint64
	k, _ := ls.Marshal(math.MaxUint64)
	if len(k) != 64 {
		t.Fatalf("setup: key length %d, want 64", len(k))
	}

	allocs := testing.AllocsPerRun(200, func() {
		if _, err := ls.Unmarshal(k); err != nil {
			t.Fatal(err)
		}
	})
	if allocs != 0 {
		t.Errorf("Unmarshal allocated %.1f times on a 64-char key, want 0", allocs)
	}
}

// --- helpers ---

func mustDyn(alphabet string) *Locksmith   { return MustNewDynamic(alphabet) }
func mustFixed(a string, n int) *Locksmith { return MustNewFixed(a, n) }

func assertPanics(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s did not panic", name)
		}
	}()
	fn()
}
