package main

import (
  "fmt"

  "github.com/estoneman/crawly/internal/http_client"
)

func main() {
	// s := `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`
  // getURLsFromHTML(s, "https://google.com")
  url := "https://google.com/"
  body := http_client.HttpGet(url)

  links, err := getURLsFromHTML(body, url)
  if err != nil {
    fmt.Printf("error getting URLs from %s\n", url)
  }

  for _, link := range links {
    fmt.Println(link)
  }
}
