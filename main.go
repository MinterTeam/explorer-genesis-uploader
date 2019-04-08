package main

import (
	"github.com/MinterTeam/explorer-genesis-uploader/core"
	"github.com/MinterTeam/minter-explorer-extender/env"
	"os"
)

func main() {
	envData := env.New()
	uploader := core.New(envData)
	err := uploader.Do()
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
