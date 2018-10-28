package main

import (
  // "fmt"
	"strings"
)

func parse_positions(positions string) [12]Bitboard {
	var result_bbs [12]Bitboard
	ranks := strings.Split(positions, "/")
	for rank, rank_string := range ranks {
		offset := 0
		for file, sq := range rank_string {
			if strings.Contains(PIECE_STRING, string(sq)) {
				index := make_piece(string(sq))
        sq_num := make_square(RANK_8-rank, file+offset)
				result_bbs[index] = result_bbs[index] ^ SQUARE_BBS[sq_num]
			} else {
				offset += int(sq-'0') - 1
			}
		}
	}
	return result_bbs
}

func bitboards_to_grid(bitboards [12]Bitboard) [8][8]string {
	var grid [8][8]string
  for _, s := range SQUARES {
    for _, piece := range PIECES {
      if occupied_at(bitboards[piece], s) {
        grid[square_rank(s)][square_file(s)] = piece.String()
      }
    }
  }
	return grid
}
