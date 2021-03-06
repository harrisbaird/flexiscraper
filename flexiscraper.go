package flexiscraper

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/antchfx/xquery/html"
	"github.com/temoto/robotstxt"
)

// DefaultUserAgent is the default user agent string. It's used in all http
// requests and during robots.txt validation.
const DefaultUserAgent = "Flexiscraper (https://github.com/harrisbaird/flexiscraper)"

// ErrDisallowedByRobots is returned when the requested URL is disallowed
// by robots.txt.
var ErrDisallowedByRobots = errors.New("HTTP request disallowed by robots.txt")

// Path used for robots.txt fetching.
var robotsTxtParsedPath, _ = url.Parse("/robots.txt")

// New initialises a new Scraper.
func New() *Scraper {
	return &Scraper{
		ObeyRobots: true,
		UserAgent:  DefaultUserAgent,
		HTTPClient: http.DefaultClient,
	}
}

// A Scraper defines the parameters for running a web scraper.
type Scraper struct {
	// The user agent string sent during http requests and when checking
	// robots.txt.
	UserAgent string

	// The http client to use when fetching, defaults to http.DefaultClient.
	HTTPClient *http.Client

	// ObeyRobots enables robot.txt policy checking.
	// Default: true
	ObeyRobots bool
}

// NewDomain initialises a new domain.
// This is used for robots.txt and to ensure absolute urls.
func (s *Scraper) NewDomain(baseDomain string) *Domain {
	domainURL, _ := url.Parse(baseDomain)

	domain := Domain{
		Scraper: s,
		Domain:  domainURL,
	}

	if s.ObeyRobots {
		robotsData, err := s.getRobots(domainURL)
		if err != nil {
			fmt.Println("flexiscraper: failed to get robots.txt for " + baseDomain)
		}
		domain.RobotsData = robotsData
	}

	return &domain
}

func (s *Scraper) getRobots(domain *url.URL) (*robotstxt.Group, error) {
	robotsPath := domain.ResolveReference(robotsTxtParsedPath)
	res, err := s.getRequest(robotsPath.String())
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	data, err := robotstxt.FromResponse(res)
	if err != nil {
		return nil, err
	}
	return data.FindGroup(s.UserAgent), nil
}

func (s *Scraper) getRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", s.UserAgent)
	return s.HTTPClient.Do(req)
}

// Domain defines implementation for scraping a single domain.
type Domain struct {
	*Scraper
	Domain     *url.URL
	RobotsData *robotstxt.Group
}

// Fetch and parse html from the given URL, checks and obeys robots.txt if
// ObeyRobots is true in the scraper.
func (d *Domain) Fetch(url string) (*Context, error) {
	context := &Context{URL: d.makeAbsoluteURL(url)}

	if d.ObeyRobots {
		if !d.RobotsData.Test(url) {
			return context, ErrDisallowedByRobots
		}
	}

	res, err := d.getRequest(context.URL)
	if err != nil {
		return context, err
	}
	defer res.Body.Close()
	err = d.Parse(context, res.Body)
	return context, err
}

// FetchRoot convinience function for fetching the current domains root URL.
func (d *Domain) FetchRoot() (*Context, error) {
	return d.Fetch(d.Domain.String())
}

// Parse html from the given reader.
func (d *Domain) Parse(context *Context, r io.Reader) error {
	node, err := htmlquery.Parse(r)
	context.Node = node
	return err
}

func (d *Domain) makeAbsoluteURL(currentURL string) string {
	parsed, err := url.Parse(currentURL)
	if err != nil {
		return currentURL
	}

	if !parsed.IsAbs() {
		return d.Domain.ResolveReference(parsed).String()
	}

	return currentURL
}
