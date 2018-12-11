package main

import (
	// "fmt"
	"strconv"
)

func (p *Position) clear_sq(sq Square) {
	for _, pc := range PIECES {
		p.remove_piece(pc, sq)
	}
}

func (p *Position) do_move(m Move, new_state StateInfo) {
	src, dst, move_type := move_src(m), move_dst(m), move_type(m)
	mover := p.piece_at(src)
	// Update new state
	new_state.castling_rights = update_castling_right(p.state.castling_rights, src)
	new_state.ep_sq = update_ep_sq(m, p.placement[pt_to_p(PAWN, opposite(p.to_move))])
	new_state.rule_50 = update_rule_50(p.state.rule_50, m, p.piece_type_at(src))

	// Update placement
	if move_type == CAPTURE {
		p.remove_piece(p.piece_at(dst), dst)
	}
	p.remove_piece(mover, src)
	p.place_piece(mover, dst)

	// Reassign state
	new_state.prev = &p.state
	p.state = new_state

	// Update position
	p.move_count += 1 * int(p.to_move&1)
	p.to_move = opposite(p.to_move)
}

func (p *Position) place_piece(pc Piece, sq Square) {
	p.placement[pc] |= SQUARE_BBS[sq]
}

func (p *Position) remove_piece(pc Piece, sq Square) {
	p.placement[pc] &^= SQUARE_BBS[sq]
}

func (p *Position) set_fen_info(positions string, color string, castles string, enps string, rule_50 string, move_count int) {
	rule_50_int, _ := strconv.Atoi(rule_50)
	// On Position
	p.move_count = move_count
	p.placement = parse_positions(positions)
	p.to_move = parse_color(color)
	// On StateInfo
	p.state.castling_rights = make_castle_state_info(castles)
	p.state.ep_sq = parse_square(enps)
	p.state.rule_50 = rule_50_int
}
