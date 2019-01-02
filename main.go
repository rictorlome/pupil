package main

import "fmt"

func init() {
	// Initialize Square Bitboards
	e := Bitboard(0)
	for _, s := range SQUARES {
		var SQUARE_BB Bitboard = 0x1 << s
		SQUARE_BBS[s] = SQUARE_BB
		NEIGHBOR_BBS[s] = neighbors(SQUARE_BB)
		KING_ATTACK_BBS[s] = precompute_king_attacks(SQUARE_BB)
		KNIGHT_ATTACK_BBS[s] = precompute_knight_attacks(SQUARE_BB)
		ROOK_ATTACK_MASKS[s] = slider_rook_attacks(Bitboard(0), s)
		BISHOP_ATTACK_MASKS[s] = slider_bishop_attacks(Bitboard(0), s)
		ROOK_OCCUPANCY_MASKS[s] = occupancy_mask(s, ROOK_DIRECTIONS)
		BISHOP_OCCUPANCY_MASKS[s] = occupancy_mask(s, BISHOP_DIRECTIONS)
		RELEVANT_ROOK_OCCUPANCY[s] = ROOK_ATTACK_MASKS[s] &^ ROOK_OCCUPANCY_MASKS[s]
		RELEVANT_BISHOP_OCCUPANCY[s] = BISHOP_ATTACK_MASKS[s] &^ BISHOP_OCCUPANCY_MASKS[s]
		for _, color := range COLORS {
			PAWN_ATTACK_BBS[s][color] = pawn_attacks(SQUARE_BB, color)
		}
	}
	// Initialize dependent BBs
	for _, fn := range [2]AttackFunc{slider_rook_attacks, slider_bishop_attacks} {
		for _, s1 := range SQUARES {
			for _, s2 := range SQUARES {
				if occupied_at_sq(fn(e, s1), s2) {
					LINE_BBS[s1][s2] = fn(e, s1)&fn(e, s2) | SQUARE_BBS[s1] | SQUARE_BBS[s2]
					BETWEEN_BBS[s1][s2] = fn(SQUARE_BBS[s2], s1) & fn(SQUARE_BBS[s1], s2)
				}
			}
		}
	}

	init_castle_sqs()
	init_castling_masks()
	init_rook_squares_for_castling()
	for _, sq := range SQUARES {
		BishopMagics[sq] = find_magic(sq, true)
		RookMagics[sq] = find_magic(sq, false)
	}
}

func main() {
	fmt.Println("OK")
	// test_fen := "4k3/8/8/8/8/8/PPPP4/1N2K3 w - - 0 1"
	// pos := parse_fen(INITIAL_FEN)

	fen := "2bqkbn1/2pppp2/np2N3/r3P1p1/p2N2B1/5Q2/PPPPKPP1/RNB2r2 w - - 0 1"
	// fen := "2bq1bn1/2pppk2/np2N3/r3P1p1/p2N2B1/8/PPPPKPP1/RNB2r2 w - - 0 2"
	pos := parse_fen(fen)
	best := pos.minimax(4)
	fmt.Println(best)
	pos.do_move(best.move, &StateInfo{})
	solution := "2bqkbn1/2pppQ2/np2N3/r3P1p1/p2N2B1/8/PPPPKPP1/RNB2r2 b - - 0 1"
	fmt.Println(best.move)
	fmt.Println(pos.String())
	fmt.Println(solution)

	// fen := "2bq1bn1/2pppk2/np2N3/r3P1pB/p2N4/8/PPPPKPP1/RNB2r2 b - - 1 2"
	// pos := parse_fen(fen)
	// fmt.Println(pos.evaluate(pos.in_checkmate()))
}
