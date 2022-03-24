package domain

type Genesis struct {
	GenesisTime     string          `json:"genesis_time"`
	ChainID         string          `json:"chain_id"`
	InitialHeight   uint64          `json:"initial_height"`
	ConsensusParams ConsensusParams `json:"consensus_params"`
	AppHash         string          `json:"app_hash"`
	AppState        AppState        `json:"app_state"`
}

type AppState struct {
	Version             string             `json:"version"`
	Note                string             `json:"note"`
	Validators          []ValidatorElement `json:"validators"`
	Candidates          []Candidate        `json:"candidates"`
	DeletedCandidates   []DeletedCandidate `json:"deleted_candidates"`
	Coins               []GenesisCoin      `json:"coins"`
	FrozenFunds         []FrozenFund       `json:"frozen_funds"`
	BlockListCandidates []string           `json:"block_list_candidates"`
	Waitlist            []Waitlist         `json:"waitlist"`
	Accounts            []Account          `json:"accounts"`
	HaltBlocks          []interface{}      `json:"halt_blocks"`
	Pools               []Pool             `json:"pools"`
	NextOrderID         uint64             `json:"next_order_id"`
	Commission          Commission         `json:"commission"`
	CommissionVotes     []interface{}      `json:"commission_votes"`
	UsedChecks          []string           `json:"used_checks"`
	MaxGas              uint64             `json:"max_gas"`
	TotalSlashed        string             `json:"total_slashed"`
}

type Account struct {
	Address      string           `json:"address"`
	Balance      []GenesisBalance `json:"balance"`
	Nonce        uint64           `json:"nonce"`
	MultisigData *MultisigData    `json:"multisig_data"`
}

type GenesisBalance struct {
	Coin  uint64 `json:"coin"`
	Value string `json:"value"`
}

type MultisigData struct {
	Threshold uint64   `json:"threshold"`
	Weights   []uint64 `json:"weights"`
	Addresses []string `json:"addresses"`
}

type Candidate struct {
	ID                       uint64         `json:"id"`
	RewardAddress            string         `json:"reward_address"`
	OwnerAddress             string         `json:"owner_address"`
	ControlAddress           string         `json:"control_address"`
	TotalBipStake            string         `json:"total_bip_stake"`
	PublicKey                string         `json:"public_key"`
	Commission               uint64         `json:"commission"`
	Stakes                   []GenesisStake `json:"stakes"`
	Updates                  []interface{}  `json:"updates"`
	Status                   int64          `json:"status"`
	JailedUntil              int64          `json:"jailed_until"`
	LastEditCommissionHeight int64          `json:"last_edit_commission_height"`
}

type Waitlist struct {
	Owner       string `json:"owner"`
	Coin        uint64 `json:"coin"`
	Value       string `json:"value"`
	BipValue    string `json:"bip_value"`
	CandidateID string `json:"candidate_id"`
}

