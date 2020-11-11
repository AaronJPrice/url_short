package http_server

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/aaronjprice/url_short/shortener"
	"github.com/aaronjprice/url_short/testlib"
	"github.com/aaronjprice/url_short/url_store"
)

const localhost = "http://localhost"

// A HTTP client that _won't_ follow redirects
// From http.Client docs: "Clients are safe for concurrent use by multiple goroutines."
var testClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
}

func TestSimple(t *testing.T) {
	port := ":8000"
	svr, errCh := Start(
		Config{Address: port},
		testEnv{t: t},
		url_store.New(shortener.NewMapHash()),
	)

	originalURL := "https://subdomain.domain.tld/path1"

	postResp, err := testClient.Post(localhost+port, "", strings.NewReader(originalURL))
	testlib.ErrCheck(t, err)
	testlib.AssertEqual(t, postResp.StatusCode, http.StatusOK)

	postRespBody, err := ioutil.ReadAll(postResp.Body)
	testlib.ErrCheck(t, err)

	getResp, err := testClient.Get(localhost + port + "/" + string(postRespBody))
	testlib.ErrCheck(t, err)
	testlib.AssertEqual(t, getResp.StatusCode, http.StatusSeeOther)

	resultURL, err := getResp.Location()
	testlib.ErrCheck(t, err)

	testlib.AssertEqual(t, resultURL.String(), originalURL)

	select {
	case err := <-errCh:
		testlib.ErrCheck(t, err)
	default:
	}
	svr.Shutdown(context.Background())
}

func TestNotFound(t *testing.T) {
	port := ":8001"
	svr, errCh := Start(
		Config{Address: port},
		testEnv{t: t},
		url_store.New(shortener.NewMapHash()),
	)

	originalURL := "https://subdomain.domain.tld/path1"

	postResp, err := testClient.Post(localhost+port, "", strings.NewReader(originalURL))
	testlib.ErrCheck(t, err)
	testlib.AssertEqual(t, postResp.StatusCode, http.StatusOK)

	getResp, err := testClient.Get(localhost + port + "/Definitely_not_the_correct_string")
	testlib.ErrCheck(t, err)
	testlib.AssertEqual(t, getResp.StatusCode, http.StatusNotFound)

	select {
	case err := <-errCh:
		testlib.ErrCheck(t, err)
	default:
	}
	svr.Shutdown(context.Background())
}

func TestSpecialChars(t *testing.T) {
	port := ":8002"
	specialChars := "abc123!@£$%^&*()_+-={}[]:\"|;'\\<>?,./~`'±§"

	var lastCompressed string
	tus := testUrlStore{
		compress: func(long string) string {
			lastCompressed = long
			return specialChars
		},
		expand: func(short string) (string, bool) {
			if short == specialChars {
				return lastCompressed, true
			}
			return "", false
		},
	}

	svr, errCh := Start(
		Config{Address: port},
		testEnv{t: t},
		tus,
	)

	originalURL := "https://subdomain.domain.tld/path1"

	postResp, err := testClient.Post(localhost+port, "", strings.NewReader(originalURL))
	testlib.ErrCheck(t, err)
	testlib.AssertEqual(t, postResp.StatusCode, http.StatusOK)

	postRespBody, err := ioutil.ReadAll(postResp.Body)
	testlib.ErrCheck(t, err)

	getResp, err := testClient.Get(localhost + port + "/" + string(postRespBody))
	testlib.ErrCheck(t, err)
	testlib.AssertEqual(t, getResp.StatusCode, http.StatusSeeOther)

	resultURL, err := getResp.Location()
	testlib.ErrCheck(t, err)

	testlib.AssertEqual(t, resultURL.String(), originalURL)

	select {
	case err := <-errCh:
		testlib.ErrCheck(t, err)
	default:
	}
	svr.Shutdown(context.Background())
}

type testEnv struct {
	t *testing.T
}

func (te testEnv) Printf(level int, format string, v ...interface{}) {
	te.t.Logf(format, v...)
}

type testUrlStore struct {
	compress func(string) string
	expand   func(string) (string, bool)
}

func (tus testUrlStore) Compress(long string) string {
	return tus.compress(long)
}

func (tus testUrlStore) Expand(short string) (string, bool) {
	return tus.expand(short)
}
