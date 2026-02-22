package arb

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

func TestCalcAmountInV2ToV30To1(t *testing.T) {
	logger.Init("", "", false)
	decimal.DivisionPrecision = 40

	fee1 := uint32(500)
	fee2 := uint32(500)
	x1 := convertToBigIntMust("1000")
	y1 := convertToBigIntMust("2508000")

	price := float64(2499.5)
	pl := createPool(price, fee2)
	mint(pl, 2500, 2502, convertToBigIntMust("10"), convertToBigIntMust("2500"))
	mint(pl, 2499, 2500, convertToBigIntMust("1"), convertToBigIntMust("250"))
	mint(pl, 2495, 2499, convertToBigIntMust("100"), convertToBigIntMust("250000"))
	mint(pl, 2490, 2496, convertToBigIntMust("100"), convertToBigIntMust("250000"))
	pl.Reload()

	bounds := BuildTickLiquidity(pl, false, 5)
	printBoundedTicks(bounds)

	amt, amtMid, amtOut, profit, profitable := CalcAmountInV2ToV3(true, nil, nil, x1, y1, bounds, fee1, fee2, true)
	t.Logf("amountIn: %v amtMid: %v amtOut: %v profit: %v profitable: %v", amt, amtMid, amtOut, profit, profitable)

	// amtOld, _ := bestAmountV2ToV3Tick(x1, y1, bounds[0].Liquidity, bounds[0].SPCurr, big.NewInt(0), uint(fee1), uint(fee2), true)
	// t.Logf("amountIn 2: %v ", amtOld)

	// for i := -5; i <= 10; i++ {
	// nAmt := new(big.Int).Div(mulmul(amt, big.NewInt(int64(10+i))), big.NewInt(10))
	logger.Info().Msgf("amtIn: %v profit: %v", amt, calcProfitV2ToV3(true, x1, y1, pl, amt, fee1))
	printPriceAfterSwapV2(true, x1, y1, amt, fee1)
	amountOut := getAmountOutV2(x1, y1, amt, new(big.Int).Sub(e6, big.NewInt(int64(fee1))))
	printPriceAfterSwapV3(false, pl, amountOut)
	// }
}

func printPriceAfterSwapV2(zeroForOne bool, x1, y1, amtIn *big.Int, fee uint32) {
	var x2, y2 *big.Int

	if zeroForOne {
		amountOut := getAmountOutV2(x1, y1, amtIn, new(big.Int).Sub(e6, big.NewInt(int64(fee))))
		x2 = new(big.Int).Add(x1, amtIn)
		y2 = new(big.Int).Sub(y1, amountOut)
	} else {
		amountOut := getAmountOutV2(y1, x1, amtIn, new(big.Int).Sub(e6, big.NewInt(int64(fee))))
		x2 = new(big.Int).Sub(x1, amountOut)
		y2 = new(big.Int).Add(y1, amtIn)
	}

	price := decimal.NewFromBigInt(y2, 0).Div(decimal.NewFromBigInt(x2, 0))
	logger.Info().Msgf("price after swap: %v zeroForOne: %v amountIn: %v", price, zeroForOne, amtIn)
}

func printPriceAfterSwapV3(zeroForOne bool, pl *pool.Pool, amtIn *big.Int) {
	var spLimit *big.Int
	if zeroForOne {
		spLimit = new(big.Int).Add(pool.MIN_SQRT_RATIO, big.NewInt(1))
	} else {
		spLimit = new(big.Int).Sub(pool.MAX_SQRT_RATIO, big.NewInt(1))
	}
	amt0, amt1, spNext, _ := pl.MockSwap(zeroForOne, amtIn, spLimit)

	logger.Info().Msgf("price after swap: %v zeroForOne: %v amountIn: %v amt0: %v amt1: %v", SqrtPriceX96ToPrice(spNext), zeroForOne, amtIn, amt0, amt1)
}

