package core

import "time"

type GenesisResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Result  struct {
		Genesis Genesis `json:"genesis"`
	} `json:"result"`
}

type Genesis struct {
	GenesisTime     time.Time `json:"genesis_time"`
	ChainID         string    `json:"chain_id"`
	ConsensusParams struct {
		Block struct {
			MaxBytes   string `json:"max_bytes"`
			MaxGas     string `json:"max_gas"`
			TimeIotaMs string `json:"time_iota_ms"`
		} `json:"block"`
		Evidence struct {
			MaxAge string `json:"max_age"`
		} `json:"evidence"`
		Validator struct {
			PubKeyTypes []string `json:"pub_key_types"`
		} `json:"validator"`
	} `json:"consensus_params"`
	AppHash  string   `json:"app_hash"`
	AppState AppState `json:"app_state"`
}

type AppState struct {
	StartHeight  string      `json:"start_height"`
	Validators   []Validator `json:"validators"`
	Candidates   []Candidate `json:"candidates"`
	Accounts     []Account   `json:"accounts"`
	Coins        []Coin      `json:"coins"`
	UsedChecks   []string    `json:"used_checks"`
	MaxGas       string      `json:"max_gas"`
	TotalSlashed string      `json:"total_slashed"`
}

type Validator struct {
	RewardAddress string `json:"reward_address"`
	TotalBipStake string `json:"total_bip_stake"`
	PubKey        string `json:"pub_key"`
	Commission    string `json:"commission"`
	AccumReward   string `json:"accum_reward"`
	AbsentTimes   string `json:"absent_times"`
}

type Candidate struct {
	RewardAddress  string  `json:"reward_address"`
	OwnerAddress   string  `json:"owner_address"`
	TotalBipStake  string  `json:"total_bip_stake"`
	PubKey         string  `json:"pub_key"`
	Commission     string  `json:"commission"`
	Stakes         []Stake `json:"stakes"`
	CreatedAtBlock string  `json:"created_at_block"`
	Status         int     `json:"status"`
}

type Stake struct {
	Owner    string `json:"owner"`
	Coin     string `json:"coin"`
	Value    string `json:"value"`
	BipValue string `json:"bip_value"`
}

type Coin struct {
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	Volume         string `json:"volume"`
	Crr            string `json:"crr"`
	ReserveBalance string `json:"reserve_balance"`
}

type Account struct {
	Address string    `json:"address"`
	Balance []Balance `json:"balance"`
	Nonce   string    `json:"nonce"`
}

type Balance struct {
	Coin  string `json:"coin"`
	Value string `json:"value"`
}
