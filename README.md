<p align="center" background="black"><img src="minter-logo.svg" width="400"></p>

<p align="center" style="text-align: center;">
    <a href="https://github.com/MinterTeam/explorer-genesis-uploader/blob/master/LICENSE">
        <img src="https://img.shields.io/packagist/l/doctrine/orm.svg" alt="License">
    </a>
    <img alt="undefined" src="https://img.shields.io/github/last-commit/MinterTeam/explorer-genesis-uploader.svg">
</p>

# Minter Explorer Genesis Uploader

The official repository of Minter Explorer Genesis Uploader service.

Minter Explorer Genesis Uploader is a service which provides to upload primary network state data to Minter Explorer database after network reset or first start.

## Requirement

- PostgresSQL

## Build

- use database migration from `database` directory

- run `go mod tidy`

- run `go build -o ./builds/explorer_genesis_uploader ./cmd/explorer_genesis_uploader.go`

## Run

- copy `.env.prod` to `.env` and fill with own values

- run `./builds/explorer-genesis-uploader` or `docker-compose up`