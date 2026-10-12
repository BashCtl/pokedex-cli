package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  hello",
			expected: []string{"hello"},
		},
		{
			input:    "  ",
			expected: []string{},
		},
		{
			input:    "  cool     dude  in  the  world  ",
			expected: []string{"cool", "dude", "in", "the", "world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Length of slice doesn't mutch expected: %v, but got: %v", c.expected, actual)
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Words in slice don't mutch, expected: %v, but got: %v", expectedWord, word)
			}
		}
	}
}
