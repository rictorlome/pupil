package main

import (
// "fmt"
// "strings"
)

var MAX_SCORE int = 32000
var TT_GLOBAL *TT = createTT()

type MoveScore struct {
	move  Move
	score int
}

func (p *Position) ab(alpha int, beta int, depth int) int {
	tt_entry, ok := TT_GLOBAL.read(p.state.key)
	if ok && tt_entry.depth == depth && tt_entry.key == p.state.key {
		return tt_entry.score
	}
	moves := p.generate_moves()
	if depth == 0 {
		score := p.evaluate(len(moves) == 0 && p.in_check())
		new_entry := TTEntry{score: score, depth: depth, key: p.state.key}
		TT_GLOBAL.write(p.state.key, &new_entry)
		return score
	}
	for _, move := range moves {
		p.do_move(move, &StateInfo{})
		score := -p.ab(-beta, -alpha, depth-1)
		p.undo_move(move)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}
	// // This causes the NegaEqualsAB test to fail
	// new_entry := TTEntry{score: alpha, depth: depth, key: p.state.key}
	// // This causes even more errors
	// if ok && depth > tt_entry.depth {
	// 		TT_GLOBAL.write(p.state.key, &new_entry)
	// }
	return alpha
}

func (p *Position) ab_root(depth int) MoveScore {
	alpha, beta, best_move := -MAX_SCORE, MAX_SCORE, Move(0)
	for _, move := range p.generate_moves() {
		p.do_move(move, &StateInfo{})
		score := -p.ab(-beta, -alpha, depth-1)
		p.undo_move(move)
		if score >= beta {
			return MoveScore{move, beta}
		}
		if score > alpha {
			best_move, alpha = move, score
		}
	}
	return MoveScore{best_move, alpha}
}
