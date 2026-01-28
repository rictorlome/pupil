package main

func init() {
	// Initialize Square Bitboards
	e := Bitboard(0)
	for _, s := range SQUARES {
		var SQUARE_BB Bitboard = 0x1 << s
		SQUARE_BBS[s] = SQUARE_BB
		NEIGHBOR_BBS[s] = neighbors(SQUARE_BB)
		KING_ATTACK_BBS[s] = precomputeKingAttacks(SQUARE_BB)
		KNIGHT_ATTACK_BBS[s] = precomputeKnightAttacks(SQUARE_BB)
		ROOK_ATTACK_MASKS[s] = sliderRookAttacks(Bitboard(0), s)
		BISHOP_ATTACK_MASKS[s] = sliderBishopAttacks(Bitboard(0), s)
		ROOK_OCCUPANCY_MASKS[s] = occupancyMask(s, ROOK_DIRECTIONS)
		BISHOP_OCCUPANCY_MASKS[s] = occupancyMask(s, BISHOP_DIRECTIONS)
		RELEVANT_ROOK_OCCUPANCY[s] = ROOK_ATTACK_MASKS[s] &^ ROOK_OCCUPANCY_MASKS[s]
		RELEVANT_BISHOP_OCCUPANCY[s] = BISHOP_ATTACK_MASKS[s] &^ BISHOP_OCCUPANCY_MASKS[s]
		for _, color := range COLORS {
			PAWN_ATTACK_BBS[s][color] = pawnAttacks(SQUARE_BB, color)
		}
	}
	// Initialize dependent BBs
	for _, fn := range [2]AttackFunc{sliderRookAttacks, sliderBishopAttacks} {
		for _, s1 := range SQUARES {
			for _, s2 := range SQUARES {
				if occupiedAtSq(fn(e, s1), s2) {
					LINE_BBS[s1][s2] = fn(e, s1)&fn(e, s2) | SQUARE_BBS[s1] | SQUARE_BBS[s2]
					BETWEEN_BBS[s1][s2] = fn(SQUARE_BBS[s2], s1) & fn(SQUARE_BBS[s1], s2)
				}
			}
		}
	}

	initCastleSqs()
	initCastlingMasks()
	initRookSquaresForCastling()
	initMagics()
	initZobrists()
	initPool()
}
