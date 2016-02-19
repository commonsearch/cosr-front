package main

import (
	"encoding/json"
	"gopkg.in/olivere/elastic.v3"
	"log"
	"net/url"
	"os"
	"time"
)

// ElasticsearchTextClient is the ES client to the main index.
var ElasticsearchTextClient *elastic.Client

// ElasticsearchDocsClient is the ES client to the document store.
var ElasticsearchDocsClient *elastic.Client

// ElasticsearchConnect sets up persistent connections to both ES servers.
func ElasticsearchConnect() {

	ElasticsearchTextClient = ElasticsearchConnectServer(Config.ElasticsearchText)

	ElasticsearchDocsClient = ElasticsearchConnectServer(Config.ElasticsearchDocs)

}

// ElasticsearchConnectServer connects one single client to its ES server.
func ElasticsearchConnectServer(url string) *elastic.Client {

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(url),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC: ", log.LstdFlags)))

	// elastic.SetTraceLog(log.New(os.Stderr, "ELASTIC: ", log.LstdFlags)),
	// elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags))

	// We must do this to allow having an unconnected client instance,
	// and throw proper errors instead of panicking at startup.
	client.Stop()
	elastic.SetHealthcheck(true)(client)
	elastic.SetSniff(true)(client)
	client.Start()

	if err != nil {
		// Don't panic, this might be a temporary failure.
		log.Println("Connection to Elasticsearch failed:", err)
	}

	return client

}

// ElasticsearchRequest sends a POST request to an ES server and parses the returned JSON.
func ElasticsearchRequest(client *elastic.Client, path string, body string) (*elastic.SearchResult, time.Duration, error) {

	if client == nil {
		return nil, 0, elastic.ErrNoClient
	}

	// This is measured on our side in addition to the ElasticSearch-provided SearchResult.TookInMillis
	// It is a good measure of the network & deserialization overhead.
	start := time.Now()

	params := make(url.Values)
	res, err := client.PerformRequest("POST", path, params, body)
	if err != nil {
		return nil, time.Since(start), err
	}

	ret := new(elastic.SearchResult)

	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, time.Since(start), err
	}
	return ret, time.Since(start), nil
}
