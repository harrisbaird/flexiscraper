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
* [xmlpath](http://gopkg.in/xmlpath.v2)
* [robotstxt](github.com/temoto/robotstxt)

## Usage

```Go
package main

import (
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
    panic(err)
  }

  items := []HackerNewsItem{}

  c.Each("//tr[@class=\"athing\"]", func(i int, c *Context) {
    item := HackerNewsItem{
      Title: c.Find(".//td[@class=\"title\"]/a"),
      URL:   c.Find(".//td[@class=\"title\"]/a/@href"),
      User:  c.Find(".//following-sibling::tr//a[@class=\"hnuser\"]"),
      Points: q.Build(
        q.XPath(c.Node, ".//following-sibling::tr//span[@class=\"age\"]"),
        q.Regexp("(\\d+)"),
      ).Int(),
    }
    items = append(items, item)
  })
}
```
