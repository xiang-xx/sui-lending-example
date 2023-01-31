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
const (
	PackageLendingPortal      = "0x725996f982c461ddb1060cb64a7b47246e7332be"
	PackageExternalInterfaces = "0xaa65494974bfa11425bfbff836db69cf7950f3ef"
	PackageWormholeBridge     = "0x10ec199c006b64d40511ca7f2f0527051577d23f"
	PackageFaucet             = "0x72b846eca3c7f91961ec3cae20441be96a21e1fe"
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
