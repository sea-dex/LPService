package pool

import (
	"math/big"
)

const FixedPoint96_RESOLUTION = 96

var (
	MIN_SQRT_RATIO = big.NewInt(4295128739)
	MAX_SQRT_RATIO = toBigIntMust("1461446703485210103287273052203988822378723970342")
)

func toBigIntMust(s string) *big.Int {
	v, ok := big.NewInt(0).SetString(s, 0)
	if !ok {
		panic("convert string to bigInt failed")
	}

	return v
}

// x >= 0, y >= 0.
func divRoundingUp(x, y *big.Int) *big.Int {
	// assembly {
	// 	z := add(div(x, y), gt(mod(x, y), 0))
	// }
	quotient := big.NewInt(0).Quo(x, y)
	mod := big.NewInt(0).Rem(x, y)

	if mod.Cmp(bigZero) == 0 {
		return quotient
	}

	return quotient.Add(quotient, big.NewInt(1))
}

func getTickAtSqrtRatio(sqrtPriceX96 *big.Int) int {
	if sqrtPriceX96.Cmp(MIN_SQRT_RATIO) < 0 || sqrtPriceX96.Cmp(MAX_SQRT_RATIO) >= 0 {
		panic("invalid sqrtPriceX96")
	}

	ratio := big.NewInt(0).Lsh(sqrtPriceX96, 32)
	r := big.NewInt(0).Set(ratio)
	msb := uint64(0)

	shlOrShr := func(n int64, mask *big.Int) {
		f := uint64(0)
		if r.Cmp(mask) > 0 {
			f = 1 << n
		}

		// msb.Or(msb, big.NewInt(int64(f)))
		msb = msb | f

		if n > 0 {
			r.Rsh(r, uint(f))
		}
	}

	shlOrShr(7, toBigIntMust("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"))
	shlOrShr(6, toBigIntMust("0xFFFFFFFFFFFFFFFF"))
	shlOrShr(5, toBigIntMust("0xFFFFFFFF"))
	shlOrShr(4, toBigIntMust("0xFFFF"))
	shlOrShr(3, toBigIntMust("0xFF"))
	shlOrShr(2, toBigIntMust("0xF"))
	shlOrShr(1, toBigIntMust("0x3"))
	shlOrShr(0, toBigIntMust("0x1"))

	if msb >= 128 {
		r = big.NewInt(0).Rsh(ratio, uint(msb-127))
	} else {
		r = big.NewInt(0).Lsh(ratio, uint(127-msb))
	}

	log_2 := big.NewInt(0).Lsh(big.NewInt(int64(msb)-128), 64) // nolint

	mulShr3 := func(n uint) {
		r.Mul(r, r)
		r.Rsh(r, 127)
		f := big.NewInt(0).Rsh(r, 128)
		log_2.Or(log_2, big.NewInt(0).Lsh(f, n))

		if n > 50 {
			r.Rsh(r, uint(f.Uint64()))
		}
	}

	mulShr3(63)
	mulShr3(62)
	mulShr3(61)
	mulShr3(60)
	mulShr3(59)
	mulShr3(58)
	mulShr3(57)
	mulShr3(56)
	mulShr3(55)
	mulShr3(54)
	mulShr3(53)
	mulShr3(52)
	mulShr3(51)
	mulShr3(50)

	log_sqrt10001 := big.NewInt(0).Mul(log_2, toBigIntMust("255738958999603826347141"))
	// int24 tickLow = int24((log_sqrt10001 - 3402992956809132418596140100660247210) >> 128);
	biTickLow := big.NewInt(0).Sub(log_sqrt10001, toBigIntMust("3402992956809132418596140100660247210"))
	biTickLow.Rsh(biTickLow, 128)
	tickLow := int(biTickLow.Int64())
	// int24 tickHi = int24((log_sqrt10001 + 291339464771989622907027621153398088495) >> 128);
	biTickHi := big.NewInt(0).Add(log_sqrt10001, toBigIntMust("291339464771989622907027621153398088495"))
	biTickHi.Rsh(biTickHi, 128)
	tickHi := int(biTickHi.Int64())

	if tickLow == tickHi {
		return tickLow
	}

	if getSqrtRatioAtTick(tickHi).Cmp(sqrtPriceX96) <= 0 {
		return tickHi
	} else {
		return tickLow
	}
}

