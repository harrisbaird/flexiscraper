package flexiscraper

import "strconv"

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
