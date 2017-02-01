package q

import (
	"errors"
	"fmt"
	"regexp"

	xmlpath "gopkg.in/xmlpath.v2"
)

var ErrNodeNotSet = errors.New("AddXPath called but ScrapedItem.Node wasn't set")

type QueryFunc func([]string) ([]string, error)

// XPath performs an xpath query on the current node.
func XPath(node *xmlpath.Node, exp string) QueryFunc {
	return func(values []string) ([]string, error) {
		p, err := xmlpath.Compile(exp)
		if err != nil {
			return values, err
		}

		nodes := p.Iter(node)
		for nodes.Next() {
			values = append(values, nodes.Node().String())
		}

		if len(values) == 0 {
			return values, errors.New("XPath didn't match: " + exp)
		}

		return values, nil
	}
}

// Use the input value with no processing.
func With(v []string) QueryFunc {
	return func(values []string) ([]string, error) {
		return v, nil
	}
}

// Replace calls sprintf using template and previous value in query chain.
func Replace(template string) QueryFunc {
	return func(values []string) ([]string, error) {
		for i, v := range values {
			values[i] = fmt.Sprintf(template, v)
		}

		return values, nil
	}
}

// Regexp performs a regular expression on the previous value in query chain,
// returning the first matched string.
func Regexp(r string) QueryFunc {
	return func(values []string) ([]string, error) {
		for i, v := range values {
			rx, err := regexp.Compile(r)
			if err != nil {
				return values, err
			}
			values[i] = rx.FindString(v)
		}
		return values, nil
	}
}
