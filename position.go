package main

func (p *Position) is_capture(m Move) bool {
	return p.occupied_at(move_dst(m)) || move_type(m) == EN_PASSANT
}

func (p *Position) occupied_at(sq Square) bool {
	return p.piece_at(sq) != NULL_PIECE
}

func (p *Position) piece_at(sq Square) Piece {
	for piece, piece_bb := range p.placement {
		if occupied_at_sq(piece_bb, sq) {
			return Piece(piece)
		}
	}
	return NULL_PIECE
}
