/*
Copyright © 2025 raashed
*/
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rshdhere/vibecheck/cmd"
)

func main() {
	// Only try to load .env file if we're not just checking version/help
	// This prevents unnecessary log messages during simple --version checks
	if !isVersionOrHelp() {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found — continuing without it")
		}
	}
	cmd.Execute()
}

// isVersionOrHelp checks if the command is just asking for version or help
func isVersionOrHelp() bool {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "--version" || arg == "-v" || arg == "version" ||
			arg == "--help" || arg == "-h" || arg == "help" {
			return true
		}
	}
	return len(args) == 0 // No args = help
}
