package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kruszczynski/blobfish-server/handlers"
)

func main() {
	r := mux.NewRouter()

	handler := handlers.NewHandler()

	r.HandleFunc("/memes", handler.MemesIndex).Methods("GET")
	r.HandleFunc("/memes", handler.MemesCreate).Methods("POST")
	r.HandleFunc("/memes/{memeID}", handler.MemeDestroy).Methods("DELETE")

	// Bind to a port and pass our router in
	http.ListenAndServe(":8000", r)
}
