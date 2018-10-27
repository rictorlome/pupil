package main

import "fmt"

func setup() {
	// Initialize Square Bitboards
	var x Bitboard = 0x1
	for s := SQ_A1; s <= SQ_H8; s++ {
		SQUARE_BBS[s] = x << s
	}
}

func main() {
	fmt.Println("OK")
	setup()

	for _, dir := range DIRECTIONS {
		fmt.Println(signed_shift(SQUARE_BBS[SQ_E4], dir))
	}
}
