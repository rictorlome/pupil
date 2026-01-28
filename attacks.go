package main

func attacksByColor(occ Bitboard, pieces []Bitboard, color Color) Bitboard {
	return pawnAttacks(pieces[ptToP(PAWN, color)], color) |
		serializeRookAttacks(pieces[ptToP(ROOK, color)], occ) |
		serializeKnightAttacks(pieces[ptToP(KNIGHT, color)], occ) |
		serializeBishopAttacks(pieces[ptToP(BISHOP, color)], occ) |
		serializeQueenAttacks(pieces[ptToP(QUEEN, color)], occ) |
		serializeKingAttacks(pieces[ptToP(KING, color)], occ)
}

// Returns bitboard of occupancy by color c whose pieces attack square sq
func attackersToSqByColor(pieces []Bitboard, sq Square, color Color) Bitboard {
	occ := occupiedSquares(pieces)
	sqBB := SQUARE_BBS[sq]
	var attackers Bitboard
	for _, piece := range pieceRangeByColor(color) {
		t := pieceToType(piece)
		switch t {
		case KING:
			continue
		case PAWN:
			attackers |= pawnAttacks(sqBB, opposite(color)) & pieces[piece]
		case KNIGHT:
			attackers |= knightAttacks(occ, sq) & pieces[piece]
		default:
			attackers |= serializeForAttacks(sqBB, occ, getAttackFunc(t)) & pieces[piece]
		}
	}
	return attackers
}

func bishopAttacks(occ Bitboard, sq Square) Bitboard {
	return BishopAttackTable[attackIndexWithOffset(&BishopMagics[sq], occ)]
}

func getAttackFunc(pt PieceType) AttackFunc {
	return AttackFuncs[pt]
}

func kingAttacks(occ Bitboard, square Square) Bitboard {
	return KING_ATTACK_BBS[square]
}

func knightAttacks(occ Bitboard, square Square) Bitboard {
	return KNIGHT_ATTACK_BBS[square]
}

func nullAttacks(occ Bitboard, sq Square) Bitboard {
	return Bitboard(0)
}

// Mask over occupancy to exclude outermost squares
// where occupancy is irrelevant to calculating attack
func occupancyMask(sq Square, directions []int) Bitboard {
	var mask Bitboard
	for _, direction := range directions {
		for cursor := shiftDirection(SQUARE_BBS[sq], direction); !empty(cursor); cursor = shiftDirection(cursor, direction) {
			if empty(shiftDirection(cursor, direction)) {
				mask |= cursor
			}
		}
	}
	return mask
}

func pawnAttacks(pawns Bitboard, color Color) Bitboard {
	if color == BLACK {
		return shiftDirection(pawns, SOUTH_EAST) | shiftDirection(pawns, SOUTH_WEST)
	}
	return shiftDirection(pawns, NORTH_EAST) | shiftDirection(pawns, NORTH_WEST)
}

func precomputeKingAttacks(b Bitboard) Bitboard {
	var attacks Bitboard
	for _, direction := range DIRECTIONS {
		attacks |= shiftDirection(b, direction)
	}
	return attacks
}

func precomputeKnightAttacks(b Bitboard) Bitboard {
	var attacks Bitboard
	for i := 0; i <= 14; i += 2 {
		attacks |= shiftDirection(shiftDirection(b, KNIGHT_DIRECTIONS[i]), KNIGHT_DIRECTIONS[i+1])
	}
	return attacks
}

func queenAttacks(occ Bitboard, square Square) Bitboard {
	return rookAttacks(occ, square) | bishopAttacks(occ, square)
}

func rookAttacks(occ Bitboard, sq Square) Bitboard {
	return RookAttackTable[attackIndexWithOffset(&RookMagics[sq], occ)]
}

func serializeForAttacks(pieceBB Bitboard, occ Bitboard, fn AttackFunc) Bitboard {
	var attacks Bitboard
	for cursor := pieceBB; cursor != 0; cursor &= cursor - 1 {
		sq := Square(lsb(cursor))
		attacks |= fn(occ, sq)
	}
	return attacks
}

func serializeKnightAttacks(pieceBB Bitboard, occ Bitboard) Bitboard {
	var attacks Bitboard
	for cursor := pieceBB; cursor != 0; cursor &= cursor - 1 {
		attacks |= knightAttacks(occ, Square(lsb(cursor)))
	}
	return attacks
}

func serializeRookAttacks(pieceBB Bitboard, occ Bitboard) Bitboard {
	var attacks Bitboard
	for cursor := pieceBB; cursor != 0; cursor &= cursor - 1 {
		attacks |= rookAttacks(occ, Square(lsb(cursor)))
	}
	return attacks
}

func serializeBishopAttacks(pieceBB Bitboard, occ Bitboard) Bitboard {
	var attacks Bitboard
	for cursor := pieceBB; cursor != 0; cursor &= cursor - 1 {
		attacks |= bishopAttacks(occ, Square(lsb(cursor)))
	}
	return attacks
}

func serializeQueenAttacks(pieceBB Bitboard, occ Bitboard) Bitboard {
	var attacks Bitboard
	for cursor := pieceBB; cursor != 0; cursor &= cursor - 1 {
		attacks |= queenAttacks(occ, Square(lsb(cursor)))
	}
	return attacks
}

func serializeKingAttacks(pieceBB Bitboard, occ Bitboard) Bitboard {
	var attacks Bitboard
	for cursor := pieceBB; cursor != 0; cursor &= cursor - 1 {
		attacks |= kingAttacks(occ, Square(lsb(cursor)))
	}
	return attacks
}

func sliderAttacks(occ Bitboard, sq Square, directions []int) Bitboard {
	var attacks Bitboard
	for _, direction := range directions {
		for cursor := shiftDirection(SQUARE_BBS[sq], direction); !empty(cursor); cursor = shiftDirection(cursor, direction) {
			attacks |= cursor
			if occupiedAtBB(occ, cursor) {
				break
			}
		}
	}
	return attacks
}

func sliderBishopAttacks(occ Bitboard, square Square) Bitboard {
	return sliderAttacks(occ, square, BISHOP_DIRECTIONS)
}

func sliderRookAttacks(occ Bitboard, square Square) Bitboard {
	return sliderAttacks(occ, square, ROOK_DIRECTIONS)
}
