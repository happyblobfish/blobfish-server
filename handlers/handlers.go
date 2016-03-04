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
var maxImageContentLength int64 = 1000000 // 1 Mb
var allowedContentTypes = []string{"image/jpeg", "image/png", "image/gif"}

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

// MemesCreate creates a meme
// TODO: function too long
func (h *Handler) MemesCreate(w http.ResponseWriter, r *http.Request) {
	meme := &models.Meme{URL: r.FormValue("url")}
	resp, err := http.Get(meme.URL)
	if err != nil {
		log.Fatalf("%s %s", "http error", meme.URL)
		return
	}
	if resp.StatusCode != 200 || resp.ContentLength > maxImageContentLength {
		log.Fatalf("%s %s", "error status or response too big", resp.StatusCode)
		return
	}
	contentTypeValid := false
	for _, contentType := range allowedContentTypes {
		if contentType == resp.Header.Get("Content-Type") {
			contentTypeValid = true
			break
		}
	}
	if !contentTypeValid {
		log.Fatalf("invalid content type header")
		return
	}
	if err := meme.SetData(resp.Body, resp.Header.Get("Content-Type")); err != nil {
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

	h.db.View(meme.GetImage)                         // fetch image
	w.Header().Set("Content-Type", meme.ContentType) // set mime
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

	ptr := new(models.Meme)
	h.db.View(models.MemeFetch(ptr, memeID))
	return ptr
}

func openDB() *bolt.DB {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
