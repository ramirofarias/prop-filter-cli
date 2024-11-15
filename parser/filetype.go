package parser

import (
	"fmt"
	"strings"
)

func ParseFiletype(s string) (string, error) {
	ext := strings.ToLower(s[strings.LastIndex(s, ".")+1:])
	if ext == "" {
		return "", fmt.Errorf("output file must have an extension (e.g., .json, .csv)")
	}

	if ext != "json" && ext != "csv" {
		return "", fmt.Errorf("invalid output type: %s", ext)
	}

	return ext, nil
}
