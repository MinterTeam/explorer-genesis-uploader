package core

import (
	"encoding/json"
	"fmt"
	"github.com/MinterTeam/explorer-genesis-uploader/domain"
	"github.com/MinterTeam/explorer-genesis-uploader/repository"
	"github.com/MinterTeam/minter-explorer-tools/v4/helpers"
	"github.com/go-pg/pg/v10"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

type ExplorerGenesisUploader struct {
	addressRepository   *repository.Address
	balanceRepository   *repository.Balance
	coinRepository      *repository.Coin
	validatorRepository *repository.Validator
	logger              *logrus.Entry
}

func New() *ExplorerGenesisUploader {
	//Init Logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(false)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	contextLogger := logger.WithFields(logrus.Fields{
		"version": "1.2",
		"app":     "Minter Explorer Explorer Genesis Uploader",
	})

	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	// Repositories
	addressRepository := repository.NewAddressRepository(db)
	coinRepository := repository.NewCoinRepository(db)
	validatorRepository := repository.NewValidatorRepository(db)
	balanceRepository := repository.NewBalanceRepository(db)

	return &ExplorerGenesisUploader{
		addressRepository:   addressRepository,
		balanceRepository:   balanceRepository,
		coinRepository:      coinRepository,
		validatorRepository: validatorRepository,
		logger:              contextLogger,
	}
}

func (egu *ExplorerGenesisUploader) Do() error {

	if !egu.isEmptyDB() {
		egu.logger.Info("DB is not empty")
		os.Exit(0)
	}

	start := time.Now()
	egu.logger.Info("Getting genesis data...")

	// Create a Resty Client
	client := resty.New().SetTimeout(time.Minute).SetHostURL(fmt.Sprintf("%s", os.Getenv("NODE_API")))
	data, err := client.R().Get("/genesis")
	helpers.HandleError(err)

	genesisResponse := new(GenesisResponse)
	err = json.Unmarshal(data.Body(), genesisResponse)
	helpers.HandleError(err)
	egu.logger.Info("Genesis data has been downloaded")

	genesis := &genesisResponse.Result.Genesis

	egu.logger.Info("Extracting addresses...")
	addresses, err := egu.extractAddresses(genesis)
	helpers.HandleError(err)
	egu.logger.Info(fmt.Sprintf("%d addresses has been extracted", len(addresses)))

	egu.saveAddresses(addresses)

	egu.logger.Info("Extracting coins...")
	coins, err := egu.extractCoins(genesis)
	helpers.HandleError(err)
	egu.logger.Info(fmt.Sprintf("%d coins has been extracted", len(coins)+1))

	egu.saveCoins(coins)

	egu.logger.Info("Extracting validators...")
	validators, err := egu.extractCandidates(genesis)
	helpers.HandleError(err)
	egu.logger.Info(fmt.Sprintf("%d validators have been extracted", len(validators)))
	err = egu.saveCandidates(validators)
	helpers.HandleError(err)
	egu.logger.Info("Validators has been uploaded")

	egu.logger.Info("Extracting balances...")
	balances, err := egu.extractBalances(genesis)
	if err != nil {
		helpers.HandleError(err)
	}
	helpers.HandleError(err)
	egu.logger.Info(fmt.Sprintf("%d balances has been extracted", len(balances)))
	err = egu.saveBalances(balances)
	helpers.HandleError(err)
	egu.logger.Info("Balances has been uploaded")

	egu.logger.Info("Extracting stakes...")
	stakes, err := egu.extractStakes(genesis)
	helpers.HandleError(err)
	egu.logger.Info(fmt.Sprintf("%d stakes have been extracted", len(stakes)))
	err = egu.saveStakes(stakes)
	helpers.HandleError(err)
	egu.logger.Info("Stakes has been uploaded")
	egu.logger.Info("Upload complete")

	elapsed := time.Since(start)
	egu.logger.Info("Processing time: ", elapsed)
	return err
}

func (egu *ExplorerGenesisUploader) extractAddresses(genesis *Genesis) ([]string, error) {
	addressesMap := make(map[string]struct{})
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
		if coin.OwnerAddress != "" {
			addressesMap[helpers.RemovePrefix(coin.OwnerAddress)] = struct{}{}
		}
	}

	var addresses = make([]string, len(addressesMap))
	i := 0
	for adr := range addressesMap {
		addresses[i] = adr
		i++
	}
	return addresses, nil
}

