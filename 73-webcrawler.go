package main

import (
    "fmt"
)

type Fetcher interface {
    // Fetch returns the body of URL and
    // a slice of URLs found on that page.
    Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
    // TODO: Fetch URLs in parallel.
    // TODO: Don't fetch the same URL twice.
    // This implementation doesn't do either:
  if depth <= 0 {
    return
  }
  type fetchResult struct {
    depth int
    urls  []string
  }
  workQueue := make(chan fetchResult)
  getPage := func(url string, currentDepth int) {
    body, urls, err := fetcher.Fetch(url)
    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Printf("found[%d:%s] %q\n", currentDepth, url, body)
    }
    workQueue <- fetchResult{currentDepth+1, urls}
  }

  // I don't like this use of this counter, but don't know channels well enough
  // to do without it.
  pending := 1
  go getPage(url, 0)

  visited := map[string]bool{url:true}
  for pending > 0 {
    next := <- workQueue
    pending--
    if next.depth > depth {
      continue
    }
    for _, url := range next.urls {
      if _, seen := visited[url]; seen {
        continue
      }
      visited[url] = true
      pending++
      go getPage(url, next.depth)
    }
  }
}

func main() {
    Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
    body string
    urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
    if res, ok := f[url]; ok {
        return res.body, res.urls, nil
    }
    return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
    "http://golang.org/": &fakeResult{
        "The Go Programming Language",
        []string{
            "http://golang.org/pkg/",
            "http://golang.org/cmd/",
        },
    },
    "http://golang.org/pkg/": &fakeResult{
        "Packages",
        []string{
            "http://golang.org/",
            "http://golang.org/cmd/",
            "http://golang.org/pkg/fmt/",
            "http://golang.org/pkg/os/",
        },
    },
    "http://golang.org/pkg/fmt/": &fakeResult{
        "Package fmt",
        []string{
            "http://golang.org/",
            "http://golang.org/pkg/",
        },
    },
    "http://golang.org/pkg/os/": &fakeResult{
        "Package os",
        []string{
            "http://golang.org/",
            "http://golang.org/pkg/",
        },
    },
}
