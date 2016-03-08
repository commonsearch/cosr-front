package main

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

// ConfigSpec contains all the possible configuration variables.
// See http://github.com/kelseyhightower/envconfig for syntax
type ConfigSpec struct {

	// Debug mode turns off some optimizations and makes local development easier.
	Debug bool

	// IsDemo controls the presence of the demo warning on the frontend.
	IsDemo bool

	// TestData controls the use of test data as mock search results
	TestData bool

	// Env sets the current environment name. Valid values are "local", "ci", "prod".
	Env string `default:"local"`

	// Port sets the port we are listening on for requests.
	Port string `envconfig:"PORT" default:"9700"`

	// Host sets the IP address we are listening on for requests. Set to 127.0.0.1 to restrict to local.
	Host string `default:"0.0.0.0"`

	// ElasticsearchDocs is the HTTP url of the Elasticsearch instance for the document store.
	ElasticsearchDocs string `default:"http://192.168.99.100:39200"`

	// ElasticsearchText is the HTTP url of the Elasticsearch instance for the text index.
	ElasticsearchText string `default:"http://192.168.99.100:39200"`

	// PathFront is the path to the base directory of cosr-front.
	PathFront string `default:""`

	// ResultPageSize controls the number of results on each page.
	ResultPageSize int `default:"25"`
}

// Config contains the current configuration values.
var Config ConfigSpec

// LoadConfig populates the global configuration from COSR_* environment variables.
func LoadConfig() {

	err := envconfig.Process("cosr", &Config)

	if err != nil {
		log.Fatal(err.Error())
	}
}
