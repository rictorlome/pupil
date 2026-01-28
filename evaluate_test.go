package main

import (
	"testing"
)

type PositionValueTest struct {
	pt  PieceType
	sq  Square
	c   Color
	val int
}

var PositionValueTests = []PositionValueTest{
	PositionValueTest{PAWN, SQ_A2, WHITE, 5},
	PositionValueTest{PAWN, SQ_A2, BLACK, 50},
	PositionValueTest{PAWN, SQ_A4, WHITE, 0},
	PositionValueTest{PAWN, SQ_A4, BLACK, 5},
	PositionValueTest{PAWN, SQ_E4, WHITE, 20},
	PositionValueTest{PAWN, SQ_E4, BLACK, 25},
	PositionValueTest{ROOK, SQ_E7, WHITE, 10},
	PositionValueTest{ROOK, SQ_E7, BLACK, 0},
	PositionValueTest{KNIGHT, SQ_E2, WHITE, 5},
	PositionValueTest{KNIGHT, SQ_E2, BLACK, 0},
	PositionValueTest{BISHOP, SQ_B2, WHITE, 5},
	PositionValueTest{BISHOP, SQ_B2, BLACK, 0},
}

func TestAddPositionVal(t *testing.T) {
	for _, test := range PositionValueTests {
		if test.val != addPositionValue(test.pt, test.c, test.sq) {
			t.Errorf("Test: %v, Got: %v, Expected %v.", test, addPositionValue(test.pt, test.c, test.sq), test.val)
		}
	}
}
