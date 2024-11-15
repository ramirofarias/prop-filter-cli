package output

import (
	"encoding/json"
	"fmt"
	"os"
)

func ToJSONFile(data interface{}, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("could not encode data to JSON: %v", err)
	}

	return nil
}

func ToJSONStdOut(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("could not encode data to JSON: %v", err)
	}

	return nil
}
