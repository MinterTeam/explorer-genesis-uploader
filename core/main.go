package core

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/MinterTeam/explorer-genesis-uploader/domain"
	"github.com/MinterTeam/explorer-genesis-uploader/env"
	"github.com/MinterTeam/explorer-genesis-uploader/helpers"
	"github.com/MinterTeam/explorer-genesis-uploader/repository"
	"github.com/MinterTeam/minter-go-sdk/v2/api/grpc_client"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
	"math"
	"math/big"
	"os"
	"sync"
	"time"
)

//-file=./tmp/genesis.json
var file = flag.String(`file`, "", `Path to genesis json file`)

type ExplorerGenesisUploader struct {
	startBlock              uint64
	addressRepository       *repository.Address
	balanceRepository       *repository.Balance
	coinRepository          *repository.Coin
	validatorRepository     *repository.Validator
	liquidityPoolRepository *repository.LiquidityPool
	logger                  *logrus.Entry
	env                     env.Config
}

func (egu *ExplorerGenesisUploader) StartBlock() uint64 {
	return egu.startBlock
}

func New(cfg env.Config) *ExplorerGenesisUploader {
	//Init Logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(false)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	contextLogger := logger.WithFields(logrus.Fields{
		"version": "1.4.0",
		"app":     "Minter Explorer Explorer Genesis Uploader",
	})

	pgOptions := &pg.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.PostgresHost, cfg.PostgresPort),
		User:     cfg.PostgresUser,
		Database: cfg.PostgresDB,
		Password: cfg.PostgresPassword,
	}

	if cfg.PostgresSSLEnabled {
		pgOptions.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	db := pg.Connect(pgOptions)

	// Repositories
	addressRepository := repository.NewAddressRepository(db)
	coinRepository := repository.NewCoinRepository(db)
	validatorRepository := repository.NewValidatorRepository(db)
	balanceRepository := repository.NewBalanceRepository(db)
	liquidityPoolRepository := repository.NewLiquidityPoolRepository(db)

	return &ExplorerGenesisUploader{
		env:                     cfg,
		addressRepository:       addressRepository,
		balanceRepository:       balanceRepository,
		coinRepository:          coinRepository,
		validatorRepository:     validatorRepository,
		liquidityPoolRepository: liquidityPoolRepository,
		logger:                  contextLogger,
	}
}

