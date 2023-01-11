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
	acc := common.GetEnvAccount()
	client := common.GetDevClient()
	contract, err := gosuilending.NewFaucet(client, common.PackageFaucet, config.Faucet)
	common.AssertNil(err)
	signer, err := types.NewHexData(acc.Address)
	common.AssertNil(err)
	ctx := context.Background()
	coins, err := client.GetSuiCoinsOwnedByAddress(ctx, *signer)
	common.AssertNil(err)
	gasCoin, err := coins.PickCoinNoLess(1000)
	common.AssertNil(err)
	tx, err := contract.Claim(context.Background(), *signer, []string{
		config.USDT,
	}, gosuilending.CallOptions{
		Gas:       &gasCoin.Reference.ObjectId,
		GasBudget: 1000,
	})
	common.AssertNil(err)

	signedTx := tx.SignWith(acc.PrivateKey)
	resp, err := client.ExecuteTransaction(ctx, *signedTx, types.TxnRequestTypeWaitForLocalExecution)
	common.AssertNil(err)
	fmt.Println(resp.EffectsCert.Certificate.TransactionDigest)
}
