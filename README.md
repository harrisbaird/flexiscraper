[![GoDoc](https://godoc.org/github.com/harrisbaird/flexiscraper?status.svg)](https://godoc.org/github.com/harrisbaird/flexiscraper)
[![Build Status](https://travis-ci.org/harrisbaird/flexiscraper.svg?branch=master)](https://travis-ci.org/harrisbaird/flexiscraper)
[![Go Report Card](https://goreportcard.com/badge/github.com/harrisbaird/flexiscraper)](https://goreportcard.com/report/github.com/harrisbaird/flexiscraper)
[![Coverage Status](https://coveralls.io/repos/github/harrisbaird/flexiscraper/badge.svg?branch=master)](https://coveralls.io/github/harrisbaird/flexiscraper?branch=master)

# Flexiscraper

A simple web scraper designed for extracting structured data from a small number of pages.

## Installation

```
go get -u github.com/harrisbaird/flexiscraper
```

### External Dependancies
* [xquery](https://github.com/antchfx/xquery)
* [robotstxt](https://github.com/temoto/robotstxt)

## Usage

```Go
package main

import (
	"log"

	"github.com/harrisbaird/flexiscraper"
	q "github.com/harrisbaird/flexiscraper/q"
)

type HackerNewsItem struct {
	Title  string
	URL    string
	User   string
	Points int
}

func main() {
	scraper := flexiscraper.New()
	hackerNews := scraper.NewDomain("https://news.ycombinator.com/")
	c, err := hackerNews.FetchRoot()
	if err != nil {
		log.Fatalln(err)
	}

	items := []HackerNewsItem{}

	// Iterate through nodes matching XPath expression
	c.Each("//tr[@class=\"athing\"]", func(i int, c *Context) {
		item := HackerNewsItem{
			// Find is a convience function for looking up an XPath expression,
			// returning the first result as a string.
			Title: c.Find(".//td[@class=\"title\"]/a"),
			URL:   c.Find(".//td[@class=\"title\"]/a/@href"),
			User:  c.Find(".//following-sibling::tr//a[@class=\"hnuser\"]"),

			// Create a more complex value and return as an int.
			Points: c.Build(
				q.XPath(".//following-sibling::tr//span[@class=\"age\"]"),
				q.Regexp("\\d+"),
			).Int(),
		}

		items = append(items, item)
	})
}
```
