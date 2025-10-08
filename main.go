package main

import (
	"fmt"
	"os"

	"github.com/davidmks/got/internal/commands"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		if err := commands.Init(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("got - A simple version control system")
	fmt.Println("\nUsage:")
	fmt.Println("  got <command> [arguments]")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  init      Initialize a new repository")
}
