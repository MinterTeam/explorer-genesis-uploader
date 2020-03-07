package repository

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"github.com/go-pg/pg/v9"
	"os"
	"sync"
)

type Coin struct {
	cache    *sync.Map
	invCache *sync.Map
	db       *pg.DB
}

func NewCoinRepository() *Coin {
	return &Coin{
		cache:    new(sync.Map),
		invCache: new(sync.Map),
		db: pg.Connect(&pg.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
			User:     os.Getenv("DB_USER"),
			Database: os.Getenv("DB_NAME"),
			Password: os.Getenv("DB_PASSWORD"),
		}),
	}
}

func (r *Coin) SaveAll(coins []*models.Coin) error {
	err := r.db.Insert(&coins)
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
	coin := new(models.Coin)
	err := r.db.Model(coin).
		Column("id").
		Where("symbol = ?", symbol).
		Select()

	if err != nil {
		return 0, err
	}
	return coin.ID, nil
}

func (r *Coin) GetCoinsCount() (int, error) {
	return r.db.Model((*models.Coin)(nil)).Where("symbol != ?", os.Getenv("MINTER_BASE_COIN")).Count()
}

func (r *Coin) Close() {
	r.db.Close()
}
