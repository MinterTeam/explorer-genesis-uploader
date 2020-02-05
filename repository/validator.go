package repository

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"github.com/go-pg/pg/v9"
	"os"
	"sync"
)

type Validator struct {
	cache *sync.Map
	db    *pg.DB
}

func NewValidatorRepository() *Validator {
	return &Validator{
		cache: new(sync.Map),
		db: pg.Connect(&pg.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
			User:     os.Getenv("DB_USER"),
			Database: os.Getenv("DB_NAME"),
			Password: os.Getenv("DB_PASSWORD"),
		}),
	}
}

func (r *Validator) SaveAll(validators []*models.Validator) error {
	err := r.db.Insert(&validators)
	for _, v := range validators {
		r.cache.Store(v.PublicKey, v.ID)
	}
	return err
}

//Find validator with public key.
//Return Validator ID
func (r *Validator) FindIdByPk(pk string) (uint64, error) {
	//First look in the cache
	id, ok := r.cache.Load(pk)
	if ok {
		return id.(uint64), nil
	}

	validator := new(models.Validator)
	err := r.db.Model(validator).Column("id").Where("public_key = ?", pk).Select()
	if err != nil {
		return 0, err
	}
	return validator.ID, nil
}

func (r *Validator) SaveAllStakes(stakes []*models.Stake) error {
	err := r.db.Insert(&stakes)
	return err
}

func (r *Validator) Close() {
	r.db.Close()
}
