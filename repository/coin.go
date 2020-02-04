package repository

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"github.com/go-pg/pg/v9"
	"os"
)

type Coin struct {
}

func NewCoinRepository() *Coin {
	return &Coin{}
}

func (r *Coin) SaveAll(coins []*models.Coin) error {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	defer db.Close()
	err := db.Insert(&coins)
	return err
}

// Find coin id by symbol
func (r *Coin) FindIdBySymbol(symbol string) (uint64, error) {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	defer db.Close()

	coin := new(models.Coin)
	err := db.Model(coin).
		Column("id").
		Where("symbol = ?", symbol).
		Select()

	if err != nil {
		return 0, err
	}
	return coin.ID, nil
}
