package pool

import "math/big"

func GetLiquidityForAmounts(
	sqrtRatioX96,
	sqrtRatioAX96,
	sqrtRatioBX96,
	amount0,
	amount1 *big.Int,
) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	if sqrtRatioX96.Cmp(sqrtRatioAX96) <= 0 {
		return getLiquidityForAmount0(sqrtRatioAX96, sqrtRatioBX96, amount0)
	} else if sqrtRatioX96.Cmp(sqrtRatioBX96) < 0 {
		var liquidity0, liquidity1 *big.Int
		if amount0 != nil && amount0.Cmp(bigZero) > 0 {
			liquidity0 = getLiquidityForAmount0(sqrtRatioX96, sqrtRatioBX96, amount0)
		}

		if amount1 != nil && amount1.Cmp(bigZero) > 0 {
			liquidity1 = getLiquidityForAmount1(sqrtRatioAX96, sqrtRatioX96, amount1)
		}
		// println("    liquidity0:", liquidity0.String())
		// println("    liquidity1:", liquidity1.String())

		if liquidity0 != nil && liquidity1 != nil {
			if liquidity0.Cmp(liquidity1) > 0 {
				return liquidity1
			}

			return liquidity0
		} else if liquidity0 != nil {
			return liquidity0
		} else {
			return liquidity1
		}
	} else {
		// println(" getLiquidityForAmount1   sqrtRatioAX96:", sqrtRatioAX96.String())
		// println(" getLiquidityForAmount1   sqrtRatioBX96:", sqrtRatioBX96.String())
		// println(" getLiquidityForAmount1   amount1:", amount1.String())
		return getLiquidityForAmount1(sqrtRatioAX96, sqrtRatioBX96, amount1)
	}
}

func getLiquidityForAmounts(tickLower, tickUpper int, sqrtPriceX96, amount0, amount1 *big.Int) *big.Int {
	sqrtRatioAX96 := getSqrtRatioAtTick(tickLower)
	sqrtRatioBX96 := getSqrtRatioAtTick(tickUpper)

	// println("sqrtRatioAX96:", sqrtRatioAX96.String())
	// println("sqrtRatioBX96:", sqrtRatioBX96.String())
	return GetLiquidityForAmounts(sqrtPriceX96, sqrtRatioAX96, sqrtRatioBX96, amount0, amount1)
}

func getLiquidityForAmount0(sqrtRatioAX96, sqrtRatioBX96, amount0 *big.Int) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	intermediate := mulDiv(sqrtRatioAX96, sqrtRatioBX96, Q96)

	return mulDiv(amount0, intermediate, big.NewInt(0).Sub(sqrtRatioBX96, sqrtRatioAX96))
}

func getLiquidityForAmount1(sqrtRatioAX96, sqrtRatioBX96, amount1 *big.Int) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	return mulDiv(amount1, Q96, big.NewInt(0).Sub(sqrtRatioBX96, sqrtRatioAX96))
}

func GetAmountsForLiquidity(
	sqrtRatioX96,
	sqrtRatioAX96,
	sqrtRatioBX96,
	liquidity *big.Int,
) (amount0, amount1 *big.Int) {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	if sqrtRatioX96.Cmp(sqrtRatioAX96) <= 0 {
		amount0 = GetAmount0ForLiquidity(sqrtRatioAX96, sqrtRatioBX96, liquidity)
		amount1 = big.NewInt(0)
	} else if sqrtRatioX96.Cmp(sqrtRatioBX96) < 0 {
		amount0 = GetAmount0ForLiquidity(sqrtRatioX96, sqrtRatioBX96, liquidity)
		amount1 = GetAmount1ForLiquidity(sqrtRatioAX96, sqrtRatioX96, liquidity)
	} else {
		amount0 = big.NewInt(0)
		amount1 = GetAmount1ForLiquidity(sqrtRatioAX96, sqrtRatioBX96, liquidity)
	}

	return
}

// L*Q96*(SPB-SPA)/SPB/SPA.
func GetAmount0ForLiquidity(
	sqrtRatioAX96,
	sqrtRatioBX96,
	liquidity *big.Int,
) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		panic("getAmount0ForLiquidity: sqrtRatioAX96 should less than sqrtRatioBX96")
	}

	amt := mulmul(liquidity, new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96), Q96)
	amt.Div(amt, sqrtRatioBX96)
	amt.Div(amt, sqrtRatioAX96)

	return amt
}

// L*(SPB-SPA)/Q96.
func GetAmount1ForLiquidity(
	sqrtRatioAX96,
	sqrtRatioBX96,
	liquidity *big.Int,
) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		panic("getAmount1ForLiquidity: sqrtRatioAX96 should less than sqrtRatioBX96")
	}

	return mulDiv(liquidity, new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96), Q96)
}
