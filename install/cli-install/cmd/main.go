package main

import (
	"installer/cmd/cli"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cliCmd := cli.NewCLI()
	cliCmd.Execute()
}