func (egu *ExplorerGenesisUploader) Do() error {

	if !egu.isEmptyDB() {
		return errors.New("genesis has not been uploaded DB is not empty")
	}

	start := time.Now()
	egu.logger.Info("Getting genesis data...")

	var genesis *domain.Genesis

	if *file != "" {

		var gf *domain.GenesisFile

		jsonFile, err := os.Open(*file)
		if err != nil {
			panic(err)
		}
		defer jsonFile.Close()

		dec := json.NewDecoder(jsonFile)
		dec.UseNumber()

		err = dec.Decode(gf)
		if err != nil {
			panic(err)
		}

	} else {
		client, err := grpc_client.New(egu.env.NodeGrpc)
		if err != nil {
			panic(err)
		}
		genesisResponse, err := client.Genesis()
		if err != nil {
			panic(err)
		}

		egu.startBlock = genesisResponse.InitialHeight

		genesis = egu.convertResponseToModel(genesisResponse)
	}

	egu.startBlock = genesis.InitialHeight

	egu.logger.Info(fmt.Sprintf("Genesis has been downloaded. Processing time %s", time.Since(start)))

	egu.logger.Info("Extracting addresses...")
	startOperation := time.Now()
	addresses, err := egu.extractAddresses(genesis)
	if err != nil {
		panic(err)
	}
	egu.logger.Info(fmt.Sprintf("%d addresses has been extracted. Processing time %s", len(addresses), time.Since(startOperation)))
	startOperation = time.Now()
	egu.saveAddresses(addresses)
	egu.logger.Info(fmt.Sprintf("Addresses has been saved. Processing time %s", time.Since(startOperation)))

	egu.logger.Info("Extracting coins...")
	startOperation = time.Now()
	coins, err := egu.extractCoins(genesis)
	if err != nil {
		panic(err)
	}
	egu.logger.Info(fmt.Sprintf("%d coins has been extracted. Processing time %s", len(coins)+1, time.Since(startOperation)))
	startOperation = time.Now()
	egu.saveCoins(coins)
	egu.logger.Info(fmt.Sprintf("Coins has been saved. Processing time %s", time.Since(startOperation)))

	egu.logger.Info("Extracting validators...")
	startOperation = time.Now()
	validators, err := egu.extractCandidates(genesis)
	if err != nil {
		panic(err)
	}
	egu.logger.Info(fmt.Sprintf("%d validators have been extracted. Processing time %s", len(validators), time.Since(startOperation)))
	startOperation = time.Now()
	err = egu.saveCandidates(validators)
	if err != nil {
		panic(err)
	}
	egu.logger.Info(fmt.Sprintf("Validators has been saved. Processing time %s", time.Since(startOperation)))

	egu.logger.Info("Extracting balances...")
	startOperation = time.Now()
	balances, err := egu.extractBalances(genesis)
	if err != nil {
		panic(err)
	}
	egu.logger.Info(fmt.Sprintf("%d balances has been extracted. Processing time %s", len(balances), time.Since(startOperation)))
	startOperation = time.Now()
	err = egu.saveBalances(balances)
	if err != nil {
		panic(err)
	}
	egu.logger.Info(fmt.Sprintf("Balances has been saved. Processing time %s", time.Since(startOperation)))

	egu.logger.Info("Extracting stakes...")
	startOperation = time.Now()
	stakes, err := egu.extractStakes(genesis)
	if err != nil {
		panic(err)
	}
	egu.logger.Info(fmt.Sprintf("%d stakes have been extracted. Processing time %s", len(stakes), time.Since(startOperation)))
	startOperation = time.Now()
	err = egu.saveStakes(stakes)
	if err != nil {
		panic(err)
	}
	egu.logger.Info(fmt.Sprintf("Stakes has been saved. Processing time %s", time.Since(startOperation)))

	egu.logger.Info("Extracting unbonds...")
	startOperation = time.Now()
	unbonds, err := egu.extractUnbonds(genesis)
	if err != nil {
		//panic(err)
		egu.logger.Error(err)
	}
	egu.logger.Info(fmt.Sprintf("%d unbonds have been extracted. Processing time %s", len(unbonds), time.Since(startOperation)))
	startOperation = time.Now()
	err = egu.saveUnbonds(unbonds)
	if err != nil {
		egu.logger.Error(fmt.Sprintf("Unbonds saving error: %s", err))
	} else {
		egu.logger.Info(fmt.Sprintf("Unbonds has been saved. Processing time %s", time.Since(startOperation)))
	}

	egu.logger.Info("Extracting liquidity pools...")
	startOperation = time.Now()
	lpList, err := egu.extractLiquidityPool(genesis)
	if err != nil {
		panic(err)
	}
	egu.logger.Info(fmt.Sprintf("%d liquidity pools have been extracted. Processing time %s", len(lpList), time.Since(startOperation)))
	startOperation = time.Now()
	err = egu.saveLiquidityPool(lpList)
	if err != nil {
		panic(err)
	}
	egu.logger.Info(fmt.Sprintf("Liquidity pools has been saved. Processing time %s", time.Since(startOperation)))

	egu.logger.Info("Extracting orders...")
	startOperation = time.Now()
	orderList, err := egu.extractOrders(genesis)
	if err != nil {
		panic(err)
	}
	egu.logger.Info(fmt.Sprintf("%d orders has been extracted. Processing time %s", len(orderList), time.Since(startOperation)))
	startOperation = time.Now()
	err = egu.saveOrders(orderList)
	if err != nil {
		panic(err)
	}
	egu.logger.Info(fmt.Sprintf("Orders has been saved. Processing time %s", time.Since(startOperation)))

	egu.logger.Info("Upload complete")
	elapsed := time.Since(start)
	egu.logger.Info("Processing time: ", elapsed)
	return err
}

