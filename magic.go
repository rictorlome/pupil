package main

import (
	"fmt"
	"math/rand"
)

type Magic struct {
	mask  Bitboard
	magic Bitboard
	// offset in the attack array
	offset uint
	// the number of indices in the attack table for this square
	size  uint
	shift uint
}

func attack_index(m Magic, occ Bitboard) uint {
	return uint(((occ & m.mask) * m.magic) >> m.shift)
}

func magic_rook_attack(occ Bitboard, sq Square) Bitboard {
	m := RookMagics[sq]
	return RookAttackTable[m.offset+attack_index(m, occ)]
}

func magic_bishop_attack(occ Bitboard, sq Square) Bitboard {
	m := BishopMagics[sq]
	return BishopAttackTable[m.offset+attack_index(m, occ)]
}

var RookAttackTable = make([]Bitboard, 0x19000)
var BishopAttackTable = make([]Bitboard, 0x1480)

var RookMagics [64]Magic
var BishopMagics [64]Magic

func find_magic(sq Square, bishop bool) Magic {
	m := Magic{}
	var occupancy, reference [4096]Bitboard
	// Rook by default
	magics := RookMagics
	attack_table := RookAttackTable
	attack_func := slider_rook_attacks
	m.mask = RELEVANT_ROOK_OCCUPANCY[sq]
	m.magic = Bitboard(RookMagicNums[sq])
	if bishop {
		magics = BishopMagics
		attack_table = BishopAttackTable
		attack_func = slider_bishop_attacks
		m.mask = RELEVANT_BISHOP_OCCUPANCY[sq]
		m.magic = Bitboard(BishopMagicNums[sq])
	}

	m.shift = uint(64 - popcount(m.mask))
	if sq == SQ_A1 {
		m.offset = 0
	} else {
		m.offset = magics[sq-1].offset + magics[sq-1].size
	}

	occ, size := Bitboard(0), uint(0)
	for first := true; first || occ != 0; occ = (occ - m.mask) & m.mask {
		first = false
		occupancy[size] = occ
		reference[size] = attack_func(occ, sq)
		size++
	}
	m.size = size

OUTER:
	for i := 0; i < 100000000; i++ {
		var used [4096]Bitboard
		if popcount((m.magic*m.mask)>>56) < 6 {
			m.magic = rand_few_bits()
			continue
		}
		for j := uint(0); j < size; j++ {
			idx := attack_index(m, occupancy[j])
			if used[idx] == 0 {
				used[idx] = reference[j]
				attack_table[m.offset+idx] = used[idx]
			} else if used[idx] != reference[j] {
				m.magic = rand_few_bits()
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
