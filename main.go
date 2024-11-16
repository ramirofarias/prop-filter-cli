package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ramirofarias/prop-filter-cli/filter"
	"github.com/ramirofarias/prop-filter-cli/input"
	"github.com/ramirofarias/prop-filter-cli/models"
	"github.com/ramirofarias/prop-filter-cli/output"
	"github.com/ramirofarias/prop-filter-cli/parser"
)

func main() {
	inputPath := flag.String("input", "", "Path to JSON or CSV input file")
	sqft := flag.String("sqft", "", `Filter by square footage. Examples: ">1500", "=1500", "<1500", "<=1500", ">=1500"`)
	bathrooms := flag.String("bathrooms", "", `Filter by amount of bathrooms. Examples: ">1", "=1", "<3", "<=3", ">=3"`)
	distance := flag.String("distance", "", `Filter by distance in km to lat and long flags. Examples: ">100", "=100", "<100", "<=100", ">=100"`)
	price := flag.String("price", "", `Filter by price. Examples: ">1000", "=1000", "<1000", "<=1000", ">=1000"`)
	lat := flag.Float64("lat", -999999, `Latitude to compare distance`)
	long := flag.Float64("long", -999999, `Longitude to compare distance`)
	lighting := flag.String("lighting", "", `Lighting type. Possible values: 'low' | 'medium' | 'high'`)
	keywords := flag.String("keywords", "", `Keywords to search in description (comma-separated). Example: "spacious,big"`)
	ammenities := flag.String("ammenities", "", `Required amenities (comma-separated). Example: "garage,yard"`)
	outputPath := flag.String("output", "", `Output file path in .csv or .json. Examples: "file.csv", "file.json"`)
	flag.Parse()
	var properties []models.Property

	if *inputPath != "" {
		fileType, err := parser.ParseFiletype(*inputPath)

		if err != nil {
			log.Fatal(err)
		}

		if fileType == "json" {
			properties, err = input.FromJSONFile(*inputPath)

			if err != nil {
				log.Fatal(err)
			}
		}

		if fileType == "csv" {
			properties, err = input.FromCSVFile(*inputPath)

			if err != nil {
				log.Fatal(err)
			}
		}

	}

	var filters filter.Filter
	var err error

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
	if *outputPath != "" {
		fileType, err := parser.ParseFiletype(*outputPath)

		if err != nil {
			log.Fatal(err)
		}

		if fileType == "json" {
			output.ToJSONFile(filteredProperties, *outputPath)
			os.Exit(0)
		}

		if fileType == "csv" {
			output.ToCSVFile(filteredProperties, *outputPath)
			os.Exit(0)
		}
	}

	if err := output.ToJSONStdOut(filteredProperties); err != nil {
		fmt.Printf("Error printing data to stdout: %v\n", err)
	}
}
