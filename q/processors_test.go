package q_test

import (
	"testing"

	. "github.com/harrisbaird/flexiscraper/q"
	"github.com/nbio/st"
)

func TestReplace(t *testing.T) {
	tests := []struct {
		name     string
		template string
		value    []string
		result   []string
		wantErr  bool
	}{
		{"blank", "", []string{}, []string{}, false},
		{"valid", "hello %s", []string{"world"}, []string{"hello world"}, false},
		{"missing", "hello", []string{"world"}, []string{"hello%!(EXTRA string=world)"}, false},
		{"extra", "hello %s, %s", []string{"world"}, []string{"hello world, %!s(MISSING)"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Replace(tt.template)(tt.value)
			st.Assert(t, err != nil, tt.wantErr)
			st.Assert(t, result, tt.result)
		})
	}
}

func TestRegexp(t *testing.T) {
	tests := []struct {
		name    string
		regex   string
		value   []string
		result  []string
		wantErr bool
	}{
		{"blank", "", []string{""}, []string{""}, false},
		{"invalid", "(.*", []string{"value"}, []string{"value"}, true},
		{"valid", "\\d+", []string{"Post 5"}, []string{"5"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Regexp(tt.regex)(tt.value)
			st.Assert(t, err != nil, tt.wantErr)
			st.Assert(t, result, tt.result)
		})
	}
}
