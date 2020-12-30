package crawler

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/lukasbob/srcset"
)

type RWMap struct {
	sync.RWMutex
	m map[string]int
}

func NewRWMap() RWMap {
	return RWMap{m: make(map[string]int)}
}

// Get is a wrapper for getting the value from the underlying map
func (r RWMap) Get(key string) int {
	r.RLock()
	defer r.RUnlock()
	return r.m[key]
}

// Set is a wrapper for setting the value of a key in the underlying map
func (r RWMap) Set(key string, val int) {
	r.Lock()
	defer r.Unlock()
	r.m[key] = val
}

// Inc increases the value in the RWMap for a key.
//   This is more pleasant than r.Set(key, r.Get(key)++)
func (r RWMap) Inc(key string) {
	r.Lock()
	defer r.Unlock()
	r.m[key]++
}

func (r RWMap) Has(key string) bool {
	r.Lock()
	defer r.Unlock()
	return r.m[key] > 0
}

func (r RWMap) List() map[string]int {
	return r.m
}

type Worker struct {
	NotFound RWMap
	Done     RWMap
}

func (w *Worker) AddDone(url string) bool {
	if w.Done.Has(url) {
		w.Done.Inc(url)
		return false
	}
	w.Done.Set(url, 1)
	return true
}

func (w *Worker) AddNotFound(url string) bool {
	if w.NotFound.Has(url) {
		w.NotFound.Inc(url)
		return false
	}
	w.NotFound.Set(url, 1)
	return true
}

func (w *Worker) Fetch(c *colly.Collector, link string) {
	if w.AddDone(link) {
		if err := c.Head(link); err != nil {
			fmt.Println(err, link)
			if err.Error() == "Not Found" {
				w.AddNotFound(link)
			}
		}
	}
}

func (w *Worker) Visit(target, allowed string, maxDepth int) {
	// create a new collector
	uri, _ := url.Parse(target)
	allowedP := strings.Split(allowed, ",")
	allowedP = append(allowedP, uri.Hostname())

	c := colly.NewCollector(
		// asychronus boolean
		colly.Async(false),
		colly.CacheDir("/tmp"),
		colly.AllowedDomains(allowedP...),
		colly.MaxDepth(maxDepth),
	)

	// search for all link tags that have a rel attribute that is equal to stylesheet - CSS
	c.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
		// hyperlink reference
		link := e.Attr("href")
		w.Fetch(c, link)
	})

	// search for all script tags with src attribute -- JS
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		w.Fetch(c, link)
	})

	// serach for all img tags with src attribute -- Images
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		w.Fetch(c, link)
		sets := e.Attr("srcset")
		if sets != "" {
			for _, img := range srcset.Parse(sets) {
				w.Fetch(c, img.URL)
			}
		}

	})

	// serach for all img tags with src attribute -- Images
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("href")

		// Skip fragments
		if strings.HasPrefix(link, "#") {
			return
		}

		if parsed, err := uri.Parse(link); err == nil {
			newURI := uri.ResolveReference(parsed)
			if uri.Hostname() == newURI.Hostname() && newURI.Hostname() != "" {
				if w.AddDone(newURI.String()) {
					c.Visit(newURI.String())
				}

			}
		}

	})
	c.OnRequest(func(c *colly.Request) {
		fmt.Println(c.URL)
	})

	// Visit each url and wait for stuff to load :)
	c.Visit(target)
	c.Wait()
}
