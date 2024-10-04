package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
)

type HashTable struct {
	table map[string]int
}

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

func GenerateState() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func NewHashTable() *HashTable {
	return &HashTable{
		table: make(map[string]int),
	}
}

// Add increments the count for the given album name (keystring)
func (h *HashTable) Add(keystring string) {
	if h.table[keystring] > 0 {
		h.table[keystring]++
	} else {
		h.table[keystring] = 1
	}
}

// Get returns the count for the given album name (keystring)
func (h *HashTable) Get(keystring string) int {
	if h.table[keystring] > 0 {
		return h.table[keystring]
	} else {
		return 0
	}
}
