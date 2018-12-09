package main

func attacks_by_color(pieces []Bitboard, color Color) Bitboard {
	occ := occupied_squares(pieces)
	var attacks Bitboard
	for _, piece := range piece_range_by_color(color) {
		t := piece_to_type(piece)
		if t == PAWN {
			attacks |= pawn_attacks(pieces[piece], color)
		} else {
			attacks |= serialize_for_attacks(pieces[piece], occ, get_attack_func(t))
		}
	}
	return attacks
}

func bishop_attacks(occ Bitboard, square Square) Bitboard {
	return slider_attacks(occ, square, BISHOP_DIRECTIONS)
}

func get_attack_func(pt PieceType) AttackFunc {
	return AttackFuncs[pt]
}

func king_attacks(occ Bitboard, square Square) Bitboard {
	return KING_ATTACK_BBS[square]
}

func knight_attacks(occ Bitboard, square Square) Bitboard {
	return KNIGHT_ATTACK_BBS[square]
}

func null_attacks(occ Bitboard, sq Square) Bitboard {
	return Bitboard(0)
}

// Mask over occupancy to exclude outermost squares
// where occupancy is irrelevant to calculating attack
func occupancy_mask(sq Square, directions []int) Bitboard {
	var mask Bitboard
	for _, direction := range directions {
		for cursor := shift_direction(SQUARE_BBS[sq], direction); !empty(cursor); cursor = shift_direction(cursor, direction) {
			if empty(shift_direction(cursor, direction)) {
				mask |= cursor
			}
		}
	}
	return mask
}

func pawn_attacks(pawns Bitboard, color Color) Bitboard {
	if color == BLACK {
		return shift_direction(pawns, SOUTH_EAST) | shift_direction(pawns, SOUTH_WEST)
	}
	return shift_direction(pawns, NORTH_EAST) | shift_direction(pawns, NORTH_WEST)
}

func precompute_king_attacks(b Bitboard) Bitboard {
	var attacks Bitboard
	for _, direction := range DIRECTIONS {
		attacks |= shift_direction(b, direction)
	}
	return attacks
}

func precompute_knight_attacks(b Bitboard) Bitboard {
	var attacks Bitboard
	for i := 0; i <= 14; i += 2 {
		attacks |= shift_direction(shift_direction(b, KNIGHT_DIRECTIONS[i]), KNIGHT_DIRECTIONS[i+1])
	}
	return attacks
}

func queen_attacks(occ Bitboard, square Square) Bitboard {
	return slider_attacks(occ, square, DIRECTIONS)
}

func rook_attacks(occ Bitboard, square Square) Bitboard {
	return slider_attacks(occ, square, ROOK_DIRECTIONS)
}

func serialize_for_attacks(piece_bb Bitboard, occ Bitboard, fn AttackFunc) Bitboard {
	var attacks Bitboard
	for cursor := piece_bb; cursor != 0; cursor &= cursor - 1 {
		sq := Square(lsb(cursor))
		attacks |= fn(occ, sq)
	}
	return attacks
}

func slider_attacks(occ Bitboard, sq Square, directions []int) Bitboard {
	var attacks Bitboard
	for _, direction := range directions {
		for cursor := shift_direction(SQUARE_BBS[sq], direction); !empty(cursor); cursor = shift_direction(cursor, direction) {
			attacks |= cursor
			if occupied_at_bb(occ, cursor) {
				break
			}
		}
	}
	return attacks
}
