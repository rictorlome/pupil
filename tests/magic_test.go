// Benchmarks:
// run with $ go test -bench=.
// (where `.` is a regex)
package main

import "testing"


func BenchmarkInitMagic(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		e := init_magics()
		if e != nil {
			println(e)
		}
	}
}

func BenchmarkBishopMagic(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, sq := range SQUARES {
			BishopMagics[sq], _ = find_magic(sq, true)
		}
	}
}

func BenchmarkRookMagic(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, sq := range SQUARES {
			RookMagics[sq], _ = find_magic(sq, false)
		}
	}
}
