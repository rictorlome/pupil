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
	initEvalBitboards()
}

func initEvalBitboards() {
	// Adjacent files
	for file := FILE_A; file <= FILE_H; file++ {
		ADJACENT_FILES[file] = Bitboard(0)
		if file > FILE_A {
			ADJACENT_FILES[file] |= FILE_BBS[file-1]
		}
		if file < FILE_H {
			ADJACENT_FILES[file] |= FILE_BBS[file+1]
		}
	}

	// Passed pawn masks and king shield masks
	for _, sq := range SQUARES {
		file := squareFile(sq)
		rank := squareRank(sq)

		// Passed pawn masks: all squares ahead on same and adjacent files
		for _, color := range COLORS {
			mask := Bitboard(0)
			if color == WHITE {
				for r := rank + 1; r <= RANK_8; r++ {
					mask |= SQUARE_BBS[makeSquare(r, file)]
					if file > FILE_A {
						mask |= SQUARE_BBS[makeSquare(r, file-1)]
					}
					if file < FILE_H {
						mask |= SQUARE_BBS[makeSquare(r, file+1)]
					}
				}
			} else {
				for r := rank - 1; r >= RANK_1; r-- {
					mask |= SQUARE_BBS[makeSquare(r, file)]
					if file > FILE_A {
						mask |= SQUARE_BBS[makeSquare(r, file-1)]
					}
					if file < FILE_H {
						mask |= SQUARE_BBS[makeSquare(r, file+1)]
					}
				}
			}
			PASSED_PAWN_MASKS[sq][color] = mask
		}

		// King shield masks: 3 squares directly in front of king
		for _, color := range COLORS {
			shield := Bitboard(0)
			var shieldRank int
			var validRank bool
			if color == WHITE {
				shieldRank = rank + 1
				validRank = rank < RANK_8
			} else {
				shieldRank = rank - 1
				validRank = rank > RANK_1
			}
			if validRank {
				for f := max(FILE_A, file-1); f <= min(FILE_H, file+1); f++ {
					shield |= SQUARE_BBS[makeSquare(shieldRank, f)]
				}
			}
			KING_SHIELD_MASKS[sq][color] = shield
		}
	}
}
