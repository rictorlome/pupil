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
	for _, s := range SQUARES {
		fmt.Println(KING_ATTACK_BBS[s])
	}
}
