package main

import (
// "fmt"
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
		p.do_move(move, StateInfo{})
		new_perft = new_perft.add(get_perft(p, depth-1))
		new_perft.update_with_move(move)
		new_perft.checks += indicator(p.in_check())
		new_perft.checkmates += indicator(p.in_checkmate())
		p.undo_move(move)
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
