package parser

import "strings"

func ParseText(s string) []string {
	s = strings.ToLower(s)
	words := strings.Split(s, ",")
	for i, word := range words {
		words[i] = strings.TrimSpace(word)
	}
	return words
}
