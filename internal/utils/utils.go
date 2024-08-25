package utils

import (
	"fmt"
	"strings"
)

// SliceOfStringsToString converts a slice of strings to a single string
// Example: ["a", "b", "c"] -> "a", "b", "c"
func SliceOfStringsToString(slice []string) string {
	o := make([]string, len(slice))
	for i, s := range slice {
		o[i] = fmt.Sprintf("\"%s\"", s)
	}
	return strings.Join(o, ", ")
}
