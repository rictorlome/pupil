package main

import (
	"fmt"
	"math/bits"
	"strings"
	// "strconv"
)

type Bitboard uint64

// COLORS
type Color bool
var WHITE Color = true
var BLACK Color = false
var COLORS = []Color{WHITE, BLACK}

// DIRECTIONS
var NORTH int = 8
var EAST int = 1
var SOUTH int = -NORTH
var WEST int = -EAST

var NORTH_WEST int = NORTH + WEST
var NORTH_EAST int = NORTH + EAST
var SOUTH_WEST int = SOUTH + WEST
var SOUTH_EAST int = SOUTH + EAST
var DIRECTIONS = []int {
	NORTH, EAST, SOUTH, WEST,
	NORTH_WEST, NORTH_EAST, SOUTH_WEST, SOUTH_EAST,
}
// BITBOARDS
var ALL_SQS Bitboard = 0xffffffffffffffff

// FILES
var FILE_ABB Bitboard = 0x101010101010101
var FILE_BBB Bitboard = shift_by(FILE_ABB, EAST, 1)
var FILE_CBB Bitboard = shift_by(FILE_ABB, EAST, 2)
var FILE_DBB Bitboard = shift_by(FILE_ABB, EAST, 3)
var FILE_EBB Bitboard = shift_by(FILE_ABB, EAST, 4)
var FILE_FBB Bitboard = shift_by(FILE_ABB, EAST, 5)
var FILE_GBB Bitboard = shift_by(FILE_ABB, EAST, 6)
var FILE_HBB Bitboard = shift_by(FILE_ABB, EAST, 7)
var FILE_BBS = []Bitboard{
		FILE_ABB, FILE_BBB, FILE_CBB, FILE_DBB,
		FILE_EBB, FILE_FBB, FILE_GBB, FILE_HBB,
	}

// RANKS
var RANK_1BB Bitboard = 0xff
var RANK_2BB Bitboard = shift_by(RANK_1BB, NORTH, 1)
var RANK_3BB Bitboard = shift_by(RANK_1BB, NORTH, 2)
var RANK_4BB Bitboard = shift_by(RANK_1BB, NORTH, 3)
var RANK_5BB Bitboard = shift_by(RANK_1BB, NORTH, 4)
var RANK_6BB Bitboard = shift_by(RANK_1BB, NORTH, 5)
var RANK_7BB Bitboard = shift_by(RANK_1BB, NORTH, 6)
var RANK_8BB Bitboard = shift_by(RANK_1BB, NORTH, 7)
var RANK_BBS = []Bitboard{
	RANK_1BB, RANK_2BB, RANK_3BB, RANK_4BB,
	RANK_5BB, RANK_6BB, RANK_7BB, RANK_8BB,
}

// SQUARES
type Square uint
const (
	SQ_A1 Square = iota
	SQ_B1
	SQ_C1
	SQ_D1
	SQ_E1
	SQ_F1
	SQ_G1
	SQ_H1
	SQ_A2
	SQ_B2
	SQ_C2
	SQ_D2
	SQ_E2
	SQ_F2
	SQ_G2
	SQ_H2
	SQ_A3
	SQ_B3
	SQ_C3
	SQ_D3
	SQ_E3
	SQ_F3
	SQ_G3
	SQ_H3
	SQ_A4
	SQ_B4
	SQ_C4
	SQ_D4
	SQ_E4
	SQ_F4
	SQ_G4
	SQ_H4
	SQ_A5
	SQ_B5
	SQ_C5
	SQ_D5
	SQ_E5
	SQ_F5
	SQ_G5
	SQ_H5
	SQ_A6
	SQ_B6
	SQ_C6
	SQ_D6
	SQ_E6
	SQ_F6
	SQ_G6
	SQ_H6
	SQ_A7
	SQ_B7
	SQ_C7
	SQ_D7
	SQ_E7
	SQ_F7
	SQ_G7
	SQ_H7
	SQ_A8
	SQ_B8
	SQ_C8
	SQ_D8
	SQ_E8
	SQ_F8
	SQ_G8
	SQ_H8
	NULL_SQ // NULL SQUARE = 64
)

var SQUARE_BBS [64]Bitboard

func binary(b Bitboard) string {
	return fmt.Sprintf("%v%b", strings.Repeat("0", leading_zeros(b)), uint64(b))
}

func leading_zeros(b Bitboard) int {
	return bits.LeadingZeros64(uint64(b))
}

func popcount(b Bitboard) int {
	return bits.OnesCount64(uint64(b))
}

// https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func trailing_zeros(b Bitboard) int {
	return bits.TrailingZeros64(uint64(b))
}

func shift(b Bitboard, d int) Bitboard {
	return shift_by(b, d, 1)
}

func shift_by(b Bitboard, d int, amount int) Bitboard {
	switch d > 0 {
	case true:
		return b << uint(d * amount)
	case false:
		return b >> uint(-d * amount)
	default:
		return b
	}
}

func make_square(r uint, f uint) Square {
	return Square(r * f - 1)
}

var RANK_1 uint = 1
var RANK_2 uint = 2
var RANK_3 uint = 3
var RANK_4 uint = 4
var RANK_5 uint = 5
var RANK_6 uint = 6
var RANK_7 uint = 7
var RANK_8 uint = 8

var FILE_A uint = 1
var FILE_B uint = 2
var FILE_C uint = 3
var FILE_D uint = 4
var FILE_E uint = 5
var FILE_F uint = 6
var FILE_G uint = 7
var FILE_H uint = 8

func (b Bitboard) String() string {
	var s string
	runes := []rune(binary(b))
	for i := 0; i <= 7; i++ {
		s += Reverse(string(runes[(8*i):(8*(i+1))]))
		s += "\n"
	}
	return s
}

// type Move uint16
