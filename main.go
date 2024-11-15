package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ramirofarias/prop-filter-cli/filter"
	parser "github.com/ramirofarias/prop-filter-cli/parser"
)

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

	var properties []filter.Property
	bytes, _ := io.ReadAll(file)
	json.Unmarshal(bytes, &properties)

	var filters filter.Filter

	if *sqft != "" {
		filters.SquareFootage, err = parser.ParseComparison(*sqft)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *bathrooms != "" {
		filters.Bathrooms, err = parser.ParseComparison(*bathrooms)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *distance != "" {
		filters.Distance, err = parser.ParseComparison(*distance)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *price != "" {
		filters.Price, err = parser.ParseComparison(*price)
		if err != nil {
			log.Fatal(err)
		}
	}
	filters.Lat = *lat
	filters.Long = *long
	filters.Lighting = *lighting
	if *keywords != "" {
		filters.Keywords = parser.ParseText(*keywords)
	}
	if *ammenities != "" {
		filters.Ammenities = parser.ParseText(*ammenities)
	}

	filteredProperties := filter.FilterProperties(properties, filters)
	printResultsAsJSON(filteredProperties)
}

// type Parser interface {
// 	Parse(data []byte) ([]Property, error)
// }

type JSONParser struct{}

// func (j JSONParser) Parse(data []byte) ([]Property, error) {
// }

type CSVParser struct{}

// func (c CSVParser) Parse(data []byte) ([]Property, error) {
// }

func printResultsAsJSON(properties []filter.Property) error {
	jsonOutput, err := json.MarshalIndent(properties, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonOutput))
	return nil
}
