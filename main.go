package main

import (
	"quickcheck/crawler"

	"github.com/davecgh/go-spew/spew"
)

func main() {

	worker := crawler.Worker{}
	worker.Visit("https://www.easyblognetworks.com/")
	spew.Dump(worker.NotFound)
}
