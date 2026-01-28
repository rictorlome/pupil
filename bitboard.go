package main

import (
	"fmt"
	"math/bits"
	"strings"
)

func aligned(sq1 Square, sq2 Square, sq3 Square) bool {
	return LINE_BBS[sq1][sq2]&SQUARE_BBS[sq3] != 0
}

func between(sq1 Square, sq2 Square, sq3 Square) bool {
	return BETWEEN_BBS[sq1][sq2]&SQUARE_BBS[sq3] != 0
}

func binary(b Bitboard) string {
	return fmt.Sprintf("%064b", uint64(b))
}

func binary16(u uint16) string {
	return fmt.Sprintf("%016b", u)
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

func indicator(b bool) int {
	if b {
		return 1
	}
	return 0
}

func lastRank(color Color) int {
	if color == WHITE {
		return 7
	}
	return 0
}

// NOTE: returns 64 for empty bitboard
func lsb(b Bitboard) int {
	return bits.TrailingZeros64(uint64(b))
}

func makePiece(s string) Piece {
	return Piece(strings.Index(PIECE_STRING, s))
}

func makeSquare(r int, f int) Square {
	return Square(r*8 + f)
}

// NOTE: returns -1 for empty bitboard
func msb(b Bitboard) int {
	return 63 - bits.LeadingZeros64(uint64(b))
}

func neighbors(b Bitboard) Bitboard {
	return shiftDirection(b, WEST) | shiftDirection(b, EAST)
}

func occupiedAtBB(b Bitboard, sq Bitboard) bool {
	return b&sq != 0
}

func occupiedAtSq(b Bitboard, sq Square) bool {
	return (b>>sq)&1 == 1
}

func occupiedSquares(pieces []Bitboard) Bitboard {
	var occupied Bitboard
	for _, piece := range pieces {
		occupied |= piece
	}
	return occupied
}

func occupiedSquaresByColor(pieces []Bitboard, color Color) Bitboard {
	var occupied Bitboard
	for _, piece := range pieceRangeByColor(color) {
		occupied |= pieces[piece]
	}
	return occupied
}

func onBoard(sq Square) bool {
	return SQ_A1 <= sq && sq <= SQ_H8
}

func opposite(c Color) Color {
	return c ^ 1
}

func pieceRangeByColor(color Color) []Piece {
	if color == WHITE {
		return WHITE_PIECES
	}
	return BLACK_PIECES
}

func popcount(b Bitboard) int {
	return bits.OnesCount64(uint64(b))
}

func possibleEnpassantSq(sq Square) bool {
	return squareRank(sq) == 2 || squareRank(sq) == 5
}

// https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func secondRank(dir int) int {
	if dir == NORTH {
		return 1
	}
	return 6
}

func setSquare(original Bitboard, sq Square) Bitboard {
	return original | SQUARE_BBS[sq]
}

// shifts bitboard in a direction, prevents wrapping
func shiftDirection(b Bitboard, direction int) Bitboard {
	if direction == NORTH_EAST || direction == EAST || direction == SOUTH_EAST {
		return signedShift(b&^FILE_HBB, direction)
	}
	if direction == NORTH_WEST || direction == WEST || direction == SOUTH_WEST {
		return signedShift(b&^FILE_ABB, direction)
	}
	return signedShift(b, direction)
}

// shifts bitboard by an amount (which can be a positive or negative number)
func signedShift(b Bitboard, amount int) Bitboard {
	if amount >= 0 {
		return b << uint(amount)
	}
	return b >> uint(-amount)
}

func squareFile(sq Square) int {
	return int(sq % 8)
}

func squareRank(sq Square) int {
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

func (c Color) String() string {
	if c == WHITE {
		return "w"
	}
	return "b"
}

func (p Piece) String() string {
	return string([]rune(PIECE_STRING)[p])
}

func (s Square) String() string {
	if s == NULL_SQ {
		return "-"
	}
	return fmt.Sprintf("%c%d", FILES[squareFile(s)], squareRank(s)+1)
}

func toSquare(s string) Square {
	if s == "-" {
		return NULL_SQ
	}
	return makeSquare(int(s[1]-'0')-1, strings.Index(FILES, s[0:1]))
}

func twoUp(src Square, color Color) Square {
	if color == WHITE {
		return Square(int(src) + 16)
	}
	return Square(int(src) - 16)
}
