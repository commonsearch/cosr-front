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

	if DetectBang("") != "" {
		t.Fatal("Empty search")
	}

	if !strings.Contains(DetectBang("!a littlebits"), "amazon") {
		t.Fatal("Amazon bang")
	}

	if strings.Contains(DetectBang("ra littlebits"), "amazon") {
		t.Fatal("No Amazon bang")
	}
}
