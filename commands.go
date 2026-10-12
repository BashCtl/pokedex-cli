package main

import (
	"fmt"
	"os"
)

func commandExit(conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("\nUsage of the Pokedex\nList of commands:\n\n")
	commands := getCommands()
	for key, value := range commands {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	fmt.Println()
	return nil
}

func commandMap(cfg *config) error {
	locations, err := cfg.pokeapiClient.ListLocationAreas(cfg.nextURL)
	if err != nil {
		return err
	}

	cfg.nextURL = locations.Next
	cfg.prevURL = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.prevURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locations, err := cfg.pokeapiClient.ListLocationAreas(cfg.prevURL)
	if err != nil {
		return err
	}

	cfg.nextURL = locations.Next
	cfg.prevURL = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	return nil
}
