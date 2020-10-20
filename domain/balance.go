package domain

type Balance struct {
	ID        uint64 `json:"id" pg:",pk"`
	AddressID uint64 `json:"address_id"`
	CoinID    uint64 `json:"coin_id"  pg:",use_zero"`
	Value     string `json:"value"    pg:"type:numeric(70)"`
}
