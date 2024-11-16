package parser

import (
	"reflect"
	"testing"

	"github.com/ramirofarias/prop-filter-cli/filter"
)

func TestParseComparison(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  []filter.Comparison
		expectErr bool
	}{
		{
			name:     "Valid gte comparison",
			input:    "gte 5.5",
			expected: []filter.Comparison{{Operator: "gte", Value: 5.5}},
		},
		{
			name:     "Valid lte comparison",
			input:    "lte 3",
			expected: []filter.Comparison{{Operator: "lte", Value: 3}},
		},
		{
			name:     "Valid eq comparison",
			input:    "eq 10",
			expected: []filter.Comparison{{Operator: "eq", Value: 10}},
		},
		{
			name:     "Valid range comparison",
			input:    "in 2.5, 5.5",
			expected: []filter.Comparison{{Operator: "gte", Value: 2.5}, {Operator: "lte", Value: 5.5}},
		},
		{
			name:      "Invalid operator",
			input:     "asd 3",
			expectErr: true,
		},
		{
			name:      "Invalid number format",
			input:     "gte asdsd",
			expectErr: true,
		},
		{
			name:      "Invalid range format",
			input:     "in 1",
			expectErr: true,
		},
		{
			name:      "More than 2 values in range",
			input:     "in 1, 2, 3",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseComparison(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error but got: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestParseRangeComparison(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  []filter.Comparison
		expectErr bool
	}{
		{
			name:     "Valid range",
			input:    "2.0, 4.0",
			expected: []filter.Comparison{{Operator: "gte", Value: 2.0}, {Operator: "lte", Value: 4.0}},
		},
		{
			name:      "Invalid range with single value",
			input:     "2.0",
			expectErr: true,
		},
		{
			name:      "Invalid range with non-numeric value",
			input:     "2.0, abc",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseRangeComparison(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error but got: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestParseSingleComparison(t *testing.T) {
	tests := []struct {
		name      string
		operator  string
		value     string
		expected  filter.Comparison
		expectErr bool
	}{
		{
			name:     "Valid gte comparison",
			operator: "gte",
			value:    "5.5",
			expected: filter.Comparison{Operator: "gte", Value: 5.5},
		},
		{
			name:      "Invalid number format",
			operator:  "lte",
			value:     "abc",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseSingleComparison(tt.operator, tt.value)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error but got: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}
