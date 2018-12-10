package main

import (
// "fmt"
)

// NOTE: tmp function does not include attacking.
func can_castle(side int, color Color, occ Bitboard, cr int) bool {
	// safe := atks&CASTLE_CHECK_SQS[int(color)*2+side] == 0
	empty := occ&CASTLE_MOVE_SQS[int(color)*2+side] == 0
	has_right := cr&CASTLING_RIGHTS[int(color)*2+side] != 0
	return empty && has_right
}

// NOTE: This should only be called on dst squares of DOUBLE_PAWN_PUSH moves
func dst_to_ep_sq(dst Square) Square {
	if dst < SQ_A5 {
		return Square(dst - 8)
	}
	return Square(dst + 8)
}

func has_right(rights int, right uint) bool {
	return (rights>>right)&1 == 1
}

func init_castling_masks() {
	for _, sq := range SQUARES {
		switch sq {
		case SQ_A1:
			CASTLING_MASK_BY_SQ[sq] = WQ_CASTLE
		case SQ_E1:
			CASTLING_MASK_BY_SQ[sq] = WHITE_CASLTES
		case SQ_H1:
			CASTLING_MASK_BY_SQ[sq] = WK_CASTLE
		case SQ_A8:
			CASTLING_MASK_BY_SQ[sq] = BQ_CASTLE
		case SQ_E8:
			CASTLING_MASK_BY_SQ[sq] = BLACK_CASTLES
		case SQ_H8:
			CASTLING_MASK_BY_SQ[sq] = BK_CASTLE
		default:
			CASTLING_MASK_BY_SQ[sq] = NO_CASTLE
		}
	}
}

func make_castle_state_info(available string) int {
	var castles int
	for _, char := range available {
		castles |= CHAR_TO_CASTLE[string(char)]
	}
	return castles
}

func update_castling_right(cr int, src Square) int {
	return cr &^ CASTLING_MASK_BY_SQ[src]
}

// NOTE: ep_sq is only set if enemy pawn can actually take it.
func update_ep_sq(m Move, enemy_pawns Bitboard) Square {
	if move_type(m) == DOUBLE_PAWN_PUSH && (enemy_pawns&NEIGHBOR_BBS[move_dst(m)] != 0) {
		return dst_to_ep_sq(move_dst(m))
	}
	return NULL_SQ
}

func update_rule_50(rule_50 int, m Move, pt PieceType) int {
	if is_capture(m) || pt == PAWN {
		return 0
	}
	return rule_50 + 1
}
