package main

import (
	"fmt"
)

func countLeaves(p Position, depth int) int {
	if depth == 0 {
		return 1
	}
	child_leaves := 0
	for _, move := range p.generate_moves() {
		duped := p.dup()
		duped.do_move(move, StateInfo{})
		child_leaves += countLeaves(duped, depth-1)
	}
	return child_leaves
}

func get_perft(p Position, depth int) perft {
	new_perft := perft{0, 1, 0, 0, 0, 0, 0, 0}
	if depth == 0 {
		return new_perft
	}
	new_perft.nodes = 0
	for _, move := range p.generate_moves() {
		duped := p.dup()
		duped.do_move(move, StateInfo{})
		new_perft = new_perft.add(get_perft(duped, depth-1))
		new_perft.update_with_move(move)
	}
	new_perft.depth = depth
	return new_perft
}

type perft struct {
	depth, nodes, captures, enpassants, castles, promotions, checks, checkmates int
}

func (p *perft) add(s perft) perft {
	return perft{
		s.depth, p.nodes + s.nodes, p.captures + s.captures,
		p.enpassants + s.enpassants, p.castles + s.castles,
		p.promotions + s.promotions, p.checks + s.checks,
		p.checkmates + s.checkmates,
	}
}

func indicator(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (p *perft) update_with_move(m Move) {
	p.captures += indicator(is_capture(m))
	p.enpassants += indicator(is_enpassant(m))
	p.castles += indicator(is_castle(m))
	p.promotions += indicator(is_promotion(m))
}

func (p perft) String() string {
	return fmt.Sprintf("At depth %v,\n%v nodes, %v captures, %v enpassants, %v castles, %v promotions, %v checks, and %v checkmates", p.depth, p.nodes, p.captures, p.enpassants, p.castles, p.promotions, p.checks, p.checkmates)
}
