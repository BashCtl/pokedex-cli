package main

import (
	"github.com/BashCtl/pokedex-cli/internal/pokeapi"
	"time"
)

type config struct {
	commands      map[string]cliCommand
	nextURL       *string
	prevURL       *string
	pokeapiClient pokeapi.Client
	pokedex       map[string]pokeapi.Pokemon
}

func getConf() *config {
	return &config{
		commands:      getCommands(),
		pokeapiClient: pokeapi.NewClient(5 * time.Second),
		pokedex:       make(map[string]pokeapi.Pokemon),
	}

}
