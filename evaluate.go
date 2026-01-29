package main

const (
	DOUBLED_PAWN_PENALTY   = -10
	ISOLATED_PAWN_PENALTY  = -20
	PASSED_PAWN_BONUS      = 20
	BISHOP_PAIR_BONUS      = 30
	ROOK_OPEN_FILE_BONUS   = 15
	ROOK_SEMI_OPEN_BONUS   = 10
	MOBILITY_WEIGHT        = 2
	KING_SHIELD_BONUS      = 10
)

// Material threshold for endgame (no queens, or limited material)
const ENDGAME_THRESHOLD = Q_VAL

func addPositionValue(pt PieceType, color Color, sq Square, isEndgame bool) int {
	idx := relativeSquare(sq, color)
	// Use endgame king table in endgame
	if pt == KING && isEndgame {
		return KING_SQUARE_VALUES_ENDGAME[idx]
	}
	valueArr := *POSITION_VALUES[pt]
	return valueArr[idx]
}

func (p *Position) evaluate() int {
	us := p.sideToMove()
	them := opposite(us)

	// Calculate total material (excluding kings) to detect endgame
	totalMaterial := 0
	for _, pt := range []PieceType{QUEEN, ROOK, BISHOP, KNIGHT} {
		totalMaterial += popcount(p.placement[ptToP(pt, WHITE)]) * MATERIAL_VALUES[pt]
		totalMaterial += popcount(p.placement[ptToP(pt, BLACK)]) * MATERIAL_VALUES[pt]
	}
	isEndgame := totalMaterial <= 2*ENDGAME_THRESHOLD

	score := 0

	// Material and position values
	for piece, pieceBB := range p.placement {
		score += serializeForEvaluate(us, Piece(piece), pieceBB, isEndgame)
	}

	// Pawn structure evaluation
	score += p.evaluatePawnStructure(us) - p.evaluatePawnStructure(them)

	// Bishop pair bonus
	if popcount(p.placement[ptToP(BISHOP, us)]) >= 2 {
		score += BISHOP_PAIR_BONUS
	}
	if popcount(p.placement[ptToP(BISHOP, them)]) >= 2 {
		score -= BISHOP_PAIR_BONUS
	}

	// Rook on open/semi-open files
	score += p.evaluateRooks(us) - p.evaluateRooks(them)

	// King safety (only in middlegame)
	if !isEndgame {
		score += p.evaluateKingSafety(us) - p.evaluateKingSafety(them)
	}

	return score
}

// Evaluate pawn structure: doubled, isolated, passed pawns
func (p *Position) evaluatePawnStructure(color Color) int {
	score := 0
	ourPawns := p.placement[ptToP(PAWN, color)]
	theirPawns := p.placement[ptToP(PAWN, opposite(color))]

	// Doubled pawns: count pawns that have another pawn on the same file behind them
	var doubled Bitboard
	if color == WHITE {
		// Pawns that have a pawn directly south of them
		doubled = ourPawns & (ourPawns >> NORTH)
	} else {
		// Pawns that have a pawn directly north of them
		doubled = ourPawns & (ourPawns << NORTH)
	}
	score += DOUBLED_PAWN_PENALTY * popcount(doubled)

	// Isolated pawns: pawns with no friendly pawns on adjacent files
	for cursor := ourPawns; !empty(cursor); cursor &= cursor - 1 {
		sq := Square(lsb(cursor))
		file := squareFile(sq)
		if empty(ourPawns & ADJACENT_FILES[file]) {
			score += ISOLATED_PAWN_PENALTY
		}

		// Passed pawns: no enemy pawns ahead on same or adjacent files
		if empty(theirPawns & PASSED_PAWN_MASKS[sq][color]) {
			// Bonus increases as pawn advances
			rank := relativeRank(sq, color)
			score += PASSED_PAWN_BONUS * rank / 2
		}
	}

	return score
}

// Relative rank from color's perspective (0-7, where 7 is promotion rank)
func relativeRank(sq Square, color Color) int {
	rank := squareRank(sq)
	if color == BLACK {
		return 7 - rank
	}
	return rank
}

// Evaluate rooks on open/semi-open files
func (p *Position) evaluateRooks(color Color) int {
	score := 0
	rooks := p.placement[ptToP(ROOK, color)]
	ourPawns := p.placement[ptToP(PAWN, color)]
	theirPawns := p.placement[ptToP(PAWN, opposite(color))]

	for cursor := rooks; !empty(cursor); cursor &= cursor - 1 {
		sq := Square(lsb(cursor))
		file := squareFile(sq)
		fileBB := FILE_BBS[file]

		ourPawnsOnFile := ourPawns & fileBB
		theirPawnsOnFile := theirPawns & fileBB

		if empty(ourPawnsOnFile) {
			if empty(theirPawnsOnFile) {
				score += ROOK_OPEN_FILE_BONUS // Open file
			} else {
				score += ROOK_SEMI_OPEN_BONUS // Semi-open file
			}
		}
	}

	return score
}

// Evaluate king safety based on pawn shield
func (p *Position) evaluateKingSafety(color Color) int {
	kingSq := p.kingSquare(color)
	ourPawns := p.placement[ptToP(PAWN, color)]
	// Count pawns in the shield area using precomputed mask
	shieldPawns := popcount(ourPawns & KING_SHIELD_MASKS[kingSq][color])
	return KING_SHIELD_BONUS * shieldPawns
}

func relativeMultiplier(color Color) int {
	if color == WHITE {
		return 1
	}
	return -1
}

func relativeSquare(sq Square, color Color) Square {
	return makeSquare(relativeSquareRank(sq, color), squareFile(sq))
}

func relativeSquareRank(sq Square, color Color) int {
	if color == BLACK {
		return squareRank(sq)
	}
	return 7 - squareRank(sq)
}

func serializeForEvaluate(us Color, piece Piece, pieceBB Bitboard, isEndgame bool) int {
	score := 0
	color, pt := pieceToColor(piece), pieceToType(piece)
	for cursor := pieceBB; !empty(cursor); cursor &= cursor - 1 {
		if color == us {
			score += MATERIAL_VALUES[pt]
			score += addPositionValue(pt, color, Square(lsb(cursor)), isEndgame)
		} else {
			score -= MATERIAL_VALUES[pt]
			score -= addPositionValue(pt, color, Square(lsb(cursor)), isEndgame)
		}
	}
	return score
}
