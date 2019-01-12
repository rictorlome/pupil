package main

import (
	"fmt"
)

func can_castle(side int, color Color, occ Bitboard, cr int, atks Bitboard) bool {
	safe := atks&CASTLE_CHECK_SQS[int(color)*2+side] == 0
	empty := occ&CASTLE_MOVE_SQS[int(color)*2+side] == 0
	has_right := cr&CASTLING_RIGHTS[int(color)*2+side] != 0
	return empty && has_right && safe
}

func cleanup_sq_for_ep_capture(capture_dst Square) Square {
	if capture_dst < SQ_A5 {
		return Square(capture_dst + 8)
	}
	return Square(capture_dst - 8)
}

// NOTE: This should only be called on dst squares of DOUBLE_PAWN_PUSH moves
func dst_to_ep_sq(push_dst Square) Square {
	if push_dst < SQ_A5 {
		return Square(push_dst - 8)
	}
	return Square(push_dst + 8)
}

func (s *StateInfo) dup() *StateInfo {
	return &StateInfo{
		s.castling_rights, s.ep_sq, s.rule_50, s.opposite_color_attacks,
		s.blockers_for_king, s.prev, s.captured,
	}
}

func has_right(rights int, right uint) bool {
	return (rights>>right)&1 == 1
}

func init_castling_masks() {
	CASTLING_MASK_BY_SQ[SQ_A1] = WQ_CASTLE
	CASTLING_MASK_BY_SQ[SQ_E1] = WHITE_CASTLES
	CASTLING_MASK_BY_SQ[SQ_H1] = WK_CASTLE
	CASTLING_MASK_BY_SQ[SQ_A8] = BQ_CASTLE
	CASTLING_MASK_BY_SQ[SQ_E8] = BLACK_CASTLES
	CASTLING_MASK_BY_SQ[SQ_H8] = BK_CASTLE
}

func init_rook_squares_for_castling() {
	m := [12]Square{
		SQ_C1, SQ_A1, SQ_D1,
		SQ_G1, SQ_H1, SQ_F1,
		SQ_C8, SQ_A8, SQ_D8,
		SQ_G8, SQ_H8, SQ_F8,
	}
	for king_dst, rook_src, rook_dst := 0, 1, 2; rook_dst < 12; {
		ROOK_SRC_DST[m[king_dst]] = [2]Square{m[rook_src], m[rook_dst]}
		king_dst += 3
		rook_src += 3
		rook_dst += 3
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
	dst := move_dst(m)
	if move_type(m) == DOUBLE_PAWN_PUSH && (enemy_pawns&NEIGHBOR_BBS[dst] != 0) {
		return dst_to_ep_sq(dst)
	}
	return NULL_SQ
}

func update_rule_50(rule_50 int, m Move, pt PieceType) int {
	if is_capture(m) || pt == PAWN {
		return 0
	}
	return rule_50 + 1
}

func (s StateInfo) String() string {
	recursive_string := "nil"
	if s.prev != nil {
		recursive_string = (*s.prev).String()
	}
	return fmt.Sprintf("[r50 = %v, prev = %v]", s.rule_50, recursive_string)
}
