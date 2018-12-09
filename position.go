package main

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

func (p *Position) occupied_at(sq Square) bool {
	return p.piece_at(sq) != NULL_PIECE
}

func (p *Position) piece_at(sq Square) Piece {
	return piece_on_sq(p.placement, sq)
}

func (p *Position) piece_type_at(sq Square) PieceType {
	return piece_to_type(p.piece_at(sq))
}
