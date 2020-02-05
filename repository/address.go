package repository

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"github.com/go-pg/pg/v9"
	"os"
	"sync"
)

type Address struct {
	db       *pg.DB
	cache    *sync.Map
	invCache *sync.Map
}

func NewAddressRepository() *Address {

	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	return &Address{
		cache:    new(sync.Map),
		invCache: new(sync.Map),
		db:       db,
	}
}

func (r *Address) SaveAll(addresses []string) error {
	list := make([]*models.Address, len(addresses))
	for i, a := range addresses {
		list[i] = &models.Address{Address: a}
	}
	err := r.db.Insert(&list)
	if err == nil {
		r.addToCache(list)
	}
	return err
}

func (r *Address) FindId(address string) (uint64, error) {
	//First look in the cache
	id, ok := r.cache.Load(address)
	if ok {
		return id.(uint64), nil
	}

	adr := new(models.Address)
	err := r.db.Model(adr).Column("id").Where("address = ?", address).Select(adr)
	if err != nil {
		return 0, err
	}
	return adr.ID, nil
}

func (r *Address) GetAddressesCount() (int, error) {
	return r.db.Model((*models.Address)(nil)).Count()
}

func (r *Address) Close() {
	r.db.Close()
}

func (r *Address) addToCache(addresses []*models.Address) {
	for _, a := range addresses {
		_, exist := r.cache.Load(a)
		if !exist {
			r.cache.Store(a.Address, a.ID)
			r.invCache.Store(a.ID, a.Address)
		}
	}
}
