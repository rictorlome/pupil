package main

import (
	"testing"
)

type TestStatelessUpdate struct {
	uci, initial, end string
}

var UpdatePositionTests = []TestStatelessUpdate{
	TestStatelessUpdate{"e2e4", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1"},
	TestStatelessUpdate{"e7e5", "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1", "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2"},
	TestStatelessUpdate{"g1f3", "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2", "rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"},
	TestStatelessUpdate{"b8c6", "rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2", "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3"},
	TestStatelessUpdate{"f1b5", "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", "r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3"},
	TestStatelessUpdate{"g8f6", "r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3", "r1bqkb1r/pppp1ppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4"},
	TestStatelessUpdate{"e1g1", "r1bqkb1r/pppp1ppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4", "r1bqkb1r/pppp1ppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4"},
	TestStatelessUpdate{"d7d6", "r1bqkb1r/pppp1ppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4", "r1bqkb1r/ppp2ppp/2np1n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 w kq - 0 5"},
	TestStatelessUpdate{"f1e1", "r1bqkb1r/ppp2ppp/2np1n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 w kq - 0 5", "r1bqkb1r/ppp2ppp/2np1n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQR1K1 b kq - 1 5"},
	TestStatelessUpdate{"a7a6", "r1bqkb1r/ppp2ppp/2np1n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQR1K1 b kq - 1 5", "r1bqkb1r/1pp2ppp/p1np1n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQR1K1 w kq - 0 6"},
	TestStatelessUpdate{"b5a4", "r1bqkb1r/1pp2ppp/p1np1n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQR1K1 w kq - 0 6", "r1bqkb1r/1pp2ppp/p1np1n2/4p3/B3P3/5N2/PPPP1PPP/RNBQR1K1 b kq - 1 6"},
	TestStatelessUpdate{"b7b5", "r1bqkb1r/1pp2ppp/p1np1n2/4p3/B3P3/5N2/PPPP1PPP/RNBQR1K1 b kq - 1 6", "r1bqkb1r/2p2ppp/p1np1n2/1p2p3/B3P3/5N2/PPPP1PPP/RNBQR1K1 w kq - 0 7"},
	TestStatelessUpdate{"a4b3", "r1bqkb1r/2p2ppp/p1np1n2/1p2p3/B3P3/5N2/PPPP1PPP/RNBQR1K1 w kq - 0 7", "r1bqkb1r/2p2ppp/p1np1n2/1p2p3/4P3/1B3N2/PPPP1PPP/RNBQR1K1 b kq - 1 7"},
	TestStatelessUpdate{"b5b4", "r1bqkb1r/2p2ppp/p1np1n2/1p2p3/4P3/1B3N2/PPPP1PPP/RNBQR1K1 b kq - 1 7", "r1bqkb1r/2p2ppp/p1np1n2/4p3/1p2P3/1B3N2/PPPP1PPP/RNBQR1K1 w kq - 0 8"},
	TestStatelessUpdate{"c2c4", "r1bqkb1r/2p2ppp/p1np1n2/4p3/1p2P3/1B3N2/PPPP1PPP/RNBQR1K1 w kq - 0 8", "r1bqkb1r/2p2ppp/p1np1n2/4p3/1pP1P3/1B3N2/PP1P1PPP/RNBQR1K1 b kq c3 0 8"},
	TestStatelessUpdate{"b4c3", "r1bqkb1r/2p2ppp/p1np1n2/4p3/1pP1P3/1B3N2/PP1P1PPP/RNBQR1K1 b kq c3 0 8", "r1bqkb1r/2p2ppp/p1np1n2/4p3/4P3/1Bp2N2/PP1P1PPP/RNBQR1K1 w kq - 0 9"},
	TestStatelessUpdate{"b1c3", "r1bqkb1r/2p2ppp/p1np1n2/4p3/4P3/1Bp2N2/PP1P1PPP/RNBQR1K1 w kq - 0 9", "r1bqkb1r/2p2ppp/p1np1n2/4p3/4P3/1BN2N2/PP1P1PPP/R1BQR1K1 b kq - 0 9"},
	TestStatelessUpdate{"f8e7", "r1bqkb1r/2p2ppp/p1np1n2/4p3/4P3/1BN2N2/PP1P1PPP/R1BQR1K1 b kq - 0 9", "r1bqk2r/2p1bppp/p1np1n2/4p3/4P3/1BN2N2/PP1P1PPP/R1BQR1K1 w kq - 1 10"},
	TestStatelessUpdate{"d2d4", "r1bqk2r/2p1bppp/p1np1n2/4p3/4P3/1BN2N2/PP1P1PPP/R1BQR1K1 w kq - 1 10", "r1bqk2r/2p1bppp/p1np1n2/4p3/3PP3/1BN2N2/PP3PPP/R1BQR1K1 b kq - 0 10"},
	TestStatelessUpdate{"e8g8", "r1bqk2r/2p1bppp/p1np1n2/4p3/3PP3/1BN2N2/PP3PPP/R1BQR1K1 b kq - 0 10", "r1bq1rk1/2p1bppp/p1np1n2/4p3/3PP3/1BN2N2/PP3PPP/R1BQR1K1 w - - 1 11"},
	TestStatelessUpdate{"d4e5", "r1bq1rk1/2p1bppp/p1np1n2/4p3/3PP3/1BN2N2/PP3PPP/R1BQR1K1 w - - 1 11", "r1bq1rk1/2p1bppp/p1np1n2/4P3/4P3/1BN2N2/PP3PPP/R1BQR1K1 b - - 0 11"},
	TestStatelessUpdate{"c6d4", "r1bq1rk1/2p1bppp/p1np1n2/4P3/4P3/1BN2N2/PP3PPP/R1BQR1K1 b - - 0 11", "r1bq1rk1/2p1bppp/p2p1n2/4P3/3nP3/1BN2N2/PP3PPP/R1BQR1K1 w - - 1 12"},
	TestStatelessUpdate{"e5f6", "r1bq1rk1/2p1bppp/p2p1n2/4P3/3nP3/1BN2N2/PP3PPP/R1BQR1K1 w - - 1 12", "r1bq1rk1/2p1bppp/p2p1P2/8/3nP3/1BN2N2/PP3PPP/R1BQR1K1 b - - 0 12"},
	TestStatelessUpdate{"d4b3", "r1bq1rk1/2p1bppp/p2p1P2/8/3nP3/1BN2N2/PP3PPP/R1BQR1K1 b - - 0 12", "r1bq1rk1/2p1bppp/p2p1P2/8/4P3/1nN2N2/PP3PPP/R1BQR1K1 w - - 0 13"},
	TestStatelessUpdate{"f6e7", "r1bq1rk1/2p1bppp/p2p1P2/8/4P3/1nN2N2/PP3PPP/R1BQR1K1 w - - 0 13", "r1bq1rk1/2p1Pppp/p2p4/8/4P3/1nN2N2/PP3PPP/R1BQR1K1 b - - 0 13"},
	TestStatelessUpdate{"b3c1", "r1bq1rk1/2p1Pppp/p2p4/8/4P3/1nN2N2/PP3PPP/R1BQR1K1 b - - 0 13", "r1bq1rk1/2p1Pppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 w - - 0 14"},
	TestStatelessUpdate{"e7d8q", "r1bq1rk1/2p1Pppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 w - - 0 14", "r1bQ1rk1/2p2ppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 b - - 0 14"},
	TestStatelessUpdate{"g8h8", "r1bQ1rk1/2p2ppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 b - - 0 14", "r1bQ1r1k/2p2ppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 w - - 1 15"},
	TestStatelessUpdate{"d8f8", "r1bQ1r1k/2p2ppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 w - - 1 15", "r1b2Q1k/2p2ppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 b - - 0 15"},
}

func StatelessMove(t *testing.T, test TestStatelessUpdate) {
	p := parse_fen(test.initial)
	move := p.parse_move(test.uci)
	p.do_move(move, &StateInfo{})
	if generate_fen(p) != test.end {
		t.Errorf("Move: %v, Expected not equal to actual:\n%v %v\n", test.uci, test.end, generate_fen(p))
	}
	p.undo_move(move)
	if generate_fen(p) != test.initial {
		t.Errorf("Move: %v, Expected not equal to actual:\n%v %v\n", test.uci, test.initial, generate_fen(p))
	}
}

var SpecificTest = []TestStatelessUpdate{
	TestStatelessUpdate{"b1c3", "r1bqkb1r/2p2ppp/p1np1n2/4p3/4P3/1Bp2N2/PP1P1PPP/RNBQR1K1 w kq - 0 9", "r1bqkb1r/2p2ppp/p1np1n2/4p3/4P3/1BN2N2/PP1P1PPP/R1BQR1K1 b kq - 0 9"},
}

func TestStatelessMove(t *testing.T) {
	for _, test := range UpdatePositionTests {
		StatelessMove(t, test)
	}
}
