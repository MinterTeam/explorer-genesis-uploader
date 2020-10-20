package domain

type Unbond struct {
	BlockId     uint   `json:"block_id"`
	AddressId   uint   `json:"address_id"`
	CoinId      uint   `json:"coin_id" pg:",use_zero"`
	ValidatorId uint   `json:"validator_id"`
	Value       string `json:"value"`
}
