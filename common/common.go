package common

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/coming-chat/go-sui/account"
	"github.com/coming-chat/go-sui/client"
	gosuilending "github.com/omnibtc/go-sui-lending"
)

const DevnetRpcUrl = "https://fullnode.devnet.sui.io"
const (
	PackageLendingPortal      = "0x28bc45ed3593846a57fd3bea2839baa406f4d666"
	PackageExternalInterfaces = "0x20c2b9cb6d88de7dcf2b6ba98900058e1d58781c"
	PackageWormholeBridge     = "0x0979dcaa8d5549a4d6cc4783fe5ef093d5e32c35"
	PackageFaucet             = "0x4023210ff781042e398c5901f39434dcd290b420"
)

//go:embed sui.json
var suicontent []byte

var config *SuiConfig

type SuiConfig struct {
	USDC   string
	USDT   string
	APT    string
	BTC    string
	DAI    string
	ETH    string
	MATIC  string
	Faucet string // facute object id

	PoolAPT   string `json:"Pool<APT>"`
	PoolBTC   string `json:"Pool<BTC>"`
	PoolDAI   string `json:"Pool<DAI>"`
	PoolETH   string `json:"Pool<ETH>"`
	PoolMatic string `json:"Pool<MATIC>"`
	PoolUSDC  string `json:"Pool<USDC>"`
	PoolUSDT  string `json:"Pool<USDT>"`

	PoolManagerInfo string
	PoolState       string
	PriceOracle     string
	Storage         string
	UserManagerInfo string
	WormholeState   string
}

func init() {
	config = &SuiConfig{}
	err := json.Unmarshal(suicontent, &config)
	AssertNil(err)

	config.USDC = fixAddress(config.USDC)
	config.USDT = fixAddress(config.USDT)
	config.APT = fixAddress(config.APT)
	config.BTC = fixAddress(config.BTC)
	config.DAI = fixAddress(config.DAI)
	config.ETH = fixAddress(config.ETH)
	config.MATIC = fixAddress(config.MATIC)
	config.Faucet = fixAddress(config.Faucet)
}

func fixAddress(address string) string {
	return "0x" + strings.TrimPrefix(address, "0x")
}

func GetSuiConfig() *SuiConfig {
	return config
}

func GetDevClient() *client.Client {
	c, err := client.Dial(DevnetRpcUrl)
	AssertNil(err)
	return c
}

func GetEnvAccount() *account.Account {
	account, err := account.NewAccountWithMnemonic(os.Getenv("m"))
	AssertNil(err)
	return account
}

func AssertNil(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func GetDefaultContract() *gosuilending.Contract {
	contract, err := gosuilending.NewContract(GetDevClient(), gosuilending.ContractConfig{
		LendingPortalPackageId:     PackageLendingPortal,
		ExternalInterfacePackageId: PackageExternalInterfaces,
		PoolManagerInfo:            config.PoolManagerInfo,
		PoolState:                  config.PoolState,
		PriceOracle:                config.PriceOracle,
		Storage:                    config.Storage,
		WormholeState:              config.WormholeState,
	})
	AssertNil(err)
	return contract
}
