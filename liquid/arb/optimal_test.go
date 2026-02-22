package arb

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/utils"
)

func TestBestAmountV2ToV2(t *testing.T) {
	logger.Init("", "", false)

	x1 := convertToBigIntMust2("999992858302806272")
	y1 := convertToBigIntMust2("83308209011427998")
	x2 := convertToBigIntMust2("11860242579412496001")
	y2 := convertToBigIntMust2("195756761335001164")
	amtIn, _ := bestAmountV2ToV2(x1, y1, x2, y2, 200, 200)
	t.Logf("0->1 amtIn: %v", amtIn)

	ramtIn, _ := bestAmountV2ToV2(y2, x2, y1, x1, 200, 200)
	t.Logf("1->0 amtIn: %v", ramtIn)
}

/*
func TestBestAmountV2ToV2(t *testing.T) {
	logger.Init("", "", false)

	fee1 := uint(100)
	fee2 := uint(500)
	x1 := convertToBigIntMust("1000")
	y1 := convertToBigIntMust("5000")
	x2 := convertToBigIntMust("2000")
	y2 := convertToBigIntMust("6000")

	amtIn, profit := BestAmountV2ToV2(x1, y1, x2, y2, true, fee1, fee2)
	t.Logf("best amountIn: %v profit: %v", amtIn.Div(amtIn, e18), profit)
}

func TestBestAmountV2ToV3ZeroForOne(t *testing.T) {
	fee1 := uint(3000)
	fee2 := uint(500)
	x1 := convertToBigIntMust("1000")
	y1 := convertToBigIntMust("2490000")
	sp2 := priceToSqrtPrice(2500)
	l2 := calculateLiquidityByX(convertToBigIntMust("1000"), 2500, 2401, 2601)

	t.Logf("v3 pool: sqrtPrice=%v liquidity=%v", sp2, l2)
	amt, _ := bestAmountV2ToV3Tick(y1, x1, l2, sp2, nil, fee1, fee2, true)
	t.Logf("v2 1->0, v3 0->1 best amount: %v", amt)
}

func TestBestAmountV2ToV3OneForZero(t *testing.T) {
	fee1 := uint(3000)
	fee2 := uint(500)
	x1 := convertToBigIntMust("1000")
	y1 := convertToBigIntMust("3000000")
	sp2 := priceToSqrtPrice(2500)
	l2 := calculateLiquidityByX(convertToBigIntMust("1000"), 2500, 2401, 2601)

	amt, _ := bestAmountV2ToV3Tick(x1, y1, l2, sp2, nil, fee1, fee2, false)
	// 88849352994656244387
	t.Logf("v2 0->1, v3 1->0 best amount: %v", amt)
}

func TestBestAmountV3ToV2ZeroForOne(t *testing.T) {
	utils.SkipCI(t)
	logger.Init("", "", false)

	fee1 := uint(3000)
	fee2 := uint(500)
	x2 := convertToBigIntMust("1000")
	y2 := convertToBigIntMust("2400000")
	sp1 := priceToSqrtPrice(2500)
	l1 := calculateLiquidityByX(convertToBigIntMust("1000"), 2500, 2401, 2601)

	logger.Info().Msgf("liquidity: %v sqrtPriceX96: %v", l1, sp1)

	amt := BestAmountV3ToV2(l1, sp1, x2, y2, fee1, fee2, true)
	t.Logf("v3 0->1, v2 1->0 best amount: %v", amt)
}

func TestBestAmountV3ToV2OneForZero(t *testing.T) {
	fee1 := uint(3000)
	fee2 := uint(500)
	x2 := convertToBigIntMust("1000")
	y2 := convertToBigIntMust("2410000") // 2410
	// l1 := calculateLiquidityByX(convertToBigIntMust("1000"), 2500, 0.1)
	// t.Logf("v3 pool y: %v", calculateYByLiquidity(l1, 2500, 0.1))
	// x1 := div(mulmul(l1, Q96), priceToSqrtPrice(2500))
	// y1 := div(mulmul(l1, priceToSqrtPrice(2500)), Q96)

	// t.Logf("v3 pool x: %v", x1)
	// t.Logf("v3 pool y: %v", y1)

	// amt := BestAmountV3ToV2(l1, sp1, x2, y2, fee1, fee2, false)
	// t.Logf("v3 1->0, v2 0->1 best amount: %v", amt1)

	// amt1 := BestAmountV2ToV3(l1, sp1, x2, y2, fee1, fee2, false)
	// t.Logf("v3 1->0, v2 0->1 best amount: %v", amt1)

	// amt2 := BestAmountV2ToV2(x1, y1, x2, y2, fee1, fee2)
	// t.Logf("v3 1->0, v2 0->1 best amount: %v", amt2)
	sp1 := priceToSqrtPrice(2500)
	tickLower := TickAtPrice(2401)
	tickUpper := TickAtPrice(2601)
	sqrtPriceLower := SqrtPriceX96AtTick(tickLower)
	sqrtPriceUpper := SqrtPriceX96AtTick(tickUpper)
	amt0 := convertToBigIntMust("1000")
	l1 := pool.GetLiquidityForAmounts(sp1, sqrtPriceLower, sqrtPriceUpper, amt0, bigZero)
	amt0In := BestAmountV3ToV2(l1, sp1, x2, y2, fee1, fee2, true)
	t.Logf("v3 0->1, v2 1->0 best amount: %v", amt0In)

	y2 = convertToBigIntMust("2600000") // 2600
	amt1In := BestAmountV3ToV2(l1, sp1, x2, y2, fee1, fee2, false)
	t.Logf("v3 1->0, v2 0->1 best amount: %v", amt1In)
}

func TestBestAmountV3ToV3ZeroToOne(t *testing.T) {
	fee1 := uint(3000)
	fee2 := uint(500)
	p1 := float64(2500)
	sp1 := priceToSqrtPrice(p1)
	l1 := calculateLiquidityByX(convertToBigIntMust("1000"), p1, 2401, 2601)

	p2 := float64(2490)
	sp2 := priceToSqrtPrice(p2)
	l2 := calculateLiquidityByX(convertToBigIntMust("1000"), p2, 2401, 2601)

	amt := BestAmountV3ToV3(l1, sp1, l2, sp2, fee1, fee2, true)
	t.Logf("v3 0->1, v3 1->0 best amount: %v", amt)
}

func TestBestAmountV3ToV3OneToZero(t *testing.T) {
	fee1 := uint(3000)
	fee2 := uint(500)
	p1 := float64(2500)
	sp1 := priceToSqrtPrice(p1)
	l1 := calculateLiquidityByX(convertToBigIntMust("1000"), p1, 2401, 2601)

	p2 := float64(2490)
	sp2 := priceToSqrtPrice(p2)
	l2 := calculateLiquidityByX(convertToBigIntMust("1000"), p2, 2401, 2601)

	amt := BestAmountV3ToV3(l2, sp2, l1, sp1, fee1, fee2, false)
	t.Logf("v3 1->0, v3 0->1 best amount: %v", amt)
}
*/

