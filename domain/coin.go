package domain

type Coin struct {
	ID             uint   `json:"id"  pg:",use_zero"`
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	Volume         string `json:"volume"`
	Crr            uint   `json:"crr"`
	Reserve        string `json:"reserve"`
	MaxSupply      string `json:"max_supply"`
	Version        uint   `json:"version"`
	OwnerAddressId uint   `json:"owner_address"`
}
