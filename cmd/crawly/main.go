package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/estoneman/crawly/pkg/util"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:\n\tcrawly <url> <max concurrent lookups> <max pages>")
}

func main() {

	positionalArgs := os.Args[1:]

	if len(positionalArgs) < 3 {
		fmt.Fprintln(os.Stderr, "not enough arguments supplied")
		usage()

		os.Exit(1)
	}

	lookupURL := positionalArgs[0]

	maxLookups, err := strconv.ParseInt(positionalArgs[1], 10, 32)
	if err != nil {
		log.Fatalf("Failed to parse max concurrent lookups: %v\n", err)
	}

	maxPages, err := strconv.ParseInt(positionalArgs[2], 10, 32)
	if err != nil {
		log.Fatalf("Failed to parse max pages: %v\n", err)
	}

	if maxLookups < 0 || maxLookups > 30 {
		fmt.Fprintln(os.Stderr, "Please provide a concurrency control value between 0 and 30\nUsage:")
		usage()

		os.Exit(1)
	}

	parsedLookupURL, err := url.Parse(lookupURL)
	if err != nil {
		log.Fatalf("failed to parse: %s: %v\n", lookupURL, err)
	}

	cfg := util.Config{
		Pages:              make(map[string]int),
		MaxPages:           maxPages,
		BaseURL:            parsedLookupURL,
		Mu:                 &sync.Mutex{},
		ConcurrencyControl: make(chan struct{}, maxLookups),
		Wg:                 &sync.WaitGroup{},
	}

	cfg.CrawlPage(lookupURL)

	cfg.Wg.Wait()

	cfg.PrintReport(lookupURL)
}
