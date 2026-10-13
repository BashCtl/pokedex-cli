package main

type cliCommand struct {
	name        string
	description string
	callback    func(conf *config, args []string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},

		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 location areas in the Pokemon world.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays last 20 location areas in the Pokemon world.",
			callback:    commandMapb,
		},

		"explore": {
			name:        "explore",
			description: "Explore an area and list the Pokemon found there",
			callback:    commandExplore,
		},

		"catch": {
			name:        "catch",
			description: "Attempts to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays details about a caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists the Pokemon you have caught",
			callback:    commandPokedex,
		},
	}
}
