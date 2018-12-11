package main

import (
	"fmt"
	"testing"
)

type AlignedBetween struct {
	a, b, c          Square
	aligned, between bool
}

var AlignedBetweens = []AlignedBetween{
	AlignedBetween{SQ_A1, SQ_A2, SQ_A3, true, false},
	AlignedBetween{SQ_A1, SQ_A2, SQ_A8, true, false},
	AlignedBetween{SQ_A1, SQ_A3, SQ_A1, true, false},
	AlignedBetween{SQ_B1, SQ_A2, SQ_A3, false, false},
	AlignedBetween{SQ_A1, SQ_H8, SQ_E5, true, true},
	AlignedBetween{SQ_A1, SQ_E5, SQ_H1, false, false},
	AlignedBetween{SQ_H1, SQ_H8, SQ_H1, true, false},
	AlignedBetween{SQ_B2, SQ_D4, SQ_C3, true, true},
	AlignedBetween{SQ_B2, SQ_C3, SQ_E5, true, false},
	AlignedBetween{SQ_B2, SQ_C3, SQ_E6, false, false},
	AlignedBetween{SQ_A1, SQ_E4, SQ_A4, false, false},
}

func TestAligned(t *testing.T) {
	for _, test := range AlignedBetweens {
		if aligned(test.a, test.b, test.c) != test.aligned {
			t.Error(fmt.Sprintf("Aligned test %v, %v, %v: Expected %v, got %v", test.a, test.b, test.c, test.aligned, aligned(test.a, test.b, test.c)))
		}
		if between(test.a, test.b, test.c) != test.between {
			t.Error(fmt.Sprintf("Between test %v, %v, %v: Expected %v, got %v", test.a, test.b, test.c, test.between, between(test.a, test.b, test.c)))
		}
	}
}

type ShiftByTest struct {
	start Bitboard
	d     int
	a     int
	end   Bitboard
}

var ShiftByTests = []ShiftByTest{
	ShiftByTest{FILE_ABB, EAST, 1, FILE_BBB},
	ShiftByTest{FILE_ABB, EAST, 7, FILE_HBB},
	ShiftByTest{FILE_BBB, WEST, 1, FILE_ABB},
	ShiftByTest{FILE_HBB, WEST, 7, FILE_ABB},
	ShiftByTest{RANK_1BB, NORTH, 1, RANK_2BB},
	ShiftByTest{RANK_1BB, NORTH, 7, RANK_8BB},
	ShiftByTest{RANK_2BB, SOUTH, 1, RANK_1BB},
	ShiftByTest{RANK_8BB, SOUTH, 7, RANK_1BB},
}

func TestShiftBy(t *testing.T) {
	for i, test := range ShiftByTests {
		if test.end != signed_shift(test.start, test.d*test.a) {
			t.Error(fmt.Sprintf("Test %v: Expected %v, got %v", i, test.end, signed_shift(test.start, test.d*test.a)))
		}
	}
}
