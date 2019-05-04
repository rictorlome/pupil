package main

import (
	"fmt"
	"syscall/js"
)

func init() {
	// Initialize Square Bitboards
	e := Bitboard(0)
	for _, s := range SQUARES {
		var SQUARE_BB Bitboard = 0x1 << s
		SQUARE_BBS[s] = SQUARE_BB
		NEIGHBOR_BBS[s] = neighbors(SQUARE_BB)
		KING_ATTACK_BBS[s] = precompute_king_attacks(SQUARE_BB)
		KNIGHT_ATTACK_BBS[s] = precompute_knight_attacks(SQUARE_BB)
		ROOK_ATTACK_MASKS[s] = slider_rook_attacks(Bitboard(0), s)
		BISHOP_ATTACK_MASKS[s] = slider_bishop_attacks(Bitboard(0), s)
		ROOK_OCCUPANCY_MASKS[s] = occupancy_mask(s, ROOK_DIRECTIONS)
		BISHOP_OCCUPANCY_MASKS[s] = occupancy_mask(s, BISHOP_DIRECTIONS)
		RELEVANT_ROOK_OCCUPANCY[s] = ROOK_ATTACK_MASKS[s] &^ ROOK_OCCUPANCY_MASKS[s]
		RELEVANT_BISHOP_OCCUPANCY[s] = BISHOP_ATTACK_MASKS[s] &^ BISHOP_OCCUPANCY_MASKS[s]
		for _, color := range COLORS {
			PAWN_ATTACK_BBS[s][color] = pawn_attacks(SQUARE_BB, color)
		}
	}
	// Initialize dependent BBs
	for _, fn := range [2]AttackFunc{slider_rook_attacks, slider_bishop_attacks} {
		for _, s1 := range SQUARES {
			for _, s2 := range SQUARES {
				if occupied_at_sq(fn(e, s1), s2) {
					LINE_BBS[s1][s2] = fn(e, s1)&fn(e, s2) | SQUARE_BBS[s1] | SQUARE_BBS[s2]
					BETWEEN_BBS[s1][s2] = fn(SQUARE_BBS[s2], s1) & fn(SQUARE_BBS[s1], s2)
				}
			}
		}
	}

	init_castle_sqs()
	init_castling_masks()
	init_rook_squares_for_castling()
	init_magics()
	init_zobrists()
	init_pool()
}

func checkGameOver(i []js.Value) {
	moves := BOARD_GLOBAL.generate_moves()
	if len(moves) == 0 {
		js.Global().Call("alert", "Game Over")
		js.Global().Call("reload")
	}
}

func makeHumanMove(i []js.Value) {
	move_string := js.ValueOf(i[0].String()).String()
	human_move := BOARD_GLOBAL.parse_move(move_string)
	BOARD_GLOBAL.do_move(human_move, &StateInfo{})
	js.Global().Get("board1").Call("position", BOARD_GLOBAL.String(), true)
	js.Global().Set("humansTurn", false)
}

func makeComputerMove(i []js.Value) {
	engine_move := BOARD_GLOBAL.ab_root(6)
	BOARD_GLOBAL.do_move(engine_move.move, &StateInfo{})
	js.Global().Get("board1").Call("position", BOARD_GLOBAL.String(), true)
	js.Global().Call("clearInterval", "firstMoveInterval")
	js.Global().Set("humansTurn", true)
}

// func validateHumanMove(i []js.Value) {
// 	move_string := js.ValueOf(i[0].String()).String()
// 	move := BOARD_GLOBAL.parse_move(move_string)
// 	moves := BOARD_GLOBAL.generate_moves()
// }

var BOARD_GLOBAL Position



func registerCallbacks() {
    js.Global().Set("makeComputerMove", js.NewCallback(makeComputerMove))
		js.Global().Set("makeHumanMove", js.NewCallback(makeHumanMove))
		js.Global().Set("checkGameOver", js.NewCallback(checkGameOver))
		// js.Global().Set("validateHumanMove", js.NewCallback(validateHumanMove))
}

func main() {
	fmt.Println("OK")
	BOARD_GLOBAL = parse_fen(INITIAL_FEN)
	// pattern cribbed from https://tutorialedge.net/golang/go-webassembly-tutorial/
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}
