package util

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/estoneman/crawly/internal/http_util"
	"golang.org/x/net/html"
)

func (url *CustomURL) print() {
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

	customUrl := CustomURL(*parsedUrl)
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

func (cfg *Config) CrawlPage(rawCurrentURL string) {
	cfg.Wg.Add(1)
	defer cfg.Wg.Done()
	defer func() {
		<-cfg.ConcurrencyControl
	}()

	cfg.Mu.Lock()
	if len(cfg.Pages) >= int(cfg.MaxPages) {
		cfg.Mu.Unlock()
		return
	}
	cfg.Mu.Unlock()

	// parse newly found URL
	parsedRawCurrentURL, err := url.Parse(rawCurrentURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to parse URL: %s: %v\n", rawCurrentURL, err)
		return
	}

	// don't crawl entire internet
	if cfg.BaseURL.Host != parsedRawCurrentURL.Host {
		return
	}

	// normalize current url
	normalizedRawCurrentUrl, err := NormalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to normalize: %s: %v\n", rawCurrentURL, err)
		return
	}

	// check if url has already been seen
	cfg.Mu.Lock()
	if !cfg.addPageVisit(normalizedRawCurrentUrl) {
		cfg.Mu.Unlock()

		return
	}
	cfg.Mu.Unlock()

	body, err := http_util.HttpGet(rawCurrentURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch: %s: %v\n", rawCurrentURL, err)
		return
	}

	links, err := GetURLsFromHTML(body, rawCurrentURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error occurred while searching for hrefs in body of %s: %v\n", rawCurrentURL, err)
		return
	}

	for _, link := range links {
		go cfg.CrawlPage(link)

		// send empty struct
		cfg.ConcurrencyControl <- struct{}{}
	}
}

func (cfg *Config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.Pages[normalizedURL] += 1

	return cfg.Pages[normalizedURL] == 1
}
