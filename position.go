package main

import (
// "fmt"
)

func (p *Position) get_color_attacks(color Color) Bitboard {
	return attacks_by_color(p.occupancy(), p.placement, color)
}

func (p *Position) opposite_color_attacks() Bitboard {
	return p.state.opposite_color_attacks
}

func (p *Position) dup() Position {
	new_placement := make([]Bitboard, len(p.placement))
	new_placement_by_square := make([]Piece, len(p.placement_by_square))
	copy(new_placement, p.placement)
	copy(new_placement_by_square, p.placement_by_square)
	return Position{
		p.ply,
		new_placement,
		new_placement_by_square,
		p.state.dup(),
	}
}

func (p *Position) generate_evasions() []Move {
	us, them := p.side_to_move(), opposite(p.side_to_move())
	king_sq := p.king_square(us)

	occ, self_occ := p.occupancy(), p.occupancy_by_color(us)
	occ_without_king := occ &^ SQUARE_BBS[king_sq]

	checkers := attackers_to_sq_by_color(p.placement, king_sq, them)
	atks := attacks_by_color(occ_without_king, p.placement, them)

	king_evasions := serialize_normal_moves(king_sq, king_attacks(occ, king_sq)&^(self_occ|atks), occ)
	// safe squares for king
	if popcount(checkers) > 1 {
		return king_evasions
	}
	checker_sq := Square(lsb(checkers))
	non_evasions := p.generate_non_evasions()
	blocks_or_captures := king_evasions
	for _, move := range non_evasions {
		dst := move_dst(move)
		if dst == checker_sq || occupied_at_sq(BETWEEN_BBS[king_sq][checker_sq], dst) {
			blocks_or_captures = append(blocks_or_captures, move)
		}
	}
	return blocks_or_captures
}

func (p *Position) generate_moves() []Move {
	if p.in_check() {
		return p.generate_evasions()
	}
	return p.generate_non_evasions()
}

func (p *Position) generate_non_evasions() []Move {
	pseudo_legals := pseudolegals_by_color(p.placement, p.side_to_move(), p.state.ep_sq, p.state.castling_rights)
	non_evasions := make([]Move, 0)
	for _, pseudo_legal := range pseudo_legals {
		if p.is_legal(pseudo_legal) {
			non_evasions = append(non_evasions, pseudo_legal)
		}
	}
	return non_evasions
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

	if is_enpassant(m) {
		// This enpassant check is temporary. Apparently, this is a tricky case.
		return true
	}
	if p.piece_type_at(src) == KING {
		// Remember to add not-through-attack check for castles
		return is_castle(m) || !occupied_at_sq(p.opposite_color_attacks(), dst)
	}
	return !occupied_at_sq(p.state.blockers_for_king, src) || aligned(src, dst, p.king_square(p.side_to_move()))
}

func (p *Position) move_count() int {
	return p.ply / 2
}

func (p *Position) occupancy() Bitboard {
	return occupied_squares(p.placement)
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

func (p *Position) side_to_move() Color {
	return Color(p.ply % 2)
}

// Returns bitboard of all pieces blocking attacks to sq from sliders of color c.
func (p *Position) slider_blockers(c Color, sq Square) Bitboard {
	var blockers Bitboard

	occ := p.occupancy()

	// QUESTION: is the attack mask sufficient here, or do we need pseudolegal moves?
	possible_snipers := (p.occupancy_by_pieces(pt_to_p(BISHOP, c), pt_to_p(QUEEN, c)) & BISHOP_ATTACK_MASKS[sq]) |
		(p.occupancy_by_pieces(pt_to_p(ROOK, c), pt_to_p(QUEEN, c)) & ROOK_ATTACK_MASKS[sq])

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

func (p Position) String() string {
	return generate_fen(p)
}
