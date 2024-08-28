package http_util

import (
  "io"
  "log"
  "net/http"
)

func HttpGet(url string) string {
  resp, err := http.Get(url)
  if err != nil {
    log.Fatalf("error retrieving web content: %v\n", err)
  }

  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    log.Fatalf("error reading response: %v\n", err)
  }

  return string(body)
}
