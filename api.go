package main

import (
	"encoding/json"
	"net/http"
)

type Game struct {
	pos *Position
}

func (g *Game) validateMove(m Move) bool {
	legalMoves := g.pos.generateMoves()
	for _, legalMove := range legalMoves {
		if m == legalMove {
			return true
		}
	}
	return false
}

func (g *Game) BeginEndpoint(w http.ResponseWriter, r *http.Request) {
	engineMove := g.pos.abRoot(6)
	g.pos.doMove(engineMove.move, &StateInfo{})
	json.NewEncoder(w).Encode(g.pos.String())
}

func (g *Game) RestartEndpoint(w http.ResponseWriter, r *http.Request) {
	p := parseFen(INITIAL_FEN)
	g.pos = &p
	TT_GLOBAL = createTT(20)
	json.NewEncoder(w).Encode("OK")
}

func (g *Game) ThinkEndpoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	clientMove := g.pos.parseMove(r.Form["move"][0])
	if g.validateMove(clientMove) {
		g.pos.doMove(clientMove, &StateInfo{})
		engineMove := g.pos.abRoot(6)
		g.pos.doMove(engineMove.move, &StateInfo{})
	}
	json.NewEncoder(w).Encode(g.pos.String())
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func startServer() {
	mux := http.NewServeMux()
	p := parseFen(INITIAL_FEN)
	g := Game{&p}

	mux.HandleFunc("POST /begin", g.BeginEndpoint)
	mux.HandleFunc("POST /restart", g.RestartEndpoint)
	mux.HandleFunc("POST /think", g.ThinkEndpoint)
	mux.Handle("/", http.FileServer(http.Dir("./static/")))

	if err := http.ListenAndServe(":8080", corsMiddleware(mux)); err != nil {
		panic(err)
	}
}
