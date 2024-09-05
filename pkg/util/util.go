package util

import (
  "fmt"
	"log"
	"net/url"
  "reflect"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"github.com/estoneman/crawly/pkg/types"
)

func (url *types.CustomURL) print() {
	fmt.Printf(`%s {
  Scheme: %s
  Host: %s
  Path: %s
}
`, reflect.TypeOf(url), url.Scheme, url.Host, url.Path)
}

func GetURLsFromHTML(htmlBody, rawBaseUrl string) ([]string, error) {
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

func NormalizeURL(s string) (string, error) {
	parsedUrl, err := url.Parse(s)
	if err != nil {
		log.Printf("%s not parseable, exiting: %v", s, err)
		return "", err
	}

	customUrl := types.CustomURL(*parsedUrl)
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

func (cfg *types.Config) crawlPage(rawCurrentURL string) {
  fmt.Println("implement me")
}

func (cfg *types.Config) addPageVisit(normalizedURL string) (isFirst bool) {
  fmt.Println("implement me")

  return false
}

