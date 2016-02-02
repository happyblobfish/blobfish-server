package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kruszczynski/blobfish-server/handlers"
)

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.

	r.HandleFunc("/memes", handlers.MemesIndex).Methods("GET")
	r.HandleFunc("/memes", handlers.MemesCreate).Methods("POST")

	// Bind to a port and pass our router in
	http.ListenAndServe(":8000", r)
}
