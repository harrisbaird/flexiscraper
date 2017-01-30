package q

import (
	"errors"
	"fmt"
	"regexp"

	xmlpath "gopkg.in/xmlpath.v2"
)

var ErrNodeNotSet = errors.New("AddXPath called but ScrapedItem.Node wasn't set")

type QueryFunc func(*xmlpath.Node, string) (string, error)

// Replace calls sprintf using template and previous value in query chain.
func Replace(template string) QueryFunc {
	return func(node *xmlpath.Node, value string) (string, error) {
		if value == "" {
			return "", nil
		}

		return fmt.Sprintf(template, value), nil
	}
}

// Regexp performs a regular expression on the previous value in query chain,
// returning the first matched string.
func Regexp(r string) QueryFunc {
	return func(node *xmlpath.Node, value string) (string, error) {
		rx, err := regexp.Compile(r)
		if err != nil {
			return "", err
		}
		return rx.FindString(value), nil
	}
}

// XPath performs an xpath query on the current node.
func XPath(exp string) QueryFunc {
	return func(node *xmlpath.Node, value string) (string, error) {
		p, err := xmlpath.Compile(exp)
		if err != nil {
			return "", err
		}

		v, ok := p.String(node)
		if !ok {
			return "", errors.New("XPath didn't match: " + exp)
		}

		return v, nil
	}
}
