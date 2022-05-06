package httpqv

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var parseTests = []struct {
	str     string
	values  []*Value
	failure bool
}{
	{
		str: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		values: []*Value{
			{Value: "text/html", Priority: 1.0},
			{Value: "application/xhtml+xml", Priority: 1.0},
			{Value: "application/xml", Priority: 0.9},
			{Value: "*/*", Priority: 0.8},
		},
	},
	{
		str: "text/html, application/xhtml+xml ,application/xml;q =0.9,*/*;q= 0.8",
		values: []*Value{
			{Value: "text/html", Priority: 1.0},
			{Value: "application/xhtml+xml", Priority: 1.0},
			{Value: "application/xml", Priority: 0.9},
			{Value: "*/*", Priority: 0.8},
		},
	},
	{
		str:    "",
		values: nil,
	},
	{
		str:     "text/html,application/xhtml+xml,application/xml;q=1.1,*/*;q=0.8",
		failure: true,
	},
	{
		str:     "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=-0.8",
		failure: true,
	},
	{
		str:     "text/html,,application/xhtml+xml,application/xml;s=0.9,*/*;q=0.8",
		failure: true,
	},
	{
		str:     "text/html,,application/xhtml+xml,application/xml;q=,*/*;q=0.8",
		failure: true,
	},
	{
		str:     "text/html,,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		failure: true,
	},
	{
		str:     "text/html,,application/xhtml+xml,application/xml;;q=0.9,*/*;q=0.8",
		failure: true,
	},
}

func TestParse(t *testing.T) {
	for i, tt := range parseTests {
		name := fmt.Sprintf("case %d", i+1)
		t.Run(name, func(t *testing.T) {
			values, err := Parse(tt.str)
			if tt.failure {
				if err == nil {
					t.Errorf("Parse(%#v): shoud return error", tt.str)
				}
			} else {
				if err != nil {
					t.Errorf("Parse(%#v): returns error: %v", tt.str, err)
				}

				if diff := cmp.Diff(values, tt.values); diff != "" {
					t.Errorf("Parse(%#v), differs: (-got +want)\n%s", tt.str, diff)
				}
			}
		})
	}
}

var sortTests = []struct {
	values []*Value
	sorted []*Value
}{
	{
		values: []*Value{
			{Value: "application/xml", Priority: 0.9},
			{Value: "*/*", Priority: 0.8},
			{Value: "text/html", Priority: 1.0},
			{Value: "application/xhtml+xml", Priority: 1.0},
		},
		sorted: []*Value{
			{Value: "text/html", Priority: 1.0},
			{Value: "application/xhtml+xml", Priority: 1.0},
			{Value: "application/xml", Priority: 0.9},
			{Value: "*/*", Priority: 0.8},
		},
	},
	{
		values: []*Value{
			{Value: "text/html", Priority: 1.0},
			{Value: "application/xhtml+xml", Priority: 1.0},
			{Value: "application/xml", Priority: 0.9},
			{Value: "*/*", Priority: 0.8},
		},
		sorted: []*Value{
			{Value: "text/html", Priority: 1.0},
			{Value: "application/xhtml+xml", Priority: 1.0},
			{Value: "application/xml", Priority: 0.9},
			{Value: "*/*", Priority: 0.8},
		},
	},
	{
		values: []*Value{
			{Value: "text/json", Priority: 0.7},
			{Value: "*/*", Priority: 0.8},
			{Value: "text/html", Priority: 1.0},
			{Value: "application/xml", Priority: 0.9},
			{Value: "application/xhtml+xml", Priority: 1.0},
			{Value: "application/json", Priority: 0.7},
		},
		sorted: []*Value{
			{Value: "text/html", Priority: 1.0},
			{Value: "application/xhtml+xml", Priority: 1.0},
			{Value: "application/xml", Priority: 0.9},
			{Value: "*/*", Priority: 0.8},
			{Value: "text/json", Priority: 0.7},
			{Value: "application/json", Priority: 0.7},
		},
	},
}

func TestSort(t *testing.T) {
	for i, tt := range sortTests {
		name := fmt.Sprintf("case %d", i+1)
		t.Run(name, func(t *testing.T) {
			Sort(tt.values)
			if diff := cmp.Diff(tt.values, tt.sorted); diff != "" {
				t.Errorf("Sort(), differs: (-got +want)\n%s", diff)
			}
		})
	}
}
