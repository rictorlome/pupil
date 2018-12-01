package main

import (
// "fmt"
)

func piece_on_sq(pieces []Bitboard, sq Square) Piece {
	for piece, bb := range pieces {
		if occupied_at_sq(bb, sq) {
			return Piece(piece)
		}
	}
	return NULL_PIECE
}

func to_move(dest Square, src Square, mt MoveType, pt PromotionType) Move {
	return Move(uint16(dest) | (uint16(src) << 6) | uint16(mt) | uint16(pt))
}
