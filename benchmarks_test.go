package key

import (
	"testing"
)

func BenchmarkNew(b *testing.B) {
	benchCases := []struct {
		name     string
		alphabet string
		size     int
	}{
		{"SmallAlphabet", "abc", 3},
		{"MediumAlphabet", "abcdefghijk", 5},
		{"LargeAlphabet", "abcdefghijklmnopqrstuvwxyz0123456789", 8},
		{"DynamicSize", "abcdefghijk", 0},
	}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = New(bc.alphabet, bc.size)
			}
		})
	}
}

func BenchmarkMarshal(b *testing.B) {
	benchCases := []struct {
		name     string
		alphabet string
		size     int
		id       uint64
	}{
		{"Small_ID_FixedSize", "abc", 3, 10},
		{"Medium_ID_FixedSize", "abcdefghijk", 5, 1000},
		{"Large_ID_FixedSize", "abcdefghijklmnopqrstuvwxyz0123456789", 8, 1000000},
		{"Small_ID_DynamicSize", "abc", 0, 10},
		{"Large_ID_DynamicSize", "abcdefghijklmnopqrstuvwxyz0123456789", 0, 1000000},
	}

	for _, bc := range benchCases {
		ls, _ := New(bc.alphabet, bc.size)
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = ls.Marshal(bc.id)
			}
		})
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	benchCases := []struct {
		name     string
		alphabet string
		size     int
		key      string
	}{
		{"Short_Key_FixedSize", "abc", 3, "bab"},
		{"Medium_Key_FixedSize", "abcdefghijk", 5, "defgh"},
		{"Long_Key_FixedSize", "abcdefghijklmnopqrstuvwxyz0123456789", 8, "12345678"},
		{"Short_Key_DynamicSize", "abc", 0, "bab"},
		{"Long_Key_DynamicSize", "abcdefghijklmnopqrstuvwxyz0123456789", 0, "12345678"},
	}

	for _, bc := range benchCases {
		ls, _ := New(bc.alphabet, bc.size)
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = ls.Unmarshal(bc.key)
			}
		})
	}
}

func BenchmarkParallelOperations(b *testing.B) {
	ls, _ := New("abcdefghijklmnopqrstuvwxyz0123456789", 8)

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
		key, _ := ls.Marshal(123456)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = ls.Unmarshal(key)
			}
		})
	})
}
