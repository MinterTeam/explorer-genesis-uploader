package domain

type Balance struct {
	AddressID uint64 `json:"address_id" pg:",pk"`
	CoinID    uint64 `json:"coin_id"    pg:",pk,use_zero"`
	Value     string `json:"value"      pg:"type:numeric(70)"`
}
