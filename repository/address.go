package repository

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"github.com/go-pg/pg/v9"
	"os"
)

type Address struct {
}

func NewAddressRepository() *Address {
	return &Address{}
}

func (r Address) SaveAll(addresses []string) error {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	defer db.Close()

	list := make([]*models.Address, len(addresses))
	for i, a := range addresses {
		list[i] = &models.Address{Address: a}
	}
	err := db.Insert(&list)
	return err
}

func (r *Address) FindId(address string) (uint64, error) {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	defer db.Close()

	adr := new(models.Address)
	err := db.Model(adr).Column("id").Where("address = ?", address).Select(adr)
	if err != nil {
		return 0, err
	}
	return adr.ID, nil
}
