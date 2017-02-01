package q

import (
	"errors"
	"strconv"
)

var ErrNoMatches = errors.New("No matching queries")

func Build(queries ...ProcessorFunc) *QueryValue {
	out := QueryValue{}
	for _, query := range queries {
		v, err := query(out.Value)
		if err != nil {
			out.Error = err
			break
		}
		out.Value = v
	}
	return &out
}

func Or(values ...*QueryValue) *QueryValue {
	for _, value := range values {
		if value.Error == nil {
			return value
		}
	}
	return &QueryValue{Error: ErrNoMatches}
}

type QueryValue struct {
	Value []string
	Error error
}

func (qv *QueryValue) String() string {
	if len(qv.Value) == 0 {
		return ""
	}

	return qv.Value[0]
}

func (qv *QueryValue) StringSlice() []string {
	return qv.Value
}

func (qv *QueryValue) Int() int {
	if len(qv.IntSlice()) == 0 {
		return 0
	}

	return qv.IntSlice()[0]
}

func (qv *QueryValue) IntSlice() (s []int) {
	for _, value := range qv.Value {
		v, _ := strconv.Atoi(value)
		s = append(s, v)
	}
	return
}
