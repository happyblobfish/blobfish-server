package models

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

var bucketName = []byte("memes")
var binBucketName = []byte("memes_bindata")
var allowedContentTypes = []string{"image/jpeg", "image/png", "image/gif"}
var maxImageContentLength int64 = 50000000 // 5 Mb

// Meme is a meme
type Meme struct {
	URL         string `json:"url"`
	ID          uint64 `json:"id"`
	OriginalURL string `json:"original-url"`
	ContentType string `json:"content-type"`
	bindata     []byte
}

// MemeFromBolt creates a meme from bolt's saved data
func MemeFromBolt(k []byte, v []byte) (*Meme, error) {
	meme := new(Meme)
	if err := json.Unmarshal(v, meme); err != nil {
		return meme, err
	}
	meme.ID = binary.BigEndian.Uint64(k)
	return meme, nil
}

// MemeFetch fetches all data from memes bucket and assign to meme
func MemeFetch(memePointer *Meme, id uint64) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		meme, err := MemeFromBolt(itob(id), tx.Bucket(bucketName).Get(itob(id)))
		if err != nil {
			log.Fatalln(err)
			return err
		}
		// assing new meme's pointer
		*memePointer = *meme
		return nil
	}
}

// GetImage retrieves binary image for a meme (from memes_bindata bucket) and assign
// to m.bindata
func (m *Meme) GetImage(tx *bolt.Tx) error {
	m.bindata = tx.Bucket(binBucketName).Get(itob(m.ID))
	return nil
}

// WriteImage writes binary image to a given io.Writer
func (m *Meme) WriteImage(w io.Writer) error {
	if _, err := w.Write(m.bindata); err != nil {
		return err
	}
	return nil
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
	m.ID = id
	m.URL = fmt.Sprintf("%s/%d", "memes", id)
	jsonString, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := bucket.Put(itob(id), jsonString); err != nil {
		return err
	}
	// write image binary to memes_bindata bucket
	return m.writeBindata(tx, id)
}

// Destroy deletes a meme from the DB
func (m *Meme) Destroy(tx *bolt.Tx) error {
	bucket, err := tx.CreateBucketIfNotExists(bucketName)
	if err != nil {
		return err
	}

	if err := bucket.Delete(itob(m.ID)); err != nil {
		return err
	}
	return nil
}

// SetData takes binary data from reader and saves it to image
func (m *Meme) SetData(r io.Reader, contentType string) error {
	m.ContentType = contentType
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	m.bindata = data
	return nil
}

// ContentLength returns image's bytes length
func (m *Meme) ContentLength() int {
	return len(m.bindata)
}

// GetImageBinary gets image's binary data from a URL
func (m *Meme) GetImageBinary() (*http.Response, error) {
	resp, err := http.Get(m.OriginalURL)
	if err != nil {
		errorMsg := fmt.Sprintf("%s %s", "http error", m.OriginalURL)
		return nil, errors.New(errorMsg)
	}
	if resp.StatusCode != 200 || resp.ContentLength > maxImageContentLength {
		errorMsg := fmt.Sprintf("%s %s", "error status or response too big", resp.StatusCode)
		return nil, errors.New(errorMsg)
	}
	contentTypeValid := false
	for _, contentType := range allowedContentTypes {
		if contentType == resp.Header.Get("Content-Type") {
			contentTypeValid = true
			break
		}
	}
	if !contentTypeValid {
		return nil, errors.New("invalid content type header")
	}
	return resp, nil
}

func (m *Meme) writeBindata(tx *bolt.Tx, id uint64) error {
	bucket, err := tx.CreateBucketIfNotExists(binBucketName)
	if err != nil {
		return err
	}
	if err := bucket.Put(itob(id), m.bindata); err != nil {
		return err
	}
	return nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
