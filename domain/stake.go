package domain

type Stake struct {
	ID             uint   `json:"id"               pg:",pk"`
	OwnerAddressID uint64 `json:"owner_address_id"`
	ValidatorID    uint   `json:"validator_id"`
	CoinID         uint64 `json:"coin_id"          pg:",use_zero"`
	Value          string `json:"value"            pg:"type:numeric(70)"`
	BipValue       string `json:"bip_value"        pg:"type:numeric(70)"`
}
