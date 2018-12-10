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

func piece_to_color(p Piece) Color {
	return Color(p / 6)
}

func piece_to_type(p Piece) PieceType {
	return PieceType(p % 6)
}
