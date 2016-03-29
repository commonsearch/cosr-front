package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/olivere/elastic.v3"
	"net/url"
	"strings"
	"time"
)

// SearchResultTiming is used to measure timings at various steps in the request, in microseconds.
type SearchResultTiming struct {

	// Query times reported by ElasticSearch
	DocsQuery uint32 `json:"dq"`
	TextQuery uint32 `json:"tq"`

	// Total request times measured on our side
	DocsRequest uint32 `json:"dr"`
	TextRequest uint32 `json:"tr"`

	// Total processing time on our end
	Total uint32 `json:"o"`
}

// Hit defines a matching document.
type Hit struct {
	ID      string `json:"i"`
	URL     string `json:"u"`
	Title   string `json:"t"`
	Summary string `json:"s"`
}

// SearchResult defines the result for a query, passed to the template.
type SearchResult struct {
	Hits     []Hit              `json:"h,omitempty"`
	Redirect string             `json:"r,omitempty"`
	HasMore  bool               `json:"m,omitempty"`
	Timing   SearchResultTiming `json:"t,omitempty"`
}

// SearchRequest entirely defines a search request.
type SearchRequest struct {
	Query string `json:"q"`
	Page  int    `json:"p"`
	Lang  string `json:"g"`
}

// Href returns the relative URL of this search.
// Same function is implemented on the JavaScript side
func (req SearchRequest) Href() string {

	var components []string

	if req.Lang != "" {
		components = append(components, fmt.Sprintf("g=%s", req.Lang))
	}

	if req.Page > 1 && req.Query != "" {
		components = append(components, fmt.Sprintf("p=%d", req.Page))
	}

	if req.Query != "" {
		components = append(components, "q="+url.QueryEscape(req.Query))
	}

	if len(components) == 0 {
		return "/"
	}

	return "/?" + strings.Join(components, "&")

}

// PreviousPageHref returns the relative URL of the previous page for this search.
func (req SearchRequest) PreviousPageHref() string {
	if req.Page < 2 {
		return req.Href()
	}
	prev := req
	prev.Page--
	return prev.Href()
}

// NextPageHref returns the relative URL of the next page for this search.
func (req SearchRequest) NextPageHref() string {
	next := req
	next.Page++
	if next.Page == 1 {
		next.Page++
	}
	return next.Href()
}

// BuildTextRequest returns a JSON-encoded Elasticsearch query body for the text index.
func (req SearchRequest) BuildTextRequest() (string, error) {

	var scoringFunctions []string

	jsonQuery, err := json.Marshal(req.Query)
	if err != nil {
		return "", err
	}

	scoringFunctions = append(scoringFunctions, `{
	  	"field_value_factor": {
	      "field": "rank",
	      "factor": 1,
	      "missing": 0
	    }
	}`)

	if req.Lang != "all" {
		scoringFunctions = append(scoringFunctions, fmt.Sprintf(`{
		  	"field_value_factor": {
                "field": "lang_%s",
                "missing": 0.002
            }
		}`, req.Lang))
	}

	// TODO: remove whitespace?
	textEsBody := fmt.Sprintf(`{
      "query": {
        "function_score": {
          "query": {
            "multi_match": {
              "query": %s,
              "minimum_should_match": "-25%%",
              "type": "cross_fields",
	          "tie_breaker": 0.5,
	          "fields": ["title^3", "body", "url_words^2", "domain_words^8"]
            }

          },
          "functions": [%s]
        }
      },
      "from": %d,
      "size": %d
    }`, jsonQuery, strings.Join(scoringFunctions, ","), (req.Page-1)*Config.ResultPageSize, Config.ResultPageSize)

	return textEsBody, nil
}

// BuildDocsRequest returns a JSON-encoded Elasticsearch query body for the docs index.
func BuildDocsRequest(textSearchResult *elastic.SearchResult) string {

	// Collect the IDs
	ids := make([]string, len(textSearchResult.Hits.Hits))
	for i, hit := range textSearchResult.Hits.Hits {
		ids[i] = hit.Id
	}

	return fmt.Sprintf(`{
      "query": {
        "filtered": {
          "filter": {
            "bool": {
              "must": [{
                "ids": {
                  "type": "page",
                  "values": ["%s"]
                }
              }]
            }
          }
        }
      }
    }`, strings.Join(ids, `","`))

}

