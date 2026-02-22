package pool

import (
	"math"
	"math/big"

	"github.com/shopspring/decimal"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
)

var (
	base1_0001    = float64(1.0001)
	logBase1_0001 = math.Log(base1_0001)
)

// GetAmountOut give exact amountIn, calculate max amountOut.
func (pool *Pool) GetAmountOut(zeroForOne bool, amountIn *big.Int) *big.Int {
	switch {
	case pool.Typ.IsAMMVariety():
		if pool.Stable && (pool.Typ == common.PoolTypeAeroAMM || pool.Typ == common.PoolTypeInfusionAMM) {
			return pool.getStableOut(zeroForOne, amountIn)
		}

		return pool.GetAmountOutV2(zeroForOne, amountIn)

	case pool.Typ.IsCAMMVariety():
		return pool.getAmountOutV3(zeroForOne, amountIn)

	default:
		logger.Fatal().Str("pool", pool.Address).Uint("type", uint(pool.Typ)).Msg("invalid pool type")
	}

	return big.NewInt(0)
}

func (pool *Pool) GetAmountOutAndPrice(zeroForOne bool, amountIn *big.Int, p0 bool) (*big.Int, decimal.Decimal) {
	switch {
	case pool.Typ.IsAMMVariety():
		if pool.Stable && (pool.Typ == common.PoolTypeAeroAMM || pool.Typ == common.PoolTypeInfusionAMM) {
			return pool.getStableOutAndPrice(zeroForOne, amountIn, p0)
		}

		return pool.GetAmountOutV2AndPrice(zeroForOne, amountIn, p0)

	case pool.Typ.IsCAMMVariety():
		return pool.getAmountOutV3AndPrice(zeroForOne, amountIn, p0)

	default:
		logger.Fatal().Str("pool", pool.Address).Uint("type", uint(pool.Typ)).Msg("invalid pool type")
	}

	return big.NewInt(0), decimal.Zero
}

func (pool *Pool) getStableOut(zeroForOne bool, amountIn *big.Int) *big.Int {
	if pool.decimals0 == nil || pool.decimals1 == nil {
		logger.Warn().Msgf("pool %s decimals is nil", pool.Address)
		return big.NewInt(0)
	}

	fee := mulDiv(amountIn, big.NewInt(int64(pool.Fee)), e6) // nolint new(big.Int).Sub(e6, big.NewInt(int64(pool.Fee)))
	amountInWithFee := new(big.Int).Sub(amountIn, fee)

	decimals0, decimals1 := pool.decimals0, pool.decimals1
	reserve0, reserve1 := pool.Reserve0, pool.Reserve1
	xy := _k(reserve0, reserve1, decimals0, decimals1)

	// Pre-calculate these values
	reserve0E18 := mulDiv(reserve0, e18, decimals0)
	reserve1E18 := mulDiv(reserve1, e18, decimals1)

	if zeroForOne {
		amountInE18 := mulDiv(amountInWithFee, e18, decimals0)

		_y := _get_y(decimals0, decimals1, new(big.Int).Add(amountInE18, reserve0E18), xy, reserve1E18)
		if _y == nil {
			return big.NewInt(0)
		}

		y := new(big.Int).Sub(reserve1E18, _y)

		return mulDivE18(y, decimals1)
	} else {
		amountInE18 := mulDiv(amountInWithFee, e18, decimals1)

		_y := _get_y(decimals0, decimals1, new(big.Int).Add(amountInE18, reserve1E18), xy, reserve0E18)
		if _y == nil {
			return big.NewInt(0)
		}

		y := new(big.Int).Sub(reserve0E18, _y)

		return mulDivE18(y, decimals0)
	}
}

