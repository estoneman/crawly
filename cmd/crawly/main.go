package main

import (
	"fmt"
  "log"
	"os"
	"sync"
  "net/url"

	"github.com/estoneman/crawly/pkg/util"
)

func usage() {
  fmt.Println("usage: crawler <base_url>")
}

func main() {

  args := os.Args[1:]
  lenArgs := len(args)

  if lenArgs == 0 {
    fmt.Println("no website provided")
    usage()

    os.Exit(1)
  } else if lenArgs > 1 {
    fmt.Println("too many arguments provided")
    usage()

    os.Exit(1)
  }

  lookupUrl := args[0]
  parsedLookupURL, err := url.Parse(lookupUrl)
  if err != nil {
    log.Fatalf("failed to parse: %s: %v\n", lookupUrl, err)
  }

  // move to using Config.CrawlPage
  cfg := util.Config{
    Pages: make(map[string]int),
    BaseURL: parsedLookupURL,
    Mu: &sync.Mutex{},
    ConcurrencyControl: make(chan struct{}, 20),
    Wg: &sync.WaitGroup{},
  }

  cfg.CrawlPage(lookupUrl)

  cfg.Wg.Wait()

  fmt.Fprintf(os.Stderr, "\n=== REPORT (%d URLs found) ===\n", len(cfg.Pages))
  for url, freq := range cfg.Pages {
    fmt.Println(url, freq)
  }
}
