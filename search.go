package main

import (
// "fmt"
// "strings"
)

var MAX_SCORE int = 32000

type MoveScore struct {
	move  Move
	score int
}

func (p *Position) negamax(depth int) int {
	moves := p.generate_moves()
	if depth == 0 {
		return p.evaluate(len(moves) == 0 && p.in_check())
	}
	best := -MAX_SCORE
	for _, move := range moves {
		p.do_move(move, &StateInfo{})
		score := -p.negamax(depth - 1)
		p.undo_move(move)
		if score > best {
			best = score
		}
	}
	return best
}

func (p *Position) negamaxRoot(depth int) MoveScore {
	best := MoveScore{Move(0), -MAX_SCORE}
	for _, move := range p.generate_moves() {
		p.do_move(move, &StateInfo{})
		score := -p.negamax(depth - 1)
		p.undo_move(move)
		if score > best.score {
			best = MoveScore{move, score}
		}
	}
	return best
}

func (p *Position) alphaBeta(alpha int, beta int, depth int) int {
	moves := p.generate_moves()
	if depth == 0 {
		return p.evaluate(len(moves) == 0 && p.in_check())
	}
	for _, move := range moves {
		p.do_move(move, &StateInfo{})
		score := -p.alphaBeta(-beta, -alpha, depth-1)
		p.undo_move(move)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}
	return alpha
}

func (p *Position) alphaBetaRoot(depth int) MoveScore {
	alpha, beta, best_move := -MAX_SCORE, MAX_SCORE, Move(0)
	for _, move := range p.generate_moves() {
		p.do_move(move, &StateInfo{})
		score := -p.alphaBeta(-beta, -alpha, depth-1)
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