func (egu *ExplorerGenesisUploader) extractCoins(genesis *Genesis) ([]*domain.Coin, error) {
	var coins = make([]*domain.Coin, len(genesis.AppState.Coins)+1)
	i := 1

	coins[0] = &domain.Coin{
		ID:        0,
		Name:      "Base coin",
		Symbol:    os.Getenv("MINTER_BASE_COIN"),
		Volume:    "0",
		Crr:       100,
		Reserve:   "0",
		MaxSupply: "0",
		Version:   0,
	}

	for _, c := range genesis.AppState.Coins {
		crr, err := strconv.ParseUint(c.Crr, 10, 64)
		if err != nil {
			egu.logger.Error(err)
		}

		if c.ID == 0 {
			egu.logger.Error(err)
		}

		coins[i] = &domain.Coin{
			ID:        c.ID,
			Name:      c.Name,
			Symbol:    c.Symbol,
			Volume:    c.Volume,
			Crr:       uint(crr),
			Reserve:   c.Reserve,
			MaxSupply: c.MaxSupply,
			Version:   c.Version,
		}
		if c.OwnerAddress != "" {
			addressId, err := egu.addressRepository.FindId(helpers.RemovePrefix(c.OwnerAddress))
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

func (egu ExplorerGenesisUploader) extractCandidates(genesis *Genesis) ([]*domain.Validator, error) {
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

		status := candidate.Status
		commission, err := strconv.ParseUint(candidate.Commission, 10, 64)
		stake := candidate.TotalBipStake
		validator, err := egu.validatorRepository.Add(&domain.Validator{
			PublicKey:       helpers.RemovePrefix(candidate.PublicKey),
			OwnerAddressID:  &ownerAddress,
			RewardAddressID: &rewardAddress,
			Status:          &status,
			Commission:      &commission,
			TotalStake:      &stake,
		})
		if err != nil {
			egu.logger.Panic(err)
		}

		_, err = egu.validatorRepository.AddPk(helpers.RemovePrefix(candidate.PublicKey), validator.ID)
		if err != nil {
			egu.logger.Panic(err)
		}

	}

	return validators, nil
}

func (egu *ExplorerGenesisUploader) saveAddresses(addresses []string) {
	AddrChunkSize := os.Getenv("EGU_ADDRESS_CHUNK_SIZE")
	chunkSize, err := strconv.ParseInt(AddrChunkSize, 10, 64)
	helpers.HandleError(err)
	egu.logger.Info("Saving addresses to DB...")
	if len(addresses) > 0 {
		wgAddresses := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(addresses)) / float64(chunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := int(chunkSize) * i
			end := start + int(chunkSize)
			if end > len(addresses) {
				end = len(addresses)
			}
			wgAddresses.Add(1)
			go func() {
				err := egu.addressRepository.SaveAll(addresses[start:end])
				helpers.HandleError(err)
				wgAddresses.Done()
			}()
		}
		wgAddresses.Wait()
	}
	egu.logger.Info("Addresses has been uploaded")
}

func (egu *ExplorerGenesisUploader) saveCoins(coins []*domain.Coin) {
	egu.logger.Info("Saving coins to DB...")
	coinsChunkSize := os.Getenv("EGU_COINS_CHUNK_SIZE")
	chunkSize, err := strconv.ParseInt(coinsChunkSize, 10, 64)
	helpers.HandleError(err)
	var list []*domain.Coin
	list = append(list, coins...)
	wgCoins := new(sync.WaitGroup)
	chunksCount := int(math.Ceil(float64(len(list)) / float64(chunkSize)))
	for i := 0; i < chunksCount; i++ {
		start := int(chunkSize) * i
		end := start + int(chunkSize)
		if end > len(list) {
			end = len(list)
		}
		wgCoins.Add(1)
		go func() {
			err := egu.coinRepository.SaveAll(list[start:end])
			helpers.HandleError(err)
			wgCoins.Done()
		}()
	}
	wgCoins.Wait()
	egu.logger.Info("Coins has been uploaded")
}

func (egu *ExplorerGenesisUploader) saveCandidates(validators []*domain.Validator) error {
	egu.logger.Info("Saving validators to DB...")
	validatorsChunkSize := os.Getenv("EGU_VALIDATORS_CHUNK_SIZE")
	chunkSize, err := strconv.ParseInt(validatorsChunkSize, 10, 64)
	helpers.HandleError(err)

	if len(validators) > 0 {
		wgCandidates := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(validators)) / float64(chunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := int(chunkSize) * i
			end := start + int(chunkSize)
			if end > len(validators) {
				end = len(validators)
			}
			wgCandidates.Add(1)
			go func() {
				err := egu.validatorRepository.SaveAll(validators[start:end])
				helpers.HandleError(err)
				wgCandidates.Done()
			}()
		}
		wgCandidates.Wait()
	}
	return nil
}

func (egu *ExplorerGenesisUploader) extractBalances(genesis *Genesis) ([]*domain.Balance, error) {
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
	balancesChunkSize := os.Getenv("EGU_BALANCES_CHUNK_SIZE")
	chunkSize, err := strconv.ParseInt(balancesChunkSize, 10, 64)
	helpers.HandleError(err)
	if len(balances) > 0 {
		wgBalances := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(balances)) / float64(chunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := int(chunkSize) * i
			end := start + int(chunkSize)
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

func (egu *ExplorerGenesisUploader) extractStakes(genesis *Genesis) ([]*domain.Stake, error) {
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
	stakesChunkSize := os.Getenv("EGU_STAKE_CHUNK_SIZE")
	chunkSize, err := strconv.ParseInt(stakesChunkSize, 10, 64)
	helpers.HandleError(err)

	if len(stakes) > 0 {
		wgStakes := new(sync.WaitGroup)
		chunksCount := int(math.Ceil(float64(len(stakes)) / float64(chunkSize)))
		for i := 0; i < chunksCount; i++ {
			start := int(chunkSize) * i
			end := start + int(chunkSize)
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
