package flexiscraper

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/harrisbaird/flexiscraper/q"
	xmlpath "gopkg.in/xmlpath.v2"
)

var ErrZeroMatches = errors.New("No matching queries")

func Fetch(url string) (*Context, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return Parse(res.Body)
}

func Parse(r io.Reader) (*Context, error) {
	node, err := xmlpath.ParseHTML(r)
	if err != nil {
		return nil, err
	}

	return &Context{Node: node}, nil
}

type Context struct {
	Node   *xmlpath.Node
	Errors []error
}

// Find looks up an xpath expression and returns the first match as a string.
func (c *Context) Find(exp string) string {
	return c.Build(q.XPath(exp)).String()
}

func (c *Context) Build(queries ...q.QueryFunc) *QueryValue {
	out := QueryValue{}
	for _, query := range queries {
		v, err := query(c.Node, out.Value)
		if err != nil {
			out.Error = err
			c.Errors = append(c.Errors, err)
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
	return &QueryValue{Error: ErrZeroMatches}
}

func (c *Context) Loop(sel string, fn func(int, *Context)) {
	list := xmlpath.MustCompile(sel)
	items := list.Iter(c.Node)
	i := 0
	for items.Next() {
		c := Context{Node: items.Node()}
		fn(i, &c)
		i++
	}
}

type QueryValue struct {
	Value string
	Error error
}

func (q *QueryValue) String() string {
	return q.Value
}

func (q *QueryValue) Int() int {
	v, _ := strconv.Atoi(q.Value)
	return v
}
