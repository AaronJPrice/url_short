package url_store

import (
	"sync"

	"github.com/aaronjprice/url_short/shortener"
)

type Map struct {
	longToShort map[string]string
	shortToLong map[string]string
	shortener   shortener.Shortener
	mutex       *sync.Mutex
}

func New(s shortener.Shortener) *Map {
	return &Map{
		longToShort: make(map[string]string),
		shortToLong: make(map[string]string),
		shortener:   s,
		mutex:       &sync.Mutex{},
	}
}

func (m *Map) Compress(long string) string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	short, exists := m.longToShort[long]
	if exists {
		return short
	}

	short = m.shortener.Shorten(long)

	m.longToShort[long] = short
	m.shortToLong[short] = long

	return short
}

func (m *Map) Expand(short string) (string, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if long, exists := m.shortToLong[short]; exists {
		return long, true
	}

	return "", false
}
