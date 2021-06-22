package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/MinterTeam/explorer-genesis-uploader/core"
	"github.com/MinterTeam/explorer-genesis-uploader/env"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var cfg = flag.String(`config`, "", `Path to config`)

func main() {
	flag.Parse()

	var environment env.Config

	if *cfg != "" {
		if _, err := toml.DecodeFile(*cfg, &environment); err != nil {
			panic(err)
		}
	} else {
		err := godotenv.Load()
		if err != nil {
			println("Error loading .env file")
		}
		addressChunkSize, err := strconv.ParseUint(os.Getenv("APP_ADDRESS_CHUNK_SIZE"), 10, 64)
		if err != nil {
			println(err)
		}
		coinsChunkSize, err := strconv.ParseUint(os.Getenv("APP_COINS_CHUNK_SIZE"), 10, 64)
		if err != nil {
			println(err)
		}
		balanceChunkSize, err := strconv.ParseUint(os.Getenv("APP_BALANCES_CHUNK_SIZE"), 10, 64)
		if err != nil {
			println(err)
		}
		stakeChunkSize, err := strconv.ParseUint(os.Getenv("APP_STAKE_CHUNK_SIZE"), 10, 64)
		if err != nil {
			println(err)
		}
		validatorChunkSize, err := strconv.ParseUint(os.Getenv("APP_VALIDATORS_CHUNK_SIZE"), 10, 64)
		if err != nil {
			println(err)
		}
		environment = env.Config{
			Debug:              os.Getenv("DEBUG") == "true",
			PostgresHost:       os.Getenv("POSTGRES_HOST"),
			PostgresPort:       os.Getenv("POSTGRES_PORT"),
			PostgresDB:         os.Getenv("POSTGRES_NAME"),
			PostgresUser:       os.Getenv("POSTGRES_USER"),
			PostgresPassword:   os.Getenv("POSTGRES_PASSWORD"),
			PostgresSSLEnabled: os.Getenv("POSTGRES_SSL_ENABLED") == "true",
			MinterBaseCoin:     os.Getenv("MINTER_BASE_COIN"),
			NodeGrpc:           os.Getenv("NODE_GRPC"),
			AddressChunkSize:   addressChunkSize,
			CoinsChunkSize:     coinsChunkSize,
			BalanceChunkSize:   balanceChunkSize,
			StakeChunkSize:     stakeChunkSize,
			ValidatorChunkSize: validatorChunkSize,
		}
	}

	uploader := core.New(environment)
	err := uploader.Do()
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
