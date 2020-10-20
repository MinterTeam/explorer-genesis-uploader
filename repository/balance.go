package repository

import (
	"github.com/MinterTeam/explorer-genesis-uploader/domain"
	"github.com/go-pg/pg/v10"
)

type Balance struct {
	db *pg.DB
}

func NewBalanceRepository(db *pg.DB) *Balance {
	return &Balance{
		db: db,
	}
}
func (r *Balance) SaveAll(balances []*domain.Balance) error {
	_, err := r.db.Model(&balances).Insert()
	return err
}

func (r *Balance) GetBalancesCount() (int, error) {
	return r.db.Model((*domain.Balance)(nil)).Count()
}
