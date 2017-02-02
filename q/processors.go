package q

import (
	"fmt"
	"regexp"
)

// ProcessorFunc is the return value for all processors
type ProcessorFunc func([]string) ([]string, error)

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
