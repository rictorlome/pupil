package main

import "fmt"

func init() {
	// Initialize Square Bitboards
	for _, s := range SQUARES {
		var SQUARE_BB Bitboard = 0x1 << s
		SQUARE_BBS[s] = SQUARE_BB
		NEIGHBOR_BBS[s] = neighbors(SQUARE_BB)
		KING_ATTACK_BBS[s] = precompute_king_attacks(SQUARE_BB)
		KNIGHT_ATTACK_BBS[s] = precompute_knight_attacks(SQUARE_BB)
		ROOK_ATTACK_MASKS[s] = rook_attacks(Bitboard(0), s)
		BISHOP_ATTACK_MASKS[s] = bishop_attacks(Bitboard(0), s)
		ROOK_OCCUPANCY_MASKS[s] = occupancy_mask(s, ROOK_DIRECTIONS)
		BISHOP_OCCUPANCY_MASKS[s] = occupancy_mask(s, BISHOP_DIRECTIONS)
		RELEVANT_ROOK_OCCUPANCY[s] = ROOK_ATTACK_MASKS[s] &^ ROOK_OCCUPANCY_MASKS[s]
		RELEVANT_BISHOP_OCCUPANCY[s] = BISHOP_ATTACK_MASKS[s] &^ BISHOP_OCCUPANCY_MASKS[s]
		for _, color := range COLORS {
			PAWN_ATTACK_BBS[s][color] = pawn_attacks(SQUARE_BB, color)
		}
	}
	// Initialize dependent BBs
	for _, atk_mask := range [2][64]Bitboard{ROOK_ATTACK_MASKS, BISHOP_ATTACK_MASKS} {
		for _, s1 := range SQUARES {
			for _, s2 := range SQUARES {
				if occupied_at_sq(atk_mask[s1], s2) {
					LINE_BBS[s1][s2] = atk_mask[s1]&atk_mask[s2] | SQUARE_BBS[s1] | SQUARE_BBS[s2]
				}
			}
		}
	}

	init_castle_sqs()
	init_castling_masks()
}

func main() {
	fmt.Println("OK")
	pos := parse_fen(INITIAL_FEN)
	fmt.Println(pos.king_square(BLACK))
}
