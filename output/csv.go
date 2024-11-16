package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ramirofarias/prop-filter-cli/models"
)

func ToCSVFile(data []models.Property, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"squareFootage", "lighting", "price", "rooms", "bathrooms", "latitude", "longitude", "description", "ammenities",
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	for _, property := range data {
		var row []string
		row = append(row, fmt.Sprintf("%f", property.SquareFootage))
		row = append(row, property.Lighting)
		row = append(row, fmt.Sprintf("%.2f", property.Price))
		row = append(row, fmt.Sprintf("%d", int(property.Rooms)))
		row = append(row, fmt.Sprintf("%d", int(property.Bathrooms)))
		row = append(row, fmt.Sprintf("%.6f", property.Location[0]))
		row = append(row, fmt.Sprintf("%.6f", property.Location[1]))
		row = append(row, property.Description)
		ammenitiesJSON, err := json.Marshal(property.Ammenities)
		if err != nil {
			return fmt.Errorf("error marshalling amenities to JSON: %v", err)
		}
		row = append(row, string(ammenitiesJSON))

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing CSV row: %v", err)
		}
	}

	return nil
}
