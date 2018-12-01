package main

import (
// "fmt"
)

func (p *Position) place_piece(pc Piece, sq Square) {
	p.placement[pc] |= SQUARE_BBS[sq]
}

func (p *Position) remove_piece(pc Piece, sq Square) {
	p.placement[pc] ^= SQUARE_BBS[sq]
}
