package main

import (
// "fmt"
)

func cap_or_quiet(occ Bitboard, dst Square) MoveType {
	// (occ>>dst)&1 returns 1 if occupied, else 0.
	return MoveType(((occ >> dst) & 1) << (2 + 12))
}

func is_capture(m Move) bool {
	return m&Move(CAPTURE) != 0
}

func is_castle(m Move) bool {
	return is_move_type(m, KING_CASTLE) || is_move_type(m, QUEEN_CASTLE)
}

func is_enpassant(m Move) bool {
	return is_move_type(m, EP_CAPTURE)
}

func is_move_type(m Move, mt MoveType) bool {
	return MoveType(m&^MOVE_TYPE_MASK) == mt
}

func is_promotion(m Move) bool {
	return m&Move(PROMOTION_MASK) != 0
}

func move_dst(m Move) Square {
	return Square(m &^ DST_MASK)
}

func move_src(m Move) Square {
	return Square((m &^ SRC_MASK) >> 6)
}

func move_type(m Move) MoveType {
	//TODO: remove this lookup, just and with the mask
	return MOVE_TYPES[int((m&^MOVE_TYPE_MASK)>>12)]
}

func move_type_to_idx(mt MoveType) int {
	return int(mt) >> 12
}

func move_type_to_promotion_type(mt MoveType) PieceType {
	return PROMOTION_PIECE_TYPES[(mt &^ MoveType(PROMOTION_MASK) &^ CAPTURE >> 12)]
}

func (m Move) String() string {
	return move_src(m).String() + move_dst(m).String() + PROMOTION_STRINGS[move_type_to_idx(move_type(m))]
}

func (mt MoveType) String() string {
	return MOVE_TYPE_STRINGS[move_type_to_idx(mt)]
}

func to_move(dest Square, src Square, mt MoveType) Move {
	return Move(uint16(dest) | (uint16(src) << 6) | uint16(mt))
}
