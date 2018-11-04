package main

import "fmt"

func init() {
	// Initialize Square Bitboards
	var x Bitboard = 0x1
	for _, s := range SQUARES {
		var SQUARE_BB Bitboard = x << s
		SQUARE_BBS[s] = SQUARE_BB
		KING_ATTACK_BBS[s] = king_attacks(SQUARE_BB)
		KNIGHT_ATTACK_BBS[s] = knight_attacks(SQUARE_BB)
	}
}

func main() {
	fmt.Println("OK")
	x := parse_positions(INITIAL_FEN_JUST_PIECES)
	fmt.Println(queen_attacks(occupied_squares(x), SQ_E4))
}
