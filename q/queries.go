package q

import (
	"errors"
	"fmt"
	"regexp"

	xmlpath "gopkg.in/xmlpath.v2"
)

var ErrNodeNotSet = errors.New("AddXPath called but ScrapedItem.Node wasn't set")

type QueryFunc func(*xmlpath.Node, string) (string, error)

func Replace(template string) QueryFunc {
	return func(node *xmlpath.Node, value string) (string, error) {
		return fmt.Sprintf(template, value), nil
	}
}

func Regexp(r string) QueryFunc {
	return func(node *xmlpath.Node, value string) (string, error) {
		rx, err := regexp.Compile(r)
		if err != nil {
			return "", err
		}
		return rx.FindString(value), nil
	}
}

func XPath(sel string) QueryFunc {
	return func(node *xmlpath.Node, value string) (string, error) {
		p, err := xmlpath.Compile(sel)
		if err != nil {
			return "", err
		}

		v, ok := p.String(node)
		if !ok {
			return "", errors.New("XPath didn't match: " + sel)
		}

		return v, nil
	}
}
