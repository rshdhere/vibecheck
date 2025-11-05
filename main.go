/*
Copyright © 2025 raashed
*/
package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/rshdhere/vibecheck/cmd"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — continuing without it")
	}
	cmd.Execute()
}