func (pool *Pool) getStableOutAndPrice(zeroForOne bool, amountIn *big.Int, p0 bool) (*big.Int, decimal.Decimal) {
	amountOut := pool.getStableOut(zeroForOne, amountIn)

	if amountOut.Cmp(bigZero) == 0 {
		return big.NewInt(0), decimal.Zero
	}

	bak := *pool
	if zeroForOne {
		bak.Reserve0 = new(big.Int).Add(bak.Reserve0, amountIn)
		bak.Reserve1 = new(big.Int).Sub(bak.Reserve1, amountOut)
	} else {
		bak.Reserve0 = new(big.Int).Sub(bak.Reserve0, amountOut)
		bak.Reserve1 = new(big.Int).Add(bak.Reserve1, amountIn)
	}

	var (
		amtIn  *big.Int
		amtOut *big.Int
		price  decimal.Decimal
	)

	if p0 {
		amtIn = new(big.Int).Div(bak.Reserve0, big.NewInt(10000))
		amtOut = bak.getStableOut(true, amtIn)
	} else {
		amtIn = new(big.Int).Div(bak.Reserve1, big.NewInt(10000))
		amtOut = bak.getStableOut(false, amtIn)
	}

	if amtOut.Cmp(bigZero) == 0 {
		return amountOut, decimal.Zero
	}

	price = decimal.NewFromBigInt(amtOut, 0).Div(decimal.NewFromBigInt(amtIn, 0))

	return amountOut, price
}

// uniswap v2.
func (pool *Pool) GetAmountOutV2(zeroForOne bool, amountIn *big.Int) *big.Int {
	var reserveIn, reserveOut *big.Int

	if pool.Reserve0.Cmp(bigZero) == 0 || pool.Reserve1.Cmp(bigZero) == 0 {
		return big.NewInt(0)
	}

	if zeroForOne {
		reserveIn, reserveOut = pool.Reserve0, pool.Reserve1
	} else {
		reserveIn, reserveOut = pool.Reserve1, pool.Reserve0
	}

	// Calculate the amount out using the constant product formula
	amountInWithFee := new(big.Int).Mul(amountIn, new(big.Int).Sub(e6, big.NewInt(int64(pool.Fee)))) // nolint
	numerator := new(big.Int).Mul(amountInWithFee, reserveOut)
	denominator := new(big.Int).Add(new(big.Int).Mul(reserveIn, e6), amountInWithFee)

	return new(big.Int).Div(numerator, denominator)
}

// uniswap v2.
func (pool *Pool) GetAmountOutV2AndPrice(zeroForOne bool, amountIn *big.Int, p0 bool) (*big.Int, decimal.Decimal) {
	var reserveIn, reserveOut *big.Int

	if pool.Reserve0.Cmp(bigZero) == 0 || pool.Reserve1.Cmp(bigZero) == 0 {
		return big.NewInt(0), decimal.Zero
	}

	if zeroForOne {
		reserveIn, reserveOut = pool.Reserve0, pool.Reserve1
	} else {
		reserveIn, reserveOut = pool.Reserve1, pool.Reserve0
	}

	// Calculate the amount out using the constant product formula
	amountInWithFee := new(big.Int).Mul(amountIn, new(big.Int).Sub(e6, big.NewInt(int64(pool.Fee)))) // nolint
	numerator := new(big.Int).Mul(amountInWithFee, reserveOut)
	denominator := new(big.Int).Add(new(big.Int).Mul(reserveIn, e6), amountInWithFee)

	amountOut := new(big.Int).Div(numerator, denominator)
	fee := d_e6.Sub(decimal.NewFromInt(int64(pool.Fee))).Div(d_e6) // nolint

	var (
		reserve0 = new(big.Int).Set(pool.Reserve0)
		reserve1 = new(big.Int).Set(pool.Reserve1)
		price    decimal.Decimal
	)

	if zeroForOne {
		reserve0 = reserve0.Add(reserve0, amountIn)
		reserve1 = reserve1.Sub(reserve1, amountOut)
	} else {
		reserve0 = reserve0.Sub(reserve0, amountOut)
		reserve1 = reserve1.Add(reserve1, amountIn)
	}

	p := decimal.NewFromBigInt(reserve1, 0).Div(decimal.NewFromBigInt(reserve0, 0))
	if p0 {
		price = p.Mul(fee)
	} else {
		price = decimal.NewFromInt(1).Div(p).Mul(fee)
	}

	return amountOut, price
}

