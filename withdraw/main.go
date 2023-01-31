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

	tx, err := contract.Withdraw(context.Background(), *signer, []string{
		usdtAddress,
	}, gosuilending.WithdrawArgs{
		WormholeMessageCoins:  []types.ObjectId{},
		WormholeMessageAmount: "0",
		Pool:                  *usdtPoolObject,
		DstChain:              "1",
		Amount:                "49999999",
	}, gosuilending.CallOptions{
		Gas:       &gasCoin.Reference.ObjectId,
		GasBudget: 10000,
	})
	common.AssertNil(err)

	signedTx := tx.SignSerializedSigWith(acc.PrivateKey)
	resp, err := client.ExecuteTransaction(ctx, *signedTx, types.TxnRequestTypeWaitForLocalExecution)
	common.AssertNil(err)
	fmt.Println(resp.EffectsCert.Certificate.TransactionDigest)
}
