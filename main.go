package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sushiqiren/pokedexcli/internal/pokecache"
)

type config struct {
	Next     string
	Previous string
}

type locationArea struct {
	Name string `json:"name"`
}

type locationAreaResponse struct {
	Results  []locationArea `json:"results"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
}

type locationAreaDetail struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{}
	cache := pokecache.NewCache(5 * time.Minute)
	caughtPokemon := make(map[string]Pokemon)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type cliCommand struct {
		name        string
		description string
		callback    func(*config, []string) error
	}

	commands := map[string]cliCommand{}

	commandExit := func(cfg *config, args []string) error {
		fmt.Println("Closing the Pokedex... Goodbye!")
		os.Exit(0)
		return nil
	}

	commandHelp := func(cfg *config, args []string) error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Println("Usage:")
		for _, cmd := range commands {
			fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
		}
		return nil
	}

	commandMap := func(cfg *config, args []string) error {
		url := "https://pokeapi.co/api/v2/location-area/"
		if cfg.Next != "" {
			url = cfg.Next
		}

		if data, found := cache.Get(url); found {
			fmt.Println("Using cached data")
			var response locationAreaResponse
			if err := json.Unmarshal(data, &response); err != nil {
				return err
			}
			for _, location := range response.Results {
				fmt.Println(location.Name)
			}
			cfg.Next = response.Next
			cfg.Previous = response.Previous
			return nil
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

		responseData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		cache.Add(url, responseData)

		return nil
	}

	commandMapBack := func(cfg *config, args []string) error {
		url := "https://pokeapi.co/api/v2/location-area/"
		if cfg.Previous != "" {
			url = cfg.Previous
		}

		if data, found := cache.Get(url); found {
			fmt.Println("Using cached data")
			var response locationAreaResponse
			if err := json.Unmarshal(data, &response); err != nil {
				return err
			}
			for _, location := range response.Results {
				fmt.Println(location.Name)
			}
			cfg.Next = response.Next
			cfg.Previous = response.Previous
			return nil
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

		responseData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		cache.Add(url, responseData)

		return nil
	}

	commandExplore := func(cfg *config, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("please provide a location area name")
		}
		areaName := args[0]
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", areaName)

		if data, found := cache.Get(url); found {
			fmt.Println("Using cached data")
			var response locationAreaDetail
			if err := json.Unmarshal(data, &response); err != nil {
				return err
			}
			for _, encounter := range response.PokemonEncounters {
				fmt.Println(encounter.Pokemon.Name)
			}
			return nil
		}

		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var data locationAreaDetail
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}

		for _, encounter := range data.PokemonEncounters {
			fmt.Println(encounter.Pokemon.Name)
		}

		responseData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		cache.Add(url, responseData)

		return nil
	}

	commandCatch := func(cfg *config, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("please provide a Pokemon name")
		}
		pokemonName := args[0]
		url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)

		fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

		if data, found := cache.Get(url); found {
			fmt.Println("Using cached data")
			var pokemon Pokemon
			if err := json.Unmarshal(data, &pokemon); err != nil {
				return err
			}
			if _, caught := caughtPokemon[pokemon.Name]; caught {
				fmt.Printf("%s was already caught!\n", pokemon.Name)
				return nil
			}
			if catchPokemon(rng, pokemon) {
				caughtPokemon[pokemon.Name] = pokemon
				fmt.Printf("%s was caught!\n", pokemon.Name)
				fmt.Println("You may now inspect it with the inspect command.")
			} else {
				fmt.Printf("%s escaped!\n", pokemon.Name)
			}
			return nil
		}

		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var pokemon Pokemon
		if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
			return err
		}

		if catchPokemon(rng, pokemon) {
			caughtPokemon[pokemon.Name] = pokemon
			fmt.Printf("%s was caught!\n", pokemon.Name)
			fmt.Println("You may now inspect it with the inspect command.")
		} else {
			fmt.Printf("%s escaped!\n", pokemon.Name)
		}

		responseData, err := json.Marshal(pokemon)
		if err != nil {
			return err
		}
		cache.Add(url, responseData)

		return nil
	}

	commandInspect := func(cfg *config, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("please provide a Pokemon name")
		}
		pokemonName := args[0]

		pokemon, caught := caughtPokemon[pokemonName]
		if !caught {
			fmt.Println("you have not caught that pokemon")
			return nil
		}

		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range pokemon.Types {
			fmt.Printf("  - %s\n", t.Type.Name)
		}

		return nil
	}

	commandPokedex := func(cfg *config, args []string) error {
		if len(caughtPokemon) == 0 {
			fmt.Println("Your Pokedex is empty.")
			return nil
		}

		fmt.Println("Your Pokedex:")
		for name := range caughtPokemon {
			fmt.Printf(" - %s\n", name)
		}

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
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Explore a specific location area",
		callback:    commandExplore,
	}
	commands["catch"] = cliCommand{
		name:        "catch",
		description: "Catch a specific Pokemon",
		callback:    commandCatch,
	}
	commands["inspect"] = cliCommand{
		name:        "inspect",
		description: "Inspect a caught Pokemon",
		callback:    commandInspect,
	}
	commands["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "List all caught Pokemon",
		callback:    commandPokedex,
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		words := strings.Fields(strings.ToLower(strings.TrimSpace(input)))

		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := words[1:]

		if command, ok := commands[commandName]; ok {
			err := command.callback(cfg, args)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func catchPokemon(rng *rand.Rand, pokemon Pokemon) bool {
	chance := 100 - pokemon.BaseExperience
	if chance < 10 {
		chance = 10 // Ensure there's always at least a 10% chance to catch
	}
	return rng.Intn(100) < chance
}
