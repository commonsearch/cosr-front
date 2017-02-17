package main

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
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
	ElasticsearchDocs string `default:"http://__local_docker_host__:39200"`

	// ElasticsearchText is the HTTP url of the Elasticsearch instance for the text index.
	ElasticsearchText string `default:"http://__local_docker_host__:39200"`

	// PathFront is the path to the base directory of cosr-front.
	PathFront string `default:""`

	// ResultPageSize controls the number of results on each page.
	ResultPageSize int `default:"25"`

	// Maximum allowed number of words for a query
	MaxQueryTerms int `default:"10"`
}

// Config contains the current configuration values.
var Config ConfigSpec

// LoadConfig populates the global configuration from COSR_* environment variables.
func LoadConfig() {

	err := envconfig.Process("cosr", &Config)

	if err != nil {
		log.Fatal(err.Error())
	}

	// Discover the IP of the local Docker host and replace it in the config values that may use it.
	localDockerHost := GetDockerHostIP()
	log.Println("Using Docker host IP: " + localDockerHost)

	Config.ElasticsearchDocs = strings.Replace(Config.ElasticsearchDocs, "__local_docker_host__", localDockerHost, 1)
	Config.ElasticsearchText = strings.Replace(Config.ElasticsearchText, "__local_docker_host__", localDockerHost, 1)

}

// GetDockerDaemonIP returns the IP of the Docker daemon visible from the host
func GetDockerDaemonIP() string {

	// When using boot2docker on Mac, DOCKER_HOST will be something like "tcp://192.168.99.100:2376"
	envDockerHost := os.Getenv("DOCKER_HOST")
	if envDockerHost != "" {
		parsedDockerHost, err := url.Parse(envDockerHost)
		if err == nil {
			return strings.Split(parsedDockerHost.Host, ":")[0]
		}
	}

	// On Linux, Docker should be running directly on localhost.
	return "127.0.0.1"
}

// GetDockerHostIP returns the IP of the Docker host, from inside a container
func GetDockerHostIP() string {

	daemonIP := GetDockerDaemonIP()

	if daemonIP == "127.0.0.1" {

		// 172.17.42.1 used to be hardcoded as a Docker host IP, now this seems to be the way to get it!
		out, err := exec.Command("sh", "-c", "/sbin/ip route | awk '/default/ { print $3 }'").Output()
		if err != nil {
			return daemonIP
		}
		outStr := strings.TrimSpace(string(out[:]))
		if outStr != "" {
			return outStr
		}

	}
	return daemonIP
}
