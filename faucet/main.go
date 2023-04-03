package main

import (
	"context"
	"fmt"

	"sui-lending-example/common"

	"github.com/coming-chat/go-sui/sui_types"
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
	gasCoin, err := coins.PickCoinNoLess(common.GasBudget)
	common.AssertNil(err)
	tx, err := contract.Claim(context.Background(), *signer, []string{
		config.USDT,
	}, gosuilending.CallOptions{
		Gas:       &gasCoin.Reference().ObjectId,
		GasBudget: common.GasBudget,
	})
	common.AssertNil(err)

	signature, err := acc.SignSecureWithoutEncode(tx.TxBytes, sui_types.DefaultIntent())
	common.AssertNil(err)
	options := types.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
	}
	resp, err := client.ExecuteTransactionBlock(
		context.TODO(), tx.TxBytes, []any{signature}, &options,
		types.TxnRequestTypeWaitForLocalExecution,
	)
	common.AssertNil(err)
	fmt.Println(resp.Effects.TransactionDigest)
}
