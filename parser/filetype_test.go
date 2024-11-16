package parser

import (
	"testing"
)

func TestParseFiletype(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      bool
	}{
		{
			input:    "file.json",
			expected: "json",
			err:      false,
		},
		{
			input:    "file.csv",
			expected: "csv",
			err:      false,
		},
		{
			input:    "file.txt",
			expected: "",
			err:      true,
		},
		{
			input:    "file",
			expected: "",
			err:      true,
		},
		{
			input:    "file.JSON",
			expected: "json",
			err:      false,
		},
		{
			input:    "file.CSV",
			expected: "csv",
			err:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ParseFiletype(test.input)
			if test.err && err == nil {
				t.Errorf("expected error, got nil for input %s", test.input)
			}
			if !test.err && err != nil {
				t.Errorf("unexpected error for input %s: %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("for input %s: expected %s, got %s", test.input, test.expected, result)
			}
		})
	}
}
