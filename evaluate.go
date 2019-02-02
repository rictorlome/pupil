package main

import (
// "fmt"
)

func relative_multiplier(color Color) int {
	if color == WHITE {
		return 1
	}
	return -1
}

func relative_square(sq Square, color Color) Square {
	return make_square(relative_square_rank(sq, color), square_file(sq))
}

func relative_square_rank(sq Square, color Color) int {
	if color == BLACK {
		return square_rank(sq)
	}
	return 7 - square_rank(sq)
}

func add_position_value(pt PieceType, color Color, sq Square) int {
	idx := relative_square(sq, color)
	value_arr := *POSITION_VALUES[pt]
	return value_arr[idx]
}

func serialize_for_evaluate(us Color, piece Piece, piece_bb Bitboard) int {
	score := 0
	color, pt := piece_to_color(piece), piece_to_type(piece)
	for cursor := piece_bb; !empty(cursor); cursor &= cursor - 1 {
		if color == us {
			score += MATERIAL_VALUES[pt]
			score += add_position_value(pt, color, Square(lsb(cursor)))
		} else {
			score -= MATERIAL_VALUES[pt]
			score -= add_position_value(pt, color, Square(lsb(cursor)))
		}
	}
	return score
}

func (p *Position) evaluate(checkmate bool) int {
	if !checkmate {
		score := 0
		for piece, piece_bb := range p.placement {
			score += serialize_for_evaluate(p.side_to_move(), Piece(piece), piece_bb)
		}
		return score
	}
	return -MAX_SCORE
}
