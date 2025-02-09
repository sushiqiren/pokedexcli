package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"	
)

type config struct {
    Next     string
    Previous string
}

type locationArea struct {
    Name string `json:"name"`
}

type locationAreaResponse struct {
    Results []locationArea `json:"results"`
    Next    string         `json:"next"`
    Previous string        `json:"previous"`
}

func main() {
	// Wait for user input using bufio.NewScanner
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{}

	type cliCommand struct {
		name        string
		description string
		callback    func(*config) error
	}

	commands := map[string]cliCommand{}

	 // Add a callback for the exit command.
    // This function should print Closing the Pokedex... Goodbye! then immediately exit the program
    commandExit := func(cfg *config) error {
        fmt.Println("Closing the Pokedex... Goodbye!")
        os.Exit(0)
        return nil
    }

    // Add a help command, its callback, and register it.
    // it should print
    // Welcome to the Pokedex!
    // Usage:
    // help: Displays a help message
    // exit: Exit the Pokedex
    commandHelp := func(cfg *config) error {
        fmt.Println("Welcome to the Pokedex!")
        fmt.Println("Usage:")
        for _, cmd := range commands {
            fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
        }
        return nil
    }

    // Add the map command. It displays the names of 20 location areas in the Pokemon world.
    // Each subsequent call to map should display the next 20 locations, and so on.
    // Use the PokeAPI location-area endpoint to get the location areas.
    // Update all commands (e.g. help, exit, map) to now accept a pointer to a “config” struct as a parameter.
    // This struct will contain the Next and Previous URLs that you’ll need to paginate through location areas
    commandMap := func(cfg *config) error {
        url := "https://pokeapi.co/api/v2/location-area/"
        if cfg.Next != "" {
            url = cfg.Next
        }

        resp, err := http.Get(url)
        if err != nil {
            return err
        }
        defer resp.Body.Close()

        var data locationAreaResponse
        if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
            return err
        }

        for _, location := range data.Results {
            fmt.Println(location.Name)
        }

        cfg.Next = data.Next
        cfg.Previous = data.Previous

        return nil
    }

	// Add the mapb (map back) command. It’s similar to the map command, however, instead of displaying the next 20 locations, it displays the previous 20 locations. It’s a way to go back.
	// Update the map command to accept a pointer to a “config” struct as a parameter. This struct will contain the Next and Previous URLs that you’ll need to paginate through location areas
	// Update the mapb command to accept a pointer to a “config” struct as a parameter. This struct will contain the Next and Previous URLs that you’ll need to paginate through location areas
	commandMapBack := func(cfg *config) error {
		url := "https://pokeapi.co/api/v2/location-area/"
		if cfg.Previous != "" {
			url = cfg.Previous
		}

		resp, err := http.Get(url)
		if err != nil {
			return err
        }
        defer resp.Body.Close()

        var data locationAreaResponse
        if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
            return err
        }

        for _, location := range data.Results {
            fmt.Println(location.Name)
        }

        cfg.Next = data.Next
        cfg.Previous = data.Previous

        return nil
    }

    commands["exit"] = cliCommand{
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
    }
    commands["help"] = cliCommand{
        name:        "help",
        description: "Displays a help message",
        callback:    commandHelp,
    }
    commands["map"] = cliCommand{
        name:        "map",
        description: "Displays location areas in the Pokemon world",
        callback:    commandMap,
    }
	commands["mapb"] = cliCommand{
        name:        "mapb",
        description: "Displays the previous 20 location areas",
        callback:    commandMapBack,
    }

	// Start an infinite for loop. This loop will execute once for every command the user types in
	// Use fmt.Print to print the prompt Pokedex > without a newline character
	// Use the scanner’s .Scan and .Text methods to get the user’s input as a string
	// Clean the users input by trimming any leading or trailing whitespace, and converting it to lowercase. use strings.ToLower and strings.Fields to split the input into a slice of words
	// Capture the first “word” of the input and use it to print: Your command was: <first word>
	for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        input := scanner.Text()

        // Clean the user's input
        words := strings.Fields(strings.ToLower(strings.TrimSpace(input)))

        if len(words) == 0 {
			continue
		}

		// register the exit command. Update your REPL loop to use the “command” the user typed in to look up the callback function in the registry. If the command is found, call the callback (and print any errors that are returned). If there isn’t a handler, just print Unknown command
		if command, ok := commands[words[0]]; ok {
			err := command.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
    }
}