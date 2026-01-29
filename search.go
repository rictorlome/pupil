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

// Maximum depth for quiescence search
const MAX_QUIESCE_DEPTH = 10

// Killer moves: 2 slots per ply, up to 64 ply deep
const MAX_PLY = 64
const KILLERS_PER_PLY = 2

var killerMoves [MAX_PLY][KILLERS_PER_PLY]Move

func clearKillers() {
	for i := range killerMoves {
		for j := range killerMoves[i] {
			killerMoves[i][j] = 0
		}
	}
}

func updateKillers(ply int, move Move) {
	if ply >= MAX_PLY {
		return
	}
	// Don't store if already in slot 0
	if killerMoves[ply][0] == move {
		return
	}
	// Shift slot 0 to slot 1, put new move in slot 0
	killerMoves[ply][1] = killerMoves[ply][0]
	killerMoves[ply][0] = move
}

func isKiller(ply int, move Move) bool {
	if ply >= MAX_PLY {
		return false
	}
	return killerMoves[ply][0] == move || killerMoves[ply][1] == move
}

// Quiescence search: continue searching captures to avoid horizon effect
func (p *Position) quiesce(alpha int, beta int, depth int) int {
	// Stand pat: evaluate the current position
	standPat := p.evaluate()

	// Beta cutoff: position is already too good
	if standPat >= beta {
		return beta
	}

	// Update alpha with stand pat score
	if standPat > alpha {
		alpha = standPat
	}

	// Depth limit reached
	if depth <= 0 {
		return alpha
	}

	// Generate all moves and filter to captures
	moves := p.generateMoves()

	// Terminal position check
	if len(*moves) == 0 {
		putMoveList(moves)
		if p.inCheck() {
			return -MAX_SCORE
		}
		return 0
	}

	// Search only captures
	for _, move := range *moves {
		if !isCapture(move) {
			continue
		}

		s := siPool.Get().(*StateInfo)
		p.doMove(move, s)
		score := -p.quiesce(-beta, -alpha, depth-1)
		p.undoMove(move)
		siPool.Put(s)

		if score >= beta {
			putMoveList(moves)
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}

	putMoveList(moves)
	return alpha
}

func (p *Position) ab(alpha int, beta int, depth uint8, ply int) int {
	ttEntry, empty := TT_GLOBAL.read(p.state.key)
	hit := !empty && ttEntry.key == p.state.key

	// TT cutoffs: use cached bounds when depth is sufficient
	if hit && ttEntry.depth >= depth {
		switch ttEntry.nodeType {
		case PV_NODE:
			// Exact score - can use directly
			return ttEntry.score
		case CUT_NODE:
			// Lower bound (score >= beta when stored)
			// If stored score >= current beta, we can cut off
			if ttEntry.score >= beta {
				return ttEntry.score
			}
			// Can also use as lower bound for alpha
			if ttEntry.score > alpha {
				alpha = ttEntry.score
			}
		case ALL_NODE:
			// Upper bound (score <= alpha when stored)
			// If stored score <= current alpha, we can cut off
			if ttEntry.score <= alpha {
				return ttEntry.score
			}
			// Can also use as upper bound for beta
			if ttEntry.score < beta {
				beta = ttEntry.score
			}
		}
	}

	// Default the newEntry to ALL_NODE
	newEntry := TTEntry{depth: depth, nodeType: ALL_NODE, key: p.state.key}
	score := 0
	moves := p.generateMoves()

	// Terminal position: no legal moves
	if len(*moves) == 0 {
		putMoveList(moves)
		if p.inCheck() {
			return -MAX_SCORE // Checkmate
		}
		return 0 // Stalemate
	}

	// Leaf node: drop into quiescence search
	if depth == 0 {
		putMoveList(moves)
		return p.quiesce(alpha, beta, MAX_QUIESCE_DEPTH)
	}

	// Check if best move was cached for this position
	best := Move(0)
	if hit && ttEntry.bestMove != best {
		best = ttEntry.bestMove
	}
	// Order ALL moves (not just first 1/3)
	p.orderMoves(moves, best, ply, len(*moves))

	// Main loop
	for _, move := range *moves {
		s := siPool.Get().(*StateInfo)
		p.doMove(move, s)
		score = -p.ab(-beta, -alpha, depth-1, ply+1)
		p.undoMove(move)
		siPool.Put(s)
		if score >= beta {
			// Update killer moves for non-captures
			if !isCapture(move) {
				updateKillers(ply, move)
			}
			if empty || p.state.key&1 == 1 || depth >= ttEntry.depth {
				newEntry.score = score
				newEntry.nodeType = CUT_NODE
				newEntry.bestMove = move
				TT_GLOBAL.write(p.state.key, &newEntry)
			}
			putMoveList(moves)
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
	putMoveList(moves)
	return alpha
}

func (p *Position) abRoot(depth uint8) MoveScore {
	clearKillers()
	alpha, beta, bestMove := -MAX_SCORE, MAX_SCORE, Move(0)
	moves := p.generateMoves()
	// Order moves at root too
	p.orderMoves(moves, Move(0), 0, len(*moves))
	for _, move := range *moves {
		p.doMove(move, &StateInfo{})
		score := -p.ab(-beta, -alpha, depth-1, 1)
		p.undoMove(move)
		if score >= beta {
			putMoveList(moves)
			return MoveScore{move, beta}
		}
		if score > alpha {
			bestMove, alpha = move, score
		}
	}
	putMoveList(moves)
	return MoveScore{bestMove, alpha}
}

// Sort first K moves in descending order based on p.value()
func (p *Position) orderMoves(movesPtr *[]Move, best Move, ply int, k int) {
	moves := *movesPtr
	if k > len(moves) {
		k = len(moves)
	}
	for j := 0; j < k; j++ {
		maxMoveIdx, maxMoveVal := j, p.value(moves[j], best, ply)
		for i := j + 1; i < len(moves); i++ {
			val := p.value(moves[i], best, ply)
			if val > maxMoveVal {
				maxMoveIdx, maxMoveVal = i, val
			}
		}
		if maxMoveIdx != j {
			moves[j], moves[maxMoveIdx] = moves[maxMoveIdx], moves[j]
		}
	}
}

// Move ordering values:
// 1. TT best move: 100000
// 2. Captures (MVV-LVA): 10000 + (victim - attacker)
// 3. Killer moves: 9000
// 4. Quiet moves: 0
func (p *Position) value(m Move, best Move, ply int) int {
	// TT best move - highest priority
	if m == best {
		return 100000
	}

	// Captures - use MVV-LVA (Most Valuable Victim - Least Valuable Attacker)
	if isCapture(m) {
		src, dst := moveSrc(m), moveDst(m)
		mover := p.pieceAt(src)
		moverVal := MATERIAL_VALUES[pieceToType(mover)]
		capVal := P_VAL // Default for en passant
		if !isEnpassant(m) {
			capVal = MATERIAL_VALUES[pieceToType(p.pieceAt(dst))]
		}
		// MVV-LVA: prioritize capturing valuable pieces with less valuable pieces
		return 10000 + capVal*10 - moverVal
	}

	// Killer moves - quiet moves that caused beta cutoffs
	if isKiller(ply, m) {
		return 9000
	}

	// Other quiet moves
	return 0
}
