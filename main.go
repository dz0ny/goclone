package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"quickcheck/crawler"
)

var mainURL string
var allowedDomains string
var reportPath string
var maxDepth int

func init() {
	flag.StringVar(&mainURL, "url", "https:/google.si", "URL to the website you want to check")
	flag.StringVar(&allowedDomains, "allowed", "", "Comma separated list of allowed domains to crawl")
	flag.IntVar(&maxDepth, "max", 100, "Max depth to crawl")
	flag.StringVar(&reportPath, "report", "/tmp/report.json", "Path to the report file")

}

func main() {
	flag.Parse()
	worker := crawler.Worker{crawler.NewRWMap(), crawler.NewRWMap()}
	worker.Visit(mainURL, allowedDomains, maxDepth)
	fmt.Println(worker.NotFound.List())
	file, _ := json.MarshalIndent(worker.NotFound.List(), "", " ")

	_ = ioutil.WriteFile(reportPath, file, 0644)
}
