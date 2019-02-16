package main

import "testing"

func benchmarkAlphaBeta(fen string, depth uint8, b *testing.B) {
	pos := parse_fen(fen)
	for n := 0; n < b.N; n++ {
		pos.ab_root(depth)
	}
}
func BenchmarkAlphaBetaKiwipete2(b *testing.B) {
	benchmarkAlphaBeta(KIWIPETE_FEN, 2, b)
}
func BenchmarkAlphaBetaKiwipete3(b *testing.B) {
	benchmarkAlphaBeta(KIWIPETE_FEN, 3, b)
}
func BenchmarkAlphaBetaKiwipete4(b *testing.B) {
	benchmarkAlphaBeta(KIWIPETE_FEN, 4, b)
}
func BenchmarkAlphaBetaKiwipete5(b *testing.B) {
	benchmarkAlphaBeta(KIWIPETE_FEN, 5, b)
}
func BenchmarkAlphaBetaInitial2(b *testing.B) {
	benchmarkAlphaBeta(INITIAL_FEN, 2, b)
}
func BenchmarkAlphaBetaInitial3(b *testing.B) {
	benchmarkAlphaBeta(INITIAL_FEN, 3, b)
}
func BenchmarkAlphaBetaInitial4(b *testing.B) {
	benchmarkAlphaBeta(INITIAL_FEN, 4, b)
}
func BenchmarkAlphaBetaInitial5(b *testing.B) {
	benchmarkAlphaBeta(INITIAL_FEN, 5, b)
}
