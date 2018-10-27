package main

import (
	"fmt"
	"math/bits"
	"strings"
	// "strconv"
)

func binary(b Bitboard) string {
	return fmt.Sprintf("%v%b", strings.Repeat("0", leading_zeros(b)), uint64(b))
}

func leading_zeros(b Bitboard) int {
	return bits.LeadingZeros64(uint64(b))
}

func make_square(r uint, f uint) Square {
	return Square(r*f - 1)
}

func popcount(b Bitboard) int {
	return bits.OnesCount64(uint64(b))
}

// https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// shifts bitboard by an amount (which can be a positive or negative number)
func signed_shift(b Bitboard, amount int) Bitboard {
	if amount >= 0 {
		return b << uint(amount)
	}
	return b >> uint(-amount)
}

func (b Bitboard) String() string {
	var s string
	runes := []rune(binary(b))
	for i := RANK_1; i <= RANK_8; i++ {
		s += reverse(string(runes[(8*i):(8*(i+1))])) + "\n"
	}
	return s
}

func trailing_zeros(b Bitboard) int {
	return bits.TrailingZeros64(uint64(b))
}
