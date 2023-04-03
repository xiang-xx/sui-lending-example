package main

import (
	"context"
	"fmt"

	"sui-lending-example/common"

	"github.com/coming-chat/go-sui/types"
	gosuilending "github.com/omnibtc/go-sui-lending"
)

func main() {
	acc := common.GetEnvAccount()
	client := common.GetDevClient()
	contract := common.GetDefaultContract()
	signer, err := types.NewHexData(acc.Address)
	common.AssertNil(err)
	ctx := context.Background()
	coins, err := client.GetSuiCoinsOwnedByAddress(ctx, *signer)
	common.AssertNil(err)
	gasCoin, err := coins.PickCoinNoLess(common.GasBudget)
	common.AssertNil(err)

	callOptions := gosuilending.CallOptions{
		Gas:       &gasCoin.Reference().ObjectId,
		GasBudget: common.GasBudget,
	}

	// for _, userId := range []string{"13", "17", "18"} {
	// 	f, err := contract.GetUserHealthFactor(ctx, *signer, userId, callOptions)
	// 	common.AssertNil(err)
	// 	d, _ := big.NewInt(0).SetString("1000000000000000000000000000", 10)
	// 	ff := decimal.NewFromBigInt(f, 0)
	// 	dd := decimal.NewFromBigInt(d, 0)
	// 	h := ff.Div(dd)
	// 	res, _ := h.Float64()
	// 	fmt.Printf("%s, %.5f\n", userId, res)
	// }
	// return

	dolaUserId, err := contract.GetDolaUserId(ctx, *signer, 0, signer.String(), callOptions)
	common.AssertNil(err)

	println(dolaUserId)

	liquid, err := contract.GetDolaTokenLiquidity(ctx, *signer, common.PoolIdUSDT, callOptions)
	common.AssertNil(err)
	println("dola USDT token liquidity:", liquid.String())

	liquid, err = contract.GetDolaTokenLiquidity(ctx, *signer, common.PoolIdBTC, callOptions)
	common.AssertNil(err)
	println("dola BTC token liquidity:", liquid.String())

	// appTokenLiquidity, err := contract.GetAppTokenLiquidity(ctx, *signer, 0, common.PoolIdUSDT, callOptions)
	// common.AssertNil(err)
	// println("app USDT token liquidity:", appTokenLiquidity.String())

	debtAmount, debtValue, err := contract.GetUserTokenDebt(ctx, *signer, dolaUserId, common.PoolIdBTC, callOptions)
	common.AssertNil(err)
	println("user btc token debt", debtAmount.String(), " ", debtValue.String())

	collateral, err := contract.GetUserCollateral(ctx, *signer, dolaUserId, common.PoolIdUSDT, callOptions)
	common.AssertNil(err)
	println("collateral: ", collateral.CollateralAmount.String(), " ", collateral.CollateralValue.String())

	userLendingInfo, err := contract.GetUserLendingInfo(ctx, *signer, dolaUserId, callOptions)
	common.AssertNil(err)
	fmt.Printf("lending info %v\n", userLendingInfo)

	reserveInfo, err := contract.GetReserveInfo(ctx, *signer, common.PoolIdUSDT, callOptions)
	common.AssertNil(err)
	fmt.Printf("%v\n", reserveInfo)

	// canBorrowAmount, err := contract.GetUserAllowedBorrow(ctx, *signer, usdtAddress, callOptions)
	// println("user can borrow usdt:", canBorrowAmount.String())
	// if err != nil {
	// 	println("reason:", err.Error())
	// }

	canBorrowAmount, err := contract.GetUserAllowedBorrow(ctx, *signer, dolaUserId, common.PoolIdBTC, callOptions)
	println("user can borrow btc:", canBorrowAmount.String())
	if err != nil {
		println("reason:", err.Error())
	}
}
