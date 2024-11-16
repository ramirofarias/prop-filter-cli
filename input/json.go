package input

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	data, _ := io.ReadAll(file)
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&properties)

	if err != nil {
		return []models.Property{}, fmt.Errorf("error unmarshaling json: %e", err)
	}
	return properties, nil
}
