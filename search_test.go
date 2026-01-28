package main

import (
	"testing"
)

type BestMoveTest struct {
	fen, best string
}

// from https://sites.google.com/site/darktemplarchess/mate-in-2-puzzles
var MateInTwos = []BestMoveTest{
	BestMoveTest{
		"2bqkbn1/2pppp2/np2N3/r3P1p1/p2N2B1/5Q2/PPPPKPP1/RNB2r2 w KQkq - 0 1",
		"f3f7",
	},
	BestMoveTest{
		"8/6K1/1p1B1RB1/8/2Q5/2n1kP1N/3b4/4n3 w - - 0 1",
		"d6a3",
	},
	BestMoveTest{
		"B7/K1B1p1Q1/5r2/7p/1P1kp1bR/3P3R/1P1NP3/2n5 w - - 0 1",
		"a8c6",
	},
	BestMoveTest{
		"rn1kr3/pBp2Q1p/8/2b3p1/3n4/5P2/PPbK2PP/RNB4R b - - 1 15",
		"c5b4",
	},
	BestMoveTest{
		"rr4k1/6P1/1q5p/2p1pP2/p2nN2P/P2PQ3/2P5/2KRR3 b - - 0 27",
		"b6b1", //b6b2 is also mate... replace test when possible.
	},
}

func TestMateInTwos(t *testing.T) {
	for i, test := range MateInTwos {
		pos := parseFen(test.fen)
		best := pos.abRoot(4)
		if best.move.String() != test.best {
			t.Errorf("MateinTwo %v. Got %v, Expected %v", i+1, best.move, test.best)
		}
	}
}

func (p *Position) negamax(depth int) int {
	moves := p.generateMoves()
	// Terminal position: no legal moves
	if len(moves) == 0 {
		if p.inCheck() {
			return -MAX_SCORE // Checkmate
		}
		return 0 // Stalemate
	}
	if depth == 0 {
		return p.evaluate()
	}
	best := -MAX_SCORE
	for _, move := range moves {
		p.doMove(move, &StateInfo{})
		score := -p.negamax(depth - 1)
		p.undoMove(move)
		if score > best {
			best = score
		}
	}
	return best
}

func (p *Position) negamaxRoot(depth int) MoveScore {
	best := MoveScore{Move(0), -MAX_SCORE}
	for _, move := range p.generateMoves() {
		p.doMove(move, &StateInfo{})
		score := -p.negamax(depth - 1)
		p.undoMove(move)
		if score > best.score {
			best = MoveScore{move, score}
		}
	}
	return best
}

func (p *Position) negamaxC(depth int, move Move, c chan MoveScore) {
	c <- MoveScore{move, -p.negamax(depth)}
}

// This doesn't seem to produce consistent results.
func (p *Position) negamaxRootP(depth int) MoveScore {
	best := MoveScore{Move(0), -MAX_SCORE}
	moves := p.generateMoves()
	c := make(chan MoveScore, len(moves))
	for _, move := range moves {
		dup := p.dup()
		dup.doMove(move, &StateInfo{})
		go dup.negamaxC(depth-1, move, c)
	}
	for i := 0; i < len(moves); i++ {
		ms := <-c
		if ms.score > best.score {
			best = ms
		}
	}
	return best
}

// This passes, but it's slow.
func TestNegaEqualsAB(t *testing.T) {
	for _, fen := range TestFens {
		pos := parseFen(fen)
		for j := 2; j <= 4; j += 1 {
			abBest := pos.abRoot(uint8(j))
			negaBest := pos.negamaxRoot(j)
			if abBest != negaBest {
				t.Errorf("Depth %v, Pos: %v\nAb not equal nega.\nAb: %v\nNega:%v", j, fen, abBest, negaBest)
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

			if len(moves) != 0 {
				t.Fatalf("Expected terminal position (0 moves), got %d moves", len(moves))
			}

			// Test at depth 0, 1, 2
			for depth := uint8(0); depth <= 2; depth++ {
				score := pos.ab(-MAX_SCORE, MAX_SCORE, depth)
				if score != test.expected {
					t.Errorf("depth=%d: expected %d, got %d", depth, test.expected, score)
				}
			}
		})
	}
}
