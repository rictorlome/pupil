package main

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

type Game struct {
	pos *Position
}

func (g *Game) validateMove(m Move) bool {
	legal_moves := g.pos.generate_moves()
	for _, legal_move := range legal_moves {
		if m == legal_move {
			return true
		}
	}
	return false
}

func (g *Game) BeginEndpoint(w http.ResponseWriter, r *http.Request) {
	engine_move := g.pos.ab_root(6)
	g.pos.do_move(engine_move.move, &StateInfo{})
	json.NewEncoder(w).Encode(g.pos.String())
}

func (g *Game) RestartEndpoint(w http.ResponseWriter, r *http.Request) {
	p := parse_fen(INITIAL_FEN)
	g.pos = &p
	TT_GLOBAL = createTT(20)
	json.NewEncoder(w).Encode("OK")
}

func (g *Game) ThinkEndpoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	client_move := g.pos.parse_move(r.Form["move"][0])
	if g.validateMove(client_move) {
		g.pos.do_move(client_move, &StateInfo{})
		engine_move := g.pos.ab_root(6)
		g.pos.do_move(engine_move.move, &StateInfo{})
	}
	json.NewEncoder(w).Encode(g.pos.String())
}

func startServer() {
	router := mux.NewRouter()
	p := parse_fen(INITIAL_FEN)
	g := Game{&p}
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/begin", g.BeginEndpoint).Methods("POST")
	router.HandleFunc("/restart", g.RestartEndpoint).Methods("POST")
	router.HandleFunc("/think", g.ThinkEndpoint).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	if err := http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(router)); err != nil {
		panic(err)
	}
}
