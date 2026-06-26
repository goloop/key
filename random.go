package key

import (
	"encoding/binary"
	"io"
	"math"
)

// Random returns a key for a uniformly random id drawn from the key space,
// reading its randomness from r. Pass crypto/rand.Reader for cryptographically
// secure keys, or any deterministic io.Reader in tests.
//
// The id is sampled without modulo bias: for a bounded space the function uses
// rejection sampling so every id in [0, Total) is equally likely; for a
// saturated space every uint64 is equally likely. Any read error from r is
// returned unchanged.
func (ls *Locksmith) Random(r io.Reader) (string, error) {
	var (
		id  uint64
		err error
	)
	if ls.full {
		id, err = readUint64(r)
	} else {
		id, err = uniformUint64(r, ls.total)
	}
	if err != nil {
		return "", err
	}

	return ls.Marshal(id)
}

// readUint64 reads eight bytes from r and assembles them into a uint64.
func readUint64(r io.Reader) (uint64, error) {
	var b [8]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b[:]), nil
}

// uniformUint64 returns a uniformly distributed value in [0, n) using rejection
// sampling, so the result carries no modulo bias. n must be at least 1.
func uniformUint64(r io.Reader, n uint64) (uint64, error) {
	if n == 1 {
		return 0, nil // only one possible outcome
	}

	// limit is the largest multiple of n that fits, so discarding the tail
	// [limit, MaxUint64] leaves each residue equally represented below it.
	limit := math.MaxUint64 - (math.MaxUint64 % n)
	for {
		v, err := readUint64(r)
		if err != nil {
			return 0, err
		}
		if v < limit {
			return v % n, nil
		}
	}
}
