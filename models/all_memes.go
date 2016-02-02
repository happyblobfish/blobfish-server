package models

import "github.com/boltdb/bolt"

// AllMemes returns all memes in the storage
func AllMemes(db *bolt.DB) []*Meme {
	memes := []*Meme{}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			meme := &Meme{}
			meme.fromBolt(k, v)
			memes = append(memes, meme)
		}
		return nil
	})
	return memes
}
