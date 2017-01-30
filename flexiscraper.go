package flexiscraper

import (
	"errors"
	"net/http"

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
	node, err := xmlpath.ParseHTML(res.Body)
	if err != nil {
		return nil, err
	}
	return &Context{Node: node}, nil
}

type Context struct {
	Node   *xmlpath.Node
	Errors []error
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
	if q.Error != nil {
		return ""
	}

	return q.Value
}
