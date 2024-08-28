package main

import (
  "fmt"
	"testing"
)

func TestParseHtml(t *testing.T) {
	tests := []struct {
		name      string
		inputURL string
		inputBody string
		expected  []string
	}{
		{
			name: "absolute and relative URLs",
      inputURL: "https://blog.boot.dev",
			inputBody: `
      <html>
      <body>
      <a href="/path/one">
      <span>Boot.dev</span>
      </a>
      <a href="https://other.com/path/one">
      <span>Boot.dev</span>
      </a>
      </body>
      </html>
      `,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
    {
			name: "absolute and relative URLs",
      inputURL: "http://localhost:8000/",
			inputBody: `
      <!DOCTYPE html>
<html lang="en" data-layout="responsive" data-local="">
  <head>
    
    <script>
      window.addEventListener('error', window.__err=function f(e){f.p=f.p||[];f.p.push(e)});
    </script>
    <script>
      (function() {
        const theme = document.cookie.match(/prefers-color-scheme=(light|dark|auto)/)?.[1]
        if (theme) {
          document.querySelector('html').setAttribute('data-theme', theme);
        }
      }())
    </script>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="Description" content="Package html implements an HTML5-compliant tokenizer and parser.">
    
    <meta class="js-gtmID" data-gtmid="GTM-W8MVQXG">
    <link rel="shortcut icon" href="/static/shared/icon/favicon.ico">
    
  
    <link rel="canonical" href="https://pkg.go.dev/golang.org/x/net/html">
  

    <link href="/static/frontend/frontend.min.css?version=prod-frontend-00088-pb4" rel="stylesheet">
    
    <link rel="search" type="application/opensearchdescription+xml" href="/opensearch.xml" title="Go Packages">
    
    
  <title>html package - golang.org/x/net/html - Go Packages</title>

    
  <link href="/static/frontend/unit/unit.min.css?version=prod-frontend-00088-pb4" rel="stylesheet">
  
  <link href="/static/frontend/unit/main/main.min.css?version=prod-frontend-00088-pb4" rel="stylesheet">

  <a href="https://google.com"></a>
  <a href="https://nike.com"></a>
  <a href="https://zyn.com/metrics"></a>
  <a href="https://localhost:8000/metrics"></a>

  </head>
</html>

      `,
      expected: []string{"https://google.com", "https://nike.com", "https://zyn.com/metrics", "https://localhost:8000/metrics"},
    },
		{
			name:      "empty",
			inputBody: "",
			expected:  []string{},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}

      fmt.Println(len(actual))
      for i := 0; i < len(actual); i++ {
        if actual[i] != tc.expected[i] {
				  t.Errorf("Test %v - '%s' FAIL: urls not equal (%s != %s)", i, tc.name, actual[i], tc.expected[i])
        }
      }
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "https",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "http",
			inputURL: "http://www.google.com/metrics",
			expected: "www.google.com/metrics",
		},
		{
			name:     "with port",
			inputURL: "https://www.google.com:8080/metrics",
			expected: "www.google.com:8080/metrics",
		},
		{
			name:     "extra '/'",
			inputURL: "https://www.google.com:8080/metrics/",
			expected: "www.google.com:8080/metrics",
		},
		{
			name:     "no schema",
			inputURL: "/v1/metrics/temp",
			expected: "/v1/metrics/temp",
		},
		{
			name:     "empty",
			inputURL: "",
			expected: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
