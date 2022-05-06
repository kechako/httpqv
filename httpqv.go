// Package httpqv provides parser for Quality values of HTTP header.
package httpqv

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

// A Value is a Quality value.
type Value struct {
	Value    string
	Priority float32
}

// Parse parses Quality values from s.
func Parse(s string) ([]*Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	strs := strings.Split(s, ",")
	values := make([]*Value, len(strs))
	for i, s := range strs {
		value, err := parseValue(s)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}

	return values, nil
}

func parseValue(s string) (*Value, error) {
	v, q, found := strings.Cut(s, ";")

	value := strings.TrimSpace(v)
	if value == "" {
		return nil, errors.New("empty value is found")
	}

	priority := float32(1.0)
	if found {
		key, value, found := strings.Cut(q, "=")
		if !found || strings.TrimSpace(key) != "q" {
			return nil, errors.New("invalid quality value")
		}
		p, err := strconv.ParseFloat(strings.TrimSpace(value), 32)
		if err != nil || p < 0.0 || p > 1.0 {
			return nil, errors.New("invalid quality value")
		}
		priority = float32(p)
	}

	return &Value{
		Value:    value,
		Priority: priority,
	}, nil
}

// Sort sorts the values in order of priority.
func Sort(values []*Value) {
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].Priority > values[j].Priority
	})
}
