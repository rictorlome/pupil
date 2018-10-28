package main

import (
	"fmt"
)

func has_bit(s StateInfo, bit uint) bool {
	return (s>>bit)&1 == 1
}

func make_castle_state_info(available string) StateInfo {
	var castles StateInfo
	for _, char := range available {
		castles |= CHAR_TO_CASTLE[string(char)]
	}
	return castles
}

func make_enpassant_square_info(sq Square) StateInfo {
	return StateInfo(sq << 4)
}

func make_rule_50(i int) StateInfo {
	return StateInfo(i << 10)
}

func (s StateInfo) String() string {
	return fmt.Sprintf("%016b", s)
}
