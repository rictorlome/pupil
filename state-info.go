package main

import (
// "fmt"
)

// NOTE: tmp function does not include attacking.
func can_castle(side int, color Color, occ Bitboard, cr int) bool {
	empty := occ&CASTLE_MOVE_SQS[int(color)*2+side] == 0
	has_right := cr&CASTLING_RIGHTS[int(color)*2+side] != 0
	return empty && has_right
}

func has_bit(s int, bit uint) bool {
	return (s>>bit)&1 == 1
}

func make_castle_state_info(available string) int {
	var castles int
	for _, char := range available {
		castles |= CHAR_TO_CASTLE[string(char)]
	}
	return castles
}
