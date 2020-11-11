package shortener

import (
	"encoding/base64"
	"hash/maphash"
)

type MapHash struct{}

func NewMapHash() *MapHash {
	return &MapHash{}
}

func (mh *MapHash) Shorten(long string) string {
	var h maphash.Hash
	h.WriteString(long)
	shortBytes := h.Sum([]byte{})
	shortB64 := base64.StdEncoding.EncodeToString(shortBytes)
	return shortB64
	// return url.PathEscape(shortB64)
}
