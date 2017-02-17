package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

var bangs = make(map[string](map[string]string))

// LoadBangs loads bang definitions from a static JSON file at startup.
func LoadBangs() {

	cnt, err := ioutil.ReadFile(path.Join(Config.PathFront, "server/bangs.json"))
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(cnt, &bangs); err != nil {
		log.Fatal(err)
	}
}

// DetectBang detects bang usage in a query string and returns a redirect URL if found.
// We support a small subset of https://duckduckgo.com/bang
func DetectBang(query string, lang string) string {

	// Bangs definitions don't support "all" yet.
	if lang == "all" {
		return query
	}

	parts := strings.Split(query, " ")

	for i, part := range parts {
		if len(part) > 1 && strings.HasPrefix(part, "!") {

			// Look for a definition in "[lang]", then in "any"
			url := bangs[part[1:]][lang]
			if url == "" {
				url = bangs[part[1:]]["any"]
			}

			if url != "" {
				leftoverSearch := strings.Join(append(parts[:i], parts[i+1:]...), " ")
				url = strings.Replace(url, "{{{s}}}", leftoverSearch, -1)
				url = strings.Replace(url, "{{{lang}}}", lang, -1)
				return url
			}
		}
	}

	return ""
}
