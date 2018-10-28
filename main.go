package main

import "fmt"

func init() {
	// Initialize Square Bitboards
	var x Bitboard = 0x1
	for _, s := range SQUARES {
		SQUARE_BBS[s] = x << s
	}
}

func main() {
	fmt.Println("OK")
	grid := bitboards_to_grid(parse_positions(INITIAL_FEN_JUST_PIECES))
	for _, row := range grid {
		fmt.Println(row)
	}
}
