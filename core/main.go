package core

import (
	"encoding/json"
	"fmt"
	"github.com/MinterTeam/minter-explorer-extender/address"
	"github.com/MinterTeam/minter-explorer-extender/balance"
	"github.com/MinterTeam/minter-explorer-extender/block"
	"github.com/MinterTeam/minter-explorer-extender/coin"
	"github.com/MinterTeam/minter-explorer-extender/validator"
	"github.com/MinterTeam/minter-explorer-tools/helpers"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

type ExplorerGenesisUploader struct {
	db                  *pg.DB
	env                 *models.ExtenderEnvironment
	addressRepository   *address.Repository
	balanceRepository   *balance.Repository
	blockRepository     *block.Repository
	coinRepository      *coin.Repository
	validatorRepository *validator.Repository
	logger              *logrus.Entry
}

func New(env *models.ExtenderEnvironment) *ExplorerGenesisUploader {
	//Init DB
	db := pg.Connect(&pg.Options{
		User:            env.DbUser,
		Password:        env.DbPassword,
		Database:        env.DbName,
		PoolSize:        20,
		MinIdleConns:    10,
		ApplicationName: env.AppName,
	})

	//Init Logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(false)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	contextLogger := logger.WithFields(logrus.Fields{
		"version": "1.0",
		"app":     "Minter Explorer Explorer Genesis Uploader",
	})

	// Repositories
	addressRepository := address.NewRepository(db)
	blockRepository := block.NewRepository(db)
	coinRepository := coin.NewRepository(db)
	validatorRepository := validator.NewRepository(db)
	balanceRepository := balance.NewRepository(db)

	return &ExplorerGenesisUploader{
		db:                  db,
		env:                 env,
		addressRepository:   addressRepository,
		balanceRepository:   balanceRepository,
		blockRepository:     blockRepository,
		coinRepository:      coinRepository,
		validatorRepository: validatorRepository,
		logger:              contextLogger,
	}
}

func (egu *ExplorerGenesisUploader) Do() error {

	egu.logger.Info("Getting genesis data...")
	url := egu.env.NodeApi + "/genesis"
	_, data, err := fasthttp.Get(nil, url)
	helpers.HandleError(err)
	genesisResponse := new(GenesisResponse)
	err = json.Unmarshal(data, genesisResponse)
	helpers.HandleError(err)
	egu.logger.Info("Genesis data has been downloading")

	genesis := &genesisResponse.Result.Genesis

	egu.logger.Info("Extracting addresses...")
	addresses, err := egu.extractAddresses(genesis)
	helpers.HandleError(err)
	msg := fmt.Sprintf("%d addresses have been extracting", len(addresses))
	egu.logger.Info(msg)

	egu.logger.Info("Extracting coins...")
	coins, err := egu.extractCoins(genesis)
	helpers.HandleError(err)
	msg = fmt.Sprintf("%d coins has been extracting", len(coins))
	egu.logger.Info(msg)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go egu.saveAddresses(addresses, wg)
	go egu.saveCoins(coins, wg)
	wg.Wait()

	egu.logger.Info("Extracting validators...")
	validators, err := egu.extractCandidates(genesis)
	helpers.HandleError(err)
	msg = fmt.Sprintf("%d validators have been extracting", len(validators))
	egu.logger.Info(msg)
	err = egu.saveCandidates(validators)
	helpers.HandleError(err)
	egu.logger.Info("Validators has been uploading")

	egu.logger.Info("Extracting balances...")
	balances, err := egu.extractBalances(genesis)
	helpers.HandleError(err)
	msg = fmt.Sprintf("%d balances has been extracting", len(balances))
	egu.logger.Info(msg)
	err = egu.saveBalances(balances)
	helpers.HandleError(err)
	egu.logger.Info("Balances has been uploading")

	egu.logger.Info("Extracting stakes...")
	stakes, err := egu.extractStakes(genesis)
	helpers.HandleError(err)
	msg = fmt.Sprintf("%d stakes have been extracting", len(stakes))
	egu.logger.Info(msg)
	err = egu.saveStakes(stakes)
	helpers.HandleError(err)
	egu.logger.Info("Stakes has been uploading")

	egu.logger.Info("Upload complete")
	return err
}

func (egu *ExplorerGenesisUploader) extractAddresses(genesis *Genesis) ([]string, error) {
	addressesMap := make(map[string]struct{})
	for _, val := range genesis.AppState.Validators {
		addressesMap[helpers.RemovePrefix(val.RewardAddress)] = struct{}{}
	}
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
	var addresses = make([]string, len(addressesMap))
	i := 0
	for adr := range addressesMap {
		addresses[i] = adr
		i++
	}
	return addresses, nil
}

func (egu *ExplorerGenesisUploader) extractCoins(genesis *Genesis) ([]*models.Coin, error) {
	var coins = make([]*models.Coin, len(genesis.AppState.Coins))
	i := 0
	for _, c := range genesis.AppState.Coins {
		crr, err := strconv.ParseUint(c.Crr, 10, 64)
		if err != nil {
			egu.logger.Error(err)
		}
		coins[i] = &models.Coin{
			Name:           c.Name,
			Symbol:         c.Symbol,
			Crr:            crr,
			Volume:         c.Volume,
			ReserveBalance: c.ReserveBalance,
			UpdatedAt:      time.Now(),
		}
		i++
	}
	return coins, nil
}

