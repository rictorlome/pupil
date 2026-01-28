package main

import (
	"testing"
)

func TestProfileSearch(t *testing.T) {
	pos := parseFen(INITIAL_FEN)
	getPerftParallel(&pos, 6)
	// buildTreeRecursive(&pos, nil, Move(0), 5)
}
