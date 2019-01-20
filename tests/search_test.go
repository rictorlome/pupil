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
		pos := parse_fen(test.fen)
		best := pos.alphaBetaRoot(4)
		if best.move.String() != test.best {
			t.Errorf("MateinTwo %v. Got %v, Expected %v", i+1, best.move, test.best)
		}
	}
}
