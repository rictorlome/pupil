package main

import (
	"fmt"
	"math/rand"
)

type Magic struct {
	mask    Bitboard
	magic   Bitboard
	// offset in the attack array
	offset  uint
	// the number of indices in the attack table for this square
	size 		uint
	shift   uint
}

func attack_index(m Magic, occ Bitboard) uint {
	return uint(((occ & m.mask) * m.magic) >> m.shift)
}

func magic_rook_attack(occ Bitboard, sq Square) Bitboard {
	m := RookMagics[sq]
	return RookAttackTable[attack_index(m, occ)]
}

var RookAttackTable = make([]Bitboard, 0x19000)
var BishopAttackTable = make([]Bitboard, 0x1480)

var RookMagics [64]Magic
var BishopMagics [64]Magic

func find_magic(sq Square, bishop bool) Magic {
	m := Magic{}
	var occupancy, reference, used [4096]Bitboard

	// Rook
	attack_table := RookAttackTable
	magics := RookMagics
	attack_func := rook_attacks
	// attack_mask := ROOK_ATTACK_MASKS[sq]
	m.mask = RELEVANT_ROOK_OCCUPANCY[sq]
	if bishop {
		magics = BishopMagics
		attack_table = BishopAttackTable
		attack_func = bishop_attacks
		// attack_mask = BISHOP_ATTACK_MASKS[sq]
		m.mask = RELEVANT_BISHOP_OCCUPANCY[sq]
	}


	m.shift = uint(64 - popcount(m.mask))
	if sq == SQ_A1 {
		m.offset = 0
	} else {
		m.offset = magics[sq - 1].offset + magics[sq - 1].size
	}

	b, size := Bitboard(0), uint(0)
	for first := true; first || b != 0; b = (b - m.mask) & m.mask {
		first = false
		occupancy[size] = b
		reference[size] = attack_func(b, sq)
		// fmt.Println("size", size, "offset", m.offset)
		// attack_table[m.offset + size] = reference[size]
		size++
	}
	m.size = size

	OUTER:
	for i := 0; i < 100000000; i++ {
		m.magic = rand_few_bits()
		// fmt.Println(uint64(m.magic))
		if popcount((m.magic * m.mask) >> 56) < 6 {
			continue
		}
		for i := uint(0); i < size; i++ {
			// find your index
			idx := attack_index(m, occupancy[i])
			// fmt.Println(idx < 4096)
			// if you haven't used, use it
			if used[idx] == 0 {
				attack_table[m.offset + idx] = reference[i]
				used[idx] = attack_table[m.offset + idx]
			} else if used[idx] != attack_table[m.offset + idx] {
				continue OUTER
			}
		}
		return m
	}
	fmt.Println("FAILED!")
	return m
}

// https://www.chessprogramming.org/Looking_for_Magics#Feeding_in_Randoms
func rand_few_bits() Bitboard {
	return Bitboard(rand.Uint64() & rand.Uint64() & rand.Uint64())
}
