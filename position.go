package main

func is_capture(m Move) bool {
	return m&Move(CAPTURE) != 0
}

func (p *Position) occupied_at(sq Square) bool {
	return p.piece_at(sq) != NULL_PIECE
}

func (p *Position) piece_at(sq Square) Piece {
	return piece_on_sq(p.placement, sq)
}
