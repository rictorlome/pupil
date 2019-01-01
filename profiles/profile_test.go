package main

import (
  "testing"
)

func TestProfileSearch(t *testing.T) {
  pos := parse_fen(INITIAL_FEN)
  get_perft_parallel(&pos, 6)
  // build_tree_recursive(&pos, nil, Move(0), 5)
}
