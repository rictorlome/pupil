package main

import (
	// "fmt"
	"strconv"
)

func (p *Position) do_enpassant(pawn_dst Square) {
	ep_sq := cleanup_sq_for_ep_capture(pawn_dst)
	p.remove_piece(p.piece_at(ep_sq), ep_sq)
}

func (p *Position) do_promotion(c Color, mt MoveType, dst Square) {
	pt := move_type_to_promotion_type(mt)
	p.place_piece(pt_to_p(pt, c), dst)
}

func (p *Position) clear_sq(sq Square) {
	for _, pc := range PIECES {
		p.remove_piece(pc, sq)
	}
}

func (p *Position) do_castle(king_dst Square) {
	rook_src_dst := ROOK_SRC_DST[king_dst]
	r := p.piece_at(rook_src_dst[0])
	p.remove_piece(r, rook_src_dst[0])
	p.place_piece(r, rook_src_dst[1])
}

func (p *Position) do_move(m Move, new_state StateInfo) {
	src, dst := move_src(m), move_dst(m)
	mover := p.piece_at(src)
	// Update new state
	new_state.castling_rights = update_castling_right(p.state.castling_rights, src)
	new_state.ep_sq = update_ep_sq(m, p.placement[pt_to_p(PAWN, opposite(p.to_move))])
	new_state.rule_50 = update_rule_50(p.state.rule_50, m, p.piece_type_at(src))

	// Update placement
	if is_castle(m) {
		p.do_castle(dst)
	} else if is_enpassant(m) {
		p.do_enpassant(dst)
	} else if is_capture(m) {
		p.remove_piece(p.piece_at(dst), dst)
	}

	p.remove_piece(mover, src)
	if is_promotion(m) {
		p.do_promotion(p.to_move, move_type(m), dst)
	} else {
		p.place_piece(mover, dst)
	}

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
