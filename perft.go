package main

import (
	"fmt"
)

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

func (p *perft) update_with_move(m Move) {
	p.captures += indicator(is_capture(m))
	p.enpassants += indicator(is_enpassant(m))
	p.castles += indicator(is_castle(m))
	p.promotions += indicator(is_promotion(m))
}

func (p perft) String() string {
	return fmt.Sprintf("At depth %v,\n%v nodes, %v captures, %v enpassants, %v castles, %v promotions, %v checks, and %v checkmates", p.depth, p.nodes, p.captures, p.enpassants, p.castles, p.promotions, p.checks, p.checkmates)
}
