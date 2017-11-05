package q

import (
	"fmt"

	"github.com/antchfx/xpath"
	"github.com/antchfx/xquery/html"
	"golang.org/x/net/html"
)

// InputFunc is the return value for all inputs.
type InputFunc func(*html.Node) ([]string, error)

// XPath performs an xpath query on the current node.
func XPath(expr string) InputFunc {
	return func(node *html.Node) ([]string, error) {
		output := []string{}

		expr, err := xpath.Compile(expr)
		if err != nil {
			return output, err
		}

		htmlquery.FindEach(node, expr.String(), func(i int, node *html.Node) {
			output = append(output, htmlquery.InnerText(node))
		})

		if len(output) == 0 {
			return output, fmt.Errorf("XPath didn't match: %s", expr)
		}

		return output, nil
	}
}

func Attr(expr string) InputFunc {
	return func(node *html.Node) ([]string, error) {
		expr, err := xpath.Compile(expr)
		if err != nil {
			return []string{}, err
		}

		return []string{htmlquery.SelectAttr(node, expr.String())}, nil
	}
}

// With uses the given input with no processing.
func With(v []string) InputFunc {
	return func(node *html.Node) ([]string, error) {
		return v, nil
	}
}
