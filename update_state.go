package main

import (
	"strconv"
)

func (s *StateInfo) initialize_from_fen(castles string, enps string, rule_50 string) {
	rule_50_int, _ := strconv.Atoi(rule_50)
	s.castling_rights = make_castle_state_info(castles)
	s.ep_sq = parse_square(enps)
	s.prev = nil
	s.rule_50 = rule_50_int
}
