/*
Copyright © 2025 raashed
*/
package main

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rshdhere/vibecheck/cmd"
)

func main() {
	// Only try to load .env file if we're not just checking version/help
	// This prevents unnecessary log messages during simple --version checks
	if !isVersionOrHelp() {
		loadDotEnvIfPresent()
	}
	cmd.Execute()
}

func loadDotEnvIfPresent() {
	if _, err := os.Stat(".env"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Printf("unable to stat .env: %v", err)
		}
		return
	}

	if err := godotenv.Load(); err != nil {
		log.Printf("unable to load .env: %v", err)
	}
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
