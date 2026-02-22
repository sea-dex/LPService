package pool

import (
	"math/big"

	"github.com/shopspring/decimal"
	"starbase.ag/liquidity/pkg/logger"
)

var (
	b3  = big.NewInt(3)
	e5  = big.NewInt(100000)
	e6  = big.NewInt(1000000)
	e18 = toBigIntMust("1000000000000000000")

	d_e6   = decimal.NewFromInt(1000000)
	d_zero = decimal.NewFromInt(0)
)

// curve stable pool algothrim.

func _d(x0, y *big.Int) *big.Int {
	return new(big.Int).Add(mulDivE18(new(big.Int).Mul(b3, x0), mulDivE18(y, y)), mulDivE18(mulDivE18(x0, x0), x0))
}

func _f(x0, y *big.Int) *big.Int {
	a := mulDivE18(x0, y)
	b := new(big.Int).Add(mulDivE18(x0, x0), mulDivE18(y, y))

	return mulDivE18(a, b)
}

// func _mulDiv(x, y, z *big.Int) *big.Int {
// 	v := big.NewInt(1).Mul(x, y)
// 	v.Div(v, z)

// 	if v.Cmp(uint256Max) > 0 {
// 		panic("_mulDiv great than uint256 max")
// 	}

// 	return v
// }

func _k(x, y, decimals0, decimals1 *big.Int) *big.Int {
	// must be stable pool
	_x := mulDiv(x, e18, decimals0)
	_y := mulDiv(y, e18, decimals1)
	_a := mulDiv(_x, _y, e18)
	_b := new(big.Int).Add(mulDivE18(_x, _x), mulDivE18(_y, _y))

	return mulDivE18(_a, _b)
}

func mulDivE18(x, y *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(x, y), e18)
}

func _get_y(decimals0, decimals1, x0, xy, y *big.Int) *big.Int {
	// logger.Info().Msgf("_get_y params: x0=%v xy=%v y=%v decimals: %v, %v", x0, xy, y, decimals0, decimals1)
	for i := 0; i < 255; i++ {
		k := _f(x0, y)
		d := _d(x0, y)

		// logger.Info().Msgf("_get_y: i=%v k=%v d=%v", i, k, d)
		if d.Cmp(bigZero) == 0 {
			logger.Warn().Msg("d is zero")
			return big.NewInt(0)
		}

		if k.Cmp(xy) < 0 {
			dy := mulDiv(new(big.Int).Sub(xy, k), e18, d)
			if dy.Cmp(bigZero) == 0 {
				if k.Cmp(xy) == 0 {
					return y
				}

				yplus1 := new(big.Int).Add(y, bigOne)
				if _k(x0, yplus1, decimals0, decimals1).Cmp(xy) > 0 {
					return yplus1
				}

				dy = bigOne
			}

			y = new(big.Int).Add(y, dy)
		} else {
			dy := mulDiv(new(big.Int).Sub(k, xy), e18, d)
			if dy.Cmp(bigZero) == 0 {
				if k.Cmp(xy) == 0 || _f(x0, new(big.Int).Sub(y, bigOne)).Cmp(xy) < 0 {
					return y
				}

				dy = bigOne
			}

			y = new(big.Int).Sub(y, dy)
		}
	}

	// msg := fmt.Sprintf("_get_y: !y x0=%v xy=%v y=%v", x0, xy, y)
	logger.Warn().Msgf("_get_y: !y x0=%v xy=%v y=%v", x0, xy, y)
	// panic(msg)
	return nil
}

func getAeroV2AmountOut(pl *Pool, amountIn *big.Int, tokenIn string) *big.Int {
	// (uint256 _reserve0, uint256 _reserve1) = (reserve0, reserve1);
	// amountIn -= (amountIn * IPoolFactory(factory).getFee(address(this), stable)) / 10000; // remove fee from amount received
	// return _getAmountOut(amountIn, tokenIn, _reserve0, _reserve1);
	fee := mulDiv(big.NewInt(int64(pl.Fee)), amountIn, e6) // nolint
	amountIn = new(big.Int).Sub(amountIn, fee)

	return _getAeroV2AmountOut(pl, amountIn, tokenIn)
}

func _getAeroV2AmountOut(pl *Pool, amountIn *big.Int, tokenIn string) *big.Int {
	var reserveA, reserveB *big.Int

	if pl.Stable {
		if pl.decimals0 == nil || pl.decimals1 == nil {
			logger.Warn().Msg("pool decimals0 or decimals1 is nil")
			return big.NewInt(0)
		}

		xy := _k(pl.Reserve0, pl.Reserve1, pl.decimals0, pl.decimals1)
		_reserve0 := mulDiv(pl.Reserve0, e18, pl.decimals0)
		_reserve1 := mulDiv(pl.Reserve1, e18, pl.decimals1)

		if tokenIn == pl.Token0 {
			reserveA, reserveB = _reserve0, _reserve1
			amountIn = mulDiv(amountIn, e18, pl.decimals0)
		} else {
			reserveA, reserveB = _reserve1, _reserve0
			amountIn = mulDiv(amountIn, e18, pl.decimals1)
		}

		_y := _get_y(pl.decimals0, pl.decimals1, new(big.Int).Add(amountIn, reserveA), xy, reserveB)
		y := new(big.Int).Sub(reserveB, _y)

		if tokenIn == pl.Token0 {
			return mulDiv(y, pl.decimals1, e18)
		} else {
			return mulDiv(y, pl.decimals0, e18)
		}
	}

	if tokenIn == pl.Token0 {
		reserveA, reserveB = pl.Reserve0, pl.Reserve1
	} else {
		reserveA, reserveB = pl.Reserve1, pl.Reserve0
	}

	return new(big.Int).Div(new(big.Int).Mul(amountIn, reserveB), new(big.Int).Add(reserveA, amountIn))
}
