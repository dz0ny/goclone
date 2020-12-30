package crawler

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/gocolly/colly"
	"github.com/lukasbob/srcset"
)

type Worker struct {
	NotFound []string
	Done     []string
}

func (w *Worker) AddDone(url string) bool {
	if w.containsURL(w.Done, url) {
		return false
	}
	w.Done = append(w.Done, url)
	return true
}

func (w *Worker) AddNotFound(url string) bool {
	if w.containsURL(w.Done, url) {
		return false
	}
	w.Done = append(w.Done, url)
	return true
}

func (w *Worker) containsURL(list []string, searchterm string) bool {
	i := sort.SearchStrings(list, searchterm)
	return i < len(list) && list[i] == searchterm
}

// Collector searches for css, js, and images within a given link
func (w *Worker) Visit(target string) {
	// create a new collector

	uri, _ := url.Parse(target)

	c := colly.NewCollector(
		// asychronus boolean
		colly.Async(true),
		colly.CacheDir("/tmp"),
		colly.AllowedDomains(uri.Hostname(), "cdn.paddle.com", "assets.calendly.com", "static.easyblognetworks.com", "fanstatic.niteo.co", "blog.easyblognetworks.com"),
	)

	// search for all link tags that have a rel attribute that is equal to stylesheet - CSS
	c.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
		// hyperlink reference
		link := e.Attr("href")
		// print css file was found
		//fmt.Println("Css found", "-->", link)
		if w.AddDone(link) {
			if err := c.Head(link); err != nil {
				fmt.Println(err, link)
			}
		}
	})

	// search for all script tags with src attribute -- JS
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		// Print link
		//fmt.Println("Js found", "-->", link)
		if w.AddDone(link) {
			if err := c.Head(link); err != nil {
				fmt.Println(err, link)
			}
		}
	})

	// serach for all img tags with src attribute -- Images
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		// Print link
		//fmt.Println("Img found", "-->", link)
		if w.AddDone(link) {
			if err := c.Head(link); err != nil {
				fmt.Println(err, link)
			}
		}
		sets := e.Attr("srcset")
		if sets != "" {
			for _, img := range srcset.Parse(sets) {

				if err := c.Head(img.URL); err != nil {
					fmt.Println(err, link)
				}

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

	c.OnError(func(res *colly.Response, err error) {
		if err.Error() == "Not Found" {
			w.AddNotFound(res.Request.URL.String())
		}
		fmt.Println(err, res.Request.URL.String())

	})
	// Visit each url and wait for stuff to load :)
	c.Visit(target)
	c.Wait()
}
