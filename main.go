package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
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
	Description   location        `json:"description"`
	Ammenities    map[string]bool `json:"ammenities"`
}

type Filter struct {
	SquareFootage Comparison
	Bathrooms     Comparison
	Distance      Comparison
	Lat           float64
	Long          float64
}

type location = [2]float64

func main() {
	input := flag.String("input", "", "Path to JSON or CSV input file")
	sqft := flag.String("sqft", "", `Filter by square footage. Examples: ">1500", "=1500", "<1500", "<=1500", ">=1500"`)
	bathrooms := flag.String("bathrooms", "", `Filter by amount of bathrooms. Examples: ">1", "=1", "<3", "<=3", ">=3"`)
	distance := flag.String("distance", "", `Filter by distance in km to lat and long flags. Examples: ">100", "=100", "<100", "<=100", ">=100"`)
	lat := flag.Float64("lat", -999999, `Latitude to compare distance`)
	long := flag.Float64("long", -999999, `Longitude to compare distance`)
	flag.Parse()
	flag.Parse()
	fmt.Println(input)

	file, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var properties []Property
	bytes, _ := io.ReadAll(file)
	json.Unmarshal(bytes, &properties)
	fmt.Printf("properties: %#v \n", properties)

	var filters Filter
	filters.SquareFootage, _ = parseComparison(*sqft)
	filters.Bathrooms, _ = parseComparison(*bathrooms)
	filters.Distance, _ = parseComparison(*distance)
	filters.Lat = *lat
	filters.Long = *long

	fmt.Printf("filtered: %+v \n", filterProperties(properties, filters))
}

func filterProperties(properties []Property, filters Filter) []Property {
	var filteredProperties []Property

	for _, property := range properties {
		if filters.SquareFootage != (Comparison{}) {
			if !compare(filters.SquareFootage, property.SquareFootage) {
				continue
			}
		}
		if filters.Bathrooms != (Comparison{}) {
			if !compare(filters.Bathrooms, property.Bathrooms) {
				continue
			}
		}

		if filters.Distance != (Comparison{}) {
			if !(filters.Long == -999999) && !(filters.Lat == -999999) {
				actualDistance := calculateDistance(filters.Lat, filters.Long, property.Location[0], property.Location[1])
				if !compare(filters.Distance, actualDistance) {
					continue
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

func parseComparison(s string) (Comparison, error) {
	operators := []string{"<=", ">=", "=", "<", ">"}
	trimmedString := strings.TrimSpace(s)
	var comparison Comparison

	for _, o := range operators {
		if strings.HasPrefix(trimmedString, o) {
			numStr := strings.TrimSpace(trimmedString[len(o):])
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return Comparison{}, fmt.Errorf("invalid number used in comparison: %s", numStr)
			}

			comparison.Value = num
			comparison.Operator = o
			break
		}
	}

	if comparison == (Comparison{}) {
		return Comparison{}, fmt.Errorf("invalid comparison operator: %s", trimmedString)
	}

	return comparison, nil
}

func compare(comparison Comparison, prop float64) bool {
	switch comparison.Operator {
	case "<":
		if !(prop < comparison.Value) {
			return false
		}
	case ">":
		if !(prop > comparison.Value) {
			return false
		}
	case ">=":
		if !(prop >= comparison.Value) {
			return false
		}
	case "<=":
		if !(prop <= comparison.Value) {
			return false
		}
	case "=":
		if !(prop == comparison.Value) {
			return false
		}
	default:
		return true
	}

	return true
}
