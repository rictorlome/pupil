package main

import (
	"encoding/json"
	// "fmt"
	"net/http"
	// "strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Game struct {
	pos *Position
}

func (g *Game) ThinkEndpoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	client_move := g.pos.parse_move(r.Form["move"][0])
	g.pos.do_move(client_move, &StateInfo{})
	engine_move := g.pos.alphaBetaRoot(6)
	g.pos.do_move(engine_move.move, &StateInfo{})
	json.NewEncoder(w).Encode(g.pos.String())

}

func startServer() {
	router := mux.NewRouter()
	p := parse_fen(INITIAL_FEN)
	g := Game{&p}
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/think", g.ThinkEndpoint).Methods("POST")
	if err := http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(router)); err != nil {
		panic(err)
	}
}
