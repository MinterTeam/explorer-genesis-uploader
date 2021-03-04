package repository

import (
	"github.com/MinterTeam/explorer-genesis-uploader/domain"
	"github.com/go-pg/pg/v10"
	"os"
	"sync"
)

type Coin struct {
	cache    *sync.Map
	invCache *sync.Map
	db       *pg.DB
}

func NewCoinRepository(db *pg.DB) *Coin {
	return &Coin{
		cache:    new(sync.Map),
		invCache: new(sync.Map),
		db:       db,
	}
}

func (r *Coin) SaveAll(coins []*domain.Coin) error {
	_, err := r.db.Model(&coins).Insert()
	for _, coin := range coins {
		r.cache.Store(coin.Symbol, coin.ID)
		r.invCache.Store(coin.ID, coin.Symbol)
	}
	return err
}

// Find coin id by symbol
func (r *Coin) FindIdBySymbol(symbol string) (uint64, error) {
	//First look in the cache
	id, ok := r.cache.Load(symbol)
	if ok {
		return id.(uint64), nil
	}
	coin := new(domain.Coin)
	err := r.db.Model(coin).
		Column("id").
		Where("symbol = ?", symbol).
		Select()

	if err != nil {
		return 0, err
	}
	return uint64(coin.ID), nil
}

// Find coin id by symbol
func (r *Coin) FindBySymbol(symbol string) (*domain.Coin, error) {
	coin := new(domain.Coin)
	err := r.db.Model(coin).
		Where("symbol = ?", symbol).
		Limit(1).
		Select()

	if err != nil {
		return nil, err
	}
	return coin, nil
}

func (r *Coin) GetCoinsCount() (int, error) {
	return r.db.Model((*domain.Coin)(nil)).Where("symbol != ?", os.Getenv("MINTER_BASE_COIN")).Count()
}

func (r *Coin) ChangeSequence(i int) error {
	_, err := r.db.Model().Exec(`
		alter sequence coins_id_seq START WITH ?;
	`, i)
	return err
}