func (egu *ExplorerGenesisUploader) extractAddresses(genesis *domain.Genesis) ([]string, error) {
	addressesMap := make(map[string]struct{})
	addressesMap["0000000000000000000000000000000000000000"] = struct{}{}

	for _, candidate := range genesis.AppState.Candidates {
		addressesMap[helpers.RemovePrefix(candidate.RewardAddress)] = struct{}{}
		addressesMap[helpers.RemovePrefix(candidate.OwnerAddress)] = struct{}{}
		for _, stake := range candidate.Stakes {
			addressesMap[helpers.RemovePrefix(stake.Owner)] = struct{}{}
		}
	}
	for _, account := range genesis.AppState.Accounts {
		addressesMap[helpers.RemovePrefix(account.Address)] = struct{}{}
	}
	for _, coin := range genesis.AppState.Coins {
		if coin.OwnerAddress != nil && *coin.OwnerAddress != "" {
			addressesMap[helpers.RemovePrefix(*coin.OwnerAddress)] = struct{}{}
		}
	}
	for _, data := range genesis.AppState.FrozenFunds {
		addressesMap[helpers.RemovePrefix(data.Address)] = struct{}{}
	}

	var addresses = make([]string, len(addressesMap))
	i := 0
	for adr := range addressesMap {
		addresses[i] = adr
		i++
	}
	return addresses, nil
}

func (egu *ExplorerGenesisUploader) extractCoins(genesis *domain.Genesis) ([]*domain.Coin, error) {
	var coins = make([]*domain.Coin, len(genesis.AppState.Coins)+1)
	i := 1

	coins[0] = &domain.Coin{
		ID:        0,
		Type:      domain.CoinTypeBase,
		Name:      "Base coin",
		Symbol:    egu.env.MinterBaseCoin,
		Volume:    "0",
		Crr:       100,
		Reserve:   "0",
		MaxSupply: "0",
		Version:   0,
	}

	for _, c := range genesis.AppState.Coins {
		if c.ID == 0 {
			continue
		}

		coins[i] = &domain.Coin{
			ID:        uint(c.ID),
			Type:      domain.CoinTypeBase,
			Name:      c.Name,
			Symbol:    c.Symbol,
			Volume:    c.Volume,
			Crr:       uint(c.Crr),
			Reserve:   c.Reserve,
			MaxSupply: c.MaxSupply,
			Version:   uint(c.Version),
		}
		if c.OwnerAddress != nil && *c.OwnerAddress != "" {
			addressId, err := egu.addressRepository.FindId(helpers.RemovePrefix(*c.OwnerAddress))
			if err != nil {
				egu.logger.Error(err)
			} else {
				coins[i].OwnerAddressId = uint(addressId)
			}
		}
		i++
	}
	return coins, nil
}

func (egu ExplorerGenesisUploader) extractCandidates(genesis *domain.Genesis) ([]*domain.Validator, error) {
	var validators []*domain.Validator
	for _, candidate := range genesis.AppState.Candidates {
		ownerAddress, err := egu.addressRepository.FindId(helpers.RemovePrefix(candidate.OwnerAddress))
		if err != nil {
			egu.logger.Error(err)
		}
		rewardAddress, err := egu.addressRepository.FindId(helpers.RemovePrefix(candidate.RewardAddress))
		if err != nil {
			egu.logger.Error(err)
		}

		status := uint8(candidate.Status)
		commission := candidate.Commission
		stake := candidate.TotalBipStake

		validator := &domain.Validator{
			ID:              uint(candidate.ID),
			PublicKey:       helpers.RemovePrefix(candidate.PublicKey),
			OwnerAddressID:  &ownerAddress,
			RewardAddressID: &rewardAddress,
			Status:          &status,
			Commission:      &commission,
			TotalStake:      &stake,
		}

		validators = append(validators, validator)
	}

	return validators, nil
}

