package q_test

import (
	"strings"
	"testing"

	"github.com/antchfx/xquery/html"
	. "github.com/harrisbaird/flexiscraper/q"
	"github.com/nbio/st"
)

func TestXPath(t *testing.T) {
	r := strings.NewReader("<!DOCTYPE html><html><head><title>Hello world</title></head><body><h1>Test</h1></body></html>")
	node, err := htmlquery.Parse(r)
	st.Assert(t, err, nil)

	tests := []struct {
		name    string
		exp     string
		result  []string
		wantErr bool
	}{
		{"blank", "", []string{}, true},
		{"valid", "//title", []string{"Hello world"}, false},
		{"invalid", "~~Test", []string{}, true},
		{"no match", "//invalid", []string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := XPath(tt.exp)(node)
			st.Assert(t, err != nil, tt.wantErr)
			st.Assert(t, result, tt.result)
		})
	}
}

func TestWith(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		result  []string
		wantErr bool
	}{
		{"blank", []string{""}, []string{""}, false},
		{"valid", []string{"value"}, []string{"value"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := With(tt.input)(nil)
			st.Assert(t, err != nil, tt.wantErr)
			st.Assert(t, result, tt.result)
		})
	}
}
