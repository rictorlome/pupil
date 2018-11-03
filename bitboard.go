package main

import (
	"fmt"
	"math/bits"
	"strings"
	// "strconv"
)

func binary(b Bitboard) string {
	return fmt.Sprintf("%064b", uint64(b))
}

func king_attacks(b Bitboard) Bitboard {
	var x Bitboard
	for _, direction := range DIRECTIONS {
		x |= shift_direction(b, direction)
	}
	return x
}

func knight_attacks(b Bitboard) Bitboard {
	return shift_direction(shift_direction(b, NORTH), NORTH_EAST) |
		shift_direction(shift_direction(b, NORTH), NORTH_WEST) |
		shift_direction(shift_direction(b, EAST), NORTH_EAST) |
		shift_direction(shift_direction(b, EAST), SOUTH_EAST) |
		shift_direction(shift_direction(b, SOUTH), SOUTH_EAST) |
		shift_direction(shift_direction(b, SOUTH), SOUTH_WEST) |
		shift_direction(shift_direction(b, WEST), NORTH_WEST) |
		shift_direction(shift_direction(b, WEST), SOUTH_WEST)
}

func leading_zeros(b Bitboard) int {
	return bits.LeadingZeros64(uint64(b))
}

func make_piece(s string) Piece {
	return Piece(strings.Index(PIECE_STRING, s))
}

func make_square(r int, f int) Square {
	return Square(r*8 + f)
}

func occupied_at(b Bitboard, sq Square) bool {
	return (b>>sq)&1 == 1
}

func popcount(b Bitboard) int {
	return bits.OnesCount64(uint64(b))
}

func possible_enpassant_sq(sq Square) bool {
	return square_rank(sq) == 2 || square_rank(sq) == 5
}

// https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// shifts bitboard in a direction, prevents wrapping
func shift_direction(b Bitboard, direction int) Bitboard {
	if direction == NORTH_EAST || direction == EAST || direction == SOUTH_EAST {
		return signed_shift(b&^FILE_HBB, direction)
	}
	if direction == NORTH_WEST || direction == WEST || direction == SOUTH_WEST {
		return signed_shift(b&^FILE_ABB, direction)
	}
	return signed_shift(b, direction)
}

// shifts bitboard by an amount (which can be a positive or negative number)
func signed_shift(b Bitboard, amount int) Bitboard {
	if amount >= 0 {
		return b << uint(amount)
	}
	return b >> uint(-amount)
}

func square_file(sq Square) int {
	return int(sq % 8)
}

func square_rank(sq Square) int {
	return int(sq / 8)
}

func (b Bitboard) String() string {
	var s string
	runes := []rune(binary(b))
	for i := RANK_1; i <= RANK_8; i++ {
		s += reverse(string(runes[(8*i):(8*(i+1))])) + "\n"
	}
	return s
}

func (p Piece) String() string {
	return string([]rune(PIECE_STRING)[p])
}

func (s Square) String() string {
	files := "abcdefgh"
	return fmt.Sprintf("%c%d", files[square_file(s)], square_rank(s)+1)
}

func trailing_zeros(b Bitboard) int {
	return bits.TrailingZeros64(uint64(b))
}