// PerformSearch performs the search itself and returns a SearchResult.
// We are doing 2 Elasticsearch requests:
//  - First to the "Text" server, to get matching docIDs
//  - Then to the "Docs" server with these IDs, to get the document titles/summaries
func (req SearchRequest) PerformSearch() (*SearchResult, error) {

	page := SearchResult{}

	redirect := DetectBang(req.Query, req.Lang)

	if redirect != "" {
		page.Redirect = redirect
		return &page, nil
	}

	if Config.TestData {
		return req.GenerateTestData(), nil
	}

	textEsBody, err := req.BuildTextRequest()
	if err != nil {
		return nil, err
	}

	// fmt.Println(textEsBody)

	textSearchResult, textRequestTime, err := ElasticsearchRequest(
		ElasticsearchTextClient,
		"/text/page/_search",
		textEsBody)

	if err != nil {
		return nil, err
	}

	page.Timing.TextRequest = uint32(textRequestTime.Seconds() * 1000000)
	page.Timing.TextQuery = uint32(textSearchResult.TookInMillis * 1000)

	// No results!
	if textSearchResult.Hits == nil || len(textSearchResult.Hits.Hits) == 0 {
		return &page, nil
	}

	// TODO: use ES count to determine that
	// TODO: also return textSearchResult.Hits.TotalHits
	page.HasMore = (len(textSearchResult.Hits.Hits) >= Config.ResultPageSize)

	docsEsBody := BuildDocsRequest(textSearchResult)

	docsSearchResult, docsRequestTime, err := ElasticsearchRequest(
		ElasticsearchDocsClient,
		"/docs/page/_search?fields=title,summary,url&size=100",
		docsEsBody)

	if err != nil {
		return nil, err
	}

	page.Timing.DocsRequest = uint32(docsRequestTime.Seconds() * 1000000)
	page.Timing.DocsQuery = uint32(docsSearchResult.TookInMillis * 1000)

	// No results! This shouldn't happen, are we missing documents?
	if docsSearchResult.Hits == nil {
		return &page, nil
	}

	hitsByIds := make(map[string]*Hit, len(docsSearchResult.Hits.Hits))

	// Iterate through results and convert them in their final struct
	for _, hit := range docsSearchResult.Hits.Hits {
		hitsByIds[hit.Id] = &Hit{
			ID:      hit.Id,
			URL:     hit.Fields["url"].([]interface{})[0].(string),
			Title:   hit.Fields["title"].([]interface{})[0].(string),
			Summary: hit.Fields["summary"].([]interface{})[0].(string)}
		
		//Call Highlighung Function	
		hitsByIds[hit.Id].AddHighlighting(req.Query)		
	}

	// Restore the original order of the text results.
	for _, hit := range textSearchResult.Hits.Hits {
		if hitsByIds[hit.Id] != nil {
			page.Hits = append(page.Hits, *hitsByIds[hit.Id])
		}
	}

	return &page, nil
}

//separate Highlighting function for Title and Summary
func (hit Hit) AddHighlighting(query string){
	hit.Title = strings.Replace(hit.Title, " "+query+" ", " <b>"+query+"</b> ",-1)
	hit.Summary = strings.Replace(hit.Summary, " "+query+" ", " <b>"+query+"</b> ",-1)
}

// PerformSearchWithTiming adds a Timing.Total to PerformSearch().
func (req SearchRequest) PerformSearchWithTiming() (*SearchResult, error) {

	start := time.Now()

	page, err := req.PerformSearch()

	if page != nil {
		page.Timing.Total = uint32(time.Since(start).Seconds() * 1000000)
	}

	return page, err
}

// GenerateTestData creates a mock data result for tests
func (req SearchRequest) GenerateTestData() *SearchResult {

	return &SearchResult{
		Hits: []Hit{
			Hit{
				ID:      "1",
				Title:   "Page 1",
				URL:     "http://www.example.com/page/1",
				Summary: "summary 1",
			},
			Hit{
				ID:      "2",
				Title:   "Page 2",
				URL:     "http://www.example.com/page/2",
				Summary: "summary 2",
			},
		},
	}
}
