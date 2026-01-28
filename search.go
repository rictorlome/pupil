package main

import (
// "fmt"
// "strings"
)

var MAX_SCORE int = 32000
var TT_GLOBAL *TT = createTT(20) // approx 120 mb space for entire program

type MoveScore struct {
	move  Move
	score int
}

func (p *Position) ab(alpha int, beta int, depth uint8) int {
	ttEntry, empty := TT_GLOBAL.read(p.state.key)
	hit := !empty && ttEntry.key == p.state.key
	// PV_NODEs have an exact score.
	// Would this be accurate for ttEntry.depth >= depth?
	// It causes the minimax=ab test to fail, but the selected move must be better.
	if hit && ttEntry.nodeType == PV_NODE && ttEntry.depth == depth {
		return ttEntry.score
	}

	// Default the newEntry to ALL_NODE
	newEntry := TTEntry{depth: depth, nodeType: ALL_NODE, key: p.state.key}
	score := 0
	moves := p.generateMoves()

	// Terminal position: no legal moves
	if len(moves) == 0 {
		if p.inCheck() {
			return -MAX_SCORE // Checkmate
		}
		return 0 // Stalemate
	}

	// Leaf node
	if depth == 0 {
		score = p.evaluate()
		if empty || p.state.key&1 == 1 || depth >= ttEntry.depth {
			newEntry.score = score
			newEntry.nodeType = PV_NODE
			TT_GLOBAL.write(p.state.key, &newEntry)
		}
		return score
	}

	// Check if best move was cached for this position
	best := Move(0)
	if hit && ttEntry.bestMove != best {
		best = ttEntry.bestMove
	}
	// Order first 3rd of the moves
	p.orderMoves(&moves, best, len(moves)/3)

	// Main loop
	for _, move := range moves {
		s := siPool.Get().(*StateInfo)
		p.doMove(move, s)
		score = -p.ab(-beta, -alpha, depth-1)
		p.undoMove(move)
		siPool.Put(s)
		if score >= beta {
			if empty || p.state.key&1 == 1 || depth >= ttEntry.depth {
				newEntry.score = score
				newEntry.nodeType = CUT_NODE
				newEntry.bestMove = move
				TT_GLOBAL.write(p.state.key, &newEntry)
			}
			return beta
		}
		if score > alpha {
			alpha = score
			newEntry.nodeType = PV_NODE
			newEntry.bestMove = move
		}
	}

	// Cache node
	if empty || p.state.key&1 == 1 || depth > ttEntry.depth {
		newEntry.score = alpha
		TT_GLOBAL.write(p.state.key, &newEntry)
	}
	return alpha
}

func (p *Position) abRoot(depth uint8) MoveScore {
	alpha, beta, bestMove := -MAX_SCORE, MAX_SCORE, Move(0)
	for _, move := range p.generateMoves() {
		p.doMove(move, &StateInfo{})
		score := -p.ab(-beta, -alpha, depth-1)
		p.undoMove(move)
		if score >= beta {
			return MoveScore{move, beta}
		}
		if score > alpha {
			bestMove, alpha = move, score
		}
	}
	return MoveScore{bestMove, alpha}
}

// Sort first K moves in descending order based on p.value()
func (p *Position) orderMoves(movesPtr *[]Move, best Move, k int) {
	moves := *movesPtr
	for j := 0; j < k; j++ {
		maxMoveIdx, maxMoveVal := j, p.value(moves[j], best)
		for i := j + 1; i < len(moves); i++ {
			val := p.value(moves[i], best)
			if val > maxMoveVal {
				maxMoveIdx, maxMoveVal = i, val
			}
		}
		if maxMoveIdx != j {
			moves[j], moves[maxMoveIdx] = moves[maxMoveIdx], moves[j]
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
	val := int(moveType(m))
	src, dst := moveSrc(m), moveDst(m)
	mover := p.pieceAt(src)
	// Captured value - capturing value
	if isCapture(m) && !isEnpassant(m) {
		capVal := MATERIAL_VALUES[pieceToType(p.pieceAt(dst))] * 10
		moverVal := MATERIAL_VALUES[pieceToType(mover)] * 10
		if pieceToType(mover) == PAWN {
			moverVal = 100
		}
		val += (capVal - moverVal)
	}
	return val
}