func (egu *ExplorerGenesisUploader) saveAddresses(addresses []string) {
	egu.logger.Info("Saving addresses to DB...")
	if len(addresses) > 0 {
		wgAddresses := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(addresses)) / float64(egu.env.AddressChunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := int(egu.env.AddressChunkSize) * i
			end := start + int(egu.env.AddressChunkSize)
			if end > len(addresses) {
				end = len(addresses)
			}
			wgAddresses.Add(1)
			go func() {
				err := egu.addressRepository.SaveAll(addresses[start:end])
				if err != nil {
					panic(err)
				}
				wgAddresses.Done()
			}()
		}
		wgAddresses.Wait()
	}
}

func (egu *ExplorerGenesisUploader) saveCoins(coins []*domain.Coin) {
	egu.logger.Info("Saving coins to DB...")
	var list []*domain.Coin
	list = append(list, coins...)
	wgCoins := new(sync.WaitGroup)
	chunksCount := int(math.Ceil(float64(len(list)) / float64(egu.env.CoinsChunkSize)))
	for i := 0; i < chunksCount; i++ {
		start := int(egu.env.CoinsChunkSize) * i
		end := start + int(egu.env.CoinsChunkSize)
		if end > len(list) {
			end = len(list)
		}
		wgCoins.Add(1)
		go func() {
			err := egu.coinRepository.SaveAll(list[start:end])
			if err != nil {
				panic(err)
			}
			wgCoins.Done()
		}()
	}
	wgCoins.Wait()
}

func (egu *ExplorerGenesisUploader) saveCandidates(validators []*domain.Validator) error {
	egu.logger.Info("Saving validators to DB...")

	if len(validators) > 0 {
		err := egu.validatorRepository.SaveAll(validators)
		if err != nil {
			panic(err)
		}

		var vpk []*domain.ValidatorPublicKeys

		for _, v := range validators {
			vpk = append(vpk, &domain.ValidatorPublicKeys{
				ValidatorId: v.ID,
				Key:         v.PublicKey,
			})
		}

		err = egu.validatorRepository.SaveAllPk(vpk)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (egu *ExplorerGenesisUploader) extractBalances(genesis *domain.Genesis) ([]*domain.Balance, error) {
	chunkSize := 1000
	var results []*domain.Balance
	ch := make(chan []*domain.Balance)

	if len(genesis.AppState.Accounts) > 0 {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			for v := range ch {
				results = append(results, v...)
			}
			wg.Done()
		}()
		wgBalances := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(genesis.AppState.Accounts)) / float64(chunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := chunkSize * i
			end := start + chunkSize
			if end > len(genesis.AppState.Accounts) {
				end = len(genesis.AppState.Accounts)
			}
			wgBalances.Add(1)
			go func() {
				var balances []*domain.Balance
				for _, account := range genesis.AppState.Accounts[start:end] {
					addressId, err := egu.addressRepository.FindId(helpers.RemovePrefix(account.Address))
					if err != nil {
						egu.logger.Error(err)
						continue
					}
					for _, bls := range account.Balance {
						balances = append(balances, &domain.Balance{
							CoinID:    uint64(bls.Coin),
							AddressID: addressId,
							Value:     bls.Value,
						})
					}
				}
				ch <- balances
				wgBalances.Done()
			}()
		}
		wgBalances.Wait()
		close(ch)
		wg.Wait()
	}
	return results, nil
}

func (egu *ExplorerGenesisUploader) saveBalances(balances []*domain.Balance) error {
	egu.logger.Info("Saving balances to DB...")
	if len(balances) > 0 {
		wgBalances := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(balances)) / float64(egu.env.BalanceChunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := int(egu.env.BalanceChunkSize) * i
			end := start + int(egu.env.BalanceChunkSize)
			if end > len(balances) {
				end = len(balances)
			}
			wgBalances.Add(1)
			go func() {
				err := egu.balanceRepository.SaveAll(balances[start:end])
				if err != nil {
					egu.logger.Error(err)
				}
				wgBalances.Done()
			}()
			wgBalances.Wait()
		}
	}
	return nil
}

