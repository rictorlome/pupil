package main

import (
	// "fmt"
	"strconv"
)

func (p *Position) doCastle(do bool, r Piece, rSrc Square, rDst Square) {
	if do {
		p.removePiece(r, rSrc)
		p.placePiece(r, rDst)
	} else {
		p.removePiece(r, rDst)
		p.placePiece(r, rSrc)
	}
}

func (p *Position) doMove(m Move, newState *StateInfo) {
	src, dst := moveSrc(m), moveDst(m)
	mover := p.pieceAt(src)
	us, them := p.sideToMove(), opposite(p.sideToMove())
	// Initialze new key (cancel current castling rights and ep sq to avoid conditional)
	newKey := p.state.key ^ ZOBRIST_SIDE ^ ZOBRIST_CSTL[p.state.castlingRights] ^ ZOBRIST_EPSQ[p.state.epSq]
	// Update new state
	newState.castlingRights = updateCastlingRight(p.state.castlingRights, src, dst)
	newState.epSq = updateEpSq(m, p.placement[ptToP(PAWN, them)])
	newState.rule50 = updateRule50(p.state.rule50, m, p.pieceTypeAt(src))
	// Update new key
	newKey ^= ZOBRIST_CSTL[newState.castlingRights] ^ ZOBRIST_EPSQ[newState.epSq]
	// Update placement
	if isCastle(m) {
		rookSrcDst := ROOK_SRC_DST[dst]
		rSrc, rDst := rookSrcDst[0], rookSrcDst[1]
		r := p.pieceAt(rSrc)
		p.doCastle(true, r, rSrc, rDst)
		newKey ^= ZOBRIST_PSQ[rSrc][r] ^ ZOBRIST_PSQ[rDst][r]
	} else if isEnpassant(m) {
		epSq := cleanupSqForEpCapture(dst)
		captured := p.pieceAt(epSq)
		newState.captured = captured
		p.removePiece(captured, epSq)
		newKey ^= ZOBRIST_PSQ[epSq][captured]
	} else if isCapture(m) {
		captured := p.pieceAt(dst)
		newState.captured = captured
		p.removePiece(captured, dst)
		newKey ^= ZOBRIST_PSQ[dst][captured]
	}

	p.removePiece(mover, src)
	newKey ^= ZOBRIST_PSQ[src][mover]
	if isPromotion(m) {
		pt := moveTypeToPromotionType(moveType(m))
		pc := ptToP(pt, us)
		p.placePiece(pc, dst)
		newKey ^= ZOBRIST_PSQ[dst][pc]
	} else {
		p.placePiece(mover, dst)
		newKey ^= ZOBRIST_PSQ[dst][mover]
	}

	// Update king blockers (for next turn)
	// our sliders, their king
	newState.key = newKey
	newState.oppositeColorAttacks = p.getColorAttacks(us)
	newState.blockersForKing = p.sliderBlockers(us, p.kingSquare(them))

	// Reassign state
	newState.prev = p.state
	p.state = newState

	// Update position
	p.ply += 1
	p.stm = opposite(p.stm)
}

func (p *Position) placePiece(pc Piece, sq Square) {
	sqBb := SQUARE_BBS[sq]
	p.occ |= sqBb
	p.placementBySquare[sq] = pc
	p.placement[pc] |= sqBb
}

func (p *Position) removePiece(pc Piece, sq Square) {
	sqBb := SQUARE_BBS[sq]
	p.occ &^= sqBb
	p.placementBySquare[sq] = NULL_PIECE
	p.placement[pc] &^= sqBb
}

func (p *Position) setFenInfo(positions string, color string, castles string, enps string, rule50 string, moveCount int) {
	rule50Int, _ := strconv.Atoi(rule50)
	// On Position
	p.placement = parsePositions(positions)
	p.placementBySquare = placementToPlacementBySquare(p.placement)
	p.occ = occupiedSquares(p.placement)
	p.ply = moveCount*2 + int(parseColor(color))
	p.stm = parseColor(color)
	// On StateInfo
	p.state.castlingRights = makeCastleStateInfo(castles)
	p.state.epSq = parseSquare(enps)
	p.state.rule50 = rule50Int
}

// Null move: just swap sides without making a move
func (p *Position) doNullMove(newState *StateInfo) {
	us := p.sideToMove()
	them := opposite(us)

	// Update key: flip side, clear ep square
	newKey := p.state.key ^ ZOBRIST_SIDE ^ ZOBRIST_EPSQ[p.state.epSq]

	// Copy state info
	newState.castlingRights = p.state.castlingRights
	newState.epSq = NULL_SQ // Clear ep square
	newState.rule50 = p.state.rule50 + 1

	// Update key with new ep square (NULL_SQ)
	newKey ^= ZOBRIST_EPSQ[newState.epSq]
	newState.key = newKey

	// Update attacks and blockers for the new side to move
	newState.oppositeColorAttacks = p.getColorAttacks(us)
	newState.blockersForKing = p.sliderBlockers(us, p.kingSquare(them))

	// Link states
	newState.prev = p.state
	p.state = newState

	// Swap sides
	p.ply += 1
	p.stm = them
}

func (p *Position) undoNullMove() {
	p.ply -= 1
	p.stm = opposite(p.stm)
	p.state = p.state.prev
}

func (p *Position) undoMove(m Move) {
	src, dst := moveSrc(m), moveDst(m)
	mover := p.pieceAt(dst)

	// turn has already been updated
	us := opposite(p.sideToMove())
	// Update position
	p.ply -= 1
	p.stm = opposite(p.stm)

	// move piece back to src
	p.removePiece(mover, dst)
	if isPromotion(m) {
		p.placePiece(ptToP(PAWN, us), src)
	} else {
		p.placePiece(mover, src)
	}

	// if capture, replace piece
	if isCapture(m) {
		capsq := dst
		if isEnpassant(m) {
			capsq = cleanupSqForEpCapture(dst)
		}
		p.placePiece(p.state.captured, capsq)
	}

	// if castle, undo castle
	if isCastle(m) {
		rookSrcDst := ROOK_SRC_DST[dst]
		rSrc, rDst := rookSrcDst[0], rookSrcDst[1]
		r := p.pieceAt(rDst)
		p.doCastle(false, r, rSrc, rDst)
	}

	// Reassign state
	p.state = p.state.prev
}
