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

func king_castles(move_list *[]Move, occ Bitboard, color Color, cr int, enemy_attacks Bitboard) {
	for _, side := range SIDES {
		if can_castle(side, color, occ, cr, enemy_attacks) {
			*move_list = append(*move_list, CASTLE_MOVES[int(color)*2+side])
		}
	}
}

func pawn_pushes(sq Square, dir int, occ Bitboard) Bitboard {
	pushes := signed_shift(SQUARE_BBS[sq], dir) &^ occ
	if square_rank(sq) == second_rank(dir) {
		pushes |= signed_shift(pushes, dir) &^ occ
	}
	return pushes
}

func serialize_for_pseudos_king(pl *[]Move, piece_bb Bitboard, occ Bitboard, self_occ Bitboard, color Color, cr int, enemy_attacks Bitboard) {
	src := Square(lsb(piece_bb))
	serialize_normal_moves(pl, src, king_attacks(occ, src)&^self_occ, occ)
	king_castles(pl, occ, color, cr, enemy_attacks)
}

func serialize_for_pseudos_other(pl *[]Move, piece_bb Bitboard, occ Bitboard, self_occ Bitboard, fn AttackFunc) {
	for cursor := piece_bb; cursor != 0; cursor &= cursor - 1 {
		src := Square(lsb(cursor))
		serialize_normal_moves(pl, src, fn(occ, src)&^self_occ, occ)
	}
}

// NOTE: if promoting, 2 moves are added (queen and knight promotions)
func serialize_for_pseudos_pawns(pl *[]Move, pawns Bitboard, occ Bitboard, self_occ Bitboard, color Color, ep_sq Square) {
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
				c_o_q := cap_or_quiet(occ, dst)
				*pl = append(*pl, to_move(dst, src, KNIGHT_PROMOTION|c_o_q))
				*pl = append(*pl, to_move(dst, src, QUEEN_PROMOTION|c_o_q))
				*pl = append(*pl, to_move(dst, src, ROOK_PROMOTION|c_o_q))
				*pl = append(*pl, to_move(dst, src, BISHOP_PROMOTION|c_o_q))
			case dst == ep_sq:
				*pl = append(*pl, to_move(dst, src, EP_CAPTURE))
			case dst == two_up(src, color):
				*pl = append(*pl, to_move(dst, src, DOUBLE_PAWN_PUSH))
			default:
				*pl = append(*pl, to_move(dst, src, cap_or_quiet(occ, dst)))
			}
		}
	}
}

func serialize_normal_moves(ml *[]Move, src Square, moves Bitboard, occ Bitboard) {
	for dst_cursor := moves; dst_cursor != 0; dst_cursor &= dst_cursor - 1 {
		dst := Square(lsb(dst_cursor))
		*ml = append(*ml, to_move(dst, src, cap_or_quiet(occ, dst)))
	}
}
