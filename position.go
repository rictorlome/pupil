package main

import (
// "fmt"
)

func (p *Position) dup() Position {
	newPlacement := make([]Bitboard, len(p.placement))
	newPlacementBySquare := make([]Piece, len(p.placementBySquare))
	copy(newPlacement, p.placement)
	copy(newPlacementBySquare, p.placementBySquare)
	return Position{
		p.occ,
		newPlacement,
		newPlacementBySquare,
		p.ply,
		p.state.dup(),
		p.stm,
	}
}

func (p *Position) generateEvasions(pl *[]Move, ml *[]Move) {
	us, them := p.sideToMove(), opposite(p.sideToMove())
	kingSq := p.kingSquare(us)

	occ, selfOcc := p.occupancy(), p.occupancyByColor(us)
	occWithoutKing := occ &^ SQUARE_BBS[kingSq]

	checkers := attackersToSqByColor(p.placement, kingSq, them)
	atks := attacksByColor(occWithoutKing, p.placement, them)

	// King evasions
	serializeNormalMoves(ml, kingSq, kingAttacks(occ, kingSq)&^(selfOcc|atks), occ)
	if popcount(checkers) > 1 {
		return
	}
	checkerSq := Square(lsb(checkers))
	// Regular moves (duplicate king evasions are excluded in pseudoLegalsByColor)
	p.generateNonEvasions(pl, ml, BETWEEN_BBS[kingSq][checkerSq]|checkers)
}

// generateMoves returns a pooled slice - caller must call putMoveList when done
func (p *Position) generateMoves() *[]Move {
	pseudoLegalMoveList := getMoveList()
	moveList := getMoveList()
	if p.inCheck() {
		p.generateEvasions(pseudoLegalMoveList, moveList)
	} else {
		p.generateNonEvasions(pseudoLegalMoveList, moveList, Bitboard(0))
	}
	putMoveList(pseudoLegalMoveList)
	return moveList
}

// NOTE: pseudolegal moves include those that cause check. these have to be filtered out in move generation
func (p *Position) generatePseudoLegals(pl *[]Move, forcedDsts Bitboard) {
	us := p.sideToMove()
	occ, selfOcc := p.occupancy(), p.occupancyByColor(us)

	serializeForPseudosPawns(pl, p.ourPtBb(PAWN), occ, selfOcc, us, p.state.epSq)
	serializeForPseudosOther(pl, p.ourPtBb(KNIGHT), occ, selfOcc, knightAttacks)
	serializeForPseudosOther(pl, p.ourPtBb(ROOK), occ, selfOcc, rookAttacks)
	serializeForPseudosOther(pl, p.ourPtBb(BISHOP), occ, selfOcc, bishopAttacks)
	serializeForPseudosOther(pl, p.ourPtBb(QUEEN), occ, selfOcc, queenAttacks)

	if empty(forcedDsts) {
		serializeForPseudosKing(pl, p.ourPtBb(KING), occ, selfOcc, us, p.state.castlingRights, p.state.oppositeColorAttacks)
	}
}

func (p *Position) generateNonEvasions(pl *[]Move, ml *[]Move, forcedDsts Bitboard) {
	p.generatePseudoLegals(pl, forcedDsts)
	for _, pseudoLegal := range *pl {
		if p.isLegal(pseudoLegal) && (empty(forcedDsts) || isGoodEvasion(forcedDsts, pseudoLegal)) {
			*ml = append(*ml, pseudoLegal)
		}
	}
}

func (p *Position) getColorAttacks(color Color) Bitboard {
	return attacksByColor(p.occupancy(), p.placement, color)
}

func isGoodEvasion(forcedDsts Bitboard, m Move) bool {
	return occupiedAtSq(forcedDsts, moveDst(m)) ||
		(isEnpassant(m) && occupiedAtSq(forcedDsts, cleanupSqForEpCapture(moveDst(m))))
}

func (p *Position) inCheck() bool {
	color := p.sideToMove()
	atks := p.oppositeColorAttacks()
	return occupiedAtSq(atks, p.kingSquare(color))
}

func (p *Position) inCheckmate() bool {
	if !p.inCheck() {
		return false
	}
	moves := p.generateMoves()
	result := len(*moves) == 0
	putMoveList(moves)
	return result
}

