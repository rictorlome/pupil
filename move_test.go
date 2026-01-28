package main

import (
	"fmt"
	"testing"
)

type MoveTypeTest struct {
	mt                                    MoveType
	capture, enPassant, promotion, castle bool
}

var MoveTypeTests = []MoveTypeTest{
	MoveTypeTest{QUIET, false, false, false, false},
	MoveTypeTest{DOUBLE_PAWN_PUSH, false, false, false, false},
	MoveTypeTest{KING_CASTLE, false, false, false, true},
	MoveTypeTest{QUEEN_CASTLE, false, false, false, true},
	MoveTypeTest{CAPTURE, true, false, false, false},
	MoveTypeTest{EP_CAPTURE, true, true, false, false},
	MoveTypeTest{QUIET, false, false, false, false},
	MoveTypeTest{QUIET, false, false, false, false},
	MoveTypeTest{KNIGHT_PROMOTION, false, false, true, false},
	MoveTypeTest{BISHOP_PROMOTION, false, false, true, false},
	MoveTypeTest{ROOK_PROMOTION, false, false, true, false},
	MoveTypeTest{QUEEN_PROMOTION, false, false, true, false},
	MoveTypeTest{KNIGHT_PROMOTION_CAPTURE, true, false, true, false},
	MoveTypeTest{BISHOP_PROMOTION_CAPTURE, true, false, true, false},
	MoveTypeTest{ROOK_PROMOTION_CAPTURE, true, false, true, false},
	MoveTypeTest{QUEEN_PROMOTION_CAPTURE, true, false, true, false},
}

func assertEquals(t *testing.T, name string, mt MoveType, a bool, b bool) {
	if a != b {
		t.Error(fmt.Sprintf("%v: Error with %v. Expected %v, got %v", mt, name, b, a))
	}
}

func TestIsMoveType(t *testing.T) {
	for _, test := range MoveTypeTests {
		mockMove := toMove(SQ_A1, SQ_B6, test.mt)
		// assertEquals(t, "capture", test.mt, isCapture(mockMove), test.capture)
		assertEquals(t, "castle", test.mt, isCastle(mockMove), test.castle)
		// assertEquals(t, "promotion", test.mt, isPromotion(mockMove), test.promotion)
		// assertEquals(t, "enPassant", test.mt, isEnpassant(mockMove), test.enPassant)
		// assertEquals(t, "move type", test.mt, isMoveType(mockMove, test.mt), true)
	}
}
