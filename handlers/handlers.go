package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/kruszczynski/blobfish-server/models"
)

var dbName = "blobfish.db"

// MemesIndex list the index of blobfish's memes
func MemesIndex(w http.ResponseWriter, r *http.Request) {
	db := openDB()
	defer db.Close()

	memes := models.AllMemes(db)

	json.NewEncoder(w).Encode(memes)
}

// MemesCreate creates a meme
func MemesCreate(w http.ResponseWriter, r *http.Request) {
	db := openDB()
	defer db.Close()

	meme := &models.Meme{URL: r.FormValue("url")}

	if err := db.Update(meme.Save); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(meme)
}

func openDB() *bolt.DB {
	db, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
