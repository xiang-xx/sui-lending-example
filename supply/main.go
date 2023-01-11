package main

import (
	"context"
	"fmt"

	"sui-lending-example/common"

	"github.com/coming-chat/go-sui/types"
	gosuilending "github.com/omnibtc/go-sui-lending"
)

const (
	usdtAddress = "0x13e8531463853d9a3ff017d140be14a9357f6b1d::coins::USDT"
	usdtPool    = "0xf4bc9117ff693bd9086ebdb28aea09c1c7256d9a"
)

func main() {
	usdtPoolObject, err := types.NewHexData(usdtPool)
	common.PanicIfError(err)
	acc := common.GetEnvAccount()
	client := common.GetDevClient()
	contract := common.GetDefaultContract()
	common.PanicIfError(err)
	signer, err := types.NewHexData(acc.Address)
	common.PanicIfError(err)
	ctx := context.Background()
	coins, err := client.GetSuiCoinsOwnedByAddress(ctx, *signer)
	common.PanicIfError(err)
	gasCoin, err := coins.PickCoinNoLess(10000)
	common.PanicIfError(err)

	// get usdt coins
	usdtCoins, err := client.GetCoinsOwnedByAddress(ctx, *signer, usdtAddress)
	common.PanicIfError(err)
	usdtCoinObjectIds := []types.ObjectId{}
	for _, coin := range usdtCoins {
		usdtCoinObjectIds = append(usdtCoinObjectIds, coin.Reference.ObjectId)
	}

	tx, err := contract.Supply(context.Background(), *signer, []string{
		usdtAddress,
	}, gosuilending.SupplyArgs{
		WormholeMessageCoins:  []types.ObjectId{},
		WormholeMessageAmount: 0,
		Pool:                  *usdtPoolObject,
		DepositCoins:          usdtCoinObjectIds,
		DepositAmount:         50000000,
	}, gosuilending.CallOptions{
		Gas:       &gasCoin.Reference.ObjectId,
		GasBudget: 10000,
	})
	common.PanicIfError(err)

	signedTx := tx.SignWith(acc.PrivateKey)
	resp, err := client.ExecuteTransaction(ctx, *signedTx, types.TxnRequestTypeWaitForLocalExecution)
	common.PanicIfError(err)
	fmt.Println(resp.EffectsCert.Certificate.TransactionDigest)
}