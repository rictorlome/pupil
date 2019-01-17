package main

import (
	"math/rand"
	"errors"
)

type Magic struct {
	magic Bitboard
	mask  Bitboard
	// offset in the attack array
	offset uint
	// the number of indices in the attack table for this square
	size  uint
	shift uint
}

func attack_index(m *Magic, occ Bitboard) uint {
	return uint(((occ & m.mask) * m.magic) >> m.shift)
}

func attack_index_with_offset(m *Magic, occ Bitboard) uint {
	return m.offset + uint(((occ&m.mask)*m.magic)>>m.shift)
}

func find_magic(sq Square, bishop bool) (Magic, error) {
	m := Magic{}
	var occupancy, reference [4096]Bitboard
	// Rook by default
	magics := RookMagics
	attack_table := RookAttackTable
	attack_func := slider_rook_attacks
	m.mask = RELEVANT_ROOK_OCCUPANCY[sq]
	// Use precomputed values for speedup
	m.magic = Bitboard(RookMagicNums[sq])
	if bishop {
		magics = BishopMagics
		attack_table = BishopAttackTable
		attack_func = slider_bishop_attacks
		m.mask = RELEVANT_BISHOP_OCCUPANCY[sq]
		// Use precomputed value for speedup
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
			idx := attack_index(&m, occupancy[j])
			if used[idx] == 0 {
				used[idx] = reference[j]
				attack_table[m.offset+idx] = used[idx]
			} else if used[idx] != reference[j] {
				m.magic = rand_few_bits()
				continue OUTER
			}
		}
		return m, nil
	}
	return m, errors.New("no magic")
}


func init_magics() (err error) {
	for _, sq := range SQUARES {
		BishopMagics[sq], err = find_magic(sq, true)
		if err != nil {
			return err
		}
		RookMagics[sq], err = find_magic(sq, false)
		if err != nil {
			return err
		}
	}
	return nil
}


// https://www.chessprogramming.org/Looking_for_Magics#Feeding_in_Randoms
func rand_few_bits() Bitboard {
	return Bitboard(rand.Uint64() & rand.Uint64() & rand.Uint64())
}
