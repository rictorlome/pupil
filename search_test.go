package main

import (
	"testing"
)

type BestMoveTest struct {
	fen  string
	best []string // Multiple valid solutions allowed
}

// from https://sites.google.com/site/darktemplarchess/mate-in-2-puzzles
var MateInTwos = []BestMoveTest{
	{
		"2bqkbn1/2pppp2/np2N3/r3P1p1/p2N2B1/5Q2/PPPPKPP1/RNB2r2 w KQkq - 0 1",
		[]string{"f3f7"},
	},
	{
		"8/6K1/1p1B1RB1/8/2Q5/2n1kP1N/3b4/4n3 w - - 0 1",
		[]string{"d6a3"},
	},
	{
		"B7/K1B1p1Q1/5r2/7p/1P1kp1bR/3P3R/1P1NP3/2n5 w - - 0 1",
		[]string{"a8c6", "c7d6"}, // Both are mate in 2
	},
	{
		"rn1kr3/pBp2Q1p/8/2b3p1/3n4/5P2/PPbK2PP/RNB4R b - - 1 15",
		[]string{"c5b4"},
	},
	{
		"rr4k1/6P1/1q5p/2p1pP2/p2nN2P/P2PQ3/2P5/2KRR3 b - - 0 27",
		[]string{"b6b1", "b6b2"}, // Both are mate
	},
}

func TestMateInTwos(t *testing.T) {
	for i, test := range MateInTwos {
		pos := parseFen(test.fen)
		best := pos.abRoot(4)
		found := false
		for _, valid := range test.best {
			if best.move.String() == valid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("MateinTwo %v. Got %v, Expected one of %v", i+1, best.move, test.best)
		}
	}
}

// Quiescence search for negamax (reference implementation with alpha-beta)
func (p *Position) negaQuiesce(alpha, beta, depth int) int {
	standPat := p.evaluate()

	if standPat >= beta {
		return beta
	}
	if standPat > alpha {
		alpha = standPat
	}

	// Depth limit reached
	if depth <= 0 {
		return alpha
	}

	moves := p.generateMoves()
	if len(*moves) == 0 {
		putMoveList(moves)
		if p.inCheck() {
			return -MAX_SCORE
		}
		return 0
	}

	for _, move := range *moves {
		if !isCapture(move) {
			continue
		}
		p.doMove(move, &StateInfo{})
		score := -p.negaQuiesce(-beta, -alpha, depth-1)
		p.undoMove(move)
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

func (p *Position) negamax(alpha, beta, depth int) int {
	moves := p.generateMoves()
	// Terminal position: no legal moves
	if len(*moves) == 0 {
		putMoveList(moves)
		if p.inCheck() {
			return -MAX_SCORE // Checkmate
		}
		return 0 // Stalemate
	}
	if depth == 0 {
		putMoveList(moves)
		return p.negaQuiesce(alpha, beta, MAX_QUIESCE_DEPTH)
	}
	for _, move := range *moves {
		p.doMove(move, &StateInfo{})
		score := -p.negamax(-beta, -alpha, depth-1)
		p.undoMove(move)
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

func (p *Position) negamaxRoot(depth int) MoveScore {
	alpha, beta, bestMove := -MAX_SCORE, MAX_SCORE, Move(0)
	moves := p.generateMoves()
	for _, move := range *moves {
		p.doMove(move, &StateInfo{})
		score := -p.negamax(-beta, -alpha, depth-1)
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

func (p *Position) negamaxC(depth int, move Move, c chan MoveScore) {
	c <- MoveScore{move, -p.negamax(-MAX_SCORE, MAX_SCORE, depth)}
}

// This doesn't seem to produce consistent results.
func (p *Position) negamaxRootP(depth int) MoveScore {
	best := MoveScore{Move(0), -MAX_SCORE}
	moves := p.generateMoves()
	c := make(chan MoveScore, len(*moves))
	for _, move := range *moves {
		dup := p.dup()
		dup.doMove(move, &StateInfo{})
		go dup.negamaxC(depth-1, move, c)
	}
	numMoves := len(*moves)
	putMoveList(moves)
	for i := 0; i < numMoves; i++ {
		ms := <-c
		if ms.score > best.score {
			best = ms
		}
	}
	return best
}

// This passes, but it's slow.
// Note: With killer moves, move ordering differs, so we only compare scores.
// Different moves with the same score are acceptable.
func TestNegaEqualsAB(t *testing.T) {
	for _, fen := range TestFens {
		// Clear TT to avoid interference from previous positions
		TT_GLOBAL.clear()
		clearKillers()
		pos := parseFen(fen)
		for j := 2; j <= 4; j += 1 {
			abBest := pos.abRoot(uint8(j))
			negaBest := pos.negamaxRoot(j)
			// Compare scores - moves may differ due to move ordering
			if abBest.score != negaBest.score {
				t.Errorf("Depth %v, Pos: %v\nScores differ.\nAb: %v\nNega:%v", j, fen, abBest, negaBest)
			}
		}
	}
}

func TestTerminalPositions(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		expected int
	}{
		// Stalemate: should return 0
		{"Stalemate", "8/8/8/8/8/5k2/5p2/5K2 w - - 0 1", 0},
		// Checkmate: should return -MAX_SCORE
		{"Checkmate", "4R1k1/5ppp/8/8/8/8/8/6K1 b - - 0 1", -MAX_SCORE},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pos := parseFen(test.fen)
			moves := pos.generateMoves()
			numMoves := len(*moves)
			putMoveList(moves)

			if numMoves != 0 {
				t.Fatalf("Expected terminal position (0 moves), got %d moves", numMoves)
			}

			// Test at depth 0, 1, 2
			for depth := uint8(0); depth <= 2; depth++ {
				score := pos.ab(-MAX_SCORE, MAX_SCORE, depth, 0)
				if score != test.expected {
					t.Errorf("depth=%d: expected %d, got %d", depth, test.expected, score)
				}
			}
		})
	}
}
