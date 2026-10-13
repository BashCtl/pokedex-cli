package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/BashCtl/pokedex-cli/internal/pokeapi"
)

func commandExit(conf *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("\nUsage of the Pokedex\nList of commands:\n\n")
	commands := getCommands()
	for key, value := range commands {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	fmt.Println()
	return nil
}

func commandMap(cfg *config, args []string) error {
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

func commandMapb(cfg *config, args []string) error {
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

func commandExplore(conf *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Usage: explore <location-area>")
	}

	payload, err := conf.pokeapiClient.GetLocationArea(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", payload.Name)
	fmt.Println("[+]Found Pokemon:")
	for _, pokemon := range payload.PokemonEncounters {
		fmt.Printf("-> %s\n", pokemon.Pokemon.Name)
	}

	return nil

}

var catchRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func calculateCatchChance(baseExperience int) int {
	if baseExperience <= 0 {
		return 90
	}

	chance := 100 - (baseExperience / 3)
	if chance < 20 {
		return 20
	}
	if chance > 90 {
		return 90
	}
	return chance
}

func commandCatch(conf *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon>")
	}

	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := conf.pokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	if catchRand.Intn(100) < calculateCatchChance(pokemon.BaseExperience) {
		if conf.pokedex == nil {
			conf.pokedex = make(map[string]pokeapi.Pokemon)
		}
		conf.pokedex[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%s escaped!\n", pokemon.Name)
	return nil
}
