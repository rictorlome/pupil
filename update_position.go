package main

import (
	// "fmt"
	"strconv"
)

func (p *Position) do_castle(do bool, r Piece, r_src Square, r_dst Square) {
	if do {
		p.remove_piece(r, r_src)
		p.place_piece(r, r_dst)
	} else {
		p.remove_piece(r, r_dst)
		p.place_piece(r, r_src)
	}
}

func (p *Position) do_move(m Move, new_state *StateInfo) {
	src, dst := move_src(m), move_dst(m)
	mover := p.piece_at(src)
	us, them := p.side_to_move(), opposite(p.side_to_move())
	// Initialze new key (cancel current castling rights and ep sq to avoid conditional)
	new_key := p.state.key ^ ZOBRIST_SIDE ^ ZOBRIST_CSTL[p.state.castling_rights] ^ ZOBRIST_EPSQ[p.state.ep_sq]
	// Update new state
	new_state.castling_rights = update_castling_right(p.state.castling_rights, src, dst)
	new_state.ep_sq = update_ep_sq(m, p.placement[pt_to_p(PAWN, them)])
	new_state.rule_50 = update_rule_50(p.state.rule_50, m, p.piece_type_at(src))
	// Update new key
	new_key ^= ZOBRIST_CSTL[new_state.castling_rights] ^ ZOBRIST_EPSQ[new_state.ep_sq]
	// Update placement
	if is_castle(m) {
		rook_src_dst := ROOK_SRC_DST[dst]
		r_src, r_dst := rook_src_dst[0], rook_src_dst[1]
		r := p.piece_at(r_src)
		p.do_castle(true, r, r_src, r_dst)
		new_key ^= ZOBRIST_PSQ[r_src][r] ^ ZOBRIST_PSQ[r_dst][r]
	} else if is_enpassant(m) {
		ep_sq := cleanup_sq_for_ep_capture(dst)
		captured := p.piece_at(ep_sq)
		new_state.captured = captured
		p.remove_piece(captured, ep_sq)
		new_key ^= ZOBRIST_PSQ[ep_sq][captured]
	} else if is_capture(m) {
		captured := p.piece_at(dst)
		new_state.captured = captured
		p.remove_piece(captured, dst)
		new_key ^= ZOBRIST_PSQ[dst][captured]
	}

	p.remove_piece(mover, src)
	new_key ^= ZOBRIST_PSQ[src][mover]
	if is_promotion(m) {
		pt := move_type_to_promotion_type(move_type(m))
		pc := pt_to_p(pt, us)
		p.place_piece(pc, dst)
		new_key ^= ZOBRIST_PSQ[dst][pc]
	} else {
		p.place_piece(mover, dst)
		new_key ^= ZOBRIST_PSQ[dst][mover]
	}

	// Update king blockers (for next turn)
	// our sliders, their king
	new_state.key = new_key
	new_state.opposite_color_attacks = p.get_color_attacks(us)
	new_state.blockers_for_king = p.slider_blockers(us, p.king_square(them))

	// Reassign state
	new_state.prev = p.state
	p.state = new_state

	// Update position
	p.ply += 1
	p.stm = opposite(p.stm)
}

func (p *Position) place_piece(pc Piece, sq Square) {
	sq_bb := SQUARE_BBS[sq]
	p.occ |= sq_bb
	p.placement_by_square[sq] = pc
	p.placement[pc] |= sq_bb
}

func (p *Position) remove_piece(pc Piece, sq Square) {
	sq_bb := SQUARE_BBS[sq]
	p.occ &^= sq_bb
	p.placement_by_square[sq] = NULL_PIECE
	p.placement[pc] &^= sq_bb
}

func (p *Position) set_fen_info(positions string, color string, castles string, enps string, rule_50 string, move_count int) {
	rule_50_int, _ := strconv.Atoi(rule_50)
	// On Position
	p.placement = parse_positions(positions)
	p.placement_by_square = placement_to_placement_by_square(p.placement)
	p.occ = occupied_squares(p.placement)
	p.ply = move_count*2 + int(parse_color(color))
	p.stm = parse_color(color)
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
	p.stm = opposite(p.stm)

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
		rook_src_dst := ROOK_SRC_DST[dst]
		r_src, r_dst := rook_src_dst[0], rook_src_dst[1]
		r := p.piece_at(r_dst)
		p.do_castle(false, r, r_src, r_dst)
	}

	// Reassign state
	p.state = p.state.prev
}
