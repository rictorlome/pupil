package main

import "testing"

func benchmarkAlphaBeta(fen string, depth uint8, b *testing.B) {
	pos := parseFen(fen)
	for n := 0; n < b.N; n++ {
		pos.abRoot(depth)
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
func BenchmarkAlphaBetaKiwipete6(b *testing.B) {
	benchmarkAlphaBeta(KIWIPETE_FEN, 6, b)
}
func BenchmarkAlphaBetaKiwipete7(b *testing.B) {
	benchmarkAlphaBeta(KIWIPETE_FEN, 7, b)
}
func BenchmarkAlphaBetaKiwipete8(b *testing.B) {
	benchmarkAlphaBeta(KIWIPETE_FEN, 8, b)
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
func BenchmarkAlphaBetaInitial6(b *testing.B) {
	benchmarkAlphaBeta(INITIAL_FEN, 6, b)
}
func BenchmarkAlphaBetaInitial7(b *testing.B) {
	benchmarkAlphaBeta(INITIAL_FEN, 7, b)
}
func BenchmarkAlphaBetaInitial8(b *testing.B) {
	benchmarkAlphaBeta(INITIAL_FEN, 8, b)
}
