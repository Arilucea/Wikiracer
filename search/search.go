package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	gowiki "github.com/trietmn/go-wiki"
)

// Explore articles links by
func Search(startOrEnd string,
	directionChan chan string,
	getLinksFunc func(string) ([]string, error),
	exploredMap *sync.Map,
	parentsMap *sync.Map) {

	currentLevel := []string{startOrEnd}
	for {
		nextLevel := []string{}
		for _, node := range currentLevel {
			exploredMap.Store(node, true)

			links, err := getLinksFunc(node)
			if err != nil {
				continue
			}

			for _, link := range links {
				if strings.Contains(link, ":") {
					continue
				}
				if _, explored := exploredMap.Load(link); !explored {
					if _, parent := parentsMap.Load(link); !parent {
						parentsMap.Store(link, node)
						nextLevel = append(nextLevel, link)
					}
				}
				directionChan <- link
			}
		}
		currentLevel = nextLevel
	}
}

func GetLinks(article string) ([]string, error) {
	println("Getting links for:", article)
	page, err := gowiki.GetPage(article, -1, false, true)
	if err != nil {
		return nil, err
	}
	links, err := page.GetLink()
	if err != nil {
		return nil, err
	}
	return links, nil
}

// JSON response
type ApiResponse struct {
	Continue struct {
		Blcontinue string `json:"blcontinue"`
	} `json:"continue"`
	Query struct {
		Backlinks []struct {
			Title string `json:"title"`
		} `json:"backlinks"`
	} `json:"query"`
}

// Function to fetch backlinks with pagination
func GetBacklinks(article string) ([]string, error) {
	println("Getting backlinks for:", article)
	var backlinks []string
	apiURL := "https://en.wikipedia.org/w/api.php"
	blcontinue := ""

	for {
		params := url.Values{}
		params.Add("action", "query")
		params.Add("format", "json")
		params.Add("list", "backlinks")
		params.Add("bltitle", article)
		params.Add("bllimit", "500")

		if blcontinue != "" {
			params.Add("blcontinue", blcontinue)
		}

		// Make the request to the Wikipedia API
		resp, err := http.Get(apiURL + "?" + params.Encode())
		if err != nil {
			return nil, fmt.Errorf("error fetching data: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		var apiResp ApiResponse
		err = json.Unmarshal(body, &apiResp)
		if err != nil {
			return nil, fmt.Errorf("error parsing JSON: %v", err)
		}

		for _, backlink := range apiResp.Query.Backlinks {
			backlinks = append(backlinks, backlink.Title)
		}

		// Check if there's a continue token; if not, break out of the loop
		blcontinue = apiResp.Continue.Blcontinue
		if blcontinue == "" {
			break
		}
	}

	return backlinks, nil
}
