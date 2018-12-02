package main

import (
// "fmt"
)

func is_slider(pt PieceType) bool {
	return pt == BISHOP || pt == ROOK || pt == QUEEN
}

func not_k_or_p(pt PieceType) bool {
	return pt != KING && pt != PAWN
}

func piece_on_sq(pieces []Bitboard, sq Square) Piece {
	for piece, bb := range pieces {
		if occupied_at_sq(bb, sq) {
			return Piece(piece)
		}
	}
	return NULL_PIECE
}

func piece_to_type(p Piece) PieceType {
	return PieceType(p % 6)
}
