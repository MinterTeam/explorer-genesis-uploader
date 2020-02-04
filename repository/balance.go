package repository

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"github.com/go-pg/pg/v9"
	"os"
)

type Balance struct {
}

func NewBalanceRepository() *Balance {
	return &Balance{}
}
func (r Balance) SaveAll(balances []*models.Balance) error {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	defer db.Close()

	err := db.Insert(&balances)
	return err
}
