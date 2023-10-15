package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Access command-line arguments
	args := os.Args

	// Get the name of the executable (the first argument, index 0)
	executableName := filepath.Base(args[0])

	// Ensure that at least one command-line argument is provided
	if len(args) < 2 {
		fmt.Printf("Usage: %s [somecommand]\n", executableName)
		return
	}

	// Access the second argument (index 1), which is "[somecommand]" in your example
	someCommand := args[1]
	fmt.Printf("You provided the command: %s\n", someCommand)

	// You can perform actions based on the command here
	if someCommand == "doSomething" {
		fmt.Println("Executing some specific command...")
		// Add your code for the specific command here
	} else {
		fmt.Println("Unknown command or no command provided.")
	}
}
