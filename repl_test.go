package main

import (
	"testing"
	"strings"
)

func cleanInput(text string) []string {
	// split the users input into “words” based on whitespace
	// create a new slice to hold the cleaned words
	// lowercase the input and trim any leading or trailing whitespace
	// return the cleaned input
	words := strings.Fields(strings.ToLower(text))
	return words
}

func TestCleanInput(t *testing.T) {
    cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		// add more cases here
		{
			input: "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice	
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) returned %d words, expected %d", c.input, len(actual), len(c.expected))
		}
		// Loop through the actual slice
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("cleanInput(%q) returned %q, expected %q", c.input, word, expectedWord)
			}
		}
	}
}



