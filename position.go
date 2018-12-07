package main

func (p *Position) is_capture(m Move) bool {
	return p.occupied_at(move_dst(m)) || move_type(m) == EN_PASSANT
}

func (p *Position) occupied_at(sq Square) bool {
	return p.piece_at(sq) != NULL_PIECE
}

func (p *Position) piece_at(sq Square) Piece {
  return piece_on_sq(p.placement, sq)
}
