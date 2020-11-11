package main

import (
	"log"
	"os"

	"github.com/aaronjprice/url_short/http_server"
	"github.com/aaronjprice/url_short/shortener"
	"github.com/aaronjprice/url_short/url_store"
)

func main() {
	shortener := shortener.NewMapHash()

	urlStore := url_store.New(shortener)

	port := ":8000"

	_, errCh := http_server.Start(
		http_server.Config{Address: port},
		environment{},
		urlStore,
	)

	err := <-errCh
	log.Printf("Fatal error: %v", err)
	os.Exit(1)
}

type environment struct{}

func (env environment) Printf(level int, format string, v ...interface{}) {
	log.Printf(format, v...)
}
