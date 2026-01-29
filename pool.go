package main

import (
	"sync"
)

var siPool *sync.Pool
var moveListPool *sync.Pool

func initPool() {
	siPool = &sync.Pool{
		New: func() interface{} {
			return new(StateInfo)
		},
	}
	moveListPool = &sync.Pool{
		New: func() interface{} {
			slice := make([]Move, 0, MAX_BRANCHING)
			return &slice
		},
	}
}

func getMoveList() *[]Move {
	ml := moveListPool.Get().(*[]Move)
	*ml = (*ml)[:0] // Reset length but keep capacity
	return ml
}

func putMoveList(ml *[]Move) {
	moveListPool.Put(ml)
}
