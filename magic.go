package main

import (
	"math/rand"
)

// https://www.chessprogramming.org/Looking_for_Magics#Feeding_in_Randoms
func rand_few_bits() uint64 {
	return rand.Uint64() & rand.Uint64() & rand.Uint64()
}

func transform(b Bitboard, magic uint64, bits int) int {
  return 0
}
var b, a, used = [4096]uint64
func find_magic(sq Square, bits int, bishop bool) Bitboard {
  var mask Bitboard

  mask = ROOK_OCCUPANCY_MASKS[sq]
  if bishop {
    mask = BISHOP_OCCUPANCY_MASKS[sq]
  }
  n = popcount(mask)

	// temp
	return 0
}
