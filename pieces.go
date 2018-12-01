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

func piece_to_type(p Piece) PieceType {
	return PieceType(p % 6)
}
