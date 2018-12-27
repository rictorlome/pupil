package main

import (
// "fmt"
)

func (p *Position) king_square(color Color) Square {
	return Square(lsb(p.placement[color*6]))
}

func (p *Position) is_legal(m Move) bool {
	src, dst := move_src(m), move_dst(m)

	if is_enpassant(m) {
		// This enpassant check is temporary. Apparently, this is a tricky case.
		return true
	}
	if p.piece_type_at(src) == KING {
		// Remember to add not-through-attack check for castles
		return is_castle(m) || !occupied_at_sq(attacks_by_color(p.placement, opposite(p.to_move)), dst)
	}
	return !occupied_at_sq(p.state.blockers_for_king, src) || aligned(src, dst, p.king_square(p.to_move))
}

func (p *Position) occupancy() Bitboard {
	return occupied_squares(p.placement)
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
	return piece_on_sq(p.placement, sq)
}

func (p *Position) piece_type_at(sq Square) PieceType {
	return piece_to_type(p.piece_at(sq))
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
