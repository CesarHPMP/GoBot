package utils

import (
	"fmt"
	"strings"
)

// TrimSlice trims all strings in a slice
func TrimSlice(slice []string) []string {
	for i := range slice {
		slice[i] = strings.TrimSpace(slice[i])
	}
	return slice
}

// PrintError prints an error message
func PrintError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}
