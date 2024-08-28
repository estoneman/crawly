package main

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type customURL url.URL // want to be able to do url.Print()

func getURLsFromHTML(htmlBody, rawBaseUrl string) ([]string, error) {
	// for prefixing plain webserver root filepaths
	parsedURL, err := url.Parse(rawBaseUrl)
	if err != nil {
		log.Fatalf("failed to parse: %s (%v)\n", rawBaseUrl, err)
	}

	parsedHtml, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		log.Fatalf("unable to parse html: %v", err)
	}

  var urlFull string
	var schemeURLMatch (*regexp.Regexp)
	var findURLs func(*html.Node)
  var urls []string

	findURLs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					// check if a.Val (a.k.a url) matches r'^\w+:'
					schemeURLMatch = regexp.MustCompile(`^\w+:`)
					match := schemeURLMatch.MatchString(a.Val)

					if !match {
            urlFull = parsedURL.Scheme + "://" + parsedURL.Host + a.Val
            urls = append(urls, urlFull)
					} else {
            urls = append(urls, a.Val)
          }

					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findURLs(c)
		}
	}

	findURLs(parsedHtml)

	return urls, nil
}

func (url *customURL) print() {
	fmt.Printf(`%s {
  Scheme: %s
  Host: %s
  Path: %s
}
`, reflect.TypeOf(url), url.Scheme, url.Host, url.Path)
}

func normalizeURL(s string) (string, error) {
	parsedUrl, err := url.Parse(s)
	if err != nil {
		log.Printf("%s not parseable, exiting: %v", s, err)
		return "", err
	}

	customUrl := customURL(*parsedUrl)
	normalized := customUrl.Host + customUrl.Path

	if len(normalized) == 0 {
		return "", nil
	}

	// check for extraneous '/'
	if normalized[len(normalized)-1] == '/' {
		normalized = normalized[:len(normalized)-1]
	}

	return normalized, nil
}
