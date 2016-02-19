package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

var bangs = make(map[string]string)

// LoadBangs loads bang definitions from a static JSON file at startup.
func LoadBangs() {
	cnt, err := ioutil.ReadFile(path.Join(Config.PathFront, "server/bangs.json"))
	if err != nil {
		log.Fatal(err)
	}

	cnt = bytes.Replace(cnt, []byte("\n"), []byte(""), 0)

	if err := json.Unmarshal(cnt, &bangs); err != nil {
		log.Fatal(err)
	}
}

// DetectBang detects bang usage in a query string and returns a redirect URL if found.
// We support a small subset of https://duckduckgo.com/bang
func DetectBang(query string) string {

	parts := strings.Split(query, " ")

	for i, part := range parts {
		if len(part) > 1 && strings.HasPrefix(part, "!") && bangs[part[1:]] != "" {
			leftoverSearch := strings.Join(append(parts[:i], parts[i+1:]...), " ")
			return strings.Replace(bangs[part[1:]], "{{{s}}}", leftoverSearch, -1)
		}
	}

	return ""
}
