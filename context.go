package flexiscraper

import (
	"github.com/harrisbaird/flexiscraper/q"
	xmlpath "gopkg.in/xmlpath.v2"
)

type Context struct {
	Node   *xmlpath.Node
	Errors []error
}

// Find looks up an xpath expression and returns the first match as a string.
func (c *Context) Find(exp string) string {
	return q.Build(q.XPath(c.Node, exp)).String()
}

// Find looks up an xpath expression and returns all matches.
func (c *Context) FindAll(exp string) []string {
	return q.Build(q.XPath(c.Node, exp)).StringSlice()
}

func (c *Context) Each(sel string, fn func(int, *Context)) {
	list := xmlpath.MustCompile(sel)
	items := list.Iter(c.Node)
	i := 0
	for items.Next() {
		c := Context{Node: items.Node()}
		fn(i, &c)
		i++
	}
}
