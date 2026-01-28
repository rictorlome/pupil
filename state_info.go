package main

import (
	"fmt"
)

func canCastle(side int, color Color, occ Bitboard, cr int, atks Bitboard) bool {
	safe := atks&CASTLE_CHECK_SQS[int(color)*2+side] == 0
	empty := occ&CASTLE_MOVE_SQS[int(color)*2+side] == 0
	hasRight := cr&CASTLING_RIGHTS[int(color)*2+side] != 0
	return empty && hasRight && safe
}

func cleanupSqForEpCapture(captureDst Square) Square {
	if captureDst < SQ_A5 {
		return Square(captureDst + 8)
	}
	return Square(captureDst - 8)
}

// NOTE: This should only be called on dst squares of DOUBLE_PAWN_PUSH moves
func dstToEpSq(pushDst Square) Square {
	if pushDst < SQ_A5 {
		return Square(pushDst - 8)
	}
	return Square(pushDst + 8)
}

func (s *StateInfo) dup() *StateInfo {
	return &StateInfo{
		s.castlingRights, s.epSq, s.rule50,
		s.key, s.oppositeColorAttacks,
		s.blockersForKing, s.prev, s.captured,
	}
}

func hasRight(rights int, right uint) bool {
	return (rights>>right)&1 == 1
}

func initCastlingMasks() {
	CASTLING_MASK_BY_SQ[SQ_A1] = WQ_CASTLE
	CASTLING_MASK_BY_SQ[SQ_E1] = WHITE_CASTLES
	CASTLING_MASK_BY_SQ[SQ_H1] = WK_CASTLE
	CASTLING_MASK_BY_SQ[SQ_A8] = BQ_CASTLE
	CASTLING_MASK_BY_SQ[SQ_E8] = BLACK_CASTLES
	CASTLING_MASK_BY_SQ[SQ_H8] = BK_CASTLE
}

func initRookSquaresForCastling() {
	m := [12]Square{
		SQ_C1, SQ_A1, SQ_D1,
		SQ_G1, SQ_H1, SQ_F1,
		SQ_C8, SQ_A8, SQ_D8,
		SQ_G8, SQ_H8, SQ_F8,
	}
	for kingDst, rookSrc, rookDst := 0, 1, 2; rookDst < 12; {
		ROOK_SRC_DST[m[kingDst]] = [2]Square{m[rookSrc], m[rookDst]}
		kingDst += 3
		rookSrc += 3
		rookDst += 3
	}
}

func makeCastleStateInfo(available string) int {
	var castles int
	for _, char := range available {
		castles |= CHAR_TO_CASTLE[string(char)]
	}
	return castles
}

func updateCastlingRight(cr int, src Square, dst Square) int {
	return cr &^ (CASTLING_MASK_BY_SQ[src] | CASTLING_MASK_BY_SQ[dst])
}

// NOTE: epSq is only set if enemy pawn can actually take it.
func updateEpSq(m Move, enemyPawns Bitboard) Square {
	dst := moveDst(m)
	if moveType(m) == DOUBLE_PAWN_PUSH && (enemyPawns&NEIGHBOR_BBS[dst] != 0) {
		return dstToEpSq(dst)
	}
	return NULL_SQ
}

func updateRule50(rule50 int, m Move, pt PieceType) int {
	if isCapture(m) || pt == PAWN {
		return 0
	}
	return rule50 + 1
}

func (s StateInfo) String() string {
	recursiveString := "nil"
	if s.prev != nil {
		recursiveString = (*s.prev).String()
	}
	return fmt.Sprintf("[r50 = %v, prev = %v]", s.rule50, recursiveString)
}
