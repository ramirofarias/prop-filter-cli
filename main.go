package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Property struct {
	SquareFootage float64         `json:"squareFootage"`
	Lighting      string          `json:"lighting"`
	Price         float64         `json:"price"`
	Rooms         float64         `json:"rooms"`
	Bathrooms     float64         `json:"bathrooms"`
	Location      [2]float64      `json:"location"`
	Description   string          `json:"description"`
	Ammenities    map[string]bool `json:"ammenities"`
}

type Filter struct {
	SquareFootage []Comparison
	Bathrooms     []Comparison
	Distance      []Comparison
	Price         []Comparison
	Lat           float64
	Long          float64
	Lighting      string
	Keywords      []string
	Ammenities    []string
}

type location = [2]float64

func main() {
	input := flag.String("input", "", "Path to JSON or CSV input file")
	sqft := flag.String("sqft", "", `Filter by square footage. Examples: ">1500", "=1500", "<1500", "<=1500", ">=1500"`)
	bathrooms := flag.String("bathrooms", "", `Filter by amount of bathrooms. Examples: ">1", "=1", "<3", "<=3", ">=3"`)
	distance := flag.String("distance", "", `Filter by distance in km to lat and long flags. Examples: ">100", "=100", "<100", "<=100", ">=100"`)
	price := flag.String("price", "", `Filter by price. Examples: ">1000", "=1000", "<1000", "<=1000", ">=1000"`)
	lat := flag.Float64("lat", -999999, `Latitude to compare distance`)
	long := flag.Float64("long", -999999, `Longitude to compare distance`)
	lighting := flag.String("lighting", "", `Lighting type. Possible values: 'low' | 'medium' | 'high'`)
	keywords := flag.String("keywords", "", `Keywords to search in description (comma-separated). Example: "spacious,big"`)
	ammenities := flag.String("ammenities", "", `Required amenities (comma-separated). Example: "garage,yard"`)
	flag.Parse()

	file, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var properties []Property
	bytes, _ := io.ReadAll(file)
	json.Unmarshal(bytes, &properties)

	var filters Filter

	if *sqft != "" {
		filters.SquareFootage, err = parseComparison(*sqft)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *bathrooms != "" {
		filters.Bathrooms, err = parseComparison(*bathrooms)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *distance != "" {
		filters.Distance, err = parseComparison(*distance)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *price != "" {
		filters.Price, err = parseComparison(*price)
		if err != nil {
			log.Fatal(err)
		}
	}
	filters.Lat = *lat
	filters.Long = *long
	filters.Lighting = *lighting
	if *keywords != "" {
		filters.Keywords = parseText(*keywords)
	}
	if *ammenities != "" {
		filters.Ammenities = parseText(*ammenities)
	}

	filteredProperties := filterProperties(properties, filters)
	printResultsAsJSON(filteredProperties)
}

func filterProperties(properties []Property, filters Filter) []Property {
	var filteredProperties []Property

Filters:
	for _, property := range properties {
		if len(filters.SquareFootage) != 0 {
			for _, comparison := range filters.SquareFootage {
				if !compare(comparison, property.SquareFootage) {
					continue Filters
				}
			}
		}
		if len(filters.Bathrooms) != 0 {
			for _, comparison := range filters.Bathrooms {
				if !compare(comparison, property.Bathrooms) {
					continue Filters
				}
			}
		}
		if len(filters.Distance) != 0 {
			if !(filters.Long == -999999) && !(filters.Lat == -999999) {
				actualDistance := calculateDistance(filters.Lat, filters.Long, property.Location[0], property.Location[1])
				for _, comparison := range filters.Bathrooms {
					if !compare(comparison, actualDistance) {
						continue Filters
					}
				}

			}
		}

		if len(filters.Price) != 0 {
			for _, comparison := range filters.Price {
				if !compare(comparison, property.Price) {
					continue Filters
				}
			}
		}

		if filters.Lighting != "" {
			if filters.Lighting != property.Lighting {
				continue Filters
			}
		}

		if len(filters.Keywords) != 0 {
			lowercaseDescription := strings.ToLower(property.Description)
			for _, keyword := range filters.Keywords {
				pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(keyword))
				regex := regexp.MustCompile(pattern)

				if !regex.MatchString(lowercaseDescription) {
					continue Filters
				}
			}

		}
		if len(filters.Ammenities) != 0 {
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

type Parser interface {
	Parse(data []byte) ([]Property, error)
}

type JSONParser struct{}

// func (j JSONParser) Parse(data []byte) ([]Property, error) {
// }

type CSVParser struct{}

// func (c CSVParser) Parse(data []byte) ([]Property, error) {
// }

func calculateDistance(lat1, long1, lat2, long2 float64) float64 {
	const EARTH_RADIUS = 6371
	const RADIAN = math.Pi / 180

	distance := 0.5 - math.Cos((lat2-lat1)*RADIAN)/2 + math.Cos(lat1*RADIAN)*math.Cos(lat2*RADIAN)*(1-math.Cos((long2-long1)*RADIAN))/2

	return 2 * EARTH_RADIUS * math.Asin(math.Sqrt(distance))
}

type Comparison struct {
	Operator string
	Value    float64
}

func parseComparison(s string) ([]Comparison, error) {
	operators := []string{"lte", "gte", "eq", "lt", "gt", "in"}
	trimmedString := strings.TrimSpace(s)
	var comparisons []Comparison

	for _, op := range operators {
		if strings.HasPrefix(trimmedString, op) {
			if op == "in" {
				rangeComparisons, err := parseRangeComparison(trimmedString[len(op):])
				if err != nil {
					return nil, err
				}
				comparisons = append(comparisons, rangeComparisons...)
				return comparisons, nil
			}

			comp, err := parseSingleComparison(op, trimmedString[len(op):])
			if err != nil {
				return nil, err
			}
			comparisons = append(comparisons, comp)
			return comparisons, nil
		}
	}

	return nil, fmt.Errorf("invalid comparison operator: %s", trimmedString)
}

func parseRangeComparison(valueRange string) ([]Comparison, error) {
	trimmedRange := strings.TrimSpace(valueRange)
	values := strings.Split(trimmedRange, ",")
	if len(values) != 2 {
		return nil, fmt.Errorf("range comparison requires two values, got: %s", trimmedRange)
	}

	firstNum, err := strconv.ParseFloat(strings.TrimSpace(values[0]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number in range: %s", values[0])
	}

	secondNum, err := strconv.ParseFloat(strings.TrimSpace(values[1]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number in range: %s", values[1])
	}

	return []Comparison{
		{Value: firstNum, Operator: "gte"},
		{Value: secondNum, Operator: "lte"},
	}, nil
}

func parseSingleComparison(operator, valueStr string) (Comparison, error) {
	numStr := strings.TrimSpace(valueStr)
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return Comparison{}, fmt.Errorf("invalid number in comparison: %s", numStr)
	}

	return Comparison{
		Value:    num,
		Operator: operator,
	}, nil
}

func compare(comparison Comparison, prop float64) bool {
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

func parseText(s string) []string {
	s = strings.ToLower(s)
	words := strings.Split(s, ",")
	for i, word := range words {
		words[i] = strings.TrimSpace(word)
	}
	return words
}

func printResultsAsJSON(properties []Property) error {
	jsonOutput, err := json.MarshalIndent(properties, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonOutput))
	return nil
}