func (p *Position) isLegal(m Move) bool {
	src, dst := moveSrc(m), moveDst(m)
	us, them := p.sideToMove(), opposite(p.sideToMove())
	ksq := p.kingSquare(us)

	if isEnpassant(m) {
		theirQueens := p.placement[ptToP(QUEEN, them)]
		capturedSq := cleanupSqForEpCapture(dst)
		// Occupancy after en passant: remove src and captured pawn, add dst
		occAfterEp := (p.occupancy() &^ (SQUARE_BBS[src] | SQUARE_BBS[capturedSq])) | SQUARE_BBS[dst]
		// No discovered slider attacks on the king.
		return empty(rookAttacks(occAfterEp, ksq)&(p.placement[ptToP(ROOK, them)]|theirQueens)) &&
			empty(bishopAttacks(occAfterEp, ksq)&(p.placement[ptToP(BISHOP, them)]|theirQueens))
	}
	if p.pieceTypeAt(src) == KING {
		return isCastle(m) || !occupiedAtSq(p.oppositeColorAttacks(), dst)
	}
	return !occupiedAtSq(p.state.blockersForKing, src) || aligned(src, dst, ksq)
}

func (p *Position) kingSquare(color Color) Square {
	return Square(lsb(p.placement[color*6]))
}

func (p *Position) moveCount() int {
	return p.ply / 2
}

func (p *Position) occupancy() Bitboard {
	return p.occ
}

func (p *Position) occupancyByColor(c Color) Bitboard {
	return occupiedSquaresByColor(p.placement, c)
}

func (p *Position) occupancyByPiece(pc Piece) Bitboard {
	return p.placement[pc]
}

func (p *Position) occupancyByPieces(pieces ...Piece) Bitboard {
	var occupancy Bitboard
	for _, piece := range pieces {
		occupancy |= p.occupancyByPiece(piece)
	}
	return occupancy
}

func (p *Position) occupancyByPieceType(pt PieceType) Bitboard {
	return p.occupancyByPieces(ptToP(pt, WHITE), ptToP(pt, BLACK))
}

func (p *Position) occupancyByPieceTypes(pts ...PieceType) Bitboard {
	var occupancy Bitboard
	for _, pt := range pts {
		occupancy |= p.occupancyByPieceType(pt)
	}
	return occupancy
}

func (p *Position) occupiedAt(sq Square) bool {
	return p.pieceAt(sq) != NULL_PIECE
}

func (p *Position) oppositeColorAttacks() Bitboard {
	return p.state.oppositeColorAttacks
}

func (p *Position) ourPtBb(pt PieceType) Bitboard {
	return p.placement[ptToP(pt, p.sideToMove())]
}

func (p *Position) parseMove(s string) Move {
	src := parseSquare(s[:2])
	dst := parseSquare(s[2:4])
	moverType, captured := pieceToType(p.pieceAt(src)), p.pieceAt(dst)
	return toMove(dst, src, parseMoveType(s[4:], p.occupancy(), src, dst, moverType, captured))
}

func (p *Position) pieceAt(sq Square) Piece {
	return p.placementBySquare[sq]
}

func (p *Position) pieceTypeAt(sq Square) PieceType {
	return pieceToType(p.pieceAt(sq))
}

func (p *Position) toZobrist() Key {
	var ResKey Key
	for sq, piece := range p.placementBySquare {
		if piece != NULL_PIECE {
			ResKey ^= ZOBRIST_PSQ[sq][piece]
		}
	}
	ResKey ^= ZOBRIST_CSTL[p.state.castlingRights]
	ResKey ^= ZOBRIST_EPSQ[p.state.epSq]
	ResKey ^= ZOBRIST_SIDE * Key(p.stm)
	return ResKey
}

func (p *Position) sideToMove() Color {
	return p.stm
}

// Returns bitboard of all pieces blocking attacks to sq from sliders of color c.
func (p *Position) sliderBlockers(c Color, sq Square) Bitboard {
	var blockers Bitboard

	occ := p.occupancy()

	// QUESTION: is the attack mask sufficient here, or do we need pseudolegal moves?
	queenOcc := p.occupancyByPiece(ptToP(QUEEN, c))
	possibleSnipers := ((p.occupancyByPiece(ptToP(BISHOP, c)) | queenOcc) & BISHOP_ATTACK_MASKS[sq]) |
		((p.occupancyByPiece(ptToP(ROOK, c)) | queenOcc) & ROOK_ATTACK_MASKS[sq])

	for cursor := possibleSnipers; cursor != 0; cursor &= cursor - 1 {
		sniperSq := Square(lsb(cursor))
		intermediatePieces := BETWEEN_BBS[sq][sniperSq] & occ
		// NOTE: can set pinners in this if as well, when I start keeping track of them.
		if popcount(intermediatePieces) == 1 {
			blockers |= intermediatePieces
		}
	}

	return blockers
}

func (p Position) String() string {
	return generateFen(p)
}
