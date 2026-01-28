package main

import (
// "fmt"
)

func pieceOnSq(pieces []Bitboard, sq Square) Piece {
	for piece, bb := range pieces {
		if occupiedAtSq(bb, sq) {
			return Piece(piece)
		}
	}
	return NULL_PIECE
}

func pieceToColor(p Piece) Color {
	return Color(p / 6)
}

func pieceToType(p Piece) PieceType {
	return PIECE_TYPES[p]
}

func placementToPlacementBySquare(placement []Bitboard) []Piece {
	placementBySquare := make([]Piece, 64)
	for _, SQ := range SQUARES {
		placementBySquare[SQ] = pieceOnSq(placement, SQ)
	}
	return placementBySquare
}

func ptToP(pt PieceType, color Color) Piece {
	return Piece(pt + PieceType(6*color))
}

func (pt PieceType) String() string {
	return PIECE_TYPE_STRINGS[pt]
}
