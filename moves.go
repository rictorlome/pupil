package main

import (
// "fmt"
)

func cap_or_quiet(occ Bitboard, dst Square) MoveType {
	// (occ>>dst)&1 returns 1 if occupied, else 0.
	return MoveType(((occ >> dst) & 1) << 2)
}

func is_capture(m Move) bool {
	return m&Move(CAPTURE) != 0
}

func is_castle(m Move) bool {
	return m&(Move(KING_CASTLE)|Move(QUEEN_CASTLE)) != 0
}

func is_enpassant(m Move) bool {
	return m&(Move(EP_CAPTURE)) != 0
}

func move_dst(m Move) Square {
	return Square(m &^ DST_MASK)
}

func move_src(m Move) Square {
	return Square((m &^ SRC_MASK) >> 6)
}

func move_type(m Move) MoveType {
	return MOVE_TYPES[int((m&^MOVE_TYPE_MASK)>>12)]
}

func (m Move) String() string {
	return move_src(m).String() + move_dst(m).String() + PROMOTION_STRINGS[move_type(m)]
}

func to_move(dest Square, src Square, mt MoveType) Move {
	return Move(uint16(dest) | (uint16(src) << 6) | uint16(mt))
}
