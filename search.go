package main

import (
// "fmt"
// "strings"
)

type MoveScore struct {
	move  Move
	score int
}

func (p *Position) evaluate(checkmate bool) int {
	if !checkmate {
		return 0
	}
	return max_score(opposite(p.side_to_move()))
}

func join(c Color, a MoveScore, b MoveScore) MoveScore {
	if c == WHITE {
		return a.max(b)
	}
	return a.min(b)
}

func (m MoveScore) max(n MoveScore) MoveScore {
	if m.score > n.score {
		return m
	}
	return n
}

func max_score(c Color) int {
	if c == WHITE {
		return 32000
	}
	return -32000
}

func (m MoveScore) min(n MoveScore) MoveScore {
	if m.score < n.score {
		return m
	}
	return n
}

func (p *Position) minimax(depth int) MoveScore {
	us, them := p.side_to_move(), opposite(p.side_to_move())
	moves := p.generate_moves()
	if depth == 0 {
		return MoveScore{Move(0), p.evaluate(len(moves) == 0 && p.in_check())}
	}
	best := MoveScore{Move(0), max_score(them)}
	for _, move := range moves {
		p.do_move(move, &StateInfo{})
		recursive := p.minimax(depth - 1)
		recursive.move = move
		best = join(us, best, recursive)
		p.undo_move(move)
	}
	return best
}
