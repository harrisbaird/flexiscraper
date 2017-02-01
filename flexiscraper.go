package flexiscraper

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/temoto/robotstxt"

	xmlpath "gopkg.in/xmlpath.v2"
)

// DefaultUserAgent is the default user agent string.
const DefaultUserAgent = "Flexiscraper (https://github.com/harrisbaird/flexiscraper)"

var ErrDisallowedByRobots = errors.New("HTTP request disallowed by robots.txt")

var robotsTxtParsedPath, _ = url.Parse("/robots.txt")

func New() *Scraper {
	return &Scraper{
		ObeyRobots: true,
		UserAgent:  DefaultUserAgent,
		HTTPClient: http.DefaultClient,
	}
}

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

type Domain struct {
	*Scraper
	Domain     *url.URL
	RobotsData *robotstxt.Group
}

func (d *Domain) Fetch(currentURL string) (*Context, error) {
	if d.ObeyRobots {
		if !d.RobotsData.Test(currentURL) {
			return nil, ErrDisallowedByRobots
		}
	}

	currentURL = d.ensureAbsolute(currentURL)

	res, err := d.getRequest(currentURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return d.Parse(res.Body)
}

func (d *Domain) Parse(r io.Reader) (*Context, error) {
	node, err := xmlpath.ParseHTML(r)
	if err != nil {
		return nil, err
	}

	return &Context{Node: node}, nil
}

func (d *Domain) ensureAbsolute(currentURL string) string {
	parsed, err := url.Parse(currentURL)
	if err != nil {
		return currentURL
	}

	if !parsed.IsAbs() {
		return d.Domain.ResolveReference(parsed).String()
	}

	return currentURL
}
