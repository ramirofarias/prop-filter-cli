package parser

import (
	"testing"
)

func TestParseText(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "test1, test2, test3",
			expected: []string{"test1", "test2", "test3"},
		},
		{
			input:    "   a,b ,   c  ",
			expected: []string{"a", "b", "c"},
		},
		{
			input:    "a",
			expected: []string{"a"},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := ParseText(test.input)
			if len(result) != len(test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
			for i, word := range result {
				if word != test.expected[i] {
					t.Errorf("at index %d: expected %s, got %s", i, test.expected[i], word)
				}
			}
		})
	}
}
