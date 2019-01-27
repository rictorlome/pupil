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
	// PV_NODEs have an exact score.
	// Would this be accurate for tt_entry.depth >= depth?
	// It causes the minimax=ab test to fail, but the selected move must be better.
	if ok && tt_entry.key == p.state.key && tt_entry.node_type == PV_NODE && tt_entry.depth == depth {
		return tt_entry.score
	}

	// Default the new_entry to ALL_NODE
	new_entry := TTEntry{depth: depth, key: p.state.key, node_type: ALL_NODE}
	score := 0
	moves := p.generate_moves()

	// Leaf node
	if depth == 0 {
		score = p.evaluate(len(moves) == 0 && p.in_check())
		if !ok || depth >= tt_entry.depth {
			new_entry.score = score
			new_entry.node_type = PV_NODE
			TT_GLOBAL.write(p.state.key, &new_entry)
		}
		return score
	}

	// Check if best move was cached for this position
	best := Move(0)
	if ok && tt_entry.key == p.state.key && tt_entry.best_move != best {
		best = tt_entry.best_move
	}
	// Order first 3rd of the moves
	p.order_moves(&moves, best, len(moves)/3)

	// Main loop
	for _, move := range moves {
		p.do_move(move, &StateInfo{})
		score = -p.ab(-beta, -alpha, depth-1)
		p.undo_move(move)
		if score >= beta {
			if !ok || depth >= tt_entry.depth {
				new_entry.score = score
				new_entry.node_type = CUT_NODE
				new_entry.best_move = move
				TT_GLOBAL.write(p.state.key, &new_entry)
			}
			return beta
		}
		if score > alpha {
			alpha = score
			new_entry.node_type = PV_NODE
			new_entry.best_move = move
		}
	}

	// Cache node
	if !ok || depth > tt_entry.depth {
		new_entry.score = alpha
		TT_GLOBAL.write(p.state.key, &new_entry)
	}
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

// Sort first K moves in descending order based on p.value()
func (p *Position) order_moves(moves_ptr *[]Move, best Move, k int) {
	moves := *moves_ptr
	for j := 0; j < k; j++ {
		max_move_idx := j
		for i := j + 1; i < len(moves); i++ {
			if p.value(moves[i], best) > p.value(moves[max_move_idx], best) {
				max_move_idx = i
			}
		}
		if max_move_idx != j {
			moves[j], moves[max_move_idx] = moves[max_move_idx], moves[j]
		}
	}
}

// For move ordering:
// Inspired by: https://www.redhotpawn.com/rival/programming/moveorder.php
func (p *Position) value(m Move, best Move) int {
	// Best move.
	if m == best {
		return 100000
	}
	val := int(move_type(m))
	src, dst := move_src(m), move_dst(m)
	mover := p.piece_at(src)
	// Captured value - capturing value
	if is_capture(m) && !is_enpassant(m) {
		cap_val := MATERIAL_VALUES[piece_to_type(p.piece_at(dst))] * 10
		mover_val := MATERIAL_VALUES[piece_to_type(mover)] * 10
		if piece_to_type(mover) == PAWN {
			mover_val = 100
		}
		val += (cap_val - mover_val)
	}
	return val
}
