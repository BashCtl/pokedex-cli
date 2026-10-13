package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	var result []string
	text = strings.ToLower(text)
	result = strings.Fields(text)
	return result
}
func startRepl(conf *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Input error:", err)
			}
			return
		}

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		cmd, ok := getCommands()[commandName]
		if !ok {
			fmt.Printf("Input not valid: %s\n", commandName)
			continue
		}
		if err := cmd.callback(conf, words[1:]); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