// uniswap v3.
func (pool *Pool) getAmountOutV3(zeroForOne bool, amountIn *big.Int) *big.Int {
	if pool.Liquidity.Cmp(bigZero) == 0 {
		return big.NewInt(0)
	}

	sqrtPriceLimitX96 := pool.getSqrtPriceLimitX96(zeroForOne)
	amt0, amt1, _, _ := pool.swap(zeroForOne, amountIn, sqrtPriceLimitX96, true)

	if zeroForOne {
		return new(big.Int).Abs(amt1)
	}

	return new(big.Int).Abs(amt0)
}

// uniswap v3.
func (pool *Pool) getAmountOutV3AndPrice(zeroForOne bool, amountIn *big.Int, p0 bool) (*big.Int, decimal.Decimal) {
	if pool.Liquidity.Cmp(bigZero) == 0 {
		return big.NewInt(0), decimal.Zero
	}

	sqrtPriceLimitX96 := pool.getSqrtPriceLimitX96(zeroForOne)
	amt0, amt1, sqrtPriceAfter, _ := pool.swap(zeroForOne, amountIn, sqrtPriceLimitX96, true)

	var amountOut *big.Int

	if zeroForOne {
		amountOut = new(big.Int).Abs(amt1)
	} else {
		amountOut = new(big.Int).Abs(amt0)
	}

	price := sqrtPrice2Price(sqrtPriceAfter)

	fee := d_e6.Sub(decimal.NewFromInt(int64(pool.Fee))).Div(d_e6) // nolint
	if p0 {
		return amountOut, price.Mul(fee)
	}

	return amountOut, decimal.NewFromInt(1).Div(price).Mul(fee)
}

func (pool *Pool) getSqrtPriceLimitX96(zeroForOne bool) *big.Int {
	if zeroForOne {
		return new(big.Int).Add(MIN_SQRT_RATIO, bigOne)
	}

	return new(big.Int).Sub(MAX_SQRT_RATIO, bigOne)
}

// sqrtPrice2Price converts a Uniswap V3 sqrtPrice (Q64.96 format) to a regular price. sqrtRatioX96 ** 2 / 2 ** 192 = price.
func sqrtPrice2Price(sqrtPriceX96 *big.Int) decimal.Decimal {
	// Create a new big.Float to hold the result
	// price := new(big.Float)
	// Convert sqrtPriceX96 to big.Float
	// Divide by 2^192 to account for the Q64.96 format
	// divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(2), big.NewInt(192), nil))
	// price.Quo(price, divisor)
	// return price
	sqrtPriceFloat := decimal.NewFromBigInt(new(big.Int).Mul(sqrtPriceX96, sqrtPriceX96), 0)

	return sqrtPriceFloat.Div(decimal.NewFromInt(2).Pow(decimal.NewFromInt(192)))
}

func (pool *Pool) UpdateBestPrice() {
	switch {
	case pool.Typ.IsAMMVariety():
		pool.updatePoolV2BestPrice()

	case pool.Typ.IsCAMMVariety():
		pool.updatePoolV3BestPrice()

	default:
		logger.Fatal().Str("pool", pool.Address).Uint("type", uint(pool.Typ)).Msg("invalid pool type")
	}
}

