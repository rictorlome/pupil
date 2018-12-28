package main

import (
	// "fmt"
	"strconv"
	"strings"
)

func bitboards_to_grid(bitboards []Bitboard) [8][8]string {
	var grid [8][8]string
	for _, s := range SQUARES {
		for _, piece := range PIECES {
			if occupied_at_sq(bitboards[piece], s) {
				grid[7-square_rank(s)][square_file(s)] = piece.String()
			}
		}
	}
	return grid
}

func generate_castle_string(s StateInfo) string {
	var castle_string string
	// Reversed because bit indices are right to left
	for idx, char := range "qkQK" {
		if has_right(s.castling_rights, uint(idx)) {
			castle_string = string(char) + castle_string
		}
	}
	if castle_string == "" {
		return "-"
	}
	return castle_string
}

func generate_color_string(pos Position) string {
	if pos.to_move() == WHITE {
		return "w"
	}
	return "b"
}

func generate_rule_50_string(s StateInfo) string {
	return strconv.Itoa(s.rule_50)
}

func generate_fen(pos Position) string {
	var fenArr []string
	fenArr = append(fenArr, grid_to_fen(bitboards_to_grid(pos.placement)))
	fenArr = append(fenArr, generate_color_string(pos))
	fenArr = append(fenArr, generate_castle_string(*pos.state))
	fenArr = append(fenArr, pos.state.ep_sq.String())
	fenArr = append(fenArr, generate_rule_50_string(*pos.state))
	fenArr = append(fenArr, strconv.Itoa(pos.move_count()))
	return strings.Join(fenArr, " ")
}

func grid_to_fen(grid [8][8]string) string {
	var fenArr []string
	for _, row := range grid {
		fenString := ""
		offset := 0
		for _, sq := range row {
			if sq == "" {
				offset += 1
			} else {
				if offset == 0 {
					fenString += sq
				} else {
					fenString += strconv.Itoa(offset) + sq
					offset = 0
				}
			}
		}
		if offset != 0 {
			fenString += strconv.Itoa(offset)
		}
		fenArr = append(fenArr, fenString)
	}
	return strings.Join(fenArr, "/")
}

func parse_color(s string) Color {
	if s == "w" {
		return WHITE
	}
	return BLACK
}

func parse_fen(fen string) Position {
	fields := strings.Split(fen, " ")
	move_count, _ := strconv.Atoi(fields[5])
	// Core fen info
	p := Position{}
	p.state = &StateInfo{}
	p.set_fen_info(fields[0], fields[1], fields[2], fields[3], fields[4], move_count)
	// Additional state info
	p.state.blockers_for_king = p.slider_blockers(opposite(p.to_move()), p.king_square(p.to_move()))
	p.state.prev = nil

	return p
}

func parse_move_type(promstring string, occ Bitboard, src Square, dst Square, mover_type PieceType, captured Piece) MoveType {
	file_diff := square_file(dst) - square_file(src)
	rank_diff := square_rank(dst) - square_rank(src)
	// Promotion
	var mt MoveType
	if promstring != "" {
		idx := strings.Index(strings.Join(PROMOTION_STRINGS, ""), promstring)
		mt = MoveType(idx<<12) | MoveType(PROMOTION_MASK)
	}
	// Capture
	mt |= cap_or_quiet(occ, dst)
	// En passant
	if mover_type == PAWN && square_file(dst) != square_file(src) && captured == NULL_PIECE {
		mt |= EP_CAPTURE
	}
	// Double Pawn Push
	if mover_type == PAWN && (rank_diff == 2 || rank_diff == -2) {
		mt |= DOUBLE_PAWN_PUSH
	}
	// Castles
	if mover_type == KING && file_diff == 2 {
		mt |= KING_CASTLE
	}
	if mover_type == KING && file_diff == -2 {
		mt |= QUEEN_CASTLE
	}
	return mt
}

func parse_positions(positions string) []Bitboard {
	result_bbs := make([]Bitboard, 12)
	ranks := strings.Split(positions, "/")
	for rank, rank_string := range ranks {
		offset := 0
		for file, sq := range rank_string {
			if strings.Contains(PIECE_STRING, string(sq)) {
				index := make_piece(string(sq))
				sq_num := make_square(RANK_8-rank, file+offset)
				result_bbs[index] |= SQUARE_BBS[sq_num]
			} else {
				offset += int(sq-'0') - 1
			}
		}
	}
	return result_bbs
}

func parse_square(sq string) Square {
	if sq == "-" {
		return NULL_SQ
	}
	rank := int(sq[1]-'0') - 1
	return make_square(rank, strings.Index(FILES, sq[0:1]))
}