func calcProfitV2ToV3(zeroForOne1 bool, rx, ry *big.Int, pl *pool.Pool, amtIn *big.Int, fee1 uint32) *big.Int {
	if zeroForOne1 {
		amountOut := getAmountOutV2(rx, ry, amtIn, new(big.Int).Sub(e6, big.NewInt(int64(fee1))))
		// 1->0
		amt0, amt1, _, _ := pl.MockSwap(!zeroForOne1, amountOut, nil)
		amt0 = new(big.Int).Neg(amt0)
		logger.Info().Msgf("pool v3 1->0: amt0: %v amt1: %v dx1: %v dy1: %v dy2: %v", amt0, amt1, amtIn, amountOut, amountOut)
		return new(big.Int).Sub(amt0, amtIn)
	} else {
		amountOut := getAmountOutV2(ry, rx, amtIn, new(big.Int).Sub(e6, big.NewInt(int64(fee1))))
		// 1->0
		amt0, amt1, _, _ := pl.MockSwap(!zeroForOne1, amountOut, nil)
		amt1 = new(big.Int).Neg(amt1)
		logger.Info().Msgf("pool v3 0->1: amt0: %v amt1: %v amtIn: %v", amt0, amt1, amtIn)
		return new(big.Int).Sub(amt1, amtIn)
	}
}

func calcProfitV3ToV2(zeroForOne1 bool, rx, ry *big.Int, pl *pool.Pool, amtIn *big.Int, fee2 uint32) *big.Int {
	amt0, amt1, _, _ := pl.MockSwap(zeroForOne1, amtIn, nil)
	if zeroForOne1 {
		// 0->1
		amt1 = new(big.Int).Neg(amt1)
		// 1->0
		dx2 := getAmountOutV2(ry, rx, amt1, new(big.Int).Sub(e6, big.NewInt(int64(fee2))))
		logger.Info().Msgf("swap v3(0->1), v2(1->0): dx1: %v dy1(dy2): %v dx2: %v", amt0, amt1, dx2)
		return new(big.Int).Sub(dx2, amtIn)
	} else {
		amt0 = new(big.Int).Neg(amt0)
		dy2 := getAmountOutV2(rx, ry, amt0, new(big.Int).Sub(e6, big.NewInt(int64(fee2))))
		logger.Info().Msgf("swap v3(1->0), v2(0->1): dy1: %v dx1(dx2): %v dy2: %v", amt0, amt1, dy2)
		// 1->0
		return new(big.Int).Sub(dy2, amtIn)
	}
}

func createPool(price float64, fee uint32) *pool.Pool {
	tick := TickAtPrice(price)
	sp := priceToSqrtPrice(price)
	logger.Info().Msgf("create pool: tick=%v price=%v sp=%v", tick, price, sp)
	pl := &pool.Pool{
		PoolInfo: pool.PoolInfo{
			TickSpacing: 1,
			Typ:         common.PoolTypeCAMM,
			Fee:         uint(fee),
		},
		Tick:         int(tick),
		SqrtPriceX96: sp,
		Ticks:        map[int]*pool.TickInfo{},
		Liquidity:    big.NewInt(0),
	}

	return pl
}

func mint(pl *pool.Pool, priceLower, priceUpper float64, x, y *big.Int) {
	tickLower := TickNorm(TickAtPrice(priceLower), pl.TickSpacing)
	tickUpper := TickNorm(TickAtPrice(priceUpper), pl.TickSpacing)

	spLower := SqrtPriceX96AtTick(tickLower)
	spUpper := SqrtPriceX96AtTick(tickUpper)
	k := pool.GetLiquidityForAmounts(pl.SqrtPriceX96, spLower, spUpper, x, y)

	logger.Info().Msgf("mint: tick: [%v, %v] liquidity: %v", tickLower, tickUpper, k)
	pl.UpdateTicks(int(tickLower), k, false)
	pl.UpdateTicks(int(tickUpper), k, true)
	if pl.Tick >= int(tickLower) && pl.Tick < int(tickUpper) {
		pl.Liquidity.Add(pl.Liquidity, k)
	}
}

