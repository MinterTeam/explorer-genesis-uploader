package repository

import (
	"github.com/MinterTeam/explorer-genesis-uploader/domain"
	"github.com/go-pg/pg/v10"
	"sync"
)

type Validator struct {
	cache *sync.Map
	db    *pg.DB
}

func NewValidatorRepository(db *pg.DB) *Validator {
	return &Validator{
		cache: new(sync.Map),
		db:    db,
	}
}

func (r *Validator) SaveAll(validators []*domain.Validator) error {
	_, err := r.db.Model(&validators).Insert()
	return err
}

//Find validator with public key.
//Return Validator ID
func (r *Validator) FindIdByPk(pk string) (uint, error) {
	//First look in the cache
	id, ok := r.cache.Load(pk)
	if ok {
		return id.(uint), nil
	}

	vpk := new(domain.ValidatorPublicKeys)
	err := r.db.Model(vpk).Where("key = ?", pk).Select()
	if err != nil {
		return 0, err
	}

	r.cache.Store(pk, vpk.ValidatorId)
	return vpk.ValidatorId, nil
}

func (r *Validator) SaveAllStakes(stakes []*domain.Stake) error {
	_, err := r.db.Model(&stakes).Insert()
	return err
}

func (r *Validator) SaveAllUnbonds(list []*domain.Unbond) error {
	_, err := r.db.Model(&list).Insert()
	return err
}

func (r *Validator) GetValidatorsCount() (int, error) {
	return r.db.Model((*domain.Validator)(nil)).Count()
}

func (r *Validator) AddPk(key string, validatorId uint) (*domain.ValidatorPublicKeys, error) {
	vpk := &domain.ValidatorPublicKeys{
		Key:         key,
		ValidatorId: validatorId,
	}
	_, err := r.db.Model(vpk).Insert()
	return vpk, err
}

func (r *Validator) SaveAllPk(vpk []*domain.ValidatorPublicKeys) error {
	_, err := r.db.Model(&vpk).Insert()
	return err
}

func (r *Validator) Add(v *domain.Validator) (*domain.Validator, error) {
	_, err := r.db.Model(v).Insert()
	return v, err
}
