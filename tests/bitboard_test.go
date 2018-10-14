package main

import (
  "testing"
  "fmt"
)

type ShiftByTest struct {
  start Bitboard
  d uint
  a uint
  end Bitboard
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
    if test.end != shift_by(test.start, test.d, test.a) {
      t.Error(fmt.Sprintf("Test %v: Expected %v, got %v", i, test.end, shift_by(test.start, test.d, test.a)))
    }
  }
}
