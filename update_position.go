package main

import (
	"fmt"
)

func (p *Position) clear_sq(sq Square) {
	for _, pc := range PIECES {
		p.remove_piece(pc, sq)
	}
}

func (p *Position) move_piece(m Move) {
	from, to := move_src(m), move_dst(m)
	moving := p.piece_at(from)

	if from == to {
		panic(fmt.Sprintf("Move: %v. ERR: src == dst", m))
	}
	if moving == NULL_PIECE {
		panic(fmt.Sprintf("Move: %v. ERR: Moving null piece", m))
	}

	p.remove_piece(moving, from)
	p.place_piece(moving, to)
}

func (p *Position) place_piece(pc Piece, sq Square) {
	p.placement[pc] |= SQUARE_BBS[sq]
}

func (p *Position) remove_piece(pc Piece, sq Square) {
	p.placement[pc] &^= SQUARE_BBS[sq]
}