func (pool *Pool) updatePoolV2BestPrice() {
	if pool.Stable && (pool.Typ == common.PoolTypeAeroAMM || pool.Typ == common.PoolTypeInfusionAMM) {
		if pool.Reserve0.Cmp(e6) <= 0 || pool.Reserve1.Cmp(e6) <= 0 {
			pool.price0 = &d_zero
			pool.price1 = &d_zero

			return
		}
		// logger.Info().Msgf("pool %s, reservers: %s %s", pool.Address, pool.Reserve0, pool.Reserve1)
		amountIn := new(big.Int).Div(pool.Reserve0, e5)
		amountOut := pool.getStableOut(true, amountIn)
		// logger.Info().Msgf("pool %s, amountIn: %s amountOut: %s", pool.Address, amountIn, amountOut)
		price0 := decimal.NewFromBigInt(amountOut, 0).Div(decimal.NewFromBigInt(amountIn, 0))
		pool.price0 = &price0

		amountIn = new(big.Int).Div(pool.Reserve1, e5)
		amountOut = pool.getStableOut(false, amountIn)
		// logger.Info().Msgf("pool %s, amountIn: %s amountOut: %s", pool.Address, amountIn, amountOut)
		price1 := decimal.NewFromBigInt(amountOut, 0).Div(decimal.NewFromBigInt(amountIn, 0))
		pool.price1 = &price1
	} else {
		if pool.Reserve0 == nil || pool.Reserve1 == nil {
			logger.Fatal().Msgf("pool %v %v reserve is null", pool.Vendor, pool.Address)
		}

		if pool.Reserve0.Cmp(bigZero) == 0 || pool.Reserve1.Cmp(bigZero) == 0 {
			pool.price0 = &decimal.Zero
			pool.price1 = &decimal.Zero

			return
		}

		price := decimal.NewFromBigInt(pool.Reserve1, 0).Div(decimal.NewFromBigInt(pool.Reserve0, 0))
		pool.updatePoolPrices(price)
	}
}

func (pool *Pool) updatePoolV3BestPrice() {
	if pool.Liquidity.Cmp(bigZero) > 0 {
		price := sqrtPrice2Price(pool.SqrtPriceX96)
		pool.updatePoolPrices(price)
	} else {
		ticks, currTickIdx := pool.GetTickListAndActiveRange()
		if len(ticks) == 0 {
			pool.price0 = &decimal.Zero
			pool.price1 = &decimal.Zero
			return
		}
		fee := d_e6.Sub(decimal.NewFromInt(int64(pool.Fee))).Div(d_e6) // nolint
		if pool.Tick <= ticks[0] {
			// idx: 0, cannot sell token0
			pool.price0 = &decimal.Zero
			p0 := decimal.NewFromFloat(PriceAtTick(ticks[0]))
			p1 := decimal.NewFromInt(1).Div(p0).Mul(fee)
			pool.price1 = &p1
		} else {
			if pool.Tick >= ticks[len(ticks)-1] {
				// p1 is zero, cannot sell token1
				pool.price1 = &decimal.Zero
				p0 := decimal.NewFromFloat(PriceAtTick(ticks[len(ticks)-1]))
				p0WithFee := p0.Mul(fee)
				pool.price0 = &p0WithFee
			} else {
				// normal
				p0 := decimal.NewFromFloat(PriceAtTick(ticks[currTickIdx]))
				p1 := decimal.NewFromFloat(PriceAtTick(ticks[currTickIdx+1]))
				p0WithFee := p0.Mul(fee)
				p1WithFee := p1.Mul(fee)
				pool.price0 = &p0WithFee
				pool.price1 = &p1WithFee
			}
		}
	}
}

func PriceAtTick(tick int) float64 {
	return math.Exp(float64(tick) * logBase1_0001)
}

// price0
// price1.
func (pool *Pool) updatePoolPrices(price decimal.Decimal) {
	if price.Cmp(decimal.Zero) == 0 {
		logger.Warn().Str("pool", pool.Address).
			Str("vender", pool.Vendor).
			Str("name", pool.name).
			Str("type", pool.Typ.String()).
			Msg("price is zero")

		pool.price0 = &decimal.Zero
		pool.price1 = &decimal.Zero

		return
	}

	fee := d_e6.Sub(decimal.NewFromInt(int64(pool.Fee))).Div(d_e6) // nolint
	p0 := price.Mul(fee)
	p1 := decimal.NewFromInt(1).Div(price).Mul(fee)
	pool.price0 = &p0
	pool.price1 = &p1
}
