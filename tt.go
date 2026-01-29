package main

import (
	// "fmt"
	"math"
	// "sync"
)

type TT struct {
	// lock sync.RWMutex
	m []*TTEntry
}

// alpha <= score <= beta
// all moves explored
// score is EXACT
var PV_NODE uint8 = 0

// score >= beta, beta cutoff
// fail-high nodes
// returned score <= exact score
var CUT_NODE uint8 = 1

// fail-low nodes
// all moves explored
// no move's score exceeds alpha
// score <= alpha
// exact score <= returned score
var ALL_NODE uint8 = 2

type TTEntry struct {
	bestMove Move
	depth    uint8
	key      Key
	nodeType uint8
	score    int
}

func createTT(capExp int) *TT {
	exp := float64(capExp)
	return &TT{m: make([]*TTEntry, int(math.Pow(2, exp)))}
}

func (t *TT) read(key Key) (*TTEntry, bool) {
	// t.lock.RLock()

	// modulo 2^n with &
	idx := key & Key((cap(t.m) - 1))
	entry := t.m[idx]
	// t.lock.RUnlock()
	return entry, entry == nil
}

func (t *TT) write(key Key, entry *TTEntry) {
	// t.lock.Lock()
	idx := key & Key((cap(t.m) - 1))
	t.m[idx] = entry
	// t.lock.Unlock()
}

func (t *TT) clear() {
	for i := range t.m {
		t.m[i] = nil
	}
}