func (egu *ExplorerGenesisUploader) extractStakes(genesis *domain.Genesis) ([]*domain.Stake, error) {
	var stakes []*domain.Stake
	for _, candidate := range genesis.AppState.Candidates {
		for _, stake := range candidate.Stakes {
			ownerId, err := egu.addressRepository.FindId(helpers.RemovePrefix(stake.Owner))
			if err != nil {
				egu.logger.Error(err)
			}
			validatorId, err := egu.validatorRepository.FindIdByPk(helpers.RemovePrefix(candidate.PublicKey))
			if err != nil {
				egu.logger.Error(err)
			}
			stakes = append(stakes, &domain.Stake{
				CoinID:         uint64(stake.Coin),
				OwnerAddressID: ownerId,
				ValidatorID:    validatorId,
				Value:          stake.Value,
				BipValue:       stake.BipValue,
			})
		}
	}
	return stakes, nil
}

func (egu *ExplorerGenesisUploader) saveStakes(stakes []*domain.Stake) error {
	egu.logger.Info("Saving stakes to DB...")
	if len(stakes) > 0 {
		wgStakes := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(stakes)) / float64(egu.env.StakeChunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := int(egu.env.StakeChunkSize) * i
			end := start + int(egu.env.StakeChunkSize)
			if end > len(stakes) {
				end = len(stakes)
			}
			wgStakes.Add(1)
			go func() {
				err := egu.validatorRepository.SaveAllStakes(stakes[start:end])
				if err != nil {
					egu.logger.Error(err)
				}
				wgStakes.Done()
			}()
			wgStakes.Wait()
		}
	}
	return nil
}

func (egu *ExplorerGenesisUploader) isEmptyDB() bool {
	addressesCount, err := egu.addressRepository.GetAddressesCount()
	if err != nil {
		panic(err)
	}
	balancesCount, err := egu.balanceRepository.GetBalancesCount()
	if err != nil {
		panic(err)
	}
	validatorCount, err := egu.validatorRepository.GetValidatorsCount()
	if err != nil {
		panic(err)
	}
	return addressesCount == 0 && balancesCount == 0 && validatorCount == 0
}

func (egu *ExplorerGenesisUploader) extractUnbonds(genesis *domain.Genesis) ([]*domain.Unbond, error) {
	var unbonds []*domain.Unbond
	for _, data := range genesis.AppState.FrozenFunds {
		addressId, err := egu.addressRepository.FindId(helpers.RemovePrefix(data.Address))
		if err != nil {
			egu.logger.WithField("address", data.Address).Error(err)
			continue
		}
		if data.CandidateKey == nil {
			continue
		}
		unbonds = append(unbonds, &domain.Unbond{
			AddressId:   uint(addressId),
			ValidatorId: uint(data.CandidateID),
			BlockId:     uint(data.Height),
			CoinId:      uint(data.Coin),
			Value:       data.Value,
		})
	}
	return unbonds, nil
}

func (egu *ExplorerGenesisUploader) saveUnbonds(unbonds []*domain.Unbond) error {
	egu.logger.Info("Saving unbonds to DB...")

	var list []*domain.Unbond
	for _, u := range unbonds {
		_, err := egu.validatorRepository.GetById(u.ValidatorId)
		if err == nil {
			list = append(list, u)
		}
	}

	if len(list) > 0 {
		wgStakes := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(list)) / float64(egu.env.StakeChunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := int(egu.env.StakeChunkSize) * i
			end := start + int(egu.env.StakeChunkSize)
			if end > len(list) {
				end = len(list)
			}
			wgStakes.Add(1)
			go func() {
				err := egu.validatorRepository.SaveAllUnbonds(list[start:end])
				if err != nil {
					egu.logger.Error(err)
				}
				wgStakes.Done()
			}()
			wgStakes.Wait()
		}
	}

	if len(unbonds) == len(list) {
		return errors.New("unbonds has not been saved")
	} else if len(unbonds)-len(list) > 0 {
		egu.logger.Warning(fmt.Sprintf("%d unbonds has been skiped", len(unbonds)-len(list)))
	}

	return nil
}

