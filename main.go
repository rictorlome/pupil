package main

import "fmt"

func init() {
	// Initialize Square Bitboards
	// Initialize En Passant Sqs
	var x Bitboard = 0x1
	for _, s := range SQUARES {
		SQUARE_BBS[s] = x << s
		ENPASSANT_SQS[s] = make_enpassant_square_info(s)
	}
}

func main() {
	fmt.Println("OK")
	fmt.Println(generate_fen(parse_fen(INITIAL_FEN)))
	// fmt.Println(BQ_CASTLE)
}
