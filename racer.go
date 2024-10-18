package main

import (
	"flag"
	"fmt"
	"net/url"
	"strings"
	"sync"

	s "github.com/arilucea/wikiracer/search"
)

var forwardExplored sync.Map
var backwardExplored sync.Map
var forwardParents sync.Map
var backwardParents sync.Map

func initBidirectionalSearch(start, end string) []string {
	if start == end {
		return []string{start}
	}

	forwardChan := make(chan string)
	backwardChan := make(chan string)

	go s.Search(start, forwardChan, s.GetLinks, &forwardExplored, &forwardParents)
	go s.Search(end, backwardChan, s.GetBacklinks, &backwardExplored, &backwardParents)

	for {
		select {
		case fwd := <-forwardChan:
			if _, ok := backwardExplored.Load(fwd); ok {
				return reconstructFullPath(fwd, &forwardParents, &backwardParents)
			}

		case bwd := <-backwardChan:
			if _, ok := forwardExplored.Load(bwd); ok {
				return reconstructFullPath(bwd, &forwardParents, &backwardParents)
			}
		}
	}
}

func reconstructFullPath(meetingPoint string, forwardParents, backwardParents *sync.Map) []string {
	var forwardPath []string
	for node := meetingPoint; node != ""; {
		if parent, ok := forwardParents.Load(node); ok {
			node = parent.(string)
			forwardPath = append(forwardPath, node)
		} else {
			break
		}
	}
	reverseSlice(forwardPath)
	forwardPath = append(forwardPath, meetingPoint)

	var backwardPath []string
	for node := meetingPoint; node != ""; {
		if parent, ok := backwardParents.Load(node); ok {
			node = parent.(string)
			forwardPath = append(forwardPath, node)
		} else {
			break
		}
	}

	return append(forwardPath, backwardPath...)
}

func reverseSlice(path []string) {
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
}

func main() {
	wikipediaURL := "https://en.wikipedia.org/wiki/"

	start := flag.String("start", "https://en.wikipedia.org/wiki/Uruguay", "a string")
	end := flag.String("end", "https://en.wikipedia.org/wiki/La_Rioja", "a string")
	flag.Parse()

	decodedStart, err := url.QueryUnescape(*start)
	if err != nil {
		fmt.Println("Error decoding URL:", err)
		return
	}
	decodedEnd, err := url.QueryUnescape(*end)
	if err != nil {
		fmt.Println("Error decoding URL:", err)
		return
	}

	start_article := strings.Split(decodedStart, wikipediaURL)
	end_article := strings.Split(decodedEnd, wikipediaURL)

	start_article_normalize := strings.ReplaceAll(start_article[1], "_", " ")
	end_article_normalize := strings.ReplaceAll(end_article[1], "_", " ")

	path := initBidirectionalSearch(start_article_normalize, end_article_normalize)

	fmt.Println("\n\n\nPath found: ")
	for _, article := range path {
		fmt.Println(wikipediaURL + article)
	}

}
