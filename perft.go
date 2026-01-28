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

func getPerft(p *Position, depth int, move Move, c chan perft) {
	c <- getPerftRecursive(p, depth, move)
}

func getPerftParallel(p *Position, depth int) perft {
	if depth < 1 {
		return perft{0, 1, 0, 0, 0, 0, 0, 0}
	}
	newPerft := perft{0, 0, 0, 0, 0, 0, 0, 0}
	c := make(chan perft)
	moves := p.generateMoves()
	for i := 0; i < len(moves); i++ {
		duped := p.dup()
		duped.doMove(moves[i], &StateInfo{})
		go getPerft(&duped, depth-1, moves[i], c)
	}
	for i := 0; i < len(moves); i++ {
		newPerft = newPerft.add(<-c)
	}
	newPerft.depth = depth
	return newPerft
}

var TT_TABLES = [10]*TTPerft{
	createTTPerft(), createTTPerft(), createTTPerft(), createTTPerft(), createTTPerft(),
	createTTPerft(), createTTPerft(), createTTPerft(), createTTPerft(), createTTPerft(),
}

func getPerftRecursive(p *Position, depth int, move Move) perft {
	newPerft := perft{0, 1, 0, 0, 0, 0, 0, 0}
	ttTable := TT_TABLES[depth]
	perft, ok := ttTable.read(p.state.key)
	if ok {
		return perft
	}
	if depth == 0 {
		newPerft.updateWithMove(move)
		newPerft.checks += indicator(p.inCheck())
		newPerft.checkmates += indicator(p.inCheckmate())
		return newPerft
	}
	newPerft.nodes = 0
	for _, move := range p.generateMoves() {
		s := siPool.Get().(*StateInfo)
		p.doMove(move, s)
		newPerft = newPerft.add(getPerftRecursive(p, depth-1, move))
		p.undoMove(move)
		siPool.Put(s)
	}
	newPerft.depth = depth
	ttTable.write(p.state.key, newPerft)
	return newPerft
}

func (p perft) String() string {
	return fmt.Sprintf("depth %v: %v nodes, %v captures, %v enpassants, %v castles, %v promotions, %v checks, and %v checkmates", p.depth, p.nodes, p.captures, p.enpassants, p.castles, p.promotions, p.checks, p.checkmates)
}

func (p *perft) updateWithMove(m Move) {
	p.captures += indicator(isCapture(m))
	p.enpassants += indicator(isEnpassant(m))
	p.castles += indicator(isCastle(m))
	p.promotions += indicator(isPromotion(m))
}
