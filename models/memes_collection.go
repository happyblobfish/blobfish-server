package models

import "github.com/boltdb/bolt"

// MemesCollection is an array of memes
type MemesCollection struct {
	Memes []*Meme
}

// NewMemesCollection returns an instance of MemesCollection struct
// that's used to fetch dem memes from bolt db
func NewMemesCollection() *MemesCollection {
	return &MemesCollection{Memes: []*Meme{}}
}

// GetMemes returns all memes in the storage
func (mc *MemesCollection) GetMemes(tx *bolt.Tx) error {
	b := tx.Bucket(bucketName)
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		meme := &Meme{}
		meme.fromBolt(k, v)
		mc.Memes = append(mc.Memes, meme)
	}
	return nil
}
