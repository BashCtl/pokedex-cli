package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		if scanner.Scan() {
			input := scanner.Text()
			words := cleanInput(input)

			if len(words) == 0 {
				continue
			}

			commandName := words[0]

			cmd, ok := getCommands()[commandName]
			if ok {
				cmd.callback()
			} else {
				fmt.Printf("Input not valid: %s\n", commandName)
			}

		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Input error:", err)
		}
	}
}
