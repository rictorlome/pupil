package main

import (
	"fmt"
	"sync"
)

// Multithreaded perft tests require a custom Transposition Table
type TTPerft struct {
	lock sync.RWMutex
	m    map[Key]perft
}

func createTTPerft() *TTPerft {
	return &TTPerft{m: make(map[Key]perft)}
}

func (t *TTPerft) read(key Key) (perft, bool) {
	t.lock.RLock()
	entry, ok := t.m[key]
	t.lock.RUnlock()
	return entry, ok
}

func (t *TTPerft) write(key Key, entry perft) {
	t.lock.Lock()
	t.m[key] = entry
	t.lock.Unlock()
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
	c <- get_perft_recursive(p, depth, move)
}

func get_perft_parallel(p *Position, depth int) perft {
	if depth < 1 {
		return perft{0, 1, 0, 0, 0, 0, 0, 0}
	}
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

var TT_TABLES = [10]*TTPerft{
	createTTPerft(), createTTPerft(), createTTPerft(), createTTPerft(), createTTPerft(),
	createTTPerft(), createTTPerft(), createTTPerft(), createTTPerft(), createTTPerft(),
}

func get_perft_recursive(p *Position, depth int, move Move) perft {
	new_perft := perft{0, 1, 0, 0, 0, 0, 0, 0}
	tt_table := TT_TABLES[depth]
	perft, ok := tt_table.read(p.state.key)
	if ok {
		return perft
	}
	if depth == 0 {
		new_perft.update_with_move(move)
		new_perft.checks += indicator(p.in_check())
		new_perft.checkmates += indicator(p.in_checkmate())
		return new_perft
	}
	new_perft.nodes = 0
	for _, move := range p.generate_moves() {
		s := si_pool.Get().(*StateInfo)
		p.do_move(move, s)
		new_perft = new_perft.add(get_perft_recursive(p, depth-1, move))
		p.undo_move(move)
		si_pool.Put(s)
	}
	new_perft.depth = depth
	tt_table.write(p.state.key, new_perft)
	return new_perft
}

func (p perft) String() string {
	return fmt.Sprintf("depth %v: %v nodes, %v captures, %v enpassants, %v castles, %v promotions, %v checks, and %v checkmates", p.depth, p.nodes, p.captures, p.enpassants, p.castles, p.promotions, p.checks, p.checkmates)
}

func (p *perft) update_with_move(m Move) {
	p.captures += indicator(is_capture(m))
	p.enpassants += indicator(is_enpassant(m))
	p.castles += indicator(is_castle(m))
	p.promotions += indicator(is_promotion(m))
}