func TestCalcAmountInV3ToV20To1(t *testing.T) {
	logger.Init("", "", false)
	decimal.DivisionPrecision = 40

	fee1 := uint32(500)
	fee2 := uint32(500)
	x1 := convertToBigIntMust("1000")
	y1 := convertToBigIntMust("2490000")

	price := float64(2499)
	pl := createPool(price, fee1)
	mint(pl, 2499, 2502, convertToBigIntMust("1"), convertToBigIntMust("2500"))
	mint(pl, 2497, 2499, convertToBigIntMust("1"), convertToBigIntMust("2500"))
	mint(pl, 2496, 2497, convertToBigIntMust("1"), convertToBigIntMust("2500"))
	mint(pl, 2495, 2496, convertToBigIntMust("1"), convertToBigIntMust("2500"))
	mint(pl, 2490, 2495, convertToBigIntMust("1"), convertToBigIntMust("2500"))
	pl.Reload()

	bounds := BuildTickLiquidity(pl, true, 5)
	printBoundedTicks(bounds)

	amt, amtMid, amtOut, profit, profitable := CalcAmountInV3ToV2(true, pl, nil, bounds, x1, y1, fee1, fee2, true)
	t.Logf("amountIn: %v amtMid: %v amtOut: %v profit: %v profitable: %v", amt, amtMid, amtOut, profit, profitable)

	logger.Info().Msgf("calculate profit: amtIn: %v profit: %v", amt, calcProfitV3ToV2(true, x1, y1, pl, amt, fee1))
	// for i := -20; i <= 10; i += 2 {
	// 	namt := new(big.Int).Mul(amt, big.NewInt(int64(100+i)))
	// 	namt.Div(namt, big.NewInt(100))
	// 	logger.Info().Msgf("calculate profit: i: %d amtIn: %v profit: %v", i, namt, calcProfitV3ToV2(true, x1, y1, pl, namt, fee1))
	// }

	printPriceAfterSwapV3(true, pl, amt)
	_, amountOut, _, _ := pl.MockSwap(true, amt, nil)
	printPriceAfterSwapV2(false, x1, y1, new(big.Int).Neg(amountOut), fee1)
}

func TestCalcAmountInV3ToV21To0(t *testing.T) {
	logger.Init("", "", false)
	decimal.DivisionPrecision = 40

	fee1 := uint32(500)
	fee2 := uint32(500)
	x1 := convertToBigIntMust("1000")
	y1 := convertToBigIntMust("2508000")

	price := float64(2499.5)
	pl := createPool(price, fee2)
	mint(pl, 2500, 2502, convertToBigIntMust("10"), convertToBigIntMust("2500"))
	mint(pl, 2499, 2500, convertToBigIntMust("1"), convertToBigIntMust("250"))
	mint(pl, 2495, 2499, convertToBigIntMust("100"), convertToBigIntMust("250000"))
	mint(pl, 2490, 2496, convertToBigIntMust("100"), convertToBigIntMust("250000"))
	pl.Reload()

	bounds := BuildTickLiquidity(pl, false, 5)
	printBoundedTicks(bounds)

	amt, amtMid, amtOut, profit, profitable := CalcAmountInV3ToV2(false, pl, nil, bounds, x1, y1, fee1, fee2, true)
	t.Logf("amountIn: %v amtMid: %v amtOut: %v profit: %v profitable: %v", amt, amtMid, amtOut, profit, profitable)

	// for i := -5; i <= 10; i++ {
	// nAmt := new(big.Int).Div(mulmul(amt, big.NewInt(int64(10+i))), big.NewInt(10))
	logger.Info().Msgf("amtIn: %v profit: %v", amt, calcProfitV3ToV2(false, x1, y1, pl, amt, fee1))
	printPriceAfterSwapV3(false, pl, amt)
	amountOut := getAmountOutV3(pl, amt, false)
	printPriceAfterSwapV2(true, x1, y1, amountOut, fee1)
	// }
}

