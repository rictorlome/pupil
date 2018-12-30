package main

import (
  "testing"
)

func TestProfileSearch(t *testing.T) {
  pos := parse_fen(INITIAL_FEN)
  build_tree_parallel(&pos, 5)
  // build_tree_recursive(&pos, nil, Move(0), 5)
}
