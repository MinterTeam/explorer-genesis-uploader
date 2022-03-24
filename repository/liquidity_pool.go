package repository

import (
	"github.com/MinterTeam/explorer-genesis-uploader/domain"
	"github.com/go-pg/pg/v10"
)

func (r *LiquidityPool) SaveAll(list []*domain.LiquidityPool) error {
	_, err := r.db.Model(&list).Insert()
	return err
}

func (r *LiquidityPool) SaveAllOrders(orders []domain.Order) error {
	_, err := r.db.Model(&orders).Insert()
	return err
}

type LiquidityPool struct {
	db *pg.DB
}

func NewLiquidityPoolRepository(db *pg.DB) *LiquidityPool {
	return &LiquidityPool{
		db: db,
	}
}
