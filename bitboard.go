package main

import (
	"fmt"
	"math/bits"
	"strings"
	// "strconv"
)

func aligned(sq1 Square, sq2 Square, sq3 Square) bool {
	return LINE_BBS[sq1][sq2]&SQUARE_BBS[sq3] != 0
}

func binary(b Bitboard) string {
	return fmt.Sprintf("%064b", uint64(b))
}

func bool_to_int(b bool) int {
	if b {
		return 0
	}
	return 1
}

func empty(b Bitboard) bool {
	return b == 0
}

func forward(color Color) int {
	if color == WHITE {
		return NORTH
	}
	return SOUTH
}

func last_rank(color Color) int {
	if color == WHITE {
		return 7
	}
	return 0
}

func leading_zeros(b Bitboard) int {
	return bits.LeadingZeros64(uint64(b))
}

// NOTE: returns 64 for empty bitboard
func lsb(b Bitboard) int {
	return trailing_zeros(b)
}

func make_piece(s string) Piece {
	return Piece(strings.Index(PIECE_STRING, s))
}

func make_square(r int, f int) Square {
	return Square(r*8 + f)
}

// NOTE: returns -1 for empty bitboard
func msb(b Bitboard) int {
	return 63 - leading_zeros(b)
}

func neighbors(b Bitboard) Bitboard {
	return shift_direction(b, WEST) | shift_direction(b, EAST)
}

func occupied_at_bb(b Bitboard, sq Bitboard) bool {
	return b&sq != 0
}

func occupied_at_sq(b Bitboard, sq Square) bool {
	// return b & SQUARE_BBS[sq] != 0
	return (b>>sq)&1 == 1
}

func occupied_squares(pieces []Bitboard) Bitboard {
	var occupied Bitboard
	for _, piece := range pieces {
		occupied |= piece
	}
	return occupied
}

func occupied_squares_by_color(pieces []Bitboard, color Color) Bitboard {
	var occupied Bitboard
	for _, piece := range piece_range_by_color(color) {
		occupied |= pieces[piece]
	}
	return occupied
}

func on_board(sq Square) bool {
	return SQ_A1 <= sq && sq <= SQ_H8
}

func opposite(c Color) Color {
	return c ^ 1
}

func piece_range_by_color(color Color) []Piece {
	if color == WHITE {
		return WHITE_PIECES
	}
	return BLACK_PIECES
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

func second_rank(dir int) int {
	if dir == NORTH {
		return 1
	}
	return 6
}

func set_square(original Bitboard, sq Square) Bitboard {
	return original | SQUARE_BBS[sq]
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
	if s == NULL_SQ {
		return "-"
	}
	return fmt.Sprintf("%c%d", FILES[square_file(s)], square_rank(s)+1)
}

func to_square(s string) Square {
	if s == "-" {
		return NULL_SQ
	}
	return make_square(int(s[1]-'0')-1, strings.Index(FILES, s[0:1]))
}

func trailing_zeros(b Bitboard) int {
	return bits.TrailingZeros64(uint64(b))
}
