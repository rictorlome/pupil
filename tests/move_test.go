package main

import (
	"fmt"
	"testing"
)

type MoveTypeTest struct {
  mt MoveType
  capture, en_passant, promotion, castle bool
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
    mock_move := to_move(SQ_A1, SQ_B6, test.mt)
    // assertEquals(t, "capture", test.mt, is_capture(mock_move), test.capture)
    assertEquals(t, "castle", test.mt, is_castle(mock_move), test.castle)
    // assertEquals(t, "promotion", test.mt, is_promotion(mock_move), test.promotion)
    // assertEquals(t, "en_passant", test.mt, is_enpassant(mock_move), test.en_passant)
    // assertEquals(t, "move type", test.mt, is_move_type(mock_move, test.mt), true)
  }
}
