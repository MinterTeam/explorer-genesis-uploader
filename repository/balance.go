package repository

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"github.com/go-pg/pg/v9"
	"os"
)

type Balance struct {
	db *pg.DB
}

func NewBalanceRepository() *Balance {
	return &Balance{
		db: pg.Connect(&pg.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
			User:     os.Getenv("DB_USER"),
			Database: os.Getenv("DB_NAME"),
			Password: os.Getenv("DB_PASSWORD"),
		}),
	}
}
func (r *Balance) SaveAll(balances []*models.Balance) error {
	err := r.db.Insert(&balances)
	return err
}

func (r *Balance) Close() {
	r.db.Close()
}
