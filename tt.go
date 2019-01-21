package main

import (
  // "fmt"
  "sync"
)

type TT struct {
  lock sync.RWMutex
  m    map[Key]*TTEntry
}



type TTEntry struct {
  // score int
  p perft
}

func createTT() *TT {
  return &TT{m: make(map[Key]*TTEntry)}
}

func (t *TT) read(key Key) (*TTEntry, bool) {
  t.lock.RLock()
  entry, ok := t.m[key]
  t.lock.RUnlock()
  return entry, ok
}

func (t *TT) write(key Key, entry *TTEntry) {
  t.lock.Lock()
  t.m[key] = entry
  t.lock.Unlock()
}
