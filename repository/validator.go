package repository

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"github.com/go-pg/pg/v9"
	"os"
)

type Validator struct {
}

func NewValidatorRepository() *Validator {
	return &Validator{}
}

func (r Validator) SaveAll(validators []*models.Validator) error {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	defer db.Close()

	err := db.Insert(&validators)
	return err
}

//Find validator with public key.
//Return Validator ID
func (r *Validator) FindIdByPk(pk string) (uint64, error) {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	defer db.Close()

	validator := new(models.Validator)
	err := db.Model(validator).Column("id").Where("public_key = ?", pk).Select()
	if err != nil {
		return 0, err
	}
	return validator.ID, nil
}

func (r Validator) SaveAllStakes(stakes []*models.Stake) error {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	defer db.Close()

	err := db.Insert(&stakes)
	return err
}
