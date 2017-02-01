package q

import (
	"errors"
	"fmt"
	"regexp"

	xmlpath "gopkg.in/xmlpath.v2"
)

// ProcessorFunc is the return value for all processors
type ProcessorFunc func([]string) ([]string, error)

// XPath performs an xpath query on the current node.
func XPath(node *xmlpath.Node, xpathExp string) ProcessorFunc {
	return func(values []string) ([]string, error) {
		p, err := xmlpath.Compile(xpathExp)
		if err != nil {
			return values, err
		}

		nodes := p.Iter(node)
		for nodes.Next() {
			values = append(values, nodes.Node().String())
		}

		if len(values) == 0 {
			return values, errors.New("XPath didn't match: " + xpathExp)
		}

		return values, nil
	}
}

// With uses the given input with no processing.
func With(v []string) ProcessorFunc {
	return func(values []string) ([]string, error) {
		return v, nil
	}
}

// Replace calls sprintf using template and previous value in query chain.
func Replace(template string) ProcessorFunc {
	return func(values []string) ([]string, error) {
		for i, v := range values {
			values[i] = fmt.Sprintf(template, v)
		}

		return values, nil
	}
}

// Regexp performs a regular expression on the previous value in query chain,
// returning the first matched string.
func Regexp(r string) ProcessorFunc {
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
