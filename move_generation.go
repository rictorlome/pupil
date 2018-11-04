package main

func bishop_attacks(occ Bitboard, square Square) Bitboard {
	return slider_attacks(occ, square, BISHOP_DIRECTIONS)
}

func king_attacks(b Bitboard) Bitboard {
	var x Bitboard
	for _, direction := range DIRECTIONS {
		x |= shift_direction(b, direction)
	}
	return x
}

func knight_attacks(b Bitboard) Bitboard {
	var result Bitboard
	for i := 0; i < 8; i++ {
		result |= shift_direction(shift_direction(b, KNIGHT_DIRECTIONS[i]), KNIGHT_DIRECTIONS[i+1])
	}
	return result
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

func slider_attacks(occ Bitboard, sq Square, directions []int) Bitboard {
	var result Bitboard
	for _, direction := range directions {
		cursor := shift_direction(SQUARE_BBS[sq], direction)
		for i := 1; !empty(cursor); i++ {
			result |= cursor
			if occupied_at_sq_bb(occ, cursor) {
				break
			}
			cursor = shift_direction(cursor, direction)
		}
	}
	return result
}
