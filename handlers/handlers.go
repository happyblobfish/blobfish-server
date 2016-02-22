package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/kruszczynski/blobfish-server/models"
)

var dbName = "db/blobfish.db"

// Handler is a struct with
type Handler struct {
	db *bolt.DB
}

// NewHandler returns an instance of DB handler
func NewHandler() *Handler {
	db := openDB()
	return &Handler{db: db}
}

// MemesIndex list the index of blobfish's memes
func (h *Handler) MemesIndex(w http.ResponseWriter, r *http.Request) {
	memes := models.AllMemes(h.db)

	json.NewEncoder(w).Encode(memes)
}

// MemesCreate creates a meme
func (h *Handler) MemesCreate(w http.ResponseWriter, r *http.Request) {
	meme := &models.Meme{URL: r.FormValue("url")}

	if err := h.db.Update(meme.Save); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(meme)
}

// MemeDestroy destroys a meme
func (h *Handler) MemeDestroy(w http.ResponseWriter, r *http.Request) {
	memeID, err := strconv.ParseUint(mux.Vars(r)["memeID"], 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	meme := &models.Meme{ID: uint64(memeID)}
	if err := h.db.Update(meme.Destroy); err != nil {
		log.Fatal(err)
	}
}

func openDB() *bolt.DB {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
