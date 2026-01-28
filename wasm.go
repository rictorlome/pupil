//go:build js && wasm

package main

import (
	"syscall/js"
)

var game *Game

func newGame(this js.Value, args []js.Value) interface{} {
	p := parseFen(INITIAL_FEN)
	game = &Game{&p}
	TT_GLOBAL = createTT(18) // Smaller TT for browser (64MB instead of 120MB)
	return game.pos.String()
}

func makeMove(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return js.ValueOf(map[string]interface{}{
			"error": "no move provided",
		})
	}

	uci := args[0].String()
	clientMove := game.pos.parseMove(uci)

	if !game.validateMove(clientMove) {
		return js.ValueOf(map[string]interface{}{
			"error": "illegal move",
			"fen":   game.pos.String(),
		})
	}

	game.pos.doMove(clientMove, &StateInfo{})

	return js.ValueOf(map[string]interface{}{
		"fen": game.pos.String(),
	})
}

func getEngineMove(this js.Value, args []js.Value) interface{} {
	depth := uint8(5) // Default depth 5 for faster response in browser
	if len(args) > 0 {
		depth = uint8(args[0].Int())
	}

	engineMove := game.pos.abRoot(depth)
	game.pos.doMove(engineMove.move, &StateInfo{})

	return js.ValueOf(map[string]interface{}{
		"move": engineMove.move.String(),
		"fen":  game.pos.String(),
	})
}

func getFen(this js.Value, args []js.Value) interface{} {
	if game == nil {
		return ""
	}
	return game.pos.String()
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("pupilNewGame", js.FuncOf(newGame))
	js.Global().Set("pupilMakeMove", js.FuncOf(makeMove))
	js.Global().Set("pupilGetEngineMove", js.FuncOf(getEngineMove))
	js.Global().Set("pupilGetFen", js.FuncOf(getFen))

	println("Pupil Chess Engine (WASM) loaded")

	<-c // Keep the Go program running
}
