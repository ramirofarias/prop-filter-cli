package main

import (
	"fmt"
	"os"

	"github.com/ramirofarias/prop-filter-cli/filter"
	"github.com/ramirofarias/prop-filter-cli/input"
	"github.com/ramirofarias/prop-filter-cli/models"
	"github.com/ramirofarias/prop-filter-cli/output"
	"github.com/ramirofarias/prop-filter-cli/parser"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "prop-filter-cli",
		Usage: "Filter property data from JSON or CSV files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Usage:    "Path to JSON or CSV input file",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "sqft",
				Usage: `Filter by square footage. Examples: "gt 1500", "eq 1500", "lt 1500", "lte 1500", "in 1500,2000"`,
			},
			&cli.StringFlag{
				Name:  "bathrooms",
				Usage: `Filter by amount of bathrooms. Examples: "gt 1", "eq 1", "lt 3", "lte 3", "gte 3", "in 1,3"`,
			},
			&cli.StringFlag{
				Name:  "distance",
				Usage: `Filter by distance in km to lat and long flags. Examples: "gt 100", "eq 100", "lt 100", "lte 100", "gte 100", "in 150,200"`,
			},
			&cli.StringFlag{
				Name:  "price",
				Usage: `Filter by price. Examples: "gt 1000", "eq 1000", "lt 1000", "lte 1000", "gte 1000"`,
			},
			&cli.Float64Flag{
				Name:  "lat",
				Value: -999999,
				Usage: `Latitude to compare distance`,
			},
			&cli.Float64Flag{
				Name:  "long",
				Value: -999999,
				Usage: `Longitude to compare distance`,
			},
			&cli.StringFlag{
				Name:  "lighting",
				Usage: `Lighting type. Possible values: 'low' | 'medium' | 'high'`,
			},
			&cli.StringFlag{
				Name:  "keywords",
				Usage: `Keywords to search in description (comma-separated). Example: "spacious,big"`,
			},
			&cli.StringFlag{
				Name:  "ammenities",
				Usage: `Required amenities (comma-separated). Example: "garage,yard"`,
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: `Output file path in .csv or .json. Examples: "file.csv", "file.json"`,
			},
		},
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			inputPath := c.String("input")
			var properties []models.Property

			fileType, err := parser.ParseFiletype(inputPath)
			if err != nil {
				return fmt.Errorf("error parsing input file type: %v", err)
			}

			switch fileType {
			case "json":
				properties, err = input.FromJSONFile(inputPath)
			case "csv":
				properties, err = input.FromCSVFile(inputPath)
			}

			if err != nil {
				return fmt.Errorf("error parsing input file: %v", err)
			}

			var filters filter.Filter
			if sqft := c.String("sqft"); sqft != "" {
				filters.SquareFootage, err = parser.ParseComparison(sqft)
				if err != nil {
					return fmt.Errorf("error parsing sqft filter: %v", err)
				}
			}
			if bathrooms := c.String("bathrooms"); bathrooms != "" {
				filters.Bathrooms, err = parser.ParseComparison(bathrooms)
				if err != nil {
					return fmt.Errorf("error parsing bathrooms filter: %v", err)
				}
			}
			filters.Lat = c.Float64("lat")
			filters.Long = c.Float64("long")
			if distance := c.String("distance"); distance != "" {
				if filters.Lat == -999999 || filters.Long == -999999 {
					return fmt.Errorf("lat and long flags are required when using distance filter")
				}
				filters.Distance, err = parser.ParseComparison(distance)
				if err != nil {
					return fmt.Errorf("error parsing distance filter: %v", err)
				}
			}
			if price := c.String("price"); price != "" {
				filters.Price, err = parser.ParseComparison(price)
				if err != nil {
					return fmt.Errorf("error parsing price filter: %v", err)
				}
			}

			filters.Lighting = c.String("lighting")
			if keywords := c.String("keywords"); keywords != "" {
				filters.Keywords = parser.ParseText(keywords)
			}
			if ammenities := c.String("ammenities"); ammenities != "" {
				filters.Ammenities = parser.ParseText(ammenities)
			}

			filteredProperties := filter.FilterProperties(properties, filters)

			outputPath := c.String("output")
			if outputPath != "" {
				fileType, err := parser.ParseFiletype(outputPath)
				if err != nil {
					return fmt.Errorf("error parsing output file type: %v", err)
				}

				switch fileType {
				case "json":
					if err := output.ToJSONFile(filteredProperties, outputPath); err != nil {
						return fmt.Errorf("error writing JSON output file: %v", err)
					}
				case "csv":
					if err := output.ToCSVFile(filteredProperties, outputPath); err != nil {
						return fmt.Errorf("error writing CSV output file: %v", err)
					}
				}

			} else {
				if err := output.ToJSONStdOut(filteredProperties); err != nil {
					return fmt.Errorf("error printing data to stdout: %v", err)
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error running app: %v\n", err)
		os.Exit(1)
	}
}