func TestCalcV3ZeroForOne(t *testing.T) {
	logger.Init("", "", false)
	decimal.DivisionPrecision = 40

	fee1 := uint32(500)
	fee2 := uint32(500)

	price1 := float64(2499.5)
	pl1 := createPool(price1, fee1)
	mint(pl1, 2500, 2502, convertToBigIntMust("10"), convertToBigIntMust("2500"))
	mint(pl1, 2499, 2500, convertToBigIntMust("1"), convertToBigIntMust("2500"))
	mint(pl1, 2495, 2499, convertToBigIntMust("2"), convertToBigIntMust("5000"))
	mint(pl1, 2490, 2495, convertToBigIntMust("1"), convertToBigIntMust("250000"))
	pl1.Reload()

	price2 := float64(2490.3)
	pl2 := createPool(price2, fee2)
	mint(pl2, 2500, 2502, convertToBigIntMust("3"), convertToBigIntMust("2500"))
	mint(pl2, 2499, 2500, convertToBigIntMust("1"), convertToBigIntMust("250"))
	mint(pl2, 2495, 2499, convertToBigIntMust("2"), convertToBigIntMust("250000"))
	mint(pl2, 2490, 2495, convertToBigIntMust("10"), convertToBigIntMust("250000"))
	pl2.Reload()

	tls1 := BuildTickLiquidity(pl1, true, 5)
	tls2 := BuildTickLiquidity(pl2, false, 5)
	printBoundedTicks(tls1)
	logger.Info().Msg("--------------------")
	printBoundedTicks(tls2)
	dx1, dy1, dx2, profit, profitable := CalcAmountInV3ToV3(true, pl1, pl2, tls1, tls2, fee1, fee2, true)
	t.Logf("amtIn: %v amtMid: %v amtOut: %v profit: %v profitable: %v", dx1, dy1, dx2, profit, profitable)
	t.Logf("after swap:")
	printPriceAfterSwapV3(true, pl1, dx1)
	printPriceAfterSwapV3(false, pl2, dy1)
}

func TestCalcV3OneForZero(t *testing.T) {
	logger.Init("", "", false)
	decimal.DivisionPrecision = 40

	fee1 := uint32(500)
	fee2 := uint32(500)

	price1 := float64(2499.5)
	pl1 := createPool(price1, fee1)
	mint(pl1, 2500, 2502, convertToBigIntMust("10"), convertToBigIntMust("2500"))
	mint(pl1, 2499, 2500, convertToBigIntMust("1"), convertToBigIntMust("2500"))
	mint(pl1, 2495, 2499, convertToBigIntMust("2"), convertToBigIntMust("5000"))
	mint(pl1, 2490, 2495, convertToBigIntMust("1"), convertToBigIntMust("250000"))
	pl1.Reload()

	price2 := float64(2490.3)
	pl2 := createPool(price2, fee2)
	mint(pl2, 2500, 2502, convertToBigIntMust("3"), convertToBigIntMust("2500"))
	mint(pl2, 2499, 2500, convertToBigIntMust("1"), convertToBigIntMust("250"))
	mint(pl2, 2495, 2499, convertToBigIntMust("2"), convertToBigIntMust("250000"))
	mint(pl2, 2490, 2495, convertToBigIntMust("10"), convertToBigIntMust("250000"))
	pl2.Reload()

	tls1 := BuildTickLiquidity(pl1, true, 5)
	tls2 := BuildTickLiquidity(pl2, false, 5)
	printBoundedTicks(tls1)
	logger.Info().Msg("--------------------")
	printBoundedTicks(tls2)
	amtIn, amtMid, amtOut, profit, profitable := CalcAmountInV3ToV3(false, pl2, pl1, tls2, tls1, fee1, fee2, true)
	t.Logf("amtIn: %v amtMid: %v amtOut: %v profit: %v profitable: %v", amtIn, amtMid, amtOut, profit, profitable)
	t.Logf("after swap:")
	printPriceAfterSwapV3(false, pl2, amtIn)
	printPriceAfterSwapV3(true, pl1, amtMid)
}
