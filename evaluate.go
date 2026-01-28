package main

func addPositionValue(pt PieceType, color Color, sq Square) int {
	idx := relativeSquare(sq, color)
	valueArr := *POSITION_VALUES[pt]
	return valueArr[idx]
}

func (p *Position) evaluate() int {
	score := 0
	for piece, pieceBB := range p.placement {
		score += serializeForEvaluate(p.sideToMove(), Piece(piece), pieceBB)
	}
	return score
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

func serializeForEvaluate(us Color, piece Piece, pieceBB Bitboard) int {
	score := 0
	color, pt := pieceToColor(piece), pieceToType(piece)
	for cursor := pieceBB; !empty(cursor); cursor &= cursor - 1 {
		if color == us {
			score += MATERIAL_VALUES[pt]
			score += addPositionValue(pt, color, Square(lsb(cursor)))
		} else {
			score -= MATERIAL_VALUES[pt]
			score -= addPositionValue(pt, color, Square(lsb(cursor)))
		}
	}
	return score
}
