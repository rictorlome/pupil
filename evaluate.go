package main

import (
// "fmt"
)

func add_material_value(piece Piece) int {
	return MATERIAL_VALUES[piece]
}

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
	if color == WHITE {
		return square_rank(sq)
	}
	return 7 - square_rank(sq)
}

func add_position_value(piece Piece, color Color, sq Square) int {
	idx := relative_square(sq, color)
  value_arr := *POSITION_VALUES[int(piece_to_type(piece))]
  return value_arr[idx] * relative_multiplier(color)
}

func serialize_for_evaluate(piece Piece, piece_bb Bitboard) int {
	sum := 0
	color := piece_to_color(piece)
	for cursor := piece_bb; !empty(cursor); cursor &= cursor - 1 {
		sum += add_material_value(piece)
		sum += add_position_value(piece, color, Square(lsb(cursor)))
	}
	return sum
}

func (p *Position) evaluate(checkmate bool) int {
	if !checkmate {
		sum := 0
		for piece, piece_bb := range p.placement {
			sum += serialize_for_evaluate(Piece(piece), piece_bb)
		}
		return sum
	}
	return max_score(opposite(p.side_to_move()))
}
