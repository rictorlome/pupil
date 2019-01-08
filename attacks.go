package main

import (
// "fmt"
)

func attacks_by_color(occ Bitboard, pieces []Bitboard, color Color) Bitboard {
	return pawn_attacks(pieces[pt_to_p(PAWN, color)], color) |
		serialize_rook_attacks(pieces[pt_to_p(ROOK, color)], occ) |
		serialize_knight_attacks(pieces[pt_to_p(KNIGHT, color)], occ) |
		serialize_bishop_attacks(pieces[pt_to_p(BISHOP, color)], occ) |
		serialize_queen_attacks(pieces[pt_to_p(QUEEN, color)], occ) |
		serialize_king_attacks(pieces[pt_to_p(KING, color)], occ)
}

// Returns bitboard of occupancy by color c whose pieces attack square sq
func attackers_to_sq_by_color(pieces []Bitboard, sq Square, color Color) Bitboard {
	occ := occupied_squares(pieces)
	sq_bb := SQUARE_BBS[sq]
	var attackers Bitboard
	for _, piece := range piece_range_by_color(color) {
		t := piece_to_type(piece)
		switch t {
		case KING:
			continue
		case PAWN:
			attackers |= pawn_attacks(sq_bb, opposite(color)) & pieces[piece]
		case KNIGHT:
			attackers |= knight_attacks(occ, sq) & pieces[piece]
		default:
			attackers |= serialize_for_attacks(sq_bb, occ, get_attack_func(t)) & pieces[piece]
		}
	}
	return attackers
}

func bishop_attacks(occ Bitboard, sq Square) Bitboard {
	return BishopAttackTable[attack_index_with_offset(&BishopMagics[sq], occ)]
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
	return rook_attacks(occ, square) | bishop_attacks(occ, square)
}

func rook_attacks(occ Bitboard, sq Square) Bitboard {
	return RookAttackTable[attack_index_with_offset(&RookMagics[sq], occ)]
}

func serialize_for_attacks(piece_bb Bitboard, occ Bitboard, fn AttackFunc) Bitboard {
	var attacks Bitboard
	for cursor := piece_bb; cursor != 0; cursor &= cursor - 1 {
		sq := Square(lsb(cursor))
		attacks |= fn(occ, sq)
	}
	return attacks
}

func serialize_knight_attacks(piece_bb Bitboard, occ Bitboard) Bitboard {
	var attacks Bitboard
	for cursor := piece_bb; cursor != 0; cursor &= cursor - 1 {
		attacks |= knight_attacks(occ, Square(lsb(cursor)))
	}
	return attacks
}

func serialize_rook_attacks(piece_bb Bitboard, occ Bitboard) Bitboard {
	var attacks Bitboard
	for cursor := piece_bb; cursor != 0; cursor &= cursor - 1 {
		attacks |= rook_attacks(occ, Square(lsb(cursor)))
	}
	return attacks
}

func serialize_bishop_attacks(piece_bb Bitboard, occ Bitboard) Bitboard {
	var attacks Bitboard
	for cursor := piece_bb; cursor != 0; cursor &= cursor - 1 {
		attacks |= bishop_attacks(occ, Square(lsb(cursor)))
	}
	return attacks
}

func serialize_queen_attacks(piece_bb Bitboard, occ Bitboard) Bitboard {
	var attacks Bitboard
	for cursor := piece_bb; cursor != 0; cursor &= cursor - 1 {
		attacks |= queen_attacks(occ, Square(lsb(cursor)))
	}
	return attacks
}

func serialize_king_attacks(piece_bb Bitboard, occ Bitboard) Bitboard {
	var attacks Bitboard
	for cursor := piece_bb; cursor != 0; cursor &= cursor - 1 {
		attacks |= king_attacks(occ, Square(lsb(cursor)))
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

func slider_bishop_attacks(occ Bitboard, square Square) Bitboard {
	return slider_attacks(occ, square, BISHOP_DIRECTIONS)
}

func slider_rook_attacks(occ Bitboard, square Square) Bitboard {
	return slider_attacks(occ, square, ROOK_DIRECTIONS)
}
