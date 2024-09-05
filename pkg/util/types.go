package util

import (
  "net/url"
  "sync"
)

type CustomURL url.URL // want to be able to do this: url.Print()

type Config struct {
	Pages              map[string]int
	BaseURL            *url.URL
	Mu                 *sync.Mutex
	ConcurrencyControl chan struct{}
	Wg                 *sync.WaitGroup
}
