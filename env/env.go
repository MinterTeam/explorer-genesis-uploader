package env

type Config struct {
	Debug              bool
	PostgresHost       string
	PostgresPort       string
	PostgresDB         string
	PostgresUser       string
	PostgresPassword   string
	PostgresSSLEnabled bool
	MinterBaseCoin     string
	NodeGrpc           string
	AddressChunkSize   uint64
	CoinsChunkSize     uint64
	BalanceChunkSize   uint64
	StakeChunkSize     uint64
	ValidatorChunkSize uint64
}
