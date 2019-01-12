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
    perft{0,1,0,0,0,0,0,0},
    perft{1,20,0,0,0,0,0,0},
    perft{2,400,0,0,0,0,0,0},
    perft{3,8902,34,0,0,0,12,0},
    perft{4,197281,1576,0,0,0,469,8},
    perft{5,4865609,82719,258,0,0,27351,347},
    perft{6,119060324,2812008,5248,0,0,809099,10823},
  },
}

var secondary_perft = perft_seq {
  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
  []perft{
    perft{0,1,0,0,0,0,0,0},
    perft{1,48,8,0,2,0,0,0},
    perft{2,2039,351,1,91,0,3,0},
    perft{3,97862,17102,45,3162,0,993,1},
    perft{4,4085603,757163,1929,128013,15172,25523,43},
    // perft{5,193690690,35043416,73365,4993637,8392,3309887,30171},
  },
}

func TestSearch(t *testing.T) {
    pos := parse_fen(secondary_perft.start_fen)
    for _, expected_pft := range secondary_perft.perfts {
      actual_pft := get_perft_parallel(&pos, expected_pft.depth)
      // actual_pft := get_perft_recursive(&pos, expected_pft.depth, Move(0))
      if actual_pft != expected_pft {
          t.Error(fmt.Sprintf("\nExpected:\t %v\n Got:\t\t %v", expected_pft, actual_pft))
      }
    }
}
