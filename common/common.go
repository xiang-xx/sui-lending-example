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
	PackageLendingPortal      = "0x0ea2ebec77a430932699c598d286cdd72cff25c3"
	PackageExternalInterfaces = "0xd302b6887e60b2018d046061f01d0d1241ad9455"
	PackageWormholeBridge     = "0x9286925aa8f12b71ca9398133dd0becedde9ee4d"
	PackageFaucet             = "0x7b03c8e330e491c2a1bad31b43b781b9d7d3114c"
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
	})
	AssertNil(err)
	return contract
}
