package main

import "strings"

func cleanInput(text string) []string {
	var result []string
	text = strings.ToLower(text)
	result = strings.Fields(text)
	return result
}