func (egu *ExplorerGenesisUploader) extractStakes(genesis *Genesis) ([]*models.Stake, error) {
	var stakes []*models.Stake
	for _, candidate := range genesis.AppState.Candidates {
		for _, stake := range candidate.Stakes {
			coinId, err := egu.coinRepository.FindIdBySymbol(stake.Coin)
			if err != nil {
				egu.logger.Error(err)
			}
			ownerId, err := egu.addressRepository.FindId(helpers.RemovePrefix(stake.Owner))
			if err != nil {
				egu.logger.Error(err)
			}
			validatorId, err := egu.validatorRepository.FindIdByPk(helpers.RemovePrefix(candidate.PubKey))
			if err != nil {
				egu.logger.Error(err)
			}
			stakes = append(stakes, &models.Stake{
				CoinID:         coinId,
				OwnerAddressID: ownerId,
				ValidatorID:    validatorId,
				Value:          stake.Value,
				BipValue:       stake.BipValue,
			})
		}
	}
	return stakes, nil
}

func (egu *ExplorerGenesisUploader) extractBalances(genesis *Genesis) ([]*models.Balance, error) {
	var balances []*models.Balance
	for _, account := range genesis.AppState.Accounts {
		addressId, err := egu.addressRepository.FindId(helpers.RemovePrefix(account.Address))
		if err != nil {
			egu.logger.Error(err)
		}
		for _, bls := range account.Balance {
			coinId, err := egu.coinRepository.FindIdBySymbol(bls.Coin)
			if err != nil {
				egu.logger.Error(err)
			}
			balances = append(balances, &models.Balance{
				CoinID:    coinId,
				AddressID: addressId,
				Value:     bls.Value,
			})
		}
	}
	return balances, nil
}

func (egu ExplorerGenesisUploader) extractCandidates(genesis *Genesis) ([]*models.Validator, error) {
	var validators []*models.Validator
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
		commission, err := strconv.ParseUint(candidate.Commission, 10, 64)
		stake := candidate.TotalBipStake
		validators = append(append(validators, &models.Validator{
			OwnerAddressID:  &ownerAddress,
			RewardAddressID: &rewardAddress,
			PublicKey:       helpers.RemovePrefix(candidate.PubKey),
			Status:          &status,
			Commission:      &commission,
			TotalStake:      &stake,
		}))
	}

	return validators, nil
}

func (egu *ExplorerGenesisUploader) saveAddresses(addresses []string, wg *sync.WaitGroup) {
	if len(addresses) > 0 {
		wgAddresses := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(addresses)) / float64(egu.env.AddrChunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := egu.env.AddrChunkSize * i
			end := start + egu.env.AddrChunkSize
			if end > len(addresses) {
				end = len(addresses)
			}
			wgAddresses.Add(1)
			go func() {
				err := egu.addressRepository.SaveAllIfNotExist(addresses[start:end])
				helpers.HandleError(err)
				wgAddresses.Done()
			}()
		}
		wgAddresses.Wait()
	}
	wg.Done()
	egu.logger.Info("Addresses has been uploading")
}

func (egu *ExplorerGenesisUploader) saveCoins(coins []*models.Coin, wg *sync.WaitGroup) {
	if len(coins) > 0 {
		wgCoins := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(coins)) / float64(egu.env.TxChunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := egu.env.TxChunkSize * i
			end := start + egu.env.TxChunkSize
			if end > len(coins) {
				end = len(coins)
			}
			wgCoins.Add(1)
			go func() {
				err := egu.coinRepository.SaveAllIfNotExist(coins[start:end])
				helpers.HandleError(err)
				wgCoins.Done()
			}()
		}
		wgCoins.Wait()
	}
	wg.Done()
	egu.logger.Info("Coins has been uploading")
}

func (egu *ExplorerGenesisUploader) saveCandidates(validators []*models.Validator) error {
	if len(validators) > 0 {
		wgCandidates := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(validators)) / float64(egu.env.TxChunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := egu.env.TxChunkSize * i
			end := start + egu.env.TxChunkSize
			if end > len(validators) {
				end = len(validators)
			}
			wgCandidates.Add(1)
			go func() {
				err := egu.validatorRepository.SaveAllIfNotExist(validators[start:end])
				helpers.HandleError(err)
				wgCandidates.Done()
			}()
		}
		wgCandidates.Wait()
	}
	return nil
}

func (egu *ExplorerGenesisUploader) saveBalances(balances []*models.Balance) error {
	if len(balances) > 0 {
		wgBalances := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(balances)) / float64(egu.env.TxChunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := egu.env.TxChunkSize * i
			end := start + egu.env.TxChunkSize
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

func (egu *ExplorerGenesisUploader) saveStakes(stakes []*models.Stake) error {
	if len(stakes) > 0 {
		wgStakes := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(stakes)) / float64(egu.env.TxChunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := egu.env.TxChunkSize * i
			end := start + egu.env.TxChunkSize
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
