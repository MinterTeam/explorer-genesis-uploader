package domain

type CoinType byte

const (
	_ CoinType = iota
	CoinTypeBase
	CoinTypeToken
	CoinTypePoolToken
)

type Coin struct {
	ID             uint     `json:"id"  pg:",use_zero"`
	Type           CoinType `json:"type"`
	Name           string   `json:"name"`
	Symbol         string   `json:"symbol"`
	Volume         string   `json:"volume"`
	Crr            uint     `json:"crr"`
	Reserve        string   `json:"reserve"`
	MaxSupply      string   `json:"max_supply"`
	Version        uint     `json:"version"    pg:",use_zero"`
	OwnerAddressId uint     `json:"owner_address"`
}
