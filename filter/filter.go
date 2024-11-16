package filter

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/ramirofarias/prop-filter-cli/models"
)

type location = [2]float64

type Comparison struct {
	Operator string
	Value    float64
}

type Filter struct {
	SquareFootage []Comparison
	Bathrooms     []Comparison
	Rooms         []Comparison
	Distance      []Comparison
	Price         []Comparison
	Lat           float64
	Long          float64
	Lighting      string
	Keywords      []string
	Ammenities    []string
}

func FilterProperties(properties []models.Property, filters Filter) []models.Property {
	var filteredProperties []models.Property

Filters:
	for _, property := range properties {
		if len(filters.SquareFootage) > 0 {
			for _, comparison := range filters.SquareFootage {
				if !matchesComparison(comparison, property.SquareFootage) {
					continue Filters
				}
			}
		}
		if len(filters.Bathrooms) > 0 {
			for _, comparison := range filters.Bathrooms {
				if !matchesComparison(comparison, property.Bathrooms) {
					continue Filters
				}
			}
		}
		if len(filters.Rooms) > 0 {
			for _, comparison := range filters.Rooms {
				if !matchesComparison(comparison, property.Rooms) {
					continue Filters
				}
			}
		}
		if len(filters.Distance) > 0 {
			if !(filters.Long == -999999) && !(filters.Lat == -999999) {
				actualDistance := calculateDistance(filters.Lat, filters.Long, property.Location[0], property.Location[1])
				for _, comparison := range filters.Distance {
					if !matchesComparison(comparison, actualDistance) {
						continue Filters
					}
				}

			}
		}

		if len(filters.Price) > 0 {
			for _, comparison := range filters.Price {
				if !matchesComparison(comparison, property.Price) {
					continue Filters
				}
			}
		}

		if filters.Lighting != "" {
			if filters.Lighting != property.Lighting {
				continue Filters
			}
		}

		if len(filters.Keywords) > 0 {
			for _, keyword := range filters.Keywords {
				if !hasKeyword(property.Description, keyword) {
					continue Filters
				}
			}

		}

		if len(filters.Ammenities) > 0 {
			for _, keyword := range filters.Ammenities {
				if !property.Ammenities[keyword] {
					continue Filters
				}
			}

		}

		filteredProperties = append(filteredProperties, property)
	}

	return filteredProperties
}

func matchesComparison(comparison Comparison, prop float64) bool {
	switch comparison.Operator {
	case "lt":
		if !(prop < comparison.Value) {
			return false
		}
	case "gt":
		if !(prop > comparison.Value) {
			return false
		}
	case "gte":
		if !(prop >= comparison.Value) {
			return false
		}
	case "lte":
		if !(prop <= comparison.Value) {
			return false
		}
	case "eq":
		if !(prop == comparison.Value) {
			return false
		}
	default:
		return true
	}

	return true
}

func hasKeyword(s string, k string) bool {
	lowercaseString := strings.ToLower(s)
	pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(k))
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(lowercaseString)
}

func calculateDistance(lat1, long1, lat2, long2 float64) float64 {
	const EARTH_RADIUS = 6371
	const RADIAN = math.Pi / 180

	distance := 0.5 - math.Cos((lat2-lat1)*RADIAN)/2 + math.Cos(lat1*RADIAN)*math.Cos(lat2*RADIAN)*(1-math.Cos((long2-long1)*RADIAN))/2

	return 2 * EARTH_RADIUS * math.Asin(math.Sqrt(distance))
}
