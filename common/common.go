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

const DevnetRpcUrl = "https://fullnode.testnet.sui.io"
const GasBudget = 10_000_000

const (
	PackageLendingPortal      = "0x5b81c31943358fcf8f20d2c9b92adf6b47062aa4b01afb2e5c901920807c3e09"
	PackageExternalInterfaces = "0x25bf584ec396b75ee3ab0367a5b6ebdb82f8fd8bf42f5bd00c281464d604a994"
	PackageWormholeBridge     = "0xe198cbf3b61678ba33be2a53965c4d68a2b55d00aea67af9038f54c4dba1ec61"
	// test_coins
	PackageFaucet = "0x54fc06a12aeed0752c6db5d949fcf4554bd320ca69676ee9d3085ba946b91af0"
)

const (
	DolaChainIdSui       = "0"
	DolaChainIdAptos     = "1"
	DolaChainIdEvm       = "2"
	DolaChainIdPolygon   = "5"
	DolaChainIdPolygonZk = "1422"
)

const (
	PoolIdBTC = iota
	PoolIdUSDT
	PoolIdUSDC
	PoolIdETH
	PoolIdDAI
	PoolIdMATIC
	PoolIdAPT // 6
	PoolIdBNB
)

//go:embed sui.json
var suicontent []byte

var config *SuiConfig

type SuiConfig struct {
	USDC   string
	USDT   string
	APT    string
	BNB    string
	BTC    string
	DAI    string
	ETH    string
	MATIC  string
	Faucet string // facute object id

	PoolAPT   string `json:"Pool<APT>"`
	PoolBNB   string `json:"Pool<BNB>"`
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
	CoreState       string
	LendingPortal   string
	Clock           string
	PoolApproval    string
}

func init() {
	config = &SuiConfig{}
	err := json.Unmarshal(suicontent, &config)
	AssertNil(err)

	config.USDC = fixAddress(config.USDC)
	config.USDT = fixAddress(config.USDT)
	config.APT = fixAddress(config.APT)
	config.BNB = fixAddress(config.BNB)
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
		UserManagerInfo:            config.UserManagerInfo,
		CoreState:                  config.CoreState,
		LendingPortal:              config.LendingPortal,
		Clock:                      config.Clock,
		PoolApproval:               config.PoolApproval,
	})
	AssertNil(err)
	return contract
}
