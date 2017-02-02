package q

import (
	"errors"

	xmlpath "gopkg.in/xmlpath.v2"
)

// InputFunc is the return value for all inputs.
type InputFunc func(*xmlpath.Node) ([]string, error)

// XPath performs an xpath query on the current node.
func XPath(xpathExp string) InputFunc {
	return func(node *xmlpath.Node) ([]string, error) {
		output := []string{}
		p, err := xmlpath.Compile(xpathExp)
		if err != nil {
			return output, err
		}

		nodes := p.Iter(node)
		for nodes.Next() {
			output = append(output, nodes.Node().String())
		}

		if len(output) == 0 {
			return output, errors.New("XPath didn't match: " + xpathExp)
		}

		return output, nil
	}
}

// With uses the given input with no processing.
func With(v []string) InputFunc {
	return func(node *xmlpath.Node) ([]string, error) {
		return v, nil
	}
}
