package search

import (
	"strings"
	"sync"

	gowiki "github.com/Arilucea/go-wiki"
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

func GetBacklinks(article string) ([]string, error) {
	println("Getting backlinks for:", article)
	backlinks, err := gowiki.GetBacklinks(article)
	if err != nil {
		return nil, err
	}

	return backlinks, nil
}
