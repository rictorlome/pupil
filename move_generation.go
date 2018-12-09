package main

import (
// "fmt"
)

// NOTE: CASTLE_MOVE_SQS is for testing non-occupancy of squares for castling.
// CASTLE_CHECK_SQS is used for checking if king is passing through attacked sqs.
// Include E1/E8 and exclude B1/B8 when validating if squares are attacked
func init_castle_sqs() {
	CASTLE_MOVE_SQS[0] = SQUARE_BBS[SQ_F1] | SQUARE_BBS[SQ_G1]
	CASTLE_MOVE_SQS[1] = SQUARE_BBS[SQ_D1] | SQUARE_BBS[SQ_C1] | SQUARE_BBS[SQ_B1]
	CASTLE_MOVE_SQS[2] = SQUARE_BBS[SQ_F8] | SQUARE_BBS[SQ_G8]
	CASTLE_MOVE_SQS[3] = SQUARE_BBS[SQ_D8] | SQUARE_BBS[SQ_C8] | SQUARE_BBS[SQ_B8]

	CASTLE_CHECK_SQS[0] = SQUARE_BBS[SQ_E1] | SQUARE_BBS[SQ_F1] | SQUARE_BBS[SQ_G1]
	CASTLE_CHECK_SQS[1] = SQUARE_BBS[SQ_E1] | SQUARE_BBS[SQ_D1] | SQUARE_BBS[SQ_C1]
	CASTLE_CHECK_SQS[2] = SQUARE_BBS[SQ_E8] | SQUARE_BBS[SQ_F8] | SQUARE_BBS[SQ_G8]
	CASTLE_CHECK_SQS[3] = SQUARE_BBS[SQ_E8] | SQUARE_BBS[SQ_D8] | SQUARE_BBS[SQ_C8]
}

func king_castles(occ Bitboard, color Color, cr int) []Move {
	var move_list []Move
	for _, side := range SIDES {
		if can_castle(side, color, occ, cr) {
			move_list = append(move_list, CASTLE_MOVES[int(color)*2+side])
		}
	}
	return move_list
}

func pawn_pushes(sq Square, dir int, occ Bitboard) Bitboard {
	pushes := signed_shift(SQUARE_BBS[sq], dir) &^ occ
	if square_rank(sq) == second_rank(dir) {
		pushes |= signed_shift(pushes, dir) &^ occ
	}
	return pushes
}

func pseudolegals_by_color(pieces []Bitboard, color Color, ep_sq Square, castling_rights int) []Move {
	occ, self_occ := occupied_squares(pieces), occupied_squares_by_color(pieces, color)
	var move_list []Move

	for _, piece := range piece_range_by_color(color) {
		t := piece_to_type(piece)
		switch t {
		case PAWN:
			move_list = append(move_list, serialize_for_pseudos_pawns(pieces[piece], occ, self_occ, color, ep_sq)...)
		case KING:
			move_list = append(move_list, serialize_for_pseudos_king(pieces[piece], occ, self_occ, color, castling_rights)...)
		default:
			move_list = append(move_list, serialize_for_pseudos_other(pieces[piece], occ, self_occ, get_attack_func(t))...)
		}
	}
	return move_list
}

// NOTE: This function assumes only one king in piece_bb.
func serialize_for_pseudos_king(piece_bb Bitboard, occ Bitboard, self_occ Bitboard, color Color, cr int) []Move {
	src := Square(lsb(piece_bb))
	return append(serialize_normal_moves(src, king_attacks(occ, src), occ), king_castles(occ, color, cr)...)
}

func serialize_for_pseudos_other(piece_bb Bitboard, occ Bitboard, self_occ Bitboard, fn AttackFunc) []Move {
	var move_list []Move
	for cursor := piece_bb; cursor != 0; cursor &= cursor - 1 {
		src := Square(lsb(cursor))
		move_list = append(serialize_normal_moves(src, fn(occ, src)&^self_occ, occ))
	}
	return move_list
}

// NOTE: if promoting, 2 moves are added (queen and knight promotions)
func serialize_for_pseudos_pawns(pawns Bitboard, occ Bitboard, self_occ Bitboard, color Color, ep_sq Square) []Move {
	var move_list []Move
	f_dir, l_rank := forward(color), last_rank(color)
	ep_sq_bb := Bitboard(0)

	// To avoid indexing problems
	if ep_sq != NULL_SQ {
		ep_sq_bb = SQUARE_BBS[ep_sq]
	}

	for cursor := pawns; cursor != 0; cursor &= cursor - 1 {
		src := Square(lsb(cursor))
		attacks := (PAWN_ATTACK_BBS[src][color] & (occ | ep_sq_bb)) &^ self_occ
		pseudos := attacks | pawn_pushes(src, f_dir, occ)
		for dst_cursor := pseudos; dst_cursor != 0; dst_cursor &= dst_cursor - 1 {
			dst := Square(lsb(dst_cursor))
			switch {
			case square_rank(dst) == l_rank:
				move_list = append(move_list, to_move(dst, src, KNIGHT_PROMOTION|cap_or_quiet(occ, dst)), to_move(dst, src, QUEEN_PROMOTION|cap_or_quiet(occ, dst)))
			case dst == ep_sq:
				move_list = append(move_list, to_move(dst, src, EP_CAPTURE))
			default:
				// NOTE: this current encoding does not include double pushes.
				move_list = append(move_list, to_move(dst, src, cap_or_quiet(occ, dst)))
			}
		}
	}
	return move_list
}

func serialize_normal_moves(src Square, moves Bitboard, occ Bitboard) []Move {
	var move_list []Move
	for dst_cursor := moves; dst_cursor != 0; dst_cursor &= dst_cursor - 1 {
		dst := Square(lsb(dst_cursor))
		move_list = append(move_list, to_move(dst, src, cap_or_quiet(occ, dst)))
	}
	return move_list
}
