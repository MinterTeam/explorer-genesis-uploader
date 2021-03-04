package domain

import "fmt"

type LiquidityPool struct {
	Id               uint64 `json:"id"`
	TokenId          uint64 `json:"token_id"`
	FirstCoinId      uint64 `json:"first_coin_id"  pg:",use_zero"`
	SecondCoinId     uint64 `json:"second_coin_id" pg:",use_zero"`
	FirstCoinVolume  string `json:"first_coin_volume"`
	SecondCoinVolume string `json:"second_coin_volume"`
	Liquidity        string `json:"liquidity"`
}

type AddressLiquidityPool struct {
	LiquidityPoolId uint64 `json:"liquidity_pool_id" pg:",pk"`
	AddressId       uint64 `json:"address_id"        pg:",pk"`
	Liquidity       string `json:"liquidity"`
}

func (lp *LiquidityPool) GetTokenSymbol() string {
	return fmt.Sprintf("P-%d", lp.Id)
}
