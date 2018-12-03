package main

import (
// "fmt"
)

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

func pawn_pushes(sq Square, dir int, occ Bitboard) Bitboard {
	pushes := signed_shift(SQUARE_BBS[sq], dir) &^ occ
	if square_rank(sq) == second_rank(dir) {
		pushes |= signed_shift(pushes, dir) &^ occ
	}
	return pushes
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

func pseudolegals_by_color(pieces []Bitboard, color Color, st StateInfo) []Move {
	occ, self_occ := occupied_squares(pieces), occupied_squares_by_color(pieces, color)
	var move_list []Move

	for _, piece := range piece_range_by_color(color) {
		t := piece_to_type(piece)
		switch t {
		case PAWN:
			move_list = append(move_list, serialize_for_pseudos_pawns(pieces[piece], occ, self_occ, color, st)...)
		default:
			move_list = append(move_list, serialize_for_pseudos(pieces[piece], occ, self_occ, get_attack_func(t))...)
		}
	}
	return move_list
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

func serialize_for_pseudos(piece_bb Bitboard, occ Bitboard, self_occ Bitboard, fn AttackFunc) []Move {
	var move_list []Move
	for cursor := piece_bb; cursor != 0; cursor &= cursor - 1 {
		src := Square(lsb(cursor))
		pseudos := fn(occ, src) &^ self_occ
		for dst_cursor := pseudos; dst_cursor != 0; dst_cursor &= dst_cursor - 1 {
			dst := Square(lsb(dst_cursor))
			move_list = append(move_list, to_move(dst, src, NORMAL, NO_PROMOTION))
		}
	}
	return move_list
}

func serialize_for_pseudos_pawns(pawns Bitboard, occ Bitboard, self_occ Bitboard, color Color, st StateInfo) []Move {
	var move_list []Move
	f_dir, l_rank, idx := forward(color), last_rank(color), color_to_int(color)
	enp_sq := get_enp_sq(st)
	for cursor := pawns; cursor != 0; cursor &= cursor - 1 {
		src := Square(lsb(cursor))
		attacks := (PAWN_ATTACK_BBS[src][idx] & (occ | SQUARE_BBS[enp_sq])) &^ self_occ
		pseudos := attacks | pawn_pushes(src, f_dir, occ)
		for dst_cursor := pseudos; dst_cursor != 0; dst_cursor &= dst_cursor - 1 {
			dst := Square(lsb(dst_cursor))
			switch {
			case square_rank(dst) == l_rank:
				move_list = append(move_list, to_move(dst, src, PROMOTION, KNIGHT_PROMOTION), to_move(dst, src, PROMOTION, QUEEN_PROMOTION))
			case dst == enp_sq:
				move_list = append(move_list, to_move(dst, src, EN_PASSANT, NO_PROMOTION))
			default:
				move_list = append(move_list, to_move(dst, src, NORMAL, NO_PROMOTION))
			}
		}
	}
	return move_list
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
