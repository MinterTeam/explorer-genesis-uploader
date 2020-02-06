package main

import (
	"github.com/MinterTeam/explorer-genesis-uploader/core"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	uploader := core.New()
	err = uploader.Do()
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
