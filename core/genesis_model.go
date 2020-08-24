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
	ID             uint    `json:"id"`
	RewardAddress  string  `json:"reward_address"`
	OwnerAddress   string  `json:"owner_address"`
	ControlAddress string  `json:"control_address"`
	TotalBipStake  string  `json:"total_bip_stake"`
	PubKey         string  `json:"pub_key"`
	Commission     string  `json:"commission"`
	Stakes         []Stake `json:"stakes"`
	Updates        []Stake `json:"updates"`
	Status         byte    `json:"status"`
}

type Stake struct {
	Owner    string `json:"owner"`
	Coin     uint   `json:"coin"`
	Value    string `json:"value"`
	BipValue string `json:"bip_value"`
}

type Coin struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Volume       string `json:"volume"`
	Crr          string `json:"crr"`
	Reserve      string `json:"reserve"`
	MaxSupply    string `json:"max_supply"`
	Version      uint   `json:"version"`
	OwnerAddress string `json:"owner_address"`
}

type Account struct {
	Address      string    `json:"address"`
	Balance      []Balance `json:"balance"`
	Nonce        string    `json:"nonce"`
	MultisigData *Multisig `json:"multisig_data,omitempty"`
}

type Balance struct {
	Coin  uint   `json:"coin"`
	Value string `json:"value"`
}

type Multisig struct {
	Weights   []string `json:"weights"`
	Threshold string   `json:"threshold"`
	Addresses []string `json:"addresses"`
}
