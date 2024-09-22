package util

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

type CustomURL url.URL // want to be able to do this: url.Print()

type Config struct {
	Pages              map[string]int
	MaxPages           int64
	BaseURL            *url.URL
	Mu                 *sync.Mutex
	ConcurrencyControl chan struct{}
	Wg                 *sync.WaitGroup
}

func (cfg *Config) PrintReport(baseURL string) {
	fmt.Fprintf(
		os.Stderr,
		"=============================\n  REPORT for %s\n=============================\n",
		baseURL,
	)

	for url, freq := range cfg.Pages {
		if freq == 1 {
			fmt.Printf("Found 1 internal link to %s\n", url)
			continue
		}

		fmt.Printf("Found %d internal links to %s\n", freq, url)
	}
}
