package main

import (
  "testing"
  // "fmt"
  // "net/http/pprof"
)

func TestProfileSearch(t *testing.T) {
  pos := parse_fen(INITIAL_FEN)
  build_tree(&pos, nil, Move(0), 5)
}
