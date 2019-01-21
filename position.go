package main

import (
// "fmt"
)

func (p *Position) dup() Position {
	new_placement := make([]Bitboard, len(p.placement))
	new_placement_by_square := make([]Piece, len(p.placement_by_square))
	copy(new_placement, p.placement)
	copy(new_placement_by_square, p.placement_by_square)
	return Position{
		p.occ,
		new_placement,
		new_placement_by_square,
		p.ply,
		p.state.dup(),
		p.stm,
	}
}

func (p *Position) generate_evasions(pl *[]Move, ml *[]Move) {
	us, them := p.side_to_move(), opposite(p.side_to_move())
	king_sq := p.king_square(us)

	occ, self_occ := p.occupancy(), p.occupancy_by_color(us)
	occ_without_king := occ &^ SQUARE_BBS[king_sq]

	checkers := attackers_to_sq_by_color(p.placement, king_sq, them)
	atks := attacks_by_color(occ_without_king, p.placement, them)

	// King evasions
	serialize_normal_moves(ml, king_sq, king_attacks(occ, king_sq)&^(self_occ|atks), occ)
	if popcount(checkers) > 1 {
		return
	}
	checker_sq := Square(lsb(checkers))
	// Regular moves (duplicate king evasions are excluded in pseudo_legals_by_color)
	p.generate_non_evasions(pl, ml, BETWEEN_BBS[king_sq][checker_sq]|checkers)
}

func (p *Position) generate_moves() []Move {
	pseudo_legal_move_list := make([]Move, 0, MAX_BRANCHING)
	move_list := make([]Move, 0, MAX_BRANCHING)
	if p.in_check() {
		p.generate_evasions(&pseudo_legal_move_list, &move_list)
	} else {
		p.generate_non_evasions(&pseudo_legal_move_list, &move_list, Bitboard(0))
	}
	return move_list
}

// NOTE: pseudolegal moves include those that cause check. these have to be filtered out in move generation
func (p *Position) generate_pseudo_legals(pl *[]Move, forced_dsts Bitboard) {
	us := p.side_to_move()
	occ, self_occ := p.occupancy(), p.occupancy_by_color(us)

	serialize_for_pseudos_pawns(pl, p.our_pt_bb(PAWN), occ, self_occ, us, p.state.ep_sq)
	serialize_for_pseudos_other(pl, p.our_pt_bb(KNIGHT), occ, self_occ, knight_attacks)
	serialize_for_pseudos_other(pl, p.our_pt_bb(ROOK), occ, self_occ, rook_attacks)
	serialize_for_pseudos_other(pl, p.our_pt_bb(BISHOP), occ, self_occ, bishop_attacks)
	serialize_for_pseudos_other(pl, p.our_pt_bb(QUEEN), occ, self_occ, queen_attacks)

	if empty(forced_dsts) {
		serialize_for_pseudos_king(pl, p.our_pt_bb(KING), occ, self_occ, us, p.state.castling_rights, p.state.opposite_color_attacks)
	}
}

func (p *Position) generate_non_evasions(pl *[]Move, ml *[]Move, forced_dsts Bitboard) {
	p.generate_pseudo_legals(pl, forced_dsts)
	for _, pseudo_legal := range *pl {
		if p.is_legal(pseudo_legal) && (empty(forced_dsts) || is_good_evasion(forced_dsts, pseudo_legal)) {
			*ml = append(*ml, pseudo_legal)
		}
	}
}

func is_good_evasion(forced_dsts Bitboard, m Move) bool {
	return occupied_at_sq(forced_dsts, move_dst(m)) ||
		(is_enpassant(m) && occupied_at_sq(forced_dsts, cleanup_sq_for_ep_capture(move_dst(m))))
}

func (p *Position) get_color_attacks(color Color) Bitboard {
	return attacks_by_color(p.occupancy(), p.placement, color)
}

func (p *Position) king_square(color Color) Square {
	return Square(lsb(p.placement[color*6]))
}

func (p *Position) in_check() bool {
	color := p.side_to_move()
	atks := p.opposite_color_attacks()
	return occupied_at_sq(atks, p.king_square(color))
}

func (p *Position) in_checkmate() bool {
	return p.in_check() && len(p.generate_moves()) == 0
}