func (egu *ExplorerGenesisUploader) extractLiquidityPool(genesis *domain.Genesis) ([]*domain.LiquidityPool, error) {
	var list []*domain.LiquidityPool
	for _, data := range genesis.AppState.Pools {

		token, err := egu.coinRepository.FindBySymbol(fmt.Sprintf("LP-%d", data.ID))
		if err != nil {
			egu.logger.WithField("pool_id", data.ID).Error(err)
			continue
		}

		list = append(list, &domain.LiquidityPool{
			Id:               data.ID,
			TokenId:          uint64(token.ID),
			FirstCoinId:      data.Coin0,
			SecondCoinId:     data.Coin1,
			FirstCoinVolume:  data.Reserve0,
			SecondCoinVolume: data.Reserve1,
			Liquidity:        token.Volume,
			UpdatedAtBlockId: genesis.InitialHeight,
		})
	}
	return list, nil
}

func (egu *ExplorerGenesisUploader) saveLiquidityPool(pools []*domain.LiquidityPool) error {
	egu.logger.Info("Saving liquidity pool to DB...")
	if len(pools) > 0 {
		err := egu.liquidityPoolRepository.SaveAll(pools)
		if err != nil {
			egu.logger.Error(err)
		}
	}
	return nil
}

func (egu *ExplorerGenesisUploader) extractOrders(genesis *domain.Genesis) ([]domain.Order, error) {
	var list []domain.Order
	var orderMap sync.Map
	var wg sync.WaitGroup

	for _, pool := range genesis.AppState.Pools {
		wg.Add(len(pool.Orders))
		for _, o := range pool.Orders {
			go func(ord domain.GenesisOrder) {
				defer wg.Done()

				addressId, err := egu.addressRepository.FindId(helpers.RemovePrefix(ord.Owner))
				if err != nil {
					egu.logger.Error(err)
					return
				}

				order := domain.Order{
					Id:              ord.Id,
					AddressId:       addressId,
					CreatedAtBlock:  ord.Height,
					LiquidityPoolId: pool.ID,
					Status:          1,
				}

				if ord.IsSale {
					order.CoinSellId = pool.Coin0
					order.CoinSellVolume = ord.Volume0

					order.CoinBuyId = pool.Coin1
					order.CoinBuyVolume = ord.Volume1
				} else {
					order.CoinSellId = pool.Coin1
					order.CoinSellVolume = ord.Volume1

					order.CoinBuyId = pool.Coin0
					order.CoinBuyVolume = ord.Volume0
				}

				sell, ok := big.NewFloat(0).SetString(order.CoinSellVolume)
				if !ok {
					egu.logger.Error("can't convert to big.int")
					return
				}
				buy, ok := big.NewFloat(0).SetString(order.CoinBuyVolume)
				if !ok {
					egu.logger.Error("can't convert to big.int")
					return
				}
				price := sell.Quo(sell, buy)
				order.Price = price.Text('f', 18)
				orderMap.Store(ord.Id, order)
			}(o)
		}
	}
	wg.Wait()
	orderMap.Range(func(k, v interface{}) bool {
		list = append(list, v.(domain.Order))
		return true
	})

	return list, nil
}

func (egu *ExplorerGenesisUploader) saveOrders(orders []domain.Order) error {
	chunkSize := 1000
	egu.logger.Info("Saving orders to DB...")

	if len(orders) > 0 {
		wgStakes := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(orders)) / float64(chunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := chunkSize * i
			end := start + chunkSize
			if end > len(orders) {
				end = len(orders)
			}
			wgStakes.Add(1)
			go func() {
				err := egu.liquidityPoolRepository.SaveAllOrders(orders[start:end])
				if err != nil {
					egu.logger.Error(err)
				}
				wgStakes.Done()
			}()
			wgStakes.Wait()
		}
	}
	return nil
}
