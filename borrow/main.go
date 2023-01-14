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
	btcPoolObject, err := types.NewHexData(config.PoolBTC)
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
	// addressBs, err := hex.DecodeString(acc.Address[2:])
	// common.AssertNil(err)

	btcLiquidity, err := contract.GetDolaTokenLiquidity(context.Background(), *signer, common.PoolIdBTC, gosuilending.CallOptions{
		Gas:       &gasCoin.Reference.ObjectId,
		GasBudget: 10000,
	})
	common.AssertNil(err)
	fmt.Printf("btc liquidity %s\n", btcLiquidity)

	tx, err := contract.Borrow(context.Background(), *signer, []string{
		config.BTC,
	}, gosuilending.BorrowArgs{
		Receiver:              acc.Address,
		WormholeMessageCoins:  []types.ObjectId{},
		WormholeMessageAmount: "0",
		Pool:                  *btcPoolObject,
		DstChain:              common.DolaChainIdSui,
		Amount:                "200",
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
