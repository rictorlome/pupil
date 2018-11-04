package main

func bishop_attacks(occ Bitboard, square Square) Bitboard {
	return slider_attacks(occ, square, BISHOP_DIRECTIONS)
}

func king_attacks(b Bitboard) Bitboard {
	var attacks Bitboard
	for _, direction := range DIRECTIONS {
		attacks |= shift_direction(b, direction)
	}
	return attacks
}

func knight_attacks(b Bitboard) Bitboard {
	var attacks Bitboard
	for i := 0; i < 8; i++ {
		attacks |= shift_direction(shift_direction(b, KNIGHT_DIRECTIONS[i]), KNIGHT_DIRECTIONS[i+1])
	}
	return attacks
}

func pawn_attacks(pawns Bitboard, color Color) Bitboard {
	if color == BLACK {
		return shift_direction(pawns, SOUTH_EAST) | shift_direction(pawns, SOUTH_WEST)
	}
	return shift_direction(pawns, NORTH_EAST) | shift_direction(pawns, NORTH_WEST)
}

func queen_attacks(occ Bitboard, square Square) Bitboard {
	return slider_attacks(occ, square, DIRECTIONS)
}

func rook_attacks(occ Bitboard, square Square) Bitboard {
	return slider_attacks(occ, square, ROOK_DIRECTIONS)
}

func slider_attacks(occ Bitboard, sq Square, directions []int) Bitboard {
	var attacks Bitboard
	for _, direction := range directions {
		cursor := shift_direction(SQUARE_BBS[sq], direction)
		for i := 0; !empty(cursor); i++ {
			attacks |= cursor
			if occupied_at_bb(occ, cursor) {
				break
			}
			cursor = shift_direction(cursor, direction)
		}
	}
	return attacks
}
