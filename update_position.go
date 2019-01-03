package main

import (
	// "fmt"
	"strconv"
)

func (p *Position) do_castle(do bool, king_dst Square) {
	rook_src_dst := ROOK_SRC_DST[king_dst]
	src, dst := rook_src_dst[0], rook_src_dst[1]
	if do {
		r := p.piece_at(src)
		p.remove_piece(r, src)
		p.place_piece(r, dst)
	} else {
		r := p.piece_at(dst)
		p.remove_piece(r, dst)
		p.place_piece(r, src)
	}
}

func (p *Position) do_move(m Move, new_state *StateInfo) {
	src, dst := move_src(m), move_dst(m)
	mover := p.piece_at(src)
	us, them := p.side_to_move(), opposite(p.side_to_move())
	// Update new state
	new_state.castling_rights = update_castling_right(p.state.castling_rights, src)
	new_state.ep_sq = update_ep_sq(m, p.placement[pt_to_p(PAWN, them)])
	new_state.rule_50 = update_rule_50(p.state.rule_50, m, p.piece_type_at(src))

	// Update placement
	if is_castle(m) {
		p.do_castle(true, dst)
	} else if is_enpassant(m) {
		ep_sq := cleanup_sq_for_ep_capture(dst)
		captured := p.piece_at(ep_sq)
		new_state.captured = captured
		p.remove_piece(captured, ep_sq)
	} else if is_capture(m) {
		captured := p.piece_at(dst)
		new_state.captured = captured
		p.remove_piece(captured, dst)
	}

	p.remove_piece(mover, src)
	if is_promotion(m) {
		p.do_promotion(us, move_type(m), dst)
	} else {
		p.place_piece(mover, dst)
	}

	// Update king blockers (for next turn)
	// our sliders, their king
	new_state.opposite_color_attacks = p.get_color_attacks(us)
	new_state.blockers_for_king = p.slider_blockers(us, p.king_square(them))

	// Reassign state
	new_state.prev = p.state
	p.state = new_state

	// Update position
	p.ply += 1
}

func (p *Position) do_promotion(c Color, mt MoveType, dst Square) {
	pt := move_type_to_promotion_type(mt)
	p.place_piece(pt_to_p(pt, c), dst)
}

func (p *Position) place_piece(pc Piece, sq Square) {
	p.placement_by_square[sq] = pc
	p.placement[pc] |= SQUARE_BBS[sq]
}

func (p *Position) remove_piece(pc Piece, sq Square) {
	p.placement_by_square[sq] = NULL_PIECE
	p.placement[pc] &^= SQUARE_BBS[sq]
}

func (p *Position) set_fen_info(positions string, color string, castles string, enps string, rule_50 string, move_count int) {
	rule_50_int, _ := strconv.Atoi(rule_50)
	// On Position
	p.placement = parse_positions(positions)
	p.placement_by_square = placement_to_placement_by_square(p.placement)
	p.ply = move_count*2 + int(parse_color(color))
	// On StateInfo
	p.state.castling_rights = make_castle_state_info(castles)
	p.state.ep_sq = parse_square(enps)
	p.state.rule_50 = rule_50_int
}

func (p *Position) undo_move(m Move) {
	src, dst := move_src(m), move_dst(m)
	mover := p.piece_at(dst)

	// turn has already been updated
	us := opposite(p.side_to_move())
	// Update position
	p.ply -= 1

	// move piece back to src
	p.remove_piece(mover, dst)
	if is_promotion(m) {
		p.place_piece(pt_to_p(PAWN, us), src)
	} else {
		p.place_piece(mover, src)
	}

	// if capture, replace piece
	if is_capture(m) {
		capsq := dst
		if is_enpassant(m) {
			capsq = cleanup_sq_for_ep_capture(dst)
		}
		p.place_piece(p.state.captured, capsq)
	}

	// if castle, undo castle
	if is_castle(m) {
		p.do_castle(false, dst)
	}

	// Reassign state
	p.state = p.state.prev
}
