package main

import "fmt"

func init() {
	// Initialize Square Bitboards
	var x Bitboard = 0x1
	for _, s := range SQUARES {
		var SQUARE_BB Bitboard = x << s
		SQUARE_BBS[s] = SQUARE_BB
		KING_ATTACK_BBS[s] = precompute_king_attacks(SQUARE_BB)
		KNIGHT_ATTACK_BBS[s] = precompute_knight_attacks(SQUARE_BB)
		ROOK_ATTACK_MASKS[s] = rook_attacks(Bitboard(0), s)
		BISHOP_ATTACK_MASKS[s] = bishop_attacks(Bitboard(0), s)
		ROOK_OCCUPANCY_MASKS[s] = occupancy_mask(s, ROOK_DIRECTIONS)
		BISHOP_OCCUPANCY_MASKS[s] = occupancy_mask(s, BISHOP_DIRECTIONS)
		RELEVANT_ROOK_OCCUPANCY[s] = ROOK_ATTACK_MASKS[s] &^ ROOK_OCCUPANCY_MASKS[s]
		RELEVANT_BISHOP_OCCUPANCY[s] = BISHOP_ATTACK_MASKS[s] &^ BISHOP_OCCUPANCY_MASKS[s]
	}
}

func main() {
	fmt.Println("OK")
	pos := parse_fen(INITIAL_FEN)
	fmt.Println(pseudolegals_by_color(pos.placement, BLACK, pos.state))
}
