package main

import (
	"fmt"
	"testing"
)

type TestUpdatePosition struct {
	uci, fen string
}

var UpdatePositionTests = []TestUpdatePosition{
	TestUpdatePosition{"e2e4", "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1"},
	TestUpdatePosition{"e7e5", "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2"},
	TestUpdatePosition{"g1f3", "rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"},
	TestUpdatePosition{"b8c6", "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3"},
	TestUpdatePosition{"f1b5", "r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3"},
	TestUpdatePosition{"g8f6", "r1bqkb1r/pppp1ppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4"},
	TestUpdatePosition{"e1g1", "r1bqkb1r/pppp1ppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4"},
	TestUpdatePosition{"d7d6", "r1bqkb1r/ppp2ppp/2np1n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 w kq - 0 5"},
	TestUpdatePosition{"f1e1", "r1bqkb1r/ppp2ppp/2np1n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQR1K1 b kq - 1 5"},
	TestUpdatePosition{"a7a6", "r1bqkb1r/1pp2ppp/p1np1n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQR1K1 w kq - 0 6"},
	TestUpdatePosition{"b5a4", "r1bqkb1r/1pp2ppp/p1np1n2/4p3/B3P3/5N2/PPPP1PPP/RNBQR1K1 b kq - 1 6"},
	TestUpdatePosition{"b7b5", "r1bqkb1r/2p2ppp/p1np1n2/1p2p3/B3P3/5N2/PPPP1PPP/RNBQR1K1 w kq - 0 7"},
	TestUpdatePosition{"a4b3", "r1bqkb1r/2p2ppp/p1np1n2/1p2p3/4P3/1B3N2/PPPP1PPP/RNBQR1K1 b kq - 1 7"},
	TestUpdatePosition{"b5b4", "r1bqkb1r/2p2ppp/p1np1n2/4p3/1p2P3/1B3N2/PPPP1PPP/RNBQR1K1 w kq - 0 8"},
	TestUpdatePosition{"c2c4", "r1bqkb1r/2p2ppp/p1np1n2/4p3/1pP1P3/1B3N2/PP1P1PPP/RNBQR1K1 b kq c3 0 8"},
	TestUpdatePosition{"b4c3", "r1bqkb1r/2p2ppp/p1np1n2/4p3/4P3/1Bp2N2/PP1P1PPP/RNBQR1K1 w kq - 0 9"},
	TestUpdatePosition{"b1c3", "r1bqkb1r/2p2ppp/p1np1n2/4p3/4P3/1BN2N2/PP1P1PPP/R1BQR1K1 b kq - 0 9"},
	TestUpdatePosition{"f8e7", "r1bqk2r/2p1bppp/p1np1n2/4p3/4P3/1BN2N2/PP1P1PPP/R1BQR1K1 w kq - 1 10"},
	TestUpdatePosition{"d2d4", "r1bqk2r/2p1bppp/p1np1n2/4p3/3PP3/1BN2N2/PP3PPP/R1BQR1K1 b kq - 0 10"},
	TestUpdatePosition{"e8g8", "r1bq1rk1/2p1bppp/p1np1n2/4p3/3PP3/1BN2N2/PP3PPP/R1BQR1K1 w - - 1 11"},
	TestUpdatePosition{"d4e5", "r1bq1rk1/2p1bppp/p1np1n2/4P3/4P3/1BN2N2/PP3PPP/R1BQR1K1 b - - 0 11"},
	TestUpdatePosition{"c6d4", "r1bq1rk1/2p1bppp/p2p1n2/4P3/3nP3/1BN2N2/PP3PPP/R1BQR1K1 w - - 1 12"},
	TestUpdatePosition{"e5f6", "r1bq1rk1/2p1bppp/p2p1P2/8/3nP3/1BN2N2/PP3PPP/R1BQR1K1 b - - 0 12"},
	TestUpdatePosition{"d4b3", "r1bq1rk1/2p1bppp/p2p1P2/8/4P3/1nN2N2/PP3PPP/R1BQR1K1 w - - 0 13"},
	TestUpdatePosition{"f6e7", "r1bq1rk1/2p1Pppp/p2p4/8/4P3/1nN2N2/PP3PPP/R1BQR1K1 b - - 0 13"},
	TestUpdatePosition{"b3c1", "r1bq1rk1/2p1Pppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 w - - 0 14"},
	TestUpdatePosition{"e7d8q", "r1bQ1rk1/2p2ppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 b - - 0 14"},
	TestUpdatePosition{"g8h8", "r1bQ1r1k/2p2ppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 w - - 1 15"},
	TestUpdatePosition{"d8f8", "r1b2Q1k/2p2ppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 b - - 0 15"},
}

func TestMove(t *testing.T) {
	p := parse_fen(INITIAL_FEN)
	for _, test := range UpdatePositionTests {
		p.do_move(p.parse_move(test.uci), StateInfo{})
		if generate_fen(p) != test.fen {
			t.Error(fmt.Sprintf("Move %v: Expected fen %v, got %v", test.uci, test.fen, generate_fen(p)))
		}
	}
}
