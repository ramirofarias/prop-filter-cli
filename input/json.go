package input

import (
	"encoding/json"
	"io"
	"os"

	"github.com/ramirofarias/prop-filter-cli/models"
)

func FromJSONFile(filename string) ([]models.Property, error) {
	var properties []models.Property

	file, err := os.Open(filename)
	if err != nil {
		return nil, err

	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)
	json.Unmarshal(bytes, &properties)
	return properties, nil
}
