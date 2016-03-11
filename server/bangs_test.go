package main

import (
	"strings"
	"testing"
)

func init() {

	SetupGlobals()

}

func TestBangs(t *testing.T) {
	t.Parallel()

	if DetectBang("", "en") != "" {
		t.Fatal("Empty search")
	}

	if DetectBang("!a littlebits", "en") != "https://www.amazon.com/s/?field-keywords=littlebits" {
		t.Fatal("Amazon bang EN")
	}

	if DetectBang("!a littlebits", "fr") != "https://www.amazon.fr/s/?field-keywords=littlebits" {
		t.Fatal("Amazon bang FR")
	}

	if DetectBang("!w littlebits", "en") != "https://en.wikipedia.org/wiki/Special:Search?search=littlebits" {
		t.Fatal("Wikipedia bang EN (via any)")
	}

	if strings.Contains(DetectBang("ra littlebits", "en"), "amazon") {
		t.Fatal("No Amazon bang")
	}
}
