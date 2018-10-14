package main

// import "fmt"

type Bitboard uint64

// COLORS
type Color bool

var WHITE Color = true
var BLACK Color = false
var COLORS = []Color{WHITE, BLACK}

// DIRECTIONS
var NORTH uint = 8
var EAST uint = 1
var SOUTH uint = -NORTH
var WEST uint = -EAST

var NORTH_WEST uint = NORTH + WEST
var NORTH_EAST uint = NORTH + EAST
var SOUTH_WEST uint = SOUTH + WEST
var SOUTH_EAST uint = SOUTH + EAST

// BITBOARDS
var ALL_SQS Bitboard = 0xffffffffffffffff

// FILES
var FILE_ABB Bitboard = 0x101010101010101
var FILE_BBB Bitboard = FILE_ABB >> 1
var FILE_CBB Bitboard = FILE_ABB >> 2
var FILE_DBB Bitboard = FILE_ABB >> 3
var FILE_EBB Bitboard = FILE_ABB >> 4
var FILE_FBB Bitboard = FILE_ABB >> 5
var FILE_GBB Bitboard = FILE_ABB >> 6
var FILE_HBB Bitboard = FILE_ABB >> 7

// RANKS
var RANK_1BB Bitboard = 0xff
var RANK_2BB Bitboard = RANK_1BB >> (NORTH * 1)
var RANK_3BB Bitboard = RANK_1BB >> (NORTH * 2)
var RANK_4BB Bitboard = RANK_1BB >> (NORTH * 3)
var RANK_5BB Bitboard = RANK_1BB >> (NORTH * 4)
var RANK_6BB Bitboard = RANK_1BB >> (NORTH * 5)
var RANK_7BB Bitboard = RANK_1BB >> (NORTH * 6)
var RANK_8BB Bitboard = RANK_1BB >> (NORTH * 7)

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

var SQUARE_BBS []Bitboard

func shift_by(b Bitboard, d uint, amount uint) Bitboard {
	return b >> (d * amount)
}

// func shift_by(b Bitboard, d uint, amount uint) Bitboard {
// 	switch d {
// 	case NORTH, EAST:
// 		return b >> (d * amount)
// 	case SOUTH, WEST:
// 		return b << (d * amount)
// 	default:
// 		return b
// 	}
// }

// type Move uint16
