package main

import (
// "fmt"
// "sync"
)

type TT struct {
	// lock sync.RWMutex
	m map[Key]*TTEntry
}

type TTEntry struct {
	depth int
	key   Key
	score int
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
