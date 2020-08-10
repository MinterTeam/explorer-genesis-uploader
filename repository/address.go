package repository

import (
	"github.com/MinterTeam/explorer-genesis-uploader/domain"
	"github.com/go-pg/pg/v10"
	"sync"
)

type Address struct {
	db       *pg.DB
	cache    *sync.Map
	invCache *sync.Map
}

func NewAddressRepository(db *pg.DB) *Address {
	return &Address{
		cache:    new(sync.Map),
		invCache: new(sync.Map),
		db:       db,
	}
}

func (r *Address) SaveAll(addresses []string) error {
	list := make([]*domain.Address, len(addresses))
	for i, a := range addresses {
		list[i] = &domain.Address{Address: a}
	}
	_, err := r.db.Model(&list).Insert()
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

	adr := new(domain.Address)
	err := r.db.Model(adr).Column("id").Where("address = ?", address).Select(adr)
	if err != nil {
		return 0, err
	}
	return adr.ID, nil
}

func (r *Address) GetAddressesCount() (int, error) {
	return r.db.Model((*domain.Address)(nil)).Count()
}

func (r *Address) addToCache(addresses []*domain.Address) {
	for _, a := range addresses {
		_, exist := r.cache.Load(a)
		if !exist {
			r.cache.Store(a.Address, a.ID)
			r.invCache.Store(a.ID, a.Address)
		}
	}
}
