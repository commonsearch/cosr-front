package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// resultPage represents the data needed to render an HTML result page or a JSON API return.
type resultPage struct {
	Search SearchRequest `json:"s"`
	Type   string        `json:"t,omitempty"`
	Result SearchResult  `json:"r"`
}

// getSearchRequest interprets the query in the URL by transforming a http.Request into a SearchRequest.
func getSearchRequest(r *http.Request) *SearchRequest {

	sr := SearchRequest{}

	sr.Query = strings.Trim(r.FormValue("q"), " ")

	sr.Lang = r.FormValue("g")

	sr.Page, _ = strconv.Atoi(r.FormValue("p"))

	if sr.Page == 0 || sr.Query == "" {
		sr.Page = 1
	}

	// Keeping Lang empty (=autodetect on client side) is acceptable
	// only if the query is empty (and we will land on the full homepage)
	// If we have a query, we unfortunately have to guess an english default
	if sr.Lang == "" && sr.Query != "" {
		sr.Lang = "en"
	}

	err := r.Body.Close()
	if err != nil {
		fmt.Printf("Warning: Could not close request Body")
	}

	return &sr
}

// sendResultPage renders a resultPage to HTML and sends it to the client.
func sendResultPage(w http.ResponseWriter, r *http.Request, page *resultPage) {

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	// If we are in debug mode, read the template from disk at each request!
	if Config.Debug {
		LoadTemplates()
	}

	err := Templates["index.html"].ExecuteTemplate(w, "index.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// TruncateQuery allows only Config.MaxQueryTerms words to enter the actual search
func TruncateQuery(q string) (string, string) {
	words := strings.Fields(q)
	allowedLength := Config.MaxQueryTerms

	var extra string // the starting word that are truncated
	if len(words) > allowedLength {
		q = strings.Join(words[:allowedLength], " ")
		extra = words[allowedLength]
	}

	return q, extra
}

// SearchHandler handles HTTP queries to home or result pages (/ or /?q=*).
func SearchHandler(w http.ResponseWriter, r *http.Request) {

	search := getSearchRequest(r)

	// Empty query: render the "home" version
	if search.Query == "" {
		page := resultPage{Type: "home", Search: *search}
		sendResultPage(w, r, &page)
		return
	}

	truncatedQuery, extra := TruncateQuery(search.Query)
	if extra != "" {
		search.Query = truncatedQuery
	}

	// Perform the search itself
	result, err := search.PerformSearchWithTiming()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If we used a !bang or asked for a redirect, do it now
	if result.Redirect != "" {
		http.Redirect(w, r, result.Redirect, 302)
		return
	}

	result.Extra = extra
	page := resultPage{Search: *search, Result: *result}
	sendResultPage(w, r, &page)
}

// APISearchHandler handles HTTP queries to our JSON API (/api/search?q=*)
func APISearchHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	search := getSearchRequest(r)
	truncatedQuery, extra := TruncateQuery(search.Query)

	if extra != "" {
		search.Query = truncatedQuery
	}

	// Perform the search itself
	result, err := search.PerformSearchWithTiming()
	result.Extra = extra
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the result to the client as JSON
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