func (p *Position) is_legal(m Move) bool {
	src, dst := move_src(m), move_dst(m)
	us, them := p.side_to_move(), opposite(p.side_to_move())
	ksq := p.king_square(us)

	if is_enpassant(m) {
		their_queens := p.placement[pt_to_p(QUEEN, them)]
		// No discovered slider attacks on the king.
		return empty(rook_attacks(p.occupancy()&^(SQUARE_BBS[src]|SQUARE_BBS[cleanup_sq_for_ep_capture(dst)]), ksq)&(p.placement[pt_to_p(ROOK, them)]|their_queens)) &&
			empty(bishop_attacks(p.occupancy()&^SQUARE_BBS[src], ksq)&(p.placement[pt_to_p(BISHOP, them)]|their_queens))
	}
	if p.piece_type_at(src) == KING {
		return is_castle(m) || !occupied_at_sq(p.opposite_color_attacks(), dst)
	}
	return !occupied_at_sq(p.state.blockers_for_king, src) || aligned(src, dst, ksq)
}

func (p *Position) move_count() int {
	return p.ply / 2
}

func (p *Position) occupancy() Bitboard {
	return p.occ
}

func (p *Position) occupancy_by_color(c Color) Bitboard {
	return occupied_squares_by_color(p.placement, c)
}

func (p *Position) occupancy_by_piece(pc Piece) Bitboard {
	return p.placement[pc]
}

func (p *Position) occupancy_by_pieces(pieces ...Piece) Bitboard {
	var occupancy Bitboard
	for _, piece := range pieces {
		occupancy |= p.occupancy_by_piece(piece)
	}
	return occupancy
}

func (p *Position) occupancy_by_piece_type(pt PieceType) Bitboard {
	return p.occupancy_by_pieces(pt_to_p(pt, WHITE), pt_to_p(pt, BLACK))
}

func (p *Position) occupancy_by_piece_types(pts ...PieceType) Bitboard {
	var occupancy Bitboard
	for _, pt := range pts {
		occupancy |= p.occupancy_by_piece_type(pt)
	}
	return occupancy
}

func (p *Position) occupied_at(sq Square) bool {
	return p.piece_at(sq) != NULL_PIECE
}

func (p *Position) opposite_color_attacks() Bitboard {
	return p.state.opposite_color_attacks
}

func (p *Position) parse_move(s string) Move {
	src := parse_square(s[:2])
	dst := parse_square(s[2:4])
	mover_type, captured := piece_to_type(p.piece_at(src)), p.piece_at(dst)
	return to_move(dst, src, parse_move_type(s[4:], p.occupancy(), src, dst, mover_type, captured))
}

func (p *Position) piece_at(sq Square) Piece {
	return p.placement_by_square[sq]
}

func (p *Position) piece_type_at(sq Square) PieceType {
	return piece_to_type(p.piece_at(sq))
}

func (p *Position) to_zobrist() Key {
	var ResKey Key
	for sq, piece := range p.placement_by_square {
		if piece != NULL_PIECE {
			ResKey ^= ZOBRIST_PSQ[sq][piece]
		}
	}
	ResKey ^= ZOBRIST_CSTL[p.state.castling_rights]
	ResKey ^= ZOBRIST_EPSQ[p.state.ep_sq]
	ResKey ^= ZOBRIST_SIDE * Key(p.stm)
	return ResKey
}

func (p *Position) side_to_move() Color {
	return p.stm
}

// Returns bitboard of all pieces blocking attacks to sq from sliders of color c.
func (p *Position) slider_blockers(c Color, sq Square) Bitboard {
	var blockers Bitboard

	occ := p.occupancy()

	// QUESTION: is the attack mask sufficient here, or do we need pseudolegal moves?
	queen_occ := p.occupancy_by_piece(pt_to_p(QUEEN, c))
	possible_snipers := ((p.occupancy_by_piece(pt_to_p(BISHOP, c)) | queen_occ) & BISHOP_ATTACK_MASKS[sq]) |
		((p.occupancy_by_piece(pt_to_p(ROOK, c)) | queen_occ) & ROOK_ATTACK_MASKS[sq])

	for cursor := possible_snipers; cursor != 0; cursor &= cursor - 1 {
		sniper_sq := Square(lsb(cursor))
		intermediate_pieces := BETWEEN_BBS[sq][sniper_sq] & occ
		// NOTE: can set pinners in this if as well, when I start keeping track of them.
		if popcount(intermediate_pieces) == 1 {
			blockers |= intermediate_pieces
		}
	}

	return blockers
}

func (p *Position) our_pt_bb(pt PieceType) Bitboard {
	return p.placement[pt_to_p(pt, p.side_to_move())]
}

func (p Position) String() string {
	return generate_fen(p)
}
