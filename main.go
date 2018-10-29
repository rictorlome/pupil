package main

import "fmt"

func init() {
	// Initialize Square Bitboards
	var x Bitboard = 0x1
	for _, s := range SQUARES {
		SQUARE_BBS[s] = x << s
	}
}

func main() {
	fmt.Println("OK")
	fmt.Println(generate_fen(parse_fen("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")))
	// fmt.Println(BQ_CASTLE)
}
