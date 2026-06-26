package key

import "testing"

func BenchmarkNew(b *testing.B) {
	cases := []struct {
		name     string
		alphabet string
		size     int
	}{
		{"Fixed_Small", "abc", 3},
		{"Fixed_Medium", "abcdefghijk", 5},
		{"Fixed_Large", base36, 8},
		{"Dynamic", "abcdefghijk", 0},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if c.size == 0 {
					_, _ = NewDynamic(c.alphabet)
				} else {
					_, _ = NewFixed(c.alphabet, c.size)
				}
			}
		})
	}
}

func BenchmarkMarshal(b *testing.B) {
	cases := []struct {
		name     string
		alphabet string
		size     int
		id       uint64
	}{
		{"Small_ID_Fixed", "abc", 3, 10},
		{"Medium_ID_Fixed", "abcdefghijk", 5, 1000},
		{"Large_ID_Fixed", base36, 8, 1000000},
		{"Small_ID_Dynamic", "abc", 0, 10},
		{"Large_ID_Dynamic", base36, 0, 1000000},
	}

	for _, c := range cases {
		ls := build(c.alphabet, c.size)
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = ls.Marshal(c.id)
			}
		})
	}
}

func BenchmarkMarshalAppend(b *testing.B) {
	ls := build(base36, 8)
	buf := make([]byte, 0, 16)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf, _ = ls.MarshalAppend(buf[:0], 1000000)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	cases := []struct {
		name     string
		alphabet string
		size     int
		key      string
	}{
		{"Short_Fixed", "abc", 3, "bab"},
		{"Medium_Fixed", "abcdefghijk", 5, "defgh"},
		{"Long_Fixed", base36, 8, "12345678"},
		{"Short_Dynamic", "abc", 0, "bab"},
		{"Long_Dynamic", base36, 0, "12345678"},
	}

	for _, c := range cases {
		ls := build(c.alphabet, c.size)
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = ls.Unmarshal(c.key)
			}
		})
	}
}

func BenchmarkParallelOperations(b *testing.B) {
	ls := build(base36, 8)

	b.Run("Parallel_Marshal", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			var id uint64
			for pb.Next() {
				id = (id + 1) % 1000000
				_, _ = ls.Marshal(id)
			}
		})
	})

	b.Run("Parallel_Unmarshal", func(b *testing.B) {
		k, _ := ls.Marshal(123456)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = ls.Unmarshal(k)
			}
		})
	})
}

// build is a test helper: size 0 selects a dynamic codec, otherwise fixed.
func build(alphabet string, size int) *Locksmith {
	if size == 0 {
		return MustNewDynamic(alphabet)
	}
	return MustNewFixed(alphabet, size)
}
