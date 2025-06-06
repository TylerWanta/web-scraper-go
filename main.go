package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"webScraper/package/urlTracker"

	"github.com/PuerkitoBio/goquery"
)

var log = false

func main() {
	log = len(os.Args) > 1 && os.Args[1] == "--log"

	var wg sync.WaitGroup
	urlTracker := urlTracker.NewURlTracker(100)

	wg.Add(1)

	start := time.Now()
	crawl("https://google.com", urlTracker, &wg)

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Println(urlTracker)
	fmt.Printf("Scraped %v URLs in %v seconds", urlTracker.Length(), elapsed.Seconds())
}

func crawl(url string, urlTracker *urlTracker.URLTracker, wg *sync.WaitGroup) {
	defer wg.Done()

	res, err := http.Get(url)

	if err != nil {
		if log {
			fmt.Printf("Error for url: %v. %v\n", url, err)
		}

		return
	}

	defer res.Body.Close()
	if ok := urlTracker.Add(url); !ok {
		return
	}

	if res.StatusCode != 200 {
		if log {
			fmt.Printf("Status code error: %v. %v - %v\n", url, res.StatusCode, res.Status)
		}

		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		if log {
			fmt.Printf("Error when getting document: %v\n", err)
		}

		return
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		linkURL := s.AttrOr("href", "")

		if linkURL == "" || !strings.HasPrefix(linkURL, "https") {
			return
		}

		wg.Add(1)
		crawl(linkURL, urlTracker, wg)
	})
}
