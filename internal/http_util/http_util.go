package http_util

import (
  "fmt"
  "io"
  "log"
  "net/http"
  "strings"
)

func HttpGet(url string) (string, error) {
  resp, err := http.Get(url)
  if err != nil {
    log.Fatalf("error retrieving web content: %v\n", err)
  }

  if resp.StatusCode >= http.StatusBadRequest {
    return "", fmt.Errorf("HTTP GET failed with code: %d\n", resp.StatusCode)
  }

  if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
    return "", fmt.Errorf("Content-Type header is not 'text/html' (%s)", resp.Header.Get("Content-Type"))
  }

  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    log.Fatalf("error reading response: %v\n", err)
  }

  return string(body), nil
}
