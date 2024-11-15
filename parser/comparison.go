package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ramirofarias/prop-filter-cli/filter"
)

func ParseComparison(s string) ([]filter.Comparison, error) {
	operators := []string{"lte", "gte", "eq", "lt", "gt", "in"}
	trimmedString := strings.TrimSpace(s)
	var comparisons []filter.Comparison

	for _, op := range operators {
		if strings.HasPrefix(trimmedString, op) {
			if op == "in" {
				rangeComparisons, err := parseRangeComparison(trimmedString[len(op):])
				if err != nil {
					return nil, err
				}
				comparisons = append(comparisons, rangeComparisons...)
				return comparisons, nil
			}

			comp, err := parseSingleComparison(op, trimmedString[len(op):])
			if err != nil {
				return nil, err
			}
			comparisons = append(comparisons, comp)
			return comparisons, nil
		}
	}

	return nil, fmt.Errorf("invalid comparison operator: %s", trimmedString)
}

func parseRangeComparison(valueRange string) ([]filter.Comparison, error) {
	trimmedRange := strings.TrimSpace(valueRange)
	values := strings.Split(trimmedRange, ",")
	if len(values) != 2 {
		return nil, fmt.Errorf("range comparison requires two values, got: %s", trimmedRange)
	}

	firstNum, err := strconv.ParseFloat(strings.TrimSpace(values[0]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number in range: %s", values[0])
	}

	secondNum, err := strconv.ParseFloat(strings.TrimSpace(values[1]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number in range: %s", values[1])
	}

	return []filter.Comparison{
		{Value: firstNum, Operator: "gte"},
		{Value: secondNum, Operator: "lte"},
	}, nil
}

func parseSingleComparison(operator, valueStr string) (filter.Comparison, error) {
	numStr := strings.TrimSpace(valueStr)
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return filter.Comparison{}, fmt.Errorf("invalid number in comparison: %s", numStr)
	}

	return filter.Comparison{
		Value:    num,
		Operator: operator,
	}, nil
}
