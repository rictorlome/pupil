package main

import (
	"reflect"
	"testing"
)

var parseFenTests = []struct {
	fen string
	pos Position
}{
	{
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		Position{
			state: &StateInfo{
				castling_rights: 15,
				ep_sq: 20,
				rule_50: 0,
				opposite_color_attacks: 1830813695870,
				prev: nil,
				captured: 0,
			},
			ply: 3,
			placement_by_square: []Piece{
				4, 3, 2, 1, 0, 2, 3, 4,
				5, 5, 5, 5, NULL_PIECE, 5, 5, 5, NULL_PIECE,
				NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE,
				NULL_PIECE, NULL_PIECE, NULL_PIECE, 5, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE,
				NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE,
				NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE,
				11, 11, 11, 11, 11, 11, 11, 11,
				10, 9, 8, 7, 6, 8, 9, 10,
			},
			placement: []Bitboard{
				16, 8, 36, 66, 129, 268496640, 1152921504606846976, 576460752303423488,
				2594073385365405696, 4755801206503243776, 9295429630892703744, 71776119061217280,
			},
			stm: BLACK,
			occ: 18446462599001337855,
		},
	},
	{
		"rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
		Position{
			ply: 4,
			stm: WHITE,
			occ: 18441959068093444095,
			state: &StateInfo{
				castling_rights: 15,
				ep_sq: 44,
				rule_50: 0,
				opposite_color_attacks: 9151313638215188480,
				prev: nil,
				captured: 0,
			},
			placement_by_square: []Piece{
				4, 3, 2, 1, 0, 2, 3, 4, 5, 5, 5, 5, NULL_PIECE, 5, 5, 5, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, 5, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, 11, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, NULL_PIECE, 11, 11, 11, 11, NULL_PIECE, 11, 11, 11, 10, 9, 8, 7, 6, 8, 9, 10,
			},
			placement: []Bitboard{
				16, 8, 36, 66, 129, 268496640, 1152921504606846976, 576460752303423488,
				2594073385365405696, 4755801206503243776, 9295429630892703744, 67272588153323520,
			},
		},
	},
}


func TestParseFen(t *testing.T) {
	for _, test := range parseFenTests {
		actual := parse_fen(test.fen)
		if !reflect.DeepEqual(actual, test.pos) {
			t.Errorf(
				"Expected Postition: \n%s, got \n%s from \n%s",
				test.pos.debug(), actual.debug(), test.fen,
			)
		}
	}
}

func TestGenerateFen(t *testing.T){
	for _, test := range parseFenTests {
		actual := generate_fen(test.pos)
		if !reflect.DeepEqual(actual, test.fen) {
			t.Errorf("Expected fen: \n%s, got \n%s", test.fen, actual,)
		}
	}
}

var tests = []string{
	"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
	"rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
	"rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
	"r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3",
	"r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3",
	"r1bqkb1r/pppp1ppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
	"r1bqkb1r/pppp1ppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4",
	"r1bqkb1r/ppp2ppp/2np1n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 w kq - 0 5",
	"r1bqkb1r/ppp2ppp/2np1n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQR1K1 b kq - 1 5",
	"r1bqkb1r/1pp2ppp/p1np1n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQR1K1 w kq - 0 6",
	"r1bqkb1r/1pp2ppp/p1np1n2/4p3/B3P3/5N2/PPPP1PPP/RNBQR1K1 b kq - 1 6",
	"r1bqkb1r/2p2ppp/p1np1n2/1p2p3/B3P3/5N2/PPPP1PPP/RNBQR1K1 w kq b6 0 7",
	"r1bqkb1r/2p2ppp/p1np1n2/1p2p3/4P3/1B3N2/PPPP1PPP/RNBQR1K1 b kq - 1 7",
	"r1bqkb1r/2p2ppp/p1np1n2/4p3/1p2P3/1B3N2/PPPP1PPP/RNBQR1K1 w kq - 0 8",
	"r1bqkb1r/2p2ppp/p1np1n2/4p3/1pP1P3/1B3N2/PP1P1PPP/RNBQR1K1 b kq c3 0 8",
	"r1bqkb1r/2p2ppp/p1np1n2/4p3/4P3/1Bp2N2/PP1P1PPP/RNBQR1K1 w kq - 0 9",
	"r1bqkb1r/2p2ppp/p1np1n2/4p3/4P3/1BN2N2/PP1P1PPP/R1BQR1K1 b kq - 0 9",
	"r1bqk2r/2p1bppp/p1np1n2/4p3/4P3/1BN2N2/PP1P1PPP/R1BQR1K1 w kq - 1 10",
	"r1bqk2r/2p1bppp/p1np1n2/4p3/3PP3/1BN2N2/PP3PPP/R1BQR1K1 b kq d3 0 10",
	"r1bq1rk1/2p1bppp/p1np1n2/4p3/3PP3/1BN2N2/PP3PPP/R1BQR1K1 w - - 1 11",
	"r1bq1rk1/2p1bppp/p1np1n2/4P3/4P3/1BN2N2/PP3PPP/R1BQR1K1 b - - 0 11",
	"r1bq1rk1/2p1bppp/p2p1n2/4P3/3nP3/1BN2N2/PP3PPP/R1BQR1K1 w - - 1 12",
	"r1bq1rk1/2p1bppp/p2p1P2/8/3nP3/1BN2N2/PP3PPP/R1BQR1K1 b - - 0 12",
	"r1bq1rk1/2p1bppp/p2p1P2/8/4P3/1nN2N2/PP3PPP/R1BQR1K1 w - - 0 13",
	"r1bq1rk1/2p1Pppp/p2p4/8/4P3/1nN2N2/PP3PPP/R1BQR1K1 b - - 0 13",
	"r1bq1rk1/2p1Pppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 w - - 0 14",
	"r1bQ1rk1/2p2ppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 b - - 0 14",
	"r1bQ1r1k/2p2ppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 w - - 1 15",
	"r1b2Q1k/2p2ppp/p2p4/8/4P3/2N2N2/PP3PPP/R1nQR1K1 b - - 0 15",
}


func TestParseAndGenerateFen(t *testing.T) {
	for _, test := range tests {
		actual := parse_fen(test)
		if test != generate_fen(actual) {
			t.Errorf("Expected fen: %s, got %s", actual, test)
		}
	}
}
