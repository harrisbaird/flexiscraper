package flexiscraper

import (
	"github.com/harrisbaird/flexiscraper/q"
	xmlpath "gopkg.in/xmlpath.v2"
)

type Context struct {
	Node   *xmlpath.Node
	Errors []error
}

// Find looks up a given xpath expression and returns the first match.
func (c *Context) Find(xpathExp string) string {
	return q.Build(q.XPath(c.Node, xpathExp)).String()
}

// FindAll looks up a given xpath expression and returns all matches.
func (c *Context) FindAll(xpathExp string) []string {
	return q.Build(q.XPath(c.Node, xpathExp)).StringSlice()
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
