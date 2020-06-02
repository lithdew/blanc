package main

import (
	"fmt"
	"github.com/ktr0731/go-fuzzyfinder/matching"
	"testing"
)

func TestAutocomplete(t *testing.T) {
	items := []string{
		"hot potato",
		"cooking oil",
		"cheese steak",
	}
	for _, match := range matching.FindAll("a", items) {
		fmt.Println("This looks similar:", items[match.Idx])
	}
}