type GenesisCoin struct {
	ID           uint64  `json:"id"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	Volume       string  `json:"volume"`
	Crr          uint64  `json:"crr"`
	Reserve      string  `json:"reserve"`
	MaxSupply    string  `json:"max_supply"`
	Version      uint64  `json:"version"`
	OwnerAddress *string `json:"owner_address"`
	Mintable     bool    `json:"mintable"`
	Burnable     bool    `json:"burnable"`
}

type Commission struct {
	Coin                    string `json:"coin"`
	PayloadByte             string `json:"payload_byte"`
	Send                    string `json:"send"`
	BuyBancor               string `json:"buy_bancor"`
	SellBancor              string `json:"sell_bancor"`
	SellAllBancor           string `json:"sell_all_bancor"`
	BuyPoolBase             string `json:"buy_pool_base"`
	BuyPoolDelta            string `json:"buy_pool_delta"`
	SellPoolBase            string `json:"sell_pool_base"`
	SellPoolDelta           string `json:"sell_pool_delta"`
	SellAllPoolBase         string `json:"sell_all_pool_base"`
	SellAllPoolDelta        string `json:"sell_all_pool_delta"`
	CreateTicker3           string `json:"create_ticker3"`
	CreateTicker4           string `json:"create_ticker4"`
	CreateTicker5           string `json:"create_ticker5"`
	CreateTicker6           string `json:"create_ticker6"`
	CreateTicker710         string `json:"create_ticker7_10"`
	CreateCoin              string `json:"create_coin"`
	CreateToken             string `json:"create_token"`
	RecreateCoin            string `json:"recreate_coin"`
	RecreateToken           string `json:"recreate_token"`
	DeclareCandidacy        string `json:"declare_candidacy"`
	Delegate                string `json:"delegate"`
	Unbond                  string `json:"unbond"`
	RedeemCheck             string `json:"redeem_check"`
	SetCandidateOn          string `json:"set_candidate_on"`
	SetCandidateOff         string `json:"set_candidate_off"`
	CreateMultisig          string `json:"create_multisig"`
	MultisendBase           string `json:"multisend_base"`
	MultisendDelta          string `json:"multisend_delta"`
	EditCandidate           string `json:"edit_candidate"`
	SetHaltBlock            string `json:"set_halt_block"`
	EditTickerOwner         string `json:"edit_ticker_owner"`
	EditMultisig            string `json:"edit_multisig"`
	EditCandidatePublicKey  string `json:"edit_candidate_public_key"`
	CreateSwapPool          string `json:"create_swap_pool"`
	AddLiquidity            string `json:"add_liquidity"`
	RemoveLiquidity         string `json:"remove_liquidity"`
	EditCandidateCommission string `json:"edit_candidate_commission"`
	MintToken               string `json:"mint_token"`
	BurnToken               string `json:"burn_token"`
	VoteCommission          string `json:"vote_commission"`
	VoteUpdate              string `json:"vote_update"`
	FailedTx                string `json:"failed_tx"`
	AddLimitOrder           string `json:"add_limit_order"`
	RemoveLimitOrder        string `json:"remove_limit_order"`
}

type DeletedCandidate struct {
	ID        string `json:"id"`
	PublicKey string `json:"public_key"`
}

type FrozenFund struct {
	Height       uint64  `json:"height"`
	Address      string  `json:"address"`
	CandidateKey *string `json:"candidate_key"`
	CandidateID  uint64  `json:"candidate_id"`
	Coin         uint64  `json:"coin"`
	Value        string  `json:"value"`
}

type GenesisStake struct {
	Owner    string `json:"owner"`
	Coin     uint64 `json:"coin,omitempty"`
	Value    string `json:"value,omitempty"`
	BipValue string `json:"bip_value,omitempty"`
}

type Pool struct {
	Coin0    uint64         `json:"coin0"`
	Coin1    uint64         `json:"coin1"`
	Reserve0 string         `json:"reserve0"`
	Reserve1 string         `json:"reserve1"`
	ID       uint64         `json:"id"`
	Orders   []GenesisOrder `json:"orders"`
}

type GenesisOrder struct {
	IsSale  bool   `json:"is_sale"`
	Volume0 string `json:"volume0"` // buy
	Volume1 string `json:"volume1"` // sell
	Id      uint64 `json:"id"`
	Owner   string `json:"owner"`
	Height  uint64 `json:"height"`
}

type ValidatorElement struct {
	TotalBipStake string `json:"total_bip_stake"`
	PublicKey     string `json:"public_key"`
	AccumReward   string `json:"accum_reward"`
	AbsentTimes   string `json:"absent_times"`
}

type ConsensusParams struct {
	Block     Block                    `json:"block"`
	Evidence  Evidence                 `json:"evidence"`
	Validator ConsensusParamsValidator `json:"validator"`
}

type Block struct {
	MaxBytes   string `json:"max_bytes"`
	MaxGas     string `json:"max_gas"`
	TimeIotaMS string `json:"time_iota_ms"`
}

type Evidence struct {
	MaxAgeNumBlocks string `json:"max_age_num_blocks"`
	MaxAgeDuration  string `json:"max_age_duration"`
}

type ConsensusParamsValidator struct {
	PubKeyTypes []string `json:"pub_key_types"`
}
