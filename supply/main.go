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
	gasCoin, err := coins.PickCoinNoLess(common.GasBudget)
	common.AssertNil(err)

	// get usdt coins
	usdtCoins, err := client.GetCoins(ctx, *signer, &config.USDT, nil, 100)
	common.AssertNil(err)
	usdtCoinObjectIds := []types.ObjectId{}
	for _, coin := range usdtCoins.Data {
		usdtCoinObjectIds = append(usdtCoinObjectIds, coin.Reference().ObjectId)
	}

	tx, err := contract.Supply(context.Background(), *signer, []string{
		config.USDT,
	}, gosuilending.SupplyArgs{
		Pool:          *usdtPoolObject,
		DepositCoins:  usdtCoinObjectIds,
		DepositAmount: "50000000",
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
