package domain

import "time"

type Validator struct {
	ID                   uint       `json:"id" pg:",pk"`
	RewardAddressID      *uint64    `json:"reward_address_id"`
	OwnerAddressID       *uint64    `json:"owner_address_id"`
	CreatedAtBlockID     *uint64    `json:"created_at_block_id"`
	Status               *uint8     `json:"status"`
	PublicKey            string     `json:"public_key"  pg:"type:varchar(64)"`
	Commission           *uint64    `json:"commission"`
	TotalStake           *string    `json:"total_stake"   pg:"type:numeric(70)"`
	Name                 *string    `json:"name"`
	SiteUrl              *string    `json:"site_url"`
	IconUrl              *string    `json:"icon_url"`
	Description          *string    `json:"description"`
	MetaUpdatedAtBlockID *uint64    `json:"meta_updated_at_block_id"`
	UpdateAt             *time.Time `json:"update_at"`
}
