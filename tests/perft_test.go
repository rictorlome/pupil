package main

import (
	"fmt"
	"testing"
)

type perft_seq struct {
	start_fen string
	perfts    []perft
}

// Tables taken from https://www.chessprogramming.org/Perft_Results
var initial_perft = perft_seq{
	"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	[]perft{
		perft{0, 1, 0, 0, 0, 0, 0, 0},
		perft{1, 20, 0, 0, 0, 0, 0, 0},
		perft{2, 400, 0, 0, 0, 0, 0, 0},
		perft{3, 8902, 34, 0, 0, 0, 12, 0},
		perft{4, 197281, 1576, 0, 0, 0, 469, 8},
		perft{5, 4865609, 82719, 258, 0, 0, 27351, 347},
		perft{6, 119060324, 2812008, 5248, 0, 0, 809099, 10828},
	},
}

// Kiwipete
var secondary_perft = perft_seq{
	"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
	[]perft{
		perft{0, 1, 0, 0, 0, 0, 0, 0},
		perft{1, 48, 8, 0, 2, 0, 0, 0},
		perft{2, 2039, 351, 1, 91, 0, 3, 0},
		perft{3, 97862, 17102, 45, 3162, 0, 993, 1},
		// perft{4, 4085603, 757163, 1929, 128013, 15172, 25523, 43},
		// perft{5, 193690690, 35043416, 73365, 4993637, 8392, 3309887, 30171},
	},
}

var tertiary_perft = perft_seq{
	"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
	[]perft{
		perft{0, 1, 0, 0, 0, 0, 0, 0},
		perft{1, 14, 1, 0, 0, 0, 2, 0},
		perft{2, 191, 14, 0, 0, 0, 10, 0},
		perft{3, 2812, 209, 2, 0, 0, 267, 0},
		perft{4, 43238, 3348, 123, 0, 0, 1680, 17},
		perft{5, 674624, 52051, 1165, 0, 0, 52950, 0},
		// perft{6, 11030083, 940350, 33325, 0, 7552, 452473, 2733},
	},
}

var fourth_perft = perft_seq{
	"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
	[]perft{
		perft{1, 6, 0, 0, 0, 0, 0, 0},
		perft{2, 264, 87, 0, 6, 48, 10, 0},
		perft{3, 9467, 1021, 4, 0, 120, 38, 22},
		perft{4, 422333, 131393, 0, 7795, 60032, 15492, 5},
		perft{5, 15833292, 2046173, 6512, 0, 329464, 200568, 50562},
		// perft{6, 706045033, 210369132, 212, 10882006, 81102984, 26973664, 81076},
	},
}

var all_perft_tests = []perft_seq{
	initial_perft,
	secondary_perft,
	tertiary_perft,
	fourth_perft,
}

func TestPerft(t *testing.T) {
	for _, test := range all_perft_tests {
		pos := parse_fen(test.start_fen)
		for _, expected_pft := range test.perfts {
			actual_pft := get_perft_parallel(&pos, expected_pft.depth)
			// actual_pft := get_perft_recursive(&pos, expected_pft.depth, Move(0))
			if actual_pft != expected_pft {
				t.Error(fmt.Sprintf("\nExpected:\t %v\n Got:\t\t %v", expected_pft, actual_pft))
			}
		}
	}
}
