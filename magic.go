package main

import (
	"math/rand"
)

type Magic struct {
	mask    Bitboard
	magic   Bitboard
	attacks *Bitboard
	shift   uint
}

func attack_index(m Magic, occ Bitboard) uint {
	return uint(((occ & m.mask) * m.magic) >> m.shift)
}

var RookAttackTable = make([]Bitboard, 0, 0x19000)
var BishopAttackTable = make([]Bitboard, 0, 0x1480)

var RookMagics [64]Magic
var BishopMagics [64]Magic

// 1. Loop over squares
// a.
// func init_magics(bishop bool, magics []Magic, dir []int) {
// 	var occupancy, reference [4096]Bitboard
// 	var b Bitboard
//
// 	table := RookAttackTable
// 	if bishop {
// 		table = BishopAttackTable
// 	}
//
//
// 	for _, s := range SQUARES {
// 		m := Magic{}
//
// 		m.mask = RELEVANT_ROOK_OCCUPANCY[s]
// 		if bishop {
// 			m.mask = RELEVANT_BISHOP_OCCUPANCY[s]
// 		}
// 		m.shift = uint(64 - popcount(m.mask))
// 		magics[s] = m
//
//
// 	}
// }

// https://www.chessprogramming.org/Looking_for_Magics#Feeding_in_Randoms
func rand_few_bits() uint64 {
	return rand.Uint64() & rand.Uint64() & rand.Uint64()
}

func transform(b Bitboard, magic uint64, bits int) int {
	return 0
}

var b, a, used [4096]uint64

func find_magic(sq Square, bits int, bishop bool) Bitboard {
	var mask Bitboard

	mask = ROOK_OCCUPANCY_MASKS[sq]
	if bishop {
		mask = BISHOP_OCCUPANCY_MASKS[sq]
	}
	// n := popcount(mask)

	// temp
	return mask
}
