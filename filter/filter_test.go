package filter

import (
	"reflect"
	"testing"

	"github.com/ramirofarias/prop-filter-cli/models"
)

func TestFilterProperties(t *testing.T) {
	properties := []models.Property{
		{
			SquareFootage: 1000,
			Bathrooms:     2,
			Location:      [2]float64{40.7128, -74.0060},
			Price:         300000,
			Lighting:      "low",
			Description:   "Spacious and bright apartment",
			Ammenities:    map[string]bool{"pool": true, "gym": true},
		},
		{
			SquareFootage: 750,
			Bathrooms:     1,
			Location:      [2]float64{-34.5749, -58.4303},
			Price:         200000,
			Lighting:      "medium",
			Description:   "Small house",
			Ammenities:    map[string]bool{"gym": true},
		},
	}

	tests := []struct {
		name     string
		filters  Filter
		expected []models.Property
	}{
		{
			name:     "Filter by square footage",
			filters:  Filter{SquareFootage: []Comparison{{Operator: "gte", Value: 800}}},
			expected: []models.Property{properties[0]},
		},
		{
			name:     "Filter by bathrooms",
			filters:  Filter{Bathrooms: []Comparison{{Operator: "gte", Value: 2}}},
			expected: []models.Property{properties[0]},
		},
		{
			name: "Filter by distance",
			filters: Filter{
				Distance: []Comparison{{Operator: "lt", Value: 100}},
				Lat:      -34.548024423566574,
				Long:     -58.70612937569411,
			},
			expected: []models.Property{properties[1]},
		},
		{
			name:     "Filter by price",
			filters:  Filter{Price: []Comparison{{Operator: "lt", Value: 250000}}},
			expected: []models.Property{properties[1]},
		},
		{
			name:     "Filter by lighting",
			filters:  Filter{Lighting: "medium"},
			expected: []models.Property{properties[1]},
		},
		{
			name:     "Filter by keyword",
			filters:  Filter{Keywords: []string{"spacious"}},
			expected: []models.Property{properties[0]},
		},
		{
			name:     "Filter by ammenities",
			filters:  Filter{Ammenities: []string{"pool"}},
			expected: []models.Property{properties[0]},
		},
		{
			name:     "No matches",
			filters:  Filter{SquareFootage: []Comparison{{Operator: "gt", Value: 5000}}},
			expected: []models.Property{},
		},
		{
			name:     "Multi filtering",
			filters:  Filter{Keywords: []string{"spacious"}, Bathrooms: []Comparison{{Operator: "gte", Value: 2}}},
			expected: []models.Property{properties[0]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterProperties(properties, tt.filters)
			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMatchesComparison(t *testing.T) {
	tests := []struct {
		comparison Comparison
		value      float64
		expected   bool
	}{
		{Comparison{Operator: "lt", Value: 10}, 5, true},
		{Comparison{Operator: "gt", Value: 10}, 5, false},
		{Comparison{Operator: "gte", Value: 5}, 5, true},
		{Comparison{Operator: "lte", Value: 5}, 6, false},
		{Comparison{Operator: "eq", Value: 10}, 10, true},
	}

	for _, tt := range tests {
		result := matchesComparison(tt.comparison, tt.value)
		if result != tt.expected {
			t.Errorf("expected %v, got %v", tt.expected, result)
		}
	}
}

func TestHasKeyword(t *testing.T) {
	tests := []struct {
		description string
		keyword     string
		expected    bool
	}{
		{"Foo bar", "foo", true},
		{"This is a test", "foo", false},
		{"This house has a gym", "gym", true},
	}

	for _, tt := range tests {
		result := hasKeyword(tt.description, tt.keyword)
		if result != tt.expected {
			t.Errorf("expected %v, got %v", tt.expected, result)
		}
	}
}
