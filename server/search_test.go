package main

import (
	"testing"
)

func TestSearchHref(t *testing.T) {
	t.Parallel()

	if (SearchRequest{}).Href() != "/" {
		t.Fatal("Empty search should be /")
	}

	if (SearchRequest{Query: "x"}).Href() != "/?q=x" {
		t.Fatal("Wrong SearchHref")
	}

	if (SearchRequest{Query: "x", Lang: "en"}).Href() != "/?g=en&q=x" {
		t.Fatal("Wrong SearchHref")
	}

	if (SearchRequest{Query: "x", Lang: "en", Page: 2}).Href() != "/?g=en&p=2&q=x" {
		t.Fatal("Wrong SearchHref")
	}

	if (SearchRequest{Query: "x", Lang: "all", Page: 1}).Href() != "/?g=all&q=x" {
		t.Fatal("Wrong SearchHref")
	}

	if (SearchRequest{Query: "", Lang: "all", Page: 2}).Href() != "/?g=all" {
		t.Fatal("Wrong SearchHref")
	}

	if (SearchRequest{Query: "", Lang: "", Page: 2}).Href() != "/" {
		t.Fatal("Wrong SearchHref")
	}

	if (SearchRequest{Query: "3", Lang: "", Page: 2}).Href() != "/?p=2&q=3" {
		t.Fatal("Wrong SearchHref")
	}

}

func TestSearchHrefPagination(t *testing.T) {
	t.Parallel()

	if (SearchRequest{}).NextPageHref() != "/" {
		t.Fatal("Empty search should be /")
	}

	if (SearchRequest{Query: "x"}).NextPageHref() != "/?p=2&q=x" {
		t.Fatal("Wrong SearchHref")
	}

	if (SearchRequest{Query: "x", Page: 1}).NextPageHref() != "/?p=2&q=x" {
		t.Fatal("Wrong SearchHref")
	}

	if (SearchRequest{Query: "x", Page: 2}).NextPageHref() != "/?p=3&q=x" {
		t.Fatal("Wrong SearchHref")
	}

	if (SearchRequest{Query: "x", Page: 2}).PreviousPageHref() != "/?q=x" {
		t.Fatal("Wrong SearchHref")
	}

	if (SearchRequest{Query: "x", Page: 3}).PreviousPageHref() != "/?p=2&q=x" {
		t.Fatal("Wrong SearchHref")
	}

	if (SearchRequest{Query: "x"}).PreviousPageHref() != "/?q=x" {
		t.Fatal("Wrong SearchHref")
	}

}
