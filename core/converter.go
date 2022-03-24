package core

import (
	"github.com/MinterTeam/explorer-genesis-uploader/domain"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

func (egu *ExplorerGenesisUploader) convertResponseToModel(response *api_pb.GenesisResponse) *domain.Genesis {
	g := new(domain.Genesis)

	var coins []domain.GenesisCoin
	for _, c := range response.AppState.Coins {
		var ownerAddress string
		if c.OwnerAddress != nil {
			ownerAddress = c.OwnerAddress.Value
		}

		coins = append(coins, domain.GenesisCoin{
			ID:           c.Id,
			Name:         c.Name,
			Symbol:       c.Symbol,
			Volume:       c.Volume,
			Crr:          c.Crr,
			Reserve:      c.Reserve,
			MaxSupply:    c.MaxSupply,
			Version:      c.Version,
			OwnerAddress: &ownerAddress,
			Mintable:     c.Mintable,
			Burnable:     c.Burnable,
		})
	}

	var frozenFunds []domain.FrozenFund
	for _, f := range response.AppState.FrozenFunds {

		var key *string
		if f.CandidateKey != nil {
			key = &f.CandidateKey.Value
		}

		frozenFunds = append(frozenFunds, domain.FrozenFund{
			Height:       f.Height,
			Address:      f.Address,
			CandidateKey: key,
			CandidateID:  f.CandidateId,
			Coin:         f.Coin,
			Value:        f.Value,
		})
	}

	appState := domain.AppState{
		Version:             response.AppState.Version,
		Note:                response.AppState.Note,
		Validators:          nil, //TODO: unuseful for now
		Candidates:          egu.convertCandidates(response),
		DeletedCandidates:   nil, //TODO: unuseful for now
		Coins:               coins,
		FrozenFunds:         frozenFunds,
		BlockListCandidates: nil, //TODO: unuseful for now
		Waitlist:            nil, //TODO: unuseful for now
		Accounts:            egu.convertAccounts(response),
		HaltBlocks:          nil, //TODO: unuseful for now
		Pools:               egu.convertPools(response),
		NextOrderID:         response.AppState.NextOrderId,
		Commission:          domain.Commission{}, //TODO: unuseful for now
		CommissionVotes:     nil,                 //TODO: unuseful for now
		UsedChecks:          nil,                 //TODO: unuseful for now
		MaxGas:              response.AppState.MaxGas,
		TotalSlashed:        response.AppState.TotalSlashed,
	}

	g.AppState = appState
	g.InitialHeight = response.InitialHeight

	return g
}

func (egu *ExplorerGenesisUploader) convertCandidates(response *api_pb.GenesisResponse) []domain.Candidate {
	var candidates []domain.Candidate
	for _, c := range response.AppState.Candidates {
		stakes := egu.convertCandidateStakes(c.Stakes)
		candidates = append(candidates, domain.Candidate{
			ID:                       c.Id,
			RewardAddress:            c.RewardAddress,
			OwnerAddress:             c.OwnerAddress,
			ControlAddress:           c.ControlAddress,
			TotalBipStake:            c.TotalBipStake,
			PublicKey:                c.PublicKey,
			Commission:               c.Commission,
			Stakes:                   stakes,
			Updates:                  nil, //TODO: unuseful for now
			Status:                   c.Status,
			JailedUntil:              c.JailedUntil,
			LastEditCommissionHeight: c.LastEditCommissionHeight,
		})
	}

	return candidates
}

func (egu *ExplorerGenesisUploader) convertCandidateStakes(list []*api_pb.GenesisResponse_AppState_Candidate_Stake) []domain.GenesisStake {
	var stakes []domain.GenesisStake
	for _, s := range list {
		stakes = append(stakes, domain.GenesisStake{
			Owner:    s.Owner,
			Coin:     s.Coin,
			Value:    s.Value,
			BipValue: s.BipValue,
		})
	}
	return stakes
}

func (egu *ExplorerGenesisUploader) convertAccounts(response *api_pb.GenesisResponse) []domain.Account {
	var accounts []domain.Account
	for _, a := range response.AppState.Accounts {

		var msd *domain.MultisigData
		if a.MultisigData != nil {
			msd = &domain.MultisigData{
				Threshold: a.MultisigData.Threshold,
				Weights:   a.MultisigData.Weights,
				Addresses: a.MultisigData.Addresses,
			}
		}

		var balances []domain.GenesisBalance
		for _, b := range a.Balance {
			balances = append(balances, domain.GenesisBalance{
				Coin:  b.Coin,
				Value: b.Value,
			})
		}

		accounts = append(accounts, domain.Account{
			Address:      a.Address,
			Balance:      balances,
			Nonce:        a.Nonce,
			MultisigData: msd,
		})
	}
	return accounts
}

func (egu *ExplorerGenesisUploader) convertPools(response *api_pb.GenesisResponse) []domain.Pool {
	var list []domain.Pool
	for _, p := range response.AppState.Pools {

		var orders []domain.GenesisOrder
		for _, o := range p.Orders {
			orders = append(orders, domain.GenesisOrder{
				IsSale:  o.IsSale,
				Volume0: o.Volume0,
				Volume1: o.Volume1,
				Id:      o.Id,
				Owner:   o.Owner,
				Height:  o.Height,
			})
		}

		list = append(list, domain.Pool{
			Coin0:    p.Coin0,
			Coin1:    p.Coin1,
			Reserve0: p.Reserve0,
			Reserve1: p.Reserve1,
			ID:       p.Id,
			Orders:   orders,
		})
	}

	return list
}
