package main

import "testing"

func TestCleanInput(t *testing.T) {
	tests := []struct {
		input string
		expected []string
	}{
		{
			input: "Hello World",
			expected: []string{"hello", "world"},
		},
	}

	for _, test := range tests {
		actual := cleanInput(test.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(test.expected) {
			t.Errorf("Expected %d words, but got %d", len(test.expected), len(actual))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := test.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Expected %s, but got %s", expectedWord, word)
			}
		}
	}
}