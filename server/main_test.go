package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server *httptest.Server
)

func init() {

	SetupGlobals()

	server = httptest.NewServer(CreateRouter())

}

func search(t testing.TB, path string) string {

	resp, err := http.Get(server.URL + path)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	return string(body)
}

func TestHomepageLooksOkay(t *testing.T) {
	t.Parallel()

	body := search(t, "")

	if !strings.Contains(body, "Common Search") {
		t.Fatal("Should contain Common Search!")
	}
	if strings.Contains(body, "xxxteststring") {
		t.Fatal("Should not contain xxxteststring!")
	}
}

func TestLocalizedHomepageLooksOkay(t *testing.T) {
	t.Parallel()

	body := search(t, "/?g=fr")

	if !strings.Contains(body, "&#34;fr&#34;") {
		t.Fatal("Should be set as french!")
	}
}

func TestTestSearchLooksOkay(t *testing.T) {
	t.Parallel()

	body := search(t, "/?q=xxxteststring")

	if !strings.Contains(body, "xxxteststring") {
		t.Fatal("Should contain xxxteststring!")
	}

	if !strings.Contains(body, "http://www.example.com/page/1") {
		t.Fatal("Should contain example!")
	}
}

func TestStaticFiles(t *testing.T) {
	t.Parallel()

	body := search(t, "/js/index.js")

	if !strings.Contains(body, ".XMLHttpRequest") {
		t.Fatal("Should be JavaScript!")
	}

	favicon := search(t, "/favicon.ico")

	if len(favicon) < 500 || len(favicon) > 50000 {
		t.Fatal("Strange favicon!")
	}
}

func BenchmarkHomepage(b *testing.B) {

	for n := 0; n < b.N; n++ {

		body := search(b, "")

		if !strings.Contains(body, "Common Search") {
			b.Fatal("Should contain Common Search!")
		}
	}
}
