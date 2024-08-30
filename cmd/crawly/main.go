package main

import (
	"fmt"
  "os"

	"github.com/estoneman/crawly/internal/http_util"
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

  url := args[0]
  body, err := http_util.HttpGet(url)
  if err != nil {
    fmt.Printf("error retrieving content of '%s': %v\n", url, err)
    os.Exit(1)
  }

  links, err := util.GetURLsFromHTML(body, url)
  if err != nil {
    fmt.Printf("error getting URLs from %s\n", url)
    os.Exit(1)
  }

  fmt.Printf("found %d URLs\n", len(links))
  for _, link := range links {
    fmt.Println(link)
  }
}
