package main

import (	
	"bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
	// Wait for user input using bufio.NewScanner
	scanner := bufio.NewScanner(os.Stdin)
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

        // Add a callback for the exit command. 
		// This function should print Closing the Pokedex... Goodbye! then immediately exit the program
		commandExit := func() error {
			fmt.Println("Closing the Pokedex... Goodbye!")
			os.Exit(0)
			return nil
		}

		
		type cliCommand struct {
			name        string
			description string
			callback    func() error
		}

		commands := map[string]cliCommand{
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},			
		}

		// Add a help command, its callback, and register it.
		// it should print 
		// Welcome to the Pokedex!
        // Usage:

		// help: Displays a help message
		// exit: Exit the Pokedex
		commandHelp := func() error {
            fmt.Println("Welcome to the Pokedex!")
            fmt.Println("Usage:")
            for _, cmd := range commands {
                fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
            }
            return nil
        }

		commands["help"] = cliCommand{
            name:        "help",
            description: "Displays a help message",
            callback:    commandHelp,
        }


		// register the exit command. Update your REPL loop to use the “command” the user typed in to look up the callback function in the registry. If the command is found, call the callback (and print any errors that are returned). If there isn’t a handler, just print Unknown command
		if command, ok := commands[words[0]]; ok {
			err := command.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}

		
	





		
    }
}