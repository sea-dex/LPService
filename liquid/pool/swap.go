package pool

import "math/big"

func computeSwapStep(
	sqrtRatioCurrentX96 *big.Int,
	sqrtRatioTargetX96 *big.Int,
	liquidity *big.Int,
	amountRemaining *big.Int,
	feePips int,
) (sqrtRatioNextX96, amountIn, amountOut, feeAmount *big.Int) {
	zeroForOne := sqrtRatioCurrentX96.Cmp(sqrtRatioTargetX96) >= 0
	exactIn := amountRemaining.Cmp(bigZero) >= 0

	_1e6 := big.NewInt(1e6)
	if exactIn {
		amountRemainingLessFee := mulDiv(amountRemaining, big.NewInt(1).Sub(_1e6, big.NewInt(int64(feePips))), _1e6)

		if zeroForOne {
			amountIn = GetAmount0Delta(sqrtRatioTargetX96, sqrtRatioCurrentX96, liquidity, true)
		} else {
			amountIn = GetAmount1Delta(sqrtRatioCurrentX96, sqrtRatioTargetX96, liquidity, true)
		}

		if amountRemainingLessFee.Cmp(amountIn) >= 0 {
			sqrtRatioNextX96 = big.NewInt(1).Set(sqrtRatioTargetX96)
		} else {
			sqrtRatioNextX96 = getNextSqrtPriceFromInput(
				sqrtRatioCurrentX96,
				liquidity,
				amountRemainingLessFee,
				zeroForOne)
		}
	} else {
		if zeroForOne {
			amountOut = GetAmount1Delta(sqrtRatioTargetX96, sqrtRatioCurrentX96, liquidity, false)
		} else {
			amountOut = GetAmount0Delta(sqrtRatioCurrentX96, sqrtRatioTargetX96, liquidity, false)
		}

		negAmountRemaining := big.NewInt(1).Neg(amountRemaining)
		// println("amountOut:", amountOut.String())
		// println("negAmountRemaining:", negAmountRemaining.String())
		if negAmountRemaining.Cmp(amountOut) >= 0 {
			sqrtRatioNextX96 = big.NewInt(1).Set(sqrtRatioTargetX96)
		} else {
			sqrtRatioNextX96 = getNextSqrtPriceFromOutput(
				sqrtRatioCurrentX96,
				liquidity,
				negAmountRemaining,
				zeroForOne,
			)
		}
	}

	maxSqrtPrice := sqrtRatioTargetX96.Cmp(sqrtRatioNextX96) == 0

	if zeroForOne {
		if !(maxSqrtPrice && exactIn) {
			amountIn = GetAmount0Delta(sqrtRatioNextX96, sqrtRatioCurrentX96, liquidity, true)
		}

		if !(maxSqrtPrice && !exactIn) {
			amountOut = GetAmount1Delta(sqrtRatioNextX96, sqrtRatioCurrentX96, liquidity, false)
		}
	} else {
		if !(maxSqrtPrice && exactIn) {
			// println(sqrtRatioTargetX96.String())
			// println(sqrtRatioCurrentX96.String(), sqrtRatioNextX96.String(), liquidity.String())
			amountIn = GetAmount1Delta(sqrtRatioCurrentX96, sqrtRatioNextX96, liquidity, true)
		}

		if !(maxSqrtPrice && !exactIn) {
			amountOut = GetAmount0Delta(sqrtRatioCurrentX96, sqrtRatioNextX96, liquidity, false)
		}
	}

	if !exactIn && amountOut.Cmp(big.NewInt(1).Neg(amountRemaining)) > 0 {
		amountOut = big.NewInt(1).Neg(amountRemaining)
	}

	if exactIn && sqrtRatioNextX96.Cmp(sqrtRatioTargetX96) != 0 {
		feeAmount = big.NewInt(1).Abs(amountRemaining)
		feeAmount.Sub(feeAmount, amountIn)
	} else {
		feeAmount = mulDivRoundingUp(amountIn, big.NewInt(int64(feePips)), big.NewInt(1).Sub(_1e6, big.NewInt(int64(feePips))))
	}

	return
}