func priceToSqrtPrice(price float64) *big.Int {
	// return mulmul(Q96, new(big.Int).Sqrt(big.NewInt(int64(price))))
	p := decimal.NewFromFloat(price).Mul(decimal.NewFromInt(2).Pow(decimal.NewFromInt(192)))
	return new(big.Int).Sqrt(p.BigInt())
}

// float to big.Int: v*(10**18).
func convertToBigIntMust(v string) *big.Int {
	d, err := decimal.NewFromString(v)
	if err != nil {
		panic(err.Error())
	}

	return d.Mul(d_10_exp_18).BigInt()
}

// float to big.Int: v*(10**18).
func convertToBigIntMust2(v string) *big.Int {
	d, err := decimal.NewFromString(v)
	if err != nil {
		panic(err.Error())
	}

	// big.NewInt(0).SetString(v, 10)
	return d.BigInt()
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func calculateLiquidityByX(x *big.Int, price, priceLower, priceUpper float64) *big.Int {
	// L = dx*(sp1*sp2)/(sp1-sp2)
	// 	sqrtp_low = price_to_sqrtp(4545)
	// sqrtp_cur = price_to_sqrtp(5000)
	// sqrtp_upp = price_to_sqrtp(5500)
	// def liquidity0(amount, pa, pb):
	//
	//	if pa > pb:
	//	    pa, pb = pb, pa
	//	return (amount * (pa * pb) / q96) / (pb - pa)
	//
	// liq0 = liquidity0(amount_eth, sqrtp_cur, sqrtp_upp)
	sqrtPriceCurr := SqrtPriceX96AtTick(TickAtPrice(price))
	sqrtPriceLower := SqrtPriceX96AtTick(TickAtPrice(priceLower))
	sqrtPriceUpper := SqrtPriceX96AtTick(TickAtPrice(priceUpper))

	return pool.GetLiquidityForAmounts(sqrtPriceCurr, sqrtPriceLower, sqrtPriceUpper, x, bigZero)
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func calculateLiquidityByY(y *big.Int, price, priceLower, priceUpper float64) *big.Int {
	sqrtPriceCurr := SqrtPriceX96AtTick(TickAtPrice(price))
	sqrtPriceLower := SqrtPriceX96AtTick(TickAtPrice(priceLower))
	sqrtPriceUpper := SqrtPriceX96AtTick(TickAtPrice(priceUpper))

	return pool.GetLiquidityForAmounts(sqrtPriceCurr, sqrtPriceLower, sqrtPriceUpper, bigZero, y)
}

func TestBoundedProduct(t *testing.T) {
	utils.SkipCI(t)
	logger.Init("", "", false)

	tickLower := TickAtPrice(2401) // 49*49
	tickCurr := TickAtPrice(2500)  // 50*50
	tickUpper := TickAtPrice(2601) // 51*51
	amt0 := convertToBigIntMust("100")
	logger.Info().Msgf("amt0: %v", amt0)
	bp := ComputeBoundedProductAtTickByAmounts(tickCurr, tickLower, tickUpper, amt0, bigZero)

	t.Logf("bounded: %v", bp)
}

func TestOutofBounded(t *testing.T) {
	utils.SkipCI(t)

	tickLower := int32(50000)
	tickCurr := tickLower + 500
	tickUpper := tickLower + 1000
	sqrtPriceCurr := SqrtPriceX96AtTick(tickCurr)
	sqrtPriceLower := SqrtPriceX96AtTick(tickLower)
	sqrtPriceUpper := SqrtPriceX96AtTick(tickUpper)

	t.Logf("sqrtPriceCurr: %v sqrtPriceLower: %v sqrtPriceUpper: %v", sqrtPriceCurr, sqrtPriceLower, sqrtPriceUpper)

	amt0 := convertToBigIntMust("1")
	l := pool.GetLiquidityForAmounts(sqrtPriceCurr, sqrtPriceLower, sqrtPriceUpper, amt0, bigZero)
	t.Logf("amt0: %v liqudity: %v", amt0, l)
	a0, a1 := pool.GetAmountsForLiquidity(sqrtPriceCurr, sqrtPriceLower, sqrtPriceUpper, l)
	t.Logf("amount for liquidity: amt0=%v amt1=%v", a0, a1)

	amt0In := pool.GetAmount0Delta(sqrtPriceLower, sqrtPriceCurr, l, false)
	t.Logf("drain amount1 amt0In: %v", amt0In)

	amtIns := []*big.Int{
		convertToBigIntMust("0.1"),
		convertToBigIntMust("0.5"),
		convertToBigIntMust("1"),
		convertToBigIntMust("1.01"),
		convertToBigIntMust("1.02"),
		convertToBigIntMust("1.03"),
		convertToBigIntMust("1.04"),
		convertToBigIntMust("1.05"),
		convertToBigIntMust("1.1"),
		convertToBigIntMust("1.2"),
		convertToBigIntMust("1.3"),
		convertToBigIntMust("1.5"),
		convertToBigIntMust("2"),
		convertToBigIntMust("2.5"),
		convertToBigIntMust("3"),
	}

	for _, amt := range amtIns {
		amtOut := getBoundedAmountOut(true, amt, l, sqrtPriceCurr, 500)
		t.Logf("0->1, amountIn: %v amountOut: %v price: %v", amt, amtOut,
			decimal.NewFromBigInt(amtOut, 0).Div(decimal.NewFromBigInt(amt, 0)))
	}
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func getBoundedAmountOut(zeroForOne bool, amtIn, liquidity, sqrtPriceX96 *big.Int, fee uint32) *big.Int {
	r := new(big.Int).Sub(e6, big.NewInt(int64(fee)))
	// dy = dx*L*SP*r2*SP / (L*Q192*e6 + dx*Q96*r2*SP)
	if zeroForOne {
		numerator := mulmul(amtIn, liquidity, sqrtPriceX96, r, sqrtPriceX96)
		denominator := new(big.Int).Add(mulmul(liquidity, Q192, e6), mulmul(amtIn, Q96, r, sqrtPriceX96))

		return new(big.Int).Div(numerator, denominator)
	}
	// dx = dy*r*Q192*L/ (Sp*Sp*e6*L + dy*r*Q96*Sp)
	numerator := mulmul(amtIn, r, Q192, liquidity)
	denominator := new(big.Int).Add(mulmul(sqrtPriceX96, sqrtPriceX96, e6, liquidity), mulmul(amtIn, r, Q96, sqrtPriceX96))

	return new(big.Int).Div(numerator, denominator)
}
