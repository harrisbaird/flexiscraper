package flexiscraper

import (
	"errors"

	"github.com/antchfx/xquery/html"
	"github.com/harrisbaird/flexiscraper/q"
	"golang.org/x/net/html"
)

var ErrNoMatches = errors.New("No matching queries")

type Context struct {
	URL    string
	Node   *html.Node
	Errors []error
}

// Find looks up a given xpath expression and returns the first match.
func (c *Context) Find(expr string) string {
	return c.Build(q.XPath(expr)).String()
}

// FindAll looks up a given xpath expression and returns all matches.
func (c *Context) FindAll(expr string) []string {
	return c.Build(q.XPath(expr)).StringSlice()
}

func (c *Context) Attr(expr string) string {
	return htmlquery.SelectAttr(c.Node, expr)
}

// Each finds nodes matching an xpath expression and calls the given function
// for each node.
func (c *Context) Each(expr string, fn func(int, *Context)) {
	htmlquery.FindEach(c.Node, expr, func(i int, node *html.Node) {
		c := Context{URL: c.URL, Node: node}
		fn(i, &c)
	})
}

func (c *Context) Build(input q.InputFunc, processors ...q.ProcessorFunc) *QueryValue {
	out := QueryValue{}

	value, err := input(c.Node)
	if err != nil {
		out.Error = err
		return &out
	}

	out.Value = value

	for _, processor := range processors {
		v, err := processor(out.Value)
		if err != nil {
			out.Error = err
			break
		}
		out.Value = v
	}
	return &out
}

func (c *Context) Or(values ...*QueryValue) *QueryValue {
	for _, value := range values {
		if value.Error == nil {
			return value
		}
	}
	return &QueryValue{Error: ErrNoMatches}
}
