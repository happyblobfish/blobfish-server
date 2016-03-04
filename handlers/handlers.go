package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/happyblobfish/server/models"
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
	memesCollection := models.NewMemesCollection()

	h.db.View(memesCollection.GetMemes)

	json.NewEncoder(w).Encode(memesCollection.Memes)
}

// MemesCreate creates a meme from URL
// Currently it's the only way to support Gifs
// TODO: function too long
func (h *Handler) MemesCreate(w http.ResponseWriter, r *http.Request) {
	meme := &models.Meme{OriginalURL: r.FormValue("url")}
	resp, err := meme.GetImageBinary()

	if err != nil {
		log.Fatal(err)
	}

	if err := meme.SetData(resp.Body, resp.Header.Get("Content-Type")); err != nil {
		log.Fatal(err)
	}

	if err := h.db.Update(meme.Save); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(meme)
}

// MemesUpload handles upload of a meme
func (h *Handler) MemesUpload(w http.ResponseWriter, r *http.Request) {
	meme := new(models.Meme) // no OriginalUrl ATM
	// TODO: DRY!!! Lines 76:81 are a great candidate for a Builder struct
	// REPLY: Yes and no.
	if err := meme.SetData(r.Body, "image/png"); err != nil { // frontend currently supports only PNGs
		log.Fatal(err)
	}
	if err := h.db.Update(meme.Save); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(meme)
}

// MemeGet serves a meme image
// TODO: return appropriate Content-Type header - may be required?
func (h *Handler) MemeGet(w http.ResponseWriter, r *http.Request) {
	meme := h.findMeme(r)

	h.db.View(meme.GetImage)                                             // fetch image
	w.Header().Set("Content-Length", strconv.Itoa(meme.ContentLength())) // set length
	w.Header().Set("Content-Type", meme.ContentType)                     // set mime
	if err := meme.WriteImage(w); err != nil {
		log.Fatal("error writing binary data")
	}
}

// MemeDestroy destroys a meme
func (h *Handler) MemeDestroy(w http.ResponseWriter, r *http.Request) {
	meme := h.findMeme(r)

	if err := h.db.Update(meme.Destroy); err != nil {
		log.Fatal(err)
	}
}

func (h *Handler) findMeme(r *http.Request) *models.Meme {
	memeID, err := strconv.ParseUint(mux.Vars(r)["memeID"], 10, 64)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	meme := new(models.Meme)
	h.db.View(models.MemeFetch(meme, memeID))
	return meme
}

func openDB() *bolt.DB {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
