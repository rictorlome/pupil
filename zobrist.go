package main

import (
	"math/rand"
)

var ZOBRIST_PSQ [64][12]Key
var ZOBRIST_CSTL [16]Key
var ZOBRIST_EPSQ [65]Key // room for null sq. this really only needs the files.
var ZOBRIST_SIDE Key

func init_zobrists() {

	for _, sq := range SQUARES {
		for _, p := range PIECES {
			ZOBRIST_PSQ[sq][p] = Key(rand.Uint64())
		}
		ZOBRIST_EPSQ[sq] = Key(rand.Uint64())
	}

	ZOBRIST_EPSQ[NULL_SQ] = Key(rand.Uint64())
	for cstl := 0; cstl <= ALL_CASTLES; cstl++ {
		ZOBRIST_CSTL[cstl] = Key(rand.Uint64())
	}
	ZOBRIST_SIDE = Key(rand.Uint64())
}
