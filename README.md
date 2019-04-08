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

## RUN

- make build

- ./builds/explorer-genesis-uploader -config=/etc/minter/config.json 

### Config file

Support JSON and YAML formats 

Example:

```
{
  "name": "Minter Explorer Genesis Uploader",
  "app": {
    "debug": true,
    "baseCoin": "MNT",
    "txChunkSize": 200,
    "addrChunkSize": 30,
    "eventsChunkSize": 200
  },
  "database": {
    "host": "localhost",
    "name": "explorer",
    "user": "minter",
    "password": "password",
    "minIdleConns": 10,
    "poolSize": 20
  },
  "minterApi": {
    "isSecure": false,
    "link": "localhost",
    "port": 8841
  }
}
```