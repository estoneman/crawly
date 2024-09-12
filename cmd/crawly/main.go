package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/estoneman/crawly/pkg/util"
)

func usage() {
	fmt.Fprintf(os.Stderr, "crawler -url <base_url> [-k <number>]\n")
}

func main() {

	lookupURL := flag.String("url", "", "Base URL")
	maxLookups := flag.Int("k", 10, "Maximum number of concurrent lookups, cannot exceed 30")

	flag.Parse()

	if *lookupURL == "" {
		fmt.Fprintln(os.Stderr, "You must provide a URL to crawl")
		usage()

		os.Exit(1)
	}

	if *maxLookups < 0 || *maxLookups > 30 {
		fmt.Fprintln(os.Stderr, "Please provide a concurrency control value between 0 and 30")
		usage()

		os.Exit(1)
	}

	parsedLookupURL, err := url.Parse(*lookupURL)
	if err != nil {
		log.Fatalf("failed to parse: %s: %v\n", *lookupURL, err)
	}

	// move to using Config.CrawlPage
	cfg := util.Config{
		Pages:              make(map[string]int),
		BaseURL:            parsedLookupURL,
		Mu:                 &sync.Mutex{},
		ConcurrencyControl: make(chan struct{}, *maxLookups),
		Wg:                 &sync.WaitGroup{},
	}

	cfg.CrawlPage(*lookupURL)

	cfg.Wg.Wait()

	fmt.Fprintf(os.Stderr, "\n=== REPORT (%d URLs found) ===\n", len(cfg.Pages))
	for url, freq := range cfg.Pages {
		fmt.Println(url, freq)
	}
}
