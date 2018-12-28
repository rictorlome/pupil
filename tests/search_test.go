package main

import (
  "testing"
  "fmt"
)

type perft_seq struct {
  start_fen string
  perfts []perft
}
//tables taken from https://chessprogramming.wikispaces.com/Perft%20Results
//note no castles or promotions
var initial_perft = perft_seq {
  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
  []perft{
    // perft{0,1,0,0,0,0,0,0},
    // perft{1,20,0,0,0,0,0,0},
    // perft{2,400,0,0,0,0,0,0},
    // perft{3,8902,34,0,0,0,12,0},
    // perft{4,197281,1576,0,0,0,469,8},
    // perft{5,4865609,82719,258,0,0,27351,347},
    // perft{6,119060324,2812008,5248,0,0,809099,10823},
  },
}

func TestSearch(t *testing.T) {
    pos := parse_fen(initial_perft.start_fen)
    for _, expected_pft := range initial_perft.perfts {
      actual_pft := get_perft(pos, expected_pft.depth)
      if actual_pft != expected_pft {
          t.Error(fmt.Sprintf("\nExpected:\n%v\nGot:\n%v", expected_pft, actual_pft))
      }
    }
}
