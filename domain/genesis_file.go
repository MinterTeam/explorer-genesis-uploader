package domain

type GenesisFile struct {
	GenesisTime   string              `json:"genesis_time"`
	ChainID       string              `json:"chain_id"`
	InitialHeight string              `json:"initial_height"`
	AppHash       string              `json:"app_hash"`
	AppState      GenesisFileAppState `json:"app_state"`
}

type GenesisFileAppState struct {
	Version      string                  `json:"version"`
	Note         string                  `json:"note"`
	Candidates   []GenesisFileCandidate  `json:"candidates"`
	Coins        []GenesisFileCoin       `json:"coins"`
	FrozenFunds  []GenesisFileFrozenFund `json:"frozen_funds"`
	Waitlist     []GenesisFileWaitlist   `json:"waitlist"`
	Accounts     []GenesisFileAccount    `json:"accounts"`
	Pools        []GenesisFilePool       `json:"pools"`
	NextOrderID  string                  `json:"next_order_id"`
	MaxGas       string                  `json:"max_gas"`
	TotalSlashed string                  `json:"total_slashed"`
}

type GenesisFileAccount struct {
	Address      string                   `json:"address"`
	Balance      []GenesisFileBalance     `json:"balance"`
	Nonce        string                   `json:"nonce"`
	MultisigData *GenesisFileMultisigData `json:"multisig_data"`
}

type GenesisFileBalance struct {
	Coin  string `json:"coin"`
	Value string `json:"value"`
}

type GenesisFileMultisigData struct {
	Threshold string   `json:"threshold"`
	Weights   []string `json:"weights"`
	Addresses []string `json:"addresses"`
}

type GenesisFileCandidate struct {
	ID                       string                `json:"id"`
	RewardAddress            string                `json:"reward_address"`
	OwnerAddress             string                `json:"owner_address"`
	ControlAddress           string                `json:"control_address"`
	TotalBipStake            string                `json:"total_bip_stake"`
	PublicKey                string                `json:"public_key"`
	Commission               string                `json:"commission"`
	Stakes                   []GenesisFileWaitlist `json:"stakes"`
	Updates                  []interface{}         `json:"updates"`
	Status                   string                `json:"status"`
	JailedUntil              string                `json:"jailed_until"`
	LastEditCommissionHeight string                `json:"last_edit_commission_height"`
}

type GenesisFileWaitlist struct {
	Owner       string  `json:"owner"`
	Coin        string  `json:"coin"`
	Value       string  `json:"value"`
	BipValue    *string `json:"bip_value,omitempty"`
	CandidateID *string `json:"candidate_id,omitempty"`
}

type GenesisFileCoin struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	Volume       string  `json:"volume"`
	Crr          string  `json:"crr"`
	Reserve      string  `json:"reserve"`
	MaxSupply    string  `json:"max_supply"`
	Version      string  `json:"version"`
	OwnerAddress *string `json:"owner_address"`
	Mintable     bool    `json:"mintable"`
	Burnable     bool    `json:"burnable"`
}

type GenesisFileFrozenFund struct {
	Height       string  `json:"height"`
	Address      string  `json:"address"`
	CandidateKey *string `json:"candidate_key"`
	CandidateID  string  `json:"candidate_id"`
	Coin         string  `json:"coin"`
	Value        string  `json:"value"`
}

type GenesisFilePool struct {
	Coin0    string             `json:"coin0"`
	Coin1    string             `json:"coin1"`
	Reserve0 string             `json:"reserve0"`
	Reserve1 string             `json:"reserve1"`
	ID       string             `json:"id"`
	Orders   []GenesisFileOrder `json:"orders"`
}

type GenesisFileOrder struct {
	IsSale  bool   `json:"is_sale"`
	Volume0 string `json:"volume0"`
	Volume1 string `json:"volume1"`
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Height  string `json:"height"`
}
