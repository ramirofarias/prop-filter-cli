package input

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/ramirofarias/prop-filter-cli/models"
)

func FromCSVFile(filename string) ([]models.Property, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV data: %v", err)
	}

	header := records[0]
	columnIndex := map[string]int{}
	for i, column := range header {
		columnIndex[column] = i
	}

	var properties []models.Property

	for _, record := range records[1:] {
		property := models.Property{}

		sqft, err := strconv.ParseFloat(record[columnIndex["squareFootage"]], 0)
		if err != nil {
			return nil, fmt.Errorf("invalid squareFootage value: %v", err)
		}
		property.SquareFootage = float64(sqft)

		property.Lighting = record[columnIndex["lighting"]]

		property.Price, err = strconv.ParseFloat(record[columnIndex["price"]], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid price value: %v", err)
		}

		rooms, err := strconv.Atoi(record[columnIndex["rooms"]])
		if err != nil {
			return nil, fmt.Errorf("invalid rooms value: %v", err)
		}
		property.Rooms = float64(rooms)

		bathrooms, err := strconv.Atoi(record[columnIndex["bathrooms"]])
		if err != nil {
			return nil, fmt.Errorf("invalid bathrooms value: %v", err)
		}
		property.Bathrooms = float64(bathrooms)

		property.Location[0], err = strconv.ParseFloat(record[columnIndex["latitude"]], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid latitude value: %v", err)
		}
		property.Location[1], err = strconv.ParseFloat(record[columnIndex["longitude"]], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid longitude value: %v", err)
		}

		property.Description = record[columnIndex["description"]]

		ammenitiesJSON := record[columnIndex["ammenities"]]
		err = json.Unmarshal([]byte(ammenitiesJSON), &property.Ammenities)
		if err != nil {
			return nil, fmt.Errorf("invalid ammenities JSON: %v", err)
		}

		properties = append(properties, property)
	}

	return properties, nil
}
