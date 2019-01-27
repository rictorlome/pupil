package main

import (
// "fmt"
// "sync"
)

type TT struct {
	// lock sync.RWMutex
	m map[Key]*TTEntry
}

// alpha <= score <= beta
// all moves explored
// score is EXACT
var PV_NODE int = 0

// score >= beta, beta cutoff
// fail-high nodes
// returned score <= exact score
var CUT_NODE int = 1

// fail-low nodes
// all moves explored
// no move's score exceeds alpha
// score <= alpha
// exact score <= returned score
var ALL_NODE int = 2

type TTEntry struct {
	best_move Move
	depth     int
	key       Key
	node_type int
	score     int
}

func createTT() *TT {
	return &TT{m: make(map[Key]*TTEntry)}
}

func createTTs(depth int) []*TT {
	TTS := make([]*TT, depth, depth)
	for i := 0; i <= depth; i++ {
		TTS[i] = createTT()
	}
	return TTS
}

func (t *TT) read(key Key) (*TTEntry, bool) {
	// t.lock.RLock()
	entry, ok := t.m[key]
	// t.lock.RUnlock()
	return entry, ok
}

func (t *TT) write(key Key, entry *TTEntry) {
	// t.lock.Lock()
	t.m[key] = entry
	// t.lock.Unlock()
}
