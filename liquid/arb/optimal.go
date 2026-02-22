package arb

import (
	"math/big"

	"github.com/shopspring/decimal"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

var (
	e6   = big.NewInt(1000000)
	e12  = big.NewInt(1000000000000)
	e18  = big.NewInt(1000000000000000000) //lint:ignore U1000 Ignore unused function temporarily for debugging
	Q96  = new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
	Q192 = new(big.Int).Exp(big.NewInt(2), big.NewInt(192), nil)
)

// return: amtIn, amtOut
func (lpservice *LPService) GetOptimalArbAmount(p0, p1 *pool.Pool,
	zeroForOne bool,
	ratio decimal.Decimal,
) (amtIn *big.Int, amtOut *big.Int, profit *big.Int) {
	if p0.Stable || p1.Stable {
		// profitable, param, err := lpservice.optimizeMaxAmountIn(p0, p1, ratio, 15)
		// if !profitable || err != nil {
		// 	return bigZero, bigZero, bigZero
		// }

		// return param.bestAmt, param.bestOut, param.bestProfit
		return bigZero, bigZero, bigZero
	}

	if p0.Typ.IsAMMVariety() && p1.Typ.IsAMMVariety() {
		amtIn, amtOut = BestAmountV2ToV2(p0, p1, zeroForOne)
		if amtOut.Cmp(amtIn) > 0 {
			profit = new(big.Int).Sub(amtOut, amtIn)
		} else {
			profit = big.NewInt(0)
		}
		return
	} else if p0.Typ.IsAMMVariety() && p1.Typ.IsCAMMVariety() {
		// V2V3
		ticks := BuildTickLiquidity(p1, !zeroForOne, 5)
		amtIn, _, amtOut, profit, _ = CalcAmountInV2ToV3(zeroForOne,
			p0, p1,
			p0.Reserve0, p0.Reserve1,
			ticks,
			uint32(p0.Fee), uint32(p1.Fee), //nolint
			false) // nolint
		return amtIn, amtOut, profit
	} else if p0.Typ.IsCAMMVariety() && p1.Typ.IsAMMVariety() {
		// v3, v2
		ticks := BuildTickLiquidity(p0, zeroForOne, 5)
		amtIn, _, amtOut, profit, _ = CalcAmountInV3ToV2(zeroForOne,
			p0, p1,
			ticks,
			p1.Reserve0, p1.Reserve1,
			uint32(p0.Fee), uint32(p1.Fee), //nolint
			false) // nolint
		return amtIn, amtOut, profit
	} else if p0.Typ.IsCAMMVariety() && p1.Typ.IsCAMMVariety() {
		// v3, v3
		tls1 := BuildTickLiquidity(p0, zeroForOne, 5)
		tls2 := BuildTickLiquidity(p1, !zeroForOne, 5)
		amtIn, _, amtOut, profit, _ = CalcAmountInV3ToV3(zeroForOne,
			p0, p1,
			tls1, tls2,
			uint32(p0.Fee), uint32(p1.Fee), //nolint
			false) // nolint
		return amtIn, amtOut, profit
	} else {
		panic("invalid pool type")
	}
}

type BoundedProduct struct {
	R0    *big.Int
	R1    *big.Int
	V2    bool
	K     *big.Int // Liquidity
	Alpha *big.Int
	Beta  *big.Int
}

func BestAmountV2ToV2(
	p1, p2 *pool.Pool,
	// x1, y1, x2, y2 *big.Int,
	zeroForOne bool,
	// fee1, fee2 uint,
) (*big.Int, *big.Int) {
	x1 := p1.Reserve0
	y1 := p1.Reserve1
	x2 := p2.Reserve0
	y2 := p2.Reserve1
	if x1.Cmp(bigZero) == 0 || y1.Cmp(bigZero) == 0 || x2.Cmp(bigZero) == 0 || y2.Cmp(bigZero) == 0 {
		logger.Warn().Msgf("pool %v or %v reserves is zero: x1=%v y1=%v x2=%v y2=%v", p1.Address, p2.Address, x1, y1, x2, y2)
		return big.NewInt(0), big.NewInt(0)
	}

	fee1 := p1.Fee
	fee2 := p2.Fee
	if zeroForOne {
		return bestAmountV2ToV2(x1, y1, x2, y2, fee1, fee2)
	}

	return bestAmountV2ToV2(y1, x1, y2, x2, fee1, fee2)
}

// r1 = (e6-f1)
// r2 = (e6-f1)
// x -> y
// k = (y2*r1*e6 + y1*r1*r2)/e12
// a = k*k
// b = 2*k*y2*x1
// c = y2*x1*(y2*x1*e12 - r1*r2*y1*x2)/e12
// r = (-b + sqrt(b**2 - 4*a*c)) / (2*a).
func bestAmountV2ToV2(
	x1, y1, x2, y2 *big.Int,
	fee1, fee2 uint,
) (*big.Int, *big.Int) {
	f1 := big.NewInt(int64(fee1)) // nolint
	f2 := big.NewInt(int64(fee2)) // nolint
	r1 := new(big.Int).Sub(e6, f1)
	r2 := new(big.Int).Sub(e6, f2)
	// k1 = f1*f2*y1*x2
	// k2 = y2*x1
	// k3 = y2*f1 + y1*f1*f2
	k1 := mulmul(r1, r2, y1, x2)
	k1.Div(k1, e12)

	k2 := mulmul(y2, x1)
	k3 := new(big.Int).Add(mulmul(y2, r1, e6), mulmul(y1, r1, r2))
	k3.Div(k3, e12)

	// k := getK(y1, y2, r1, r2)
	// a := mulmul(k, k)
	// b := mulmul(big.NewInt(2), k, y2, x1)
	// c := mulmul(y2, x1,
	// 	new(big.Int).Sub(mulmul(y2, x1, e12), mulmul(r1, r2, x2, y1)))
	// c.Div(c, e12)
	// if c.Cmp(bigZero) >= 0 {
	// 	logger.Warn().Msg("c should less than 0")
	// 	return bigZero
	// }

	amt := quadraticOptimalByK(k1, k2, k3, nil)
	if amt.Cmp(bigZero) <= 0 {
		logger.Warn().Msgf("v2->v2 x1: %v y1: %v x2: %v y2: %v", x1, y1, x2, y2)
		return big.NewInt(0), big.NewInt(0)
	}
	// dx2 = dx1*k1 / (k2 + dx1*k3)
	return amt, calcAmountOut(k1, k2, k3, amt)
}

// func DoBestAmountV2ToV3(x1, y1 *big.Int,
// 	sp2 *big.Int, l2s []*big.Int, spLimits []*big.Int,
// 	fee1, fee2 uint,
// 	zeroForOne2 bool,
// ) {
// 	if zeroForOne2 {
// 		// bestAmountV2ToV3Tick()
// 	} else {
// 	}
// }

// v2 -> v3(0->1)
// dy1 = dx1 * f1 * y1 / (x1 + dx1*f1)
// spNext = L * Q96 * SP / (L * Q96 + dy1 * f2 * SP)
// dx2: L * (Sp-spNext) / Q96
// = (L/Q96) * (Sp - [L * Q96 * SP / (L * Q96 + dy1 * f2 * SP)])
// = (L/Q96) * Sp*dy1*f2*Sp/(L * Q96 + dy1 * f2 * SP)
// = (L*Sp*dy1*f2*Sp)/ (Q96* L * Q96 + dy1 * f2 * SP * Q96)
// = (dx1*L*Sp*f1*y1*f2*Sp) / [(x1 + dx1*f1) * (Q96 * L * Q96 + dy1 * f2 * SP * Q96)]
// = (dx1*L*Sp*f1*y1*f2*Sp) / [(x1 + dx1*f1) * Q96 * L * Q96 + (x1 + dx1*f1) * dy1 * f2 * SP * Q96]
// = (dx1*L*Sp*f1*y1*f2*Sp) / [(x1 + dx1*f1) * Q96 * L * Q96 + dx1 * f1 * y1 * f2 * SP * Q96]
// = (dx1*L*Sp*f1*y1*f2*Sp) / [x1 * Q96 * L * Q96 + dx1*(f1 * Q96 * L * Q96 +  f1 * Q96 * y1 * f2 * SP)]
//
// 0->1:
// k1 = L*Sp*f1*y1*f2*Sp
// k2 = x1 * Q192 * L
// k3 = (f1 * Q192 * L +  f1 * Q96 * y1 * f2 * SP)
// dx2 = dx1*k1 / (k2 + dx1*k3) = f/g
// f = dx1*k1
// g = k2 + dx1*k3
// dx2' = f'g - fg'/g**2
// f'g - fg' = k1 * (k2+dx1*k3) - dx1*k1* k3 = k1*k2
// dx2' = k1*k2 / (k2 + dx1*k3)**2 = 1
// (k2 + dx1*k3)**2 - k1*k2 = 0
// k3**2 * dx1**2 + 2*k2*k3*dx1 + k2*(k2-k1) = 0
//
// 1->0:
// k1 = f1*y1*f2*L*Q192
// k2 = x1*Sp*Sp*L
// k3 = f1*Sp*Sp*L + f1*Sp*y1*f2*Q96
//
// a = k3**2
// b = 2*k2*k3
// c = k2*(k2-k1)
// r = (-b + sqrt(b**2 - 4*a*c)) / (2*a)
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
// func bestAmountV2ToV3Tick(x1, y1 *big.Int, l2, sp2, n *big.Int,
// 	fee1, fee2 uint,
// 	zeroForOne1 bool,
// ) (*big.Int, *big.Int) {
// 	var k1, k2, k3, k4 *big.Int

// 	r1 := new(big.Int).Sub(e6, big.NewInt(int64(fee1))) // nolint
// 	r2 := new(big.Int).Sub(e6, big.NewInt(int64(fee2))) // nolint

// 	if !zeroForOne1 {
// 		k1 = new(big.Int).Div(mulmul(l2, sp2, r1, y1, r2, sp2), e12)
// 		k2 = mulmul(x1, Q192, l2)

// 		if k2.Cmp(k1) >= 0 {
// 			logger.Warn().Msgf("BestAmountV2ToV3(0->1): k2 should less than k1")
// 			return big.NewInt(0), big.NewInt(0)
// 		}

// 		k3 = new(big.Int).Add(mulmul(r1, Q192, l2, e6), mulmul(r1, r2, y1, Q96, sp2))
// 		k3.Div(k3, e12)
// 	} else {
// 		// y1, x1 = x1, y1
// 		// k1 = f1*y1*f2*L*Q192
// 		k1 = new(big.Int).Div(mulmul(r1, y1, r2, l2, Q192), e12)
// 		// k2 = x1*Sp*Sp*L
// 		k2 = mulmul(x1, sp2, l2, sp2)

// 		if k2.Cmp(k1) >= 0 {
// 			logger.Warn().Msgf("BestAmountV2ToV3(1->0): k2 should less than k1")
// 			return big.NewInt(0), big.NewInt(0)
// 		}

// 		// k3 = f1*Sp*Sp*L + f1*Sp*y1*f2*Q96
// 		k3 = new(big.Int).Add(mulmul(r1, sp2, l2, sp2, e6), mulmul(r1, r2, y1, Q96, sp2))
// 		k3.Div(k3, e12)
// 	}

// 	if n != nil && n.Cmp(bigZero) != 0 {
// 		// k2' = k2 - N*k3
// 		k2 = new(big.Int).Sub(k2, mulmul(n, k3))
// 	}

// 	amt := quadraticOptimalByK(k1, k2, k3, k4)

// 	return amt, calcAmountOut(k1, k2, k3, amt)
// }

//lint:ignore U1000 Ignore unused function temporarily for debugging
// func bestAmountV2ToV3Ticks(x1, y1 *big.Int, tls []*TickLiquidity,
// 	fee1, fee2 uint,
// 	zeroForOne1 bool,
// ) (*big.Int, *big.Int) {
// 	amtIn := big.NewInt(0)
// 	amtOut := big.NewInt(0)
// 	var midOut *big.Int = nil

// 	for _, tl := range tls {
// 		amt, out := bestAmountV2ToV3Tick(x1, y1, tl.Liquidity, tl.SPCurr, midOut, fee1, fee2, zeroForOne1)

// 		if out.Cmp(tl.Amount1) < 0 {
// 			amtIn.Add(amtIn, amt)
// 			amtOut.Add(amtOut, out)
// 			break
// 		} else {
// 			amtIn.Add(amtIn, tl.Amount0)
// 			amtOut.Add(amtOut, tl.Amount1)

// 			if midOut == nil {
// 				midOut = new(big.Int).Set(tl.Amount0)
// 			} else {
// 				midOut.Add(midOut, tl.Amount1)
// 			}
// 		}
// 	}

// 	return amtIn, amtOut
// }

// func BestAmountV3ToV2(l1, sp1 *big.Int, x2, y2 *big.Int,
// 	fee1, fee2 uint,
// 	zeroForOne1 bool,
// ) *big.Int {
// 	var k1, k2, k3 *big.Int

// 	r1 := new(big.Int).Sub(e6, big.NewInt(int64(fee1))) // nolint
// 	r2 := new(big.Int).Sub(e6, big.NewInt(int64(fee2))) // nolint

// 	if zeroForOne1 {
// 		// k1 = L*Sp*f1*SP*f2*x2
// 		// k2 = Q192*L*y2
// 		// k3 = (Q96*f1*SP*y2 + L*Sp*f1*SP*f2)
// 		k1 = mulmul(l1, sp1, sp1, r1, r2, x2)
// 		k1.Div(k1, e12)

// 		k2 = mulmul(Q192, l1, y2)
// 		k3 = new(big.Int).Add(mulmul(Q96, r1, sp1, y2, e6), mulmul(l1, sp1, r1, sp1, r2))
// 		k3.Div(k3, e12)
// 	} else {
// 		// k1 = Q192*f1*f2*y2*L
// 		// k2 = x2*Sp*Sp*L
// 		// k3 = f1*Q96*Sp*x2 + f1*Q192*f2*L
// 		k1 = mulmul(Q192, r1, r2, y2, l1)
// 		k1.Div(k1, e12)

// 		k2 = mulmul(x2, sp1, sp1, l1)
// 		k3 = new(big.Int).Add(mulmul(Q96, r1, sp1, x2, e6), mulmul(Q192, r1, r2, l1))
// 		k3.Div(k3, e12)
// 	}

// 	return quadraticOptimalByK(k1, k2, k3, nil)
// }

// func BestAmountV3ToV3(l1, sp1, l2, sp2 *big.Int,
// 	fee1, fee2 uint,
// 	zeroForOne1 bool,
// ) *big.Int {
// 	var k1, k2, k3 *big.Int

// 	r1 := new(big.Int).Sub(e6, big.NewInt(int64(fee1))) // nolint
// 	r2 := new(big.Int).Sub(e6, big.NewInt(int64(fee2))) // nolint

// 	if zeroForOne1 {
// 		// k1 = (L1*Sp1*f1*Sp1*f2*Q96*L2)
// 		// k2 = Sp2*Sp2*L2* Q96* L1
// 		// k3 = (Sp2*Sp2*L2*f1*SP1  + L1*Sp1*f1*Sp1*f2*Sp2)
// 		k1 = mulmul(l1, sp1, r1, sp1, r2, Q96, l2)
// 		k1.Div(k1, e12)

// 		k2 = mulmul(sp2, sp2, l2, Q96, l1)
// 		k3 = new(big.Int).Add(mulmul(sp2, sp2, l2, r1, sp1, e6), mulmul(l1, sp1, r1, sp1, r2, sp2))
// 		k3.Div(k3, e12)
// 	} else {
// 		// k1 = f1*L2*Sp2*f2*Sp2*L1
// 		// k2 = Sp1*Sp1*L2*L1
// 		// k3 = (f1*Q96*Sp1*L2 + L1*f1*f2*SP2*Q96)
// 		k1 = mulmul(r1, l2, sp2, r2, sp2, l1)
// 		k1.Div(k1, e12)

// 		k2 = mulmul(sp1, sp1, l2, l1)
// 		k3 = new(big.Int).Add(mulmul(r1, Q96, sp1, l2, e6), mulmul(l1, r1, r2, sp2, Q96))
// 		k3.Div(k3, e12)
// 	}

// 	return quadraticOptimalByK(k1, k2, k3, nil)
// }

// k = (y2*r1*e6 + y1*r1*r2)/e6**2
// func getK(y1, y2 *big.Int, r1, r2 *big.Int) *big.Int {
// 	k1 := mulmul(y2, r1, e6)
// 	k2 := mulmul(y1, r1, r2)

// 	return new(big.Int).Div(new(big.Int).Add(k1, k2), e12)
// }

// r = (-b + sqrt(b**2 - 4*a*c)) / (2*a).
func quadraticOptimal(a, b, c *big.Int) *big.Int {
	x := new(big.Int).Sub(mulmul(b, b), mulmul(a, c, big.NewInt(4)))
	if x.Cmp(bigZero) <= 0 {
		logger.Warn().Msg("sqrt less than 0")
		return big.NewInt(0)
	}

	r := new(big.Int).Sqrt(x)
	if r.Cmp(b) <= 0 {
		logger.Warn().Msgf("numerator less than 0: %v %v", r, b)
		return big.NewInt(0)
	}

	r.Sub(r, b)

	return r.Div(r, mulmul(a, big.NewInt(2)))
}

// y = (k1x-k4)/(k2+k3x)
// a = k3*k3
// b = 2*k2*k3
// c = k2*k2-k2*k1
// x = (-b + sqrt(b**2 - 4*a*c)) / (2*a)
func quadraticOptimalByK(k1, k2, k3, k4 *big.Int) *big.Int {
	a := mulmul(k3, k3)
	b := mulmul(big.NewInt(2), k2, k3)
	c := mulmul(k2, new(big.Int).Sub(k2, k1))
	if k4 != nil {
		// c = K2*k2 - k2*k1 - k3*k4
		c.Sub(c, mulmul(k3, k4))
	}

	return quadraticOptimal(a, b, c)
}

// dx2 = dx1*k1 / (k2 + dx1*k3)
func calcAmountOut(k1, k2, k3 *big.Int, dx1 *big.Int) *big.Int {
	return div(mulmul(dx1, k1), new(big.Int).Add(k2, mulmul(dx1, k3)))
}

// compute BoundedProduct
// k: liquidity
func ComputeBoundedProductAtTick(tickCurr, tickLower, tickUpper int32, k *big.Int) BoundedProduct {
	sqrtPriceLower := SqrtPriceX96AtTick(tickLower)
	sqrtPriceUpper := SqrtPriceX96AtTick(tickUpper)
	sqrtPriceCurr := SqrtPriceX96AtTick(tickCurr)
	amt0Desired, amt1Desired := pool.GetAmountsForLiquidity(sqrtPriceCurr, sqrtPriceLower, sqrtPriceUpper, k)
	alpha, beta := BoundedAlphaBeta(k, sqrtPriceLower, sqrtPriceUpper)

	return BoundedProduct{
		R0:    amt0Desired,
		R1:    amt1Desired,
		V2:    false,
		K:     k, // Liquidity
		Alpha: alpha,
		Beta:  beta,
	}
}

func ComputeBoundedProductAtTickByAmounts(tickCurr, tickLower, tickUpper int32, amt0, amt1 *big.Int) BoundedProduct {
	sqrtPriceLower := SqrtPriceX96AtTick(tickLower)
	sqrtPriceUpper := SqrtPriceX96AtTick(tickUpper)
	sqrtPriceCurr := SqrtPriceX96AtTick(tickCurr)

	logger.Info().Msgf("sqrtPriceLower: %v sqrtPriceUpper: %v sqrtPriceCurr: %v", sqrtPriceLower, sqrtPriceUpper, sqrtPriceCurr)
	if amt0 == nil {
		amt0 = big.NewInt(0)
	}

	if amt1 == nil {
		amt1 = big.NewInt(0)
	}

	k := pool.GetLiquidityForAmounts(sqrtPriceCurr, sqrtPriceLower, sqrtPriceUpper, amt0, amt1)
	amt0Desired, amt1Desired := pool.GetAmountsForLiquidity(sqrtPriceCurr, sqrtPriceLower, sqrtPriceUpper, k)
	logger.Info().Msgf("liquidity: %v amount0: %v  amount1: %v", k, amt0Desired, amt1Desired)
	alpha, beta := BoundedAlphaBeta(k, sqrtPriceLower, sqrtPriceUpper)
	logger.Info().Msgf("alpha: %v beta: %v", alpha, beta)
	r0, r1 := BoundedSqrtReserves(k, sqrtPriceCurr, alpha, beta)
	return BoundedProduct{
		R0:    r0,
		R1:    r1,
		V2:    false,
		K:     k, // Liquidity
		Alpha: alpha,
		Beta:  beta,
	}
}

func BoundedAlphaBeta(k *big.Int, sqrtPriceLower, sqrtPriceUpper *big.Int) (*big.Int, *big.Int) {
	// sqrt(k*Q96/sqrtPriceUpper)
	alpha := new(big.Int).Sqrt(div(mulmul(k, Q96), sqrtPriceUpper))
	// sqrt(k*sqrtPriceLower/Q96)
	logger.Info().Msgf("BoundedAlphaBeta: k=%v sqrtPriceLower=%v", k, sqrtPriceLower)
	beta := new(big.Int).Sqrt(div(mulmul(k, sqrtPriceLower), Q96))

	return alpha, beta
}

func BoundedSqrtReserves(k, sqrtPrice *big.Int, alpha, beta *big.Int) (*big.Int, *big.Int) {
	// sqrt(k*Q96/p)
	r0 := new(big.Int).Sqrt(div(mulmul(k, Q96), sqrtPrice))
	r1 := new(big.Int).Sqrt(div(mulmul(k, sqrtPrice), Q96))

	logger.Info().Msgf("r0=%v r1=%v", r0, r1)
	r0.Sub(r0, alpha)
	r1.Sub(r1, beta)
	return r0, r1
}
