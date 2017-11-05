package flexiscraper

import (
	"errors"

	"github.com/harrisbaird/flexiscraper/q"
	xmlpath "gopkg.in/xmlpath.v2"
)

var ErrNoMatches = errors.New("No matching queries")

type Context struct {
	URL    string
	Node   *xmlpath.Node
	Errors []error
}

// Find looks up a given xpath expression and returns the first match.
func (c *Context) Find(xpathExp string) string {
	return c.Build(q.XPath(xpathExp)).String()
}

// FindAll looks up a given xpath expression and returns all matches.
func (c *Context) FindAll(xpathExp string) []string {
	return c.Build(q.XPath(xpathExp)).StringSlice()
}

// Each finds nodes matching an xpath expression and calls the given function
// for each node.
func (c *Context) Each(xpathExp string, fn func(int, *Context)) {
	list := xmlpath.MustCompile(xpathExp)
	items := list.Iter(c.Node)
	i := 0
	for items.Next() {
		c := Context{Node: items.Node()}
		fn(i, &c)
		i++
	}
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
