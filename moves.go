package main

import (
// "fmt"
)

func move_dst(m Move) Square {
	return Square(m &^ DST_MASK)
}

func move_prom(m Move) PromotionType {
	return PROMOTION_TYPES[int((m&^PROMOTION_MASK)>>12)]
}

func move_src(m Move) Square {
	return Square((m &^ SRC_MASK) >> 6)
}

func (m Move) String() string {
	return move_src(m).String() + move_dst(m).String() + move_prom(m).String()
}

func (p PromotionType) String() string {
	return PROMOTION_STRINGS[int(p>>12)]
}

func to_move(dest Square, src Square, mt MoveType, pt PromotionType) Move {
	return Move(uint16(dest) | (uint16(src) << 6) | uint16(mt) | uint16(pt))
}
