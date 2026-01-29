package main

import (
	"fmt"
	"testing"
)

type perftSeq struct {
	startFen string
	perfts   []perft
}

// go test -v -run TestPerft -timeout 24h
// Tables taken from https://www.chessprogramming.org/Perft_Results
var initialPerft = perftSeq{
	INITIAL_FEN,
	[]perft{
		perft{0, 1, 0, 0, 0, 0, 0, 0},
		perft{1, 20, 0, 0, 0, 0, 0, 0},
		perft{2, 400, 0, 0, 0, 0, 0, 0},
		perft{3, 8902, 34, 0, 0, 0, 12, 0},
		perft{4, 197281, 1576, 0, 0, 0, 469, 8},
		perft{5, 4865609, 82719, 258, 0, 0, 27351, 347},
		perft{6, 119060324, 2812008, 5248, 0, 0, 809099, 10828},
		// perft{7, 3195901860, 108329926, 319617, 883453, 0, 33103848, 435767},
		// perft{8, 84998978956, 3523740106, 7187977, 23605205, 0, 968981593, 9852036},
	},
}

// Kiwipete
var secondaryPerft = perftSeq{
	KIWIPETE_FEN,
	[]perft{
		perft{0, 1, 0, 0, 0, 0, 0, 0},
		perft{1, 48, 8, 0, 2, 0, 0, 0},
		perft{2, 2039, 351, 1, 91, 0, 3, 0},
		perft{3, 97862, 17102, 45, 3162, 0, 993, 1},
		perft{4, 4085603, 757163, 1929, 128013, 15172, 25523, 43},
		perft{5, 193690690, 35043416, 73365, 4993637, 8392, 3309887, 30171},
		// perft{6, 8031647685, 1558445089, 3577504, 184513607, 56627920, 92238050, 360003},
	},
}

var tertiaryPerft = perftSeq{
	"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
	[]perft{
		perft{0, 1, 0, 0, 0, 0, 0, 0},
		perft{1, 14, 1, 0, 0, 0, 2, 0},
		perft{2, 191, 14, 0, 0, 0, 10, 0},
		perft{3, 2812, 209, 2, 0, 0, 267, 0},
		perft{4, 43238, 3348, 123, 0, 0, 1680, 17},
		perft{5, 674624, 52051, 1165, 0, 0, 52950, 0},
		perft{6, 11030083, 940350, 33325, 0, 7552, 452473, 2733},
	},
}

var fourthPerft = perftSeq{
	"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
	[]perft{
		perft{1, 6, 0, 0, 0, 0, 0, 0},
		perft{2, 264, 87, 0, 6, 48, 10, 0},
		perft{3, 9467, 1021, 4, 0, 120, 38, 22},
		perft{4, 422333, 131393, 0, 7795, 60032, 15492, 5},
		perft{5, 15833292, 2046173, 6512, 0, 329464, 200568, 50562},
		perft{6, 706045033, 210369132, 212, 10882006, 81102984, 26973664, 81076}, // passes, too slow
	},
}

var allPerftTests = []perftSeq{
	initialPerft,
	secondaryPerft,
	tertiaryPerft,
	fourthPerft,
}

func TestPerft(t *testing.T) {
	for _, test := range allPerftTests {
		pos := parseFen(test.startFen)
		for _, expectedPft := range test.perfts {
			actualPft := getPerftParallel(&pos, expectedPft.depth)
			// actualPft := getPerftRecursive(&pos, expectedPft.depth, Move(0))
			if actualPft != expectedPft {
				t.Error(fmt.Sprintf("\nExpected:\t %v\n Got:\t\t %v", expectedPft, actualPft))
			}
		}
	}
}
