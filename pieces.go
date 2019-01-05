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
	return PIECE_TYPES[p]
}

func placement_to_placement_by_square(placement []Bitboard) []Piece {
	placement_by_square := make([]Piece, 64)
	for _, SQ := range SQUARES {
		placement_by_square[SQ] = piece_on_sq(placement, SQ)
	}
	return placement_by_square
}

func pt_to_p(pt PieceType, color Color) Piece {
	return Piece(pt + PieceType(6*color))
}

func (pt PieceType) String() string {
	return PIECE_TYPE_STRINGS[pt]
}
