package main

import (
// "fmt"
)

func piece_to_type(p Piece) PieceType {
	return PieceType(p % 6)
}

func to_move(dest Square, src Square, mt MoveType, pt PromotionType) Move {
	return Move(uint16(dest) | (uint16(src) << 6) | uint16(mt) | uint16(pt))
}
