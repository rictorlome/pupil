package main

import (
// "fmt"
)

func capOrQuiet(occ Bitboard, dst Square) MoveType {
	// (occ>>dst)&1 returns 1 if occupied, else 0.
	return MoveType(((occ >> dst) & 1) << (2 + 12))
}

func isCapture(m Move) bool {
	return m&Move(CAPTURE) != 0
}

func isCastle(m Move) bool {
	return isMoveType(m, KING_CASTLE) || isMoveType(m, QUEEN_CASTLE)
}

func isEnpassant(m Move) bool {
	return isMoveType(m, EP_CAPTURE)
}

func isMoveType(m Move, mt MoveType) bool {
	return MoveType(m&^MOVE_TYPE_MASK) == mt
}

func isPromotion(m Move) bool {
	return m&Move(PROMOTION_MASK) != 0
}

func moveDst(m Move) Square {
	return Square(m &^ DST_MASK)
}

func moveSrc(m Move) Square {
	return Square((m &^ SRC_MASK) >> 6)
}

func moveType(m Move) MoveType {
	return MoveType(m &^ MOVE_TYPE_MASK)
}

func moveTypeToIdx(mt MoveType) int {
	return int(mt) >> 12
}

func moveTypeToPromotionType(mt MoveType) PieceType {
	return PROMOTION_PIECE_TYPES[(mt &^ MoveType(PROMOTION_MASK) &^ CAPTURE >> 12)]
}

func (m Move) String() string {
	return moveSrc(m).String() + moveDst(m).String() + PROMOTION_STRINGS[moveTypeToIdx(moveType(m))]
}

func (mt MoveType) String() string {
	return MOVE_TYPE_STRINGS[moveTypeToIdx(mt)]
}

func toMove(dest Square, src Square, mt MoveType) Move {
	return Move(uint16(dest) | (uint16(src) << 6) | uint16(mt))
}
