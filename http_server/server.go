package http_server

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aaronjprice/url_short/url_store"
)

type Config struct {
	Address string
}

type Environment interface {
	Printf(int, string, ...interface{})
}

func Start(conf Config, env Environment, urlStore url_store.UrlStore) (*http.Server, <-chan error) {
	env.Printf(0, "Starting HTTP server on %s", conf.Address)

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, env, urlStore)
	})

	svr := http.Server{
		Addr:    conf.Address,
		Handler: serveMux,
	}

	errCh := make(chan error)

	go func() {
		err := svr.ListenAndServe()
		errCh <- err
	}()

	return &svr, errCh
}

func handle(w http.ResponseWriter, r *http.Request, env Environment, urlStore url_store.UrlStore) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r, env, urlStore)
	case http.MethodPost:
		handlePost(w, r, env, urlStore)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request, env Environment, urlStore url_store.UrlStore) {
	env.Printf(0, "Handling GET %v", r.URL.Path)

	// No need to path unescape, that's done automatically, but we do need to trim leading "/"
	short := strings.TrimPrefix(r.URL.Path, "/")

	long, found := urlStore.Expand(short)
	if !found {
		env.Printf(0, "Not found: %s", short)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	env.Printf(0, "Redirecting to : %s", short)
	w.Header().Set("Location", long)
	w.WriteHeader(http.StatusSeeOther)
}

func handlePost(w http.ResponseWriter, r *http.Request, env Environment, urlStore url_store.UrlStore) {
	env.Printf(0, "Handling POST %v", r.URL.Path)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		env.Printf(0, "ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	long := string(body)
	short := urlStore.Compress(long)
	env.Printf(0, "%s compressed to %s", long, short)

	// We should path escape the short string as we're expecting users to put it in the path
	short = url.PathEscape(short)

	_, err = w.Write([]byte(short))
	if err != nil {
		env.Printf(0, "ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