func getSqrtRatioAtTick(tick int) *big.Int {
	absTick := tick
	if absTick < 0 {
		absTick = -tick
	}

	ratio := big.NewInt(0)
	if absTick&0x1 != 0 {
		ratio.SetString("0xfffcb933bd6fad37aa2d162d1a594001", 0)
	} else {
		ratio.SetString("0x100000000000000000000000000000000", 0)
	}

	if absTick&0x2 != 0 {
		v, _ := big.NewInt(0).SetString("0xfff97272373d413259a46990580e213a", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x4 != 0 {
		v, _ := big.NewInt(0).SetString("0xfff2e50f5f656932ef12357cf3c7fdcc", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x8 != 0 {
		v, _ := big.NewInt(0).SetString("0xffe5caca7e10e4e61c3624eaa0941cd0", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x10 != 0 {
		v, _ := big.NewInt(0).SetString("0xffcb9843d60f6159c9db58835c926644", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x20 != 0 {
		v, _ := big.NewInt(0).SetString("0xff973b41fa98c081472e6896dfb254c0", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x40 != 0 {
		v, _ := big.NewInt(0).SetString("0xff2ea16466c96a3843ec78b326b52861", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x80 != 0 {
		v, _ := big.NewInt(0).SetString("0xfe5dee046a99a2a811c461f1969c3053", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x100 != 0 {
		v, _ := big.NewInt(0).SetString("0xfcbe86c7900a88aedcffc83b479aa3a4", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x200 != 0 {
		v, _ := big.NewInt(0).SetString("0xf987a7253ac413176f2b074cf7815e54", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x400 != 0 {
		v, _ := big.NewInt(0).SetString("0xf3392b0822b70005940c7a398e4b70f3", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x800 != 0 {
		v, _ := big.NewInt(0).SetString("0xe7159475a2c29b7443b29c7fa6e889d9", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x1000 != 0 {
		v, _ := big.NewInt(0).SetString("0xd097f3bdfd2022b8845ad8f792aa5825", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x2000 != 0 {
		v, _ := big.NewInt(0).SetString("0xa9f746462d870fdf8a65dc1f90e061e5", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x4000 != 0 {
		v, _ := big.NewInt(0).SetString("0x70d869a156d2a1b890bb3df62baf32f7", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x8000 != 0 {
		v, _ := big.NewInt(0).SetString("0x31be135f97d08fd981231505542fcfa6", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x10000 != 0 {
		v, _ := big.NewInt(0).SetString("0x9aa508b5b7a84e1c677de54f3e99bc9", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x20000 != 0 {
		v, _ := big.NewInt(0).SetString("0x5d6af8dedb81196699c329225ee604", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x40000 != 0 {
		v, _ := big.NewInt(0).SetString("0x2216e584f5fa1ea926041bedfe98", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if absTick&0x80000 != 0 {
		v, _ := big.NewInt(0).SetString("0x48a170391f7dc42444e8fa2", 0)
		ratio = ratio.Mul(ratio, v)
		ratio = ratio.Rsh(ratio, 128)
	}

	if tick > 0 {
		ratio = ratio.Div(uint256Max, ratio)
	}

	sqrtPrice := big.NewInt(0).Rsh(ratio, 32)
	_1_lsh32 := big.NewInt(1).Lsh(big.NewInt(1), 32)

	mod := big.NewInt(0).Mod(ratio, _1_lsh32)
	if mod.Cmp(bigZero) != 0 {
		sqrtPrice.Add(sqrtPrice, big.NewInt(1))
	}

	return sqrtPrice
}

func mulDiv(x, y, z *big.Int) *big.Int {
	v := big.NewInt(1).Mul(x, y)
	v.Quo(v, z)

	if v.Cmp(uint256Max) > 0 {
		panic("mulDiv great than uint256 max")
	}

	return v
}

func mulmul(args ...*big.Int) *big.Int {
	a := big.NewInt(1)
	for _, arg := range args {
		a.Mul(a, arg)
	}

	return a
}

func mulDivRoundingUp(x, y, z *big.Int) *big.Int {
	v := big.NewInt(1).Mul(x, y)
	vv := big.NewInt(1).Set(v)
	vv.Quo(vv, z)

	rem := big.NewInt(1).Rem(v, z)
	// println("x:", x.String())
	// println("y:", y.String())
	// println("z:", z.String())
	// println("mod:", rem.String())
	// println("mod:", big.NewInt(1).Mod(v, z).String())
	if rem.Cmp(bigZero) > 0 {
		vv.Add(vv, big.NewInt(1))
	}

	return vv
}

// / @notice Gets the amount0 delta between two prices
// / @dev Calculates liquidity / sqrt(lower) - liquidity / sqrt(upper),
// / i.e. liquidity * (sqrt(upper) - sqrt(lower)) / (sqrt(upper) * sqrt(lower))
// / @param sqrtRatioAX96 A sqrt price
// / @param sqrtRatioBX96 Another sqrt price
// / @param liquidity The amount of usable liquidity
// / @param roundUp Whether to round the amount up or down
// / @return amount0 Amount of token0 required to cover a position of size liquidity between the two passed prices.
func GetAmount0Delta(
	sqrtRatioAX96 *big.Int,
	sqrtRatioBX96 *big.Int,
	liquidity *big.Int,
	roundUp bool,
) (amount0 *big.Int) {
	// println("sqrtRatioBX96:", sqrtRatioBX96.String())
	// println("sqrtRatioAX96:", sqrtRatioAX96.String())
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	numerator1 := big.NewInt(1).Lsh(liquidity, FixedPoint96_RESOLUTION) // FixedPoint96.RESOLUTION)
	numerator2 := big.NewInt(1).Sub(sqrtRatioBX96, sqrtRatioAX96)

	if sqrtRatioAX96.Cmp(bigZero) <= 0 {
		panic("sqrtRatioAX96 should great than 0")
	}

	// println("numerator1:", numerator1.String())
	// println("numerator2:", numerator2.String())
	// println("sqrtRatioBX96:", sqrtRatioBX96.String())
	// println("sqrtRatioAX96:", sqrtRatioAX96.String())
	if roundUp {
		amount0 = divRoundingUp(mulDivRoundingUp(numerator1, numerator2, sqrtRatioBX96), sqrtRatioAX96)
	} else {
		amount0 = mulDiv(numerator1, numerator2, sqrtRatioBX96)
		amount0.Div(amount0, sqrtRatioAX96)
	}

	return
}

// / @notice Gets the amount1 delta between two prices
// / @dev Calculates liquidity * (sqrt(upper) - sqrt(lower))
// / @param sqrtRatioAX96 A sqrt price
// / @param sqrtRatioBX96 Another sqrt price
// / @param liquidity The amount of usable liquidity
// / @param roundUp Whether to round the amount up, or down
// / @return amount1 Amount of token1 required to cover a position of size liquidity between the two passed prices.
func GetAmount1Delta(
	sqrtRatioAX96 *big.Int,
	sqrtRatioBX96 *big.Int,
	liquidity *big.Int,
	roundUp bool,
) (amount1 *big.Int) {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	if roundUp {
		amount1 = mulDivRoundingUp(liquidity, big.NewInt(1).Sub(sqrtRatioBX96, sqrtRatioAX96), Q96)
	} else {
		amount1 = mulDiv(liquidity, big.NewInt(1).Sub(sqrtRatioBX96, sqrtRatioAX96), Q96)
	}

	return
}

// / @notice Gets the next sqrt price given an input amount of token0 or token1
// / @dev Throws if price or liquidity are 0, or if the next price is out of bounds
// / @param sqrtPX96 The starting price, i.e., before accounting for the input amount
// / @param liquidity The amount of usable liquidity
// / @param amountIn How much of token0, or token1, is being swapped in
// / @param zeroForOne Whether the amount in is token0 or token1
// / @return sqrtQX96 The price after adding the input amount to token0 or token1.
func getNextSqrtPriceFromInput(
	sqrtPX96 *big.Int,
	liquidity *big.Int,
	amountIn *big.Int,
	zeroForOne bool,
) *big.Int {
	if sqrtPX96.Cmp(bigZero) <= 0 {
		panic("sqrtPX96 should great than 0")
	}

	if liquidity.Cmp(bigZero) <= 0 {
		panic("liquidity should great than 0")
	}

	if zeroForOne {
		return getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amountIn, true)
	}

	return getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amountIn, true)
}

// / @notice Gets the next sqrt price given an output amount of token0 or token1
// / @dev Throws if price or liquidity are 0 or the next price is out of bounds
// / @param sqrtPX96 The starting price before accounting for the output amount
// / @param liquidity The amount of usable liquidity
// / @param amountOut How much of token0, or token1, is being swapped out
// / @param zeroForOne Whether the amount out is token0 or token1
// / @return sqrtQX96 The price after removing the output amount of token0 or token1.
func getNextSqrtPriceFromOutput(
	sqrtPX96 *big.Int,
	liquidity *big.Int,
	amountOut *big.Int,
	zeroForOne bool,
) *big.Int {
	if sqrtPX96.Cmp(bigZero) <= 0 {
		panic("sqrtPX96 should great than 0")
	}

	if liquidity.Cmp(bigZero) <= 0 {
		panic("liquidity should great than 0")
	}

	if zeroForOne {
		return getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amountOut, false)
	}

	return getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amountOut, false)
}

// / @notice Gets the next sqrt price given a delta of token0
// / @dev Always rounds up, because in the exact output case (increasing price) we need to move the price at least
// / far enough to get the desired output amount, and in the exact input case (decreasing price) we need to move the
// / price less in order to not send too much output.
// / The most precise formula for this is liquidity * sqrtPX96 / (liquidity +- amount * sqrtPX96),
// / if this is impossible because of overflow, we calculate liquidity / (liquidity / sqrtPX96 +- amount).
// / @param sqrtPX96 The starting price, i.e. before accounting for the token0 delta
// / @param liquidity The amount of usable liquidity
// / @param amount How much of token0 to add or remove from virtual reserves
// / @param add Whether to add or remove the amount of token0
// / @return The price after adding or removing amount, depending on add.
func getNextSqrtPriceFromAmount0RoundingUp(
	sqrtPX96 *big.Int,
	liquidity *big.Int,
	amount *big.Int,
	add bool,
) *big.Int {
	if amount.Cmp(bigZero) == 0 {
		return big.NewInt(1).Set(sqrtPX96)
	}

	// println("getNextSqrtPriceFromAmount0RoundingUp:", add)
	// println("sqrtPX96:", sqrtPX96.String())
	// println("liquidity:", liquidity.String())
	// println("amount:", amount.String())

	numerator1 := big.NewInt(1).Lsh(liquidity, FixedPoint96_RESOLUTION) // FixedPoint96.RESOLUTION
	product := big.NewInt(1).Mul(amount, sqrtPX96)

	if add {
		denominator := big.NewInt(1).Add(numerator1, product)
		return mulDivRoundingUp(numerator1, sqrtPX96, denominator)
	} else {
		if numerator1.Cmp(product) <= 0 {
			panic("numerator1 should great than product")
		}

		denominator := big.NewInt(1).Sub(numerator1, product)
		// println("numerator1:", numerator1.String())
		// println("sqrtPX96:", sqrtPX96.String())
		// println("denominator:", denominator.String())
		return mulDivRoundingUp(numerator1, sqrtPX96, denominator)
	}
}

// / @notice Gets the next sqrt price given a delta of token1
// / @dev Always rounds down, because in the exact output case (decreasing price) we need to move the price at least
// / far enough to get the desired output amount, and in the exact input case (increasing price) we need to move the
// / price less in order to not send too much output.
// / The formula we compute is within <1 wei of the lossless version: sqrtPX96 +- amount / liquidity
// / @param sqrtPX96 The starting price, i.e., before accounting for the token1 delta
// / @param liquidity The amount of usable liquidity
// / @param amount How much of token1 to add, or remove, from virtual reserves
// / @param add Whether to add, or remove, the amount of token1
// / @return The price after adding or removing `amount`.
func getNextSqrtPriceFromAmount1RoundingDown(
	sqrtPX96 *big.Int,
	liquidity *big.Int,
	amount *big.Int,
	add bool,
) *big.Int {
	quotient := big.NewInt(1)

	if add {
		if amount.Cmp(uint160Max) <= 0 {
			quotient.Lsh(amount, 96)
			quotient.Div(quotient, liquidity)
		} else {
			quotient = mulDiv(amount, Q96, liquidity)
		}

		return big.NewInt(1).Add(sqrtPX96, quotient)
	}

	if amount.Cmp(uint160Max) <= 0 {
		quotient = divRoundingUp(big.NewInt(1).Lsh(amount, 96), liquidity)
	} else {
		quotient = mulDivRoundingUp(amount, Q96, liquidity)
	}

	if sqrtPX96.Cmp(quotient) <= 0 {
		panic("sqrtPX96 should great than quotient")
	}

	return big.NewInt(1).Sub(sqrtPX96, quotient)
}
