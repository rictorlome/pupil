package main

func (p *Position) king_square(color Color) Square {
	return Square(lsb(p.placement[color * 6]))
}

func (p *Position) occupied_at(sq Square) bool {
	return p.piece_at(sq) != NULL_PIECE
}

func (p *Position) piece_at(sq Square) Piece {
	return piece_on_sq(p.placement, sq)
}
