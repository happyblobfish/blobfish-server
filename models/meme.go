package models

import (
	"encoding/binary"

	"github.com/boltdb/bolt"
)

var bucketName = []byte("memes")

// Meme is a meme
type Meme struct {
	URL string `json:"url"`
	ID  uint64 `json:"id"`
}

// Save saves meme in bolt storage
func (m *Meme) Save(tx *bolt.Tx) error {
	bucket, err := tx.CreateBucketIfNotExists(bucketName)
	if err != nil {
		return err
	}

	// Generate ID for the user.
	// This returns an error only if the Tx is closed or not writeable.
	// That can't happen in an Update() call so I ignore the error check.
	id, _ := bucket.NextSequence()
	if err := bucket.Put(itob(id), []byte(m.URL)); err != nil {
		return err
	}
	m.ID = id
	return nil
}

func (m *Meme) fromBolt(id []byte, url []byte) {
	m.ID = binary.BigEndian.Uint64(id)
	m.URL = string(url)
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
