package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"text/template"
)

// Remove most of the useless spaces in HTML.
var spacesRegexp = regexp.MustCompile("(>|\\})\\s+(<|\\{)")

// Replace a CSS tag with inline content.
var cssRegexp = regexp.MustCompile("<!-- CSS([\\s\\S]+?)ENDCSS -->")

// Remove protocol and trailing slash in homepages.
var simplifyURLRegexp = regexp.MustCompile("(.*?://)(([^/]+)(/.+)?).*")

// Templates is a map of all the parsed templates.
var Templates = make(map[string]*template.Template)

// preprocessTemplate does basic HTML minification before compiling templates.
func preprocessTemplate(s string) string {

	// We are in production mode, inline the CSS!
	if Config.Env == "prod" {

		css, err := ioutil.ReadFile(path.Join(Config.PathFront, "build/static/css/index.css"))

		if err != nil {
			log.Fatal("built index.css not found!")
		}

		// Trim the BOM if present.
		css = bytes.Trim(css, "\xef\xbb\xbf")

		s = cssRegexp.ReplaceAllString(s, fmt.Sprintf(`<style type="text/css">%s</style>`, css))
	}

	return spacesRegexp.ReplaceAllString(s, "$1$2")
}

// truncate truncates a string to 'max' characters.
func truncate(s string, max int) string {
	var numRunes = 0
	for index := range s {
		numRunes++
		if numRunes > max {
			return s[:index]
		}
	}
	return s
}

// Simplify a hit URL for display.
func simplifyURL(s string) string {
	return truncate(simplifyURLRegexp.ReplaceAllString(s, "$2"), 100)
}

// toJSON returns an object as a JSON string.
func toJSON(v interface{}) string {
	a, _ := json.Marshal(v)
	return string(a)
}

// getConfig returns the current config.
func getConfig() ConfigSpec {
	return Config
}

// ParseTemplate returns a pre-processed and parsed template.
func ParseTemplate(filepath string) *template.Template {
	cnt, err := ioutil.ReadFile(path.Join(Config.PathFront, "templates/"+filepath))

	if err != nil {
		log.Fatal(err)
	}

	funcMap := template.FuncMap{
		"simplifyURL": simplifyURL,
		"toJSON":      toJSON,
		"getConfig":   getConfig,
	}

	t, err := template.New(filepath).Funcs(funcMap).Parse(preprocessTemplate(string(cnt)))
	if err != nil {
		log.Fatal(err)
	}

	return t
}

// LoadTemplates loads and parses all known templates at startup.
func LoadTemplates() {
	Templates["index.html"] = ParseTemplate("index.html")
}
