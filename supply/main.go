package main

import (
	"context"
	"fmt"

	"sui-lending-example/common"

	"github.com/coming-chat/go-sui/types"
	gosuilending "github.com/omnibtc/go-sui-lending"
)

func main() {
	config := common.GetSuiConfig()
	usdtPoolObject, err := types.NewHexData(config.PoolUSDT)
	common.AssertNil(err)
	acc := common.GetEnvAccount()
	client := common.GetDevClient()
	contract := common.GetDefaultContract()
	common.AssertNil(err)
	signer, err := types.NewHexData(acc.Address)
	common.AssertNil(err)
	ctx := context.Background()
	coins, err := client.GetSuiCoinsOwnedByAddress(ctx, *signer)
	common.AssertNil(err)
	gasCoin, err := coins.PickCoinNoLess(10000)
	common.AssertNil(err)

	// get usdt coins
	usdtCoins, err := client.GetCoinsOwnedByAddress(ctx, *signer, config.USDT)
	common.AssertNil(err)
	usdtCoinObjectIds := []types.ObjectId{}
	for _, coin := range usdtCoins {
		usdtCoinObjectIds = append(usdtCoinObjectIds, coin.Reference.ObjectId)
	}

	tx, err := contract.Supply(context.Background(), *signer, []string{
		config.USDT,
	}, gosuilending.SupplyArgs{
		WormholeMessageCoins:  []types.ObjectId{},
		WormholeMessageAmount: "0",
		Pool:                  *usdtPoolObject,
		DepositCoins:          usdtCoinObjectIds,
		DepositAmount:         "50000000",
	}, gosuilending.CallOptions{
		Gas:       &gasCoin.Reference.ObjectId,
		GasBudget: 10000,
	})
	common.AssertNil(err)

	signedTx := tx.SignWith(acc.PrivateKey)
	resp, err := client.ExecuteTransaction(ctx, *signedTx, types.TxnRequestTypeWaitForLocalExecution)
	common.AssertNil(err)
	fmt.Println(resp.EffectsCert.Certificate.TransactionDigest)
}
