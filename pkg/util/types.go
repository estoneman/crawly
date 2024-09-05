package util

import (
  "net/url"
  "sync"
)

type CustomURL url.URL // want to be able to do this: url.Print()

type Config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}
