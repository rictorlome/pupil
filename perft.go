package main

import (
	"fmt"
	"sync"
)

var pool *sync.Pool

func initPool() {
	pool = &sync.Pool{
		New: func() interface{} {
			return new(StateInfo)
		},
	}
}

type thread_local struct {
	move_list, pseudo_legal_move_list *[]Move
}

func (l *thread_local) init() {
	ml, pl := make([]Move, 0, MAX_BRANCHING), make([]Move, 0, MAX_BRANCHING)
	l.move_list = &ml
	l.pseudo_legal_move_list = &pl
}

func (l *thread_local) clear() {
	*l.move_list = (*l.move_list)[:0]
	*l.pseudo_legal_move_list = (*l.pseudo_legal_move_list)[:0]
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

func get_perft(p *Position, depth int, move Move, c chan perft) {
	tl := thread_local{}
	tl.init()
	c <- get_perft_recursive(&tl, p, depth, move)
}

func get_perft_parallel(p *Position, depth int) perft {
	if depth < 1 {
		return perft{0, 1, 0, 0, 0, 0, 0, 0}
	}
	initPool()
	new_perft := perft{0, 0, 0, 0, 0, 0, 0, 0}
	c := make(chan perft)
	moves := p.generate_moves()
	for i := 0; i < len(moves); i++ {
		duped := p.dup()
		duped.do_move(moves[i], &StateInfo{})
		go get_perft(&duped, depth-1, moves[i], c)
	}
	for i := 0; i < len(moves); i++ {
		new_perft = new_perft.add(<-c)
	}
	new_perft.depth = depth
	return new_perft
}

func get_perft_recursive(tl *thread_local, p *Position, depth int, move Move) perft {
	new_perft := perft{0, 1, 0, 0, 0, 0, 0, 0}

	if depth == 0 {
		new_perft.update_with_move(move)
		new_perft.checks += indicator(p.in_check())
		new_perft.checkmates += indicator(p.in_checkmate_local(tl))
		return new_perft
	}
	new_perft.nodes = 0
	p.generate_moves_local(tl)
	for _, move := range *tl.move_list {
		s := pool.Get().(*StateInfo)
		p.do_move(move, s)
		new_perft = new_perft.add(get_perft_recursive(tl, p, depth-1, move))
		p.undo_move(move)
		pool.Put(s)
	}
	new_perft.depth = depth
	return new_perft
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
	return fmt.Sprintf("depth %v: %v nodes, %v captures, %v enpassants, %v castles, %v promotions, %v checks, and %v checkmates", p.depth, p.nodes, p.captures, p.enpassants, p.castles, p.promotions, p.checks, p.checkmates)
}
