package arb

import (
	"math/big"

	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

// return: amtIn, amtMid, amtOut, profit, profitable

// v3 0->1:
// dx = dy*f*L*Q192 / (Sp*Sp*L + dy*Q96*f*Sp)
// v3 1->0:
// dy = dx*f*L*SP*SP / (Q192*L + dx*Q96*f*SP)

/*
1. v2 0->1, v3 1->0
dx2 = (dx1*f1*f2*L2*Q192*(y1-N) - N*x1*f2*L2*Q192) / (Sp2*Sp2*L2*x1 - N*x1*Q96*f2*Sp2 + dx1*f1*Sp2*(Sp2*L2 + (y1-N)*Q96*f2))
k1 = f1*f2*L2*Q192*(y1-N)
k2 = Sp2*Sp2*L2*x1 - N*x1*Q96*f2*Sp2
k3 = f1*Sp2*(Sp2*L2 + (y1-N)*Q96*f2)
k4 = N*x1*f2*L2*Q192

2. v2 1->0, v3 0->1
dy2 = (dy1*f1*f2*L2*SP2*SP2*(x1-N)-N*y1*f2*L2*SP2*SP2) / (Q192*L2*y1-N*y1*Q96*f2*Sp2 + dy1*f1*Q96(Q96*L2 + (x1-N)*f2*Sp2))
k1 = f1*f2*L2*SP2*SP2*(x1-N)
k2 = Q192*L2*y1-N*y1*Q96*f2*Sp2
k3 = f1*Q96(Q96*L2 + (x1-N)*f2*Sp2)
k4 = N*y1*f2*L2*SP2*SP2
*/
func CalcAmountInV2ToV3(
	zeroForOne1 bool,
	pl0, pl1 *pool.Pool,
	rx, ry *big.Int,
	tls []*TickLiquidity,
	fee1, fee2 uint32,
	debug bool,
) (*big.Int, *big.Int, *big.Int, *big.Int, bool) {
	n := new(big.Int)
	m := new(big.Int)
	f1 := new(big.Int).Sub(e6, big.NewInt(int64(fee1)))
	f2 := new(big.Int).Sub(e6, big.NewInt(int64(fee2)))

	var dx1, dy1, dx2, dy2 *big.Int = big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	if zeroForOne1 {
		for _, tl := range tls {
			dx1, dy1, dy2, dx2 = calcV2ToV3BoundZeroForOne(rx, ry, tl.Liquidity, tl.SPCurr, n, f1, f2, debug)
			if dx1.Cmp(bigZero) == 0 {
				return new(big.Int), new(big.Int), new(big.Int), new(big.Int), false
			}
			if debug {
				logger.Info().Msgf("calcV2ToV3Bound(0->1): dx1: %v dy2: %v dx2: %v n: %v", dx1, dy2, dx2, n)
			}

			if dy2.Cmp(tl.Amount1) > 0 {
				n.Add(n, tl.Amount1)
				m.Add(m, tl.Amount0)
				nx := getAmountOutV2(ry, rx, n, f1)
				if debug {
					logger.Info().Msgf("CalcAmountInV2ToV3(1->0): tick liquidity used up, modify dx1 %v to %v", dx1, nx)
				}
				dx1 = nx
				dx2 = bigZero
				continue
			} else {
				break
			}
		}
		amtOut := new(big.Int).Add(m, dx2)
		profit := new(big.Int).Sub(amtOut, dx1)
		return dx1, dy1, amtOut, profit, profit.Cmp(bigZero) > 0
	} else {
		// v2 1->0, v3 0->1
		for _, tl := range tls {
			dy1, dx1, dx2, dy2 = calcV2ToV3BoundOneForZero(rx, ry, tl.Liquidity, tl.SPCurr, n, f1, f2, debug)
			if dy1.Cmp(bigZero) == 0 {
				return new(big.Int), new(big.Int), new(big.Int), new(big.Int), false
			}

			if debug {
				logger.Info().Msgf("calcV2ToV3Bound(1->0): dy1: %v dx2: %v dy2: %v n: %v", dy1, dx2, dy2, n)
			}
			if dx2.Cmp(tl.Amount0) > 0 {
				n.Add(n, tl.Amount0)
				m.Add(m, tl.Amount1)
				ny := getAmountOutV2(rx, ry, n, f1)
				if debug {
					logger.Info().Msgf("CalcAmountInV2ToV3(0->1): tick liquidity used up, modify dy1 %v to %v", dy1, ny)
				}
				dy1 = ny
				dy2 = bigZero
				continue
			} else {
				break
			}
		}
		amtOut := new(big.Int).Add(m, dy2)
		profit := new(big.Int).Sub(amtOut, dy1)
		return dy1, dx1, amtOut, profit, profit.Cmp(bigZero) > 0
	}
}

// sqrt(ry*Q192/rx) > sp
func checkSqrtPrice(rx, ry *big.Int, sp2 *big.Int, greatThan bool) bool {
	sp1 := div(mulmul(ry, Q192), rx)
	sp1.Sqrt(sp1)

	return sp1.Cmp(sp2) > 0 == greatThan
}

// sqrt(ry/rx) > sp
func calcV2ToV3BoundZeroForOne(rx, ry, l, sp, n, f1, f2 *big.Int, debug bool) (*big.Int, *big.Int, *big.Int, *big.Int) {
	if debug {
		logger.Info().Msgf("v2->v3, 0->1 rx: %v ry: %v liquidity: %v sp: %v n: %v", rx, ry, l, sp, n)
	}

	if !checkSqrtPrice(rx, ry, sp, true) {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	}

	y_n := new(big.Int).Sub(ry, n)
	// k1 = f1*f2*L2*Q192*(y1-N)
	k1 := mulmul(f1, f2, l, Q192, y_n)
	k1.Div(k1, e12)
	// k2 = Sp2*Sp2*L2*x1 - N*x1*Q96*f2*Sp2
	k2 := mulmul(sp, sp, l, rx)
	if n.Cmp(bigZero) != 0 {
		k2.Sub(k2, div(mulmul(n, rx, Q96, f2, sp), e6))
	}

	// k3 = f1*Sp2*(Sp2*L2 + (y1-N)*Q96*f2)
	k31 := new(big.Int).Add(mulmul(sp, l, e6), mulmul(y_n, Q96, f2))
	k3 := mulmul(f1, sp, k31)
	k3.Div(k3, e12)

	// k4 = N*x1*f2*L2*Q192
	k4 := new(big.Int)
	if n.Cmp(bigZero) != 0 {
		k4 = mulmul(n, rx, f2, l, Q192)
		k4.Div(k4, e6)
	}

	dx1 := quadraticOptimalByK(k1, k2, k3, k4)
	if dx1.Cmp(bigZero) == 0 {
		logger.Warn().Msgf("v2->v3, 0->1 rx: %v ry: %v liquidity: %v sp: %v n: %v", rx, ry, l, sp, n)
		return big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	}

	dy1 := getAmountOutV2(rx, ry, dx1, f1)
	if dy1.Cmp(n) < 0 {
		logger.Warn().Msgf("dy1 less than n: dx1=%v dy1=%v n=%v", dx1, dy1, n)
	}
	dy2 := new(big.Int).Sub(dy1, n)
	// dx2 = dy2*f2*L2*Q192 / (Sp2*Sp2*L + dy2*Q96*f2*Sp2)
	dx2 := div(mulmul(dy2, f2, l, Q192), new(big.Int).Add(mulmul(sp, sp, l, e6), mulmul(dy2, Q96, f2, sp)))
	if debug {
		logger.Info().Msgf("dy2: %v dx2: %v L: %v sp: %v", dy2, dx2, l, sp)
	}

	return dx1, dy2, dy2, dx2
}

// k1 = f1*f2*L2*SP2*SP2*(x1-N)
// k2 = Q192*L2*y1-N*y1*Q96*f2*Sp2
// k3 = f1*Q96*(Q96*L2 + (x1-N)*f2*Sp2)
// k4 = N*y1*f2*L2*SP2*SP2
func calcV2ToV3BoundOneForZero(rx, ry, l, sp, n, f1, f2 *big.Int, debug bool) (*big.Int, *big.Int, *big.Int, *big.Int) {
	x_n := new(big.Int).Sub(rx, n)

	k1 := mulmul(f1, f2, l, sp, sp, x_n)
	k1.Div(k1, e12)

	k2 := mulmul(Q192, l, ry)
	if n.Cmp(bigZero) != 0 {
		k2.Sub(k2, div(mulmul(n, ry, Q96, f2, sp), e6))
	}

	k31 := new(big.Int).Add(mulmul(Q96, l, e6), mulmul(x_n, sp, f2))
	k3 := mulmul(f1, Q96, k31)
	k3.Div(k3, e12)

	k4 := new(big.Int)
	if n.Cmp(bigZero) != 0 {
		k4 = mulmul(n, ry, f2, l, sp, sp)
		k4.Div(k4, e6)
	}

	dy1 := quadraticOptimalByK(k1, k2, k3, k4)
	if dy1.Cmp(bigZero) == 0 {
		logger.Warn().Msgf("v2->v3, 1->0 rx: %v ry: %v liquidity: %v sp: %v n: %v", rx, ry, l, sp, n)
		return big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	}

	dx1 := getAmountOutV2(ry, rx, dy1, f1)
	if dx1.Cmp(n) < 0 {
		logger.Warn().Msgf("dx1 less than n: dx1=%v dy1=%v n=%v", dx1, dy1, n)
	}
	dx2 := new(big.Int).Sub(dx1, n)
	// dy2 = dx2*f2*L2*SP2*SP2 / (Q192*L2 + dx2*Q96*f2*SP2)
	dy2 := div(mulmul(dx2, f2, l, sp, sp), new(big.Int).Add(mulmul(Q192, l, e6), mulmul(dx2, Q96, f2, sp)))

	return dy1, dx1, dx2, dy2
}

// return: amtIn, amtMid, amtOut, profit, profitable
// first swap v3 pool, then swap v2 pool
func CalcAmountInV3ToV2(
	zeroForOne1 bool,
	pl0, pl1 *pool.Pool,
	tls []*TickLiquidity,
	rx, ry *big.Int,
	fee1, fee2 uint32,
	debug bool,
) (*big.Int, *big.Int, *big.Int, *big.Int, bool) {
	n := new(big.Int)
	m := new(big.Int)
	f1 := new(big.Int).Sub(e6, big.NewInt(int64(fee1)))
	f2 := new(big.Int).Sub(e6, big.NewInt(int64(fee2)))

	var dx1, dy1, dx2, dy2 *big.Int = big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	if zeroForOne1 {
		for _, tl := range tls {
			dx1, dy1 = calcV3BoundToV2ZeroForOne(rx, ry, tl.Liquidity, tl.SPCurr, n, f1, f2, debug)
			if dx1.Cmp(bigZero) == 0 {
				return new(big.Int), new(big.Int), new(big.Int), new(big.Int), false
			}
			if debug {
				logger.Info().Msgf("calcV3BoundToV2(0->1): dx1: %v dy1: %v", dx1, dy1)
			}

			if dy1.Cmp(tl.Amount1) > 0 {
				n.Add(n, tl.Amount1)
				m.Add(m, tl.Amount0)
				// 1->0
				// nx := getAmountOutV2(ry, rx, n, f1)
				if debug {
					logger.Info().Msgf("calcV3BoundToV2(0->1): tick liquidity used up, dx1: %v cumulate dy: %v", dx1, n)
				}
				continue
			} else {
				dy2 = new(big.Int).Add(n, dy1)
				dx2 = getAmountOutV2(ry, rx, dy2, f2)
				if debug {
					logger.Info().Msgf("dy1: %v dy2: %v dx2: %v", dy1, dy2, dx2)
				}
				break
			}
		}
		dx1.Add(dx1, m)
		profit := new(big.Int).Sub(dx2, dx1)
		return dx1, dy2, dx2, profit, profit.Cmp(bigZero) > 0
	} else {
		// v3 1->0, v2 0->1
		for _, tl := range tls {
			dy1, dx1 = calcV3BoundToV2OneForZero(rx, ry, tl.Liquidity, tl.SPCurr, n, f1, f2, debug)
			if dy1.Cmp(bigZero) == 0 {
				return new(big.Int), new(big.Int), new(big.Int), new(big.Int), false
			}

			if debug {
				logger.Info().Msgf("calcV3BoundToV2(1->0): dy1: %v dx1: %v n: %v", dy1, dx1, n)
			}
			if dx1.Cmp(tl.Amount0) > 0 {
				n.Add(n, tl.Amount0)
				m.Add(m, tl.Amount1)
				if debug {
					logger.Info().Msgf("calcV3BoundToV2(1->0): tick liquidity used up, dy1 %v cumulative dx %v cumulative dy: %v", dy1, n, m)
				}
				continue
			} else {
				dx2 = new(big.Int).Add(n, dx1)
				dy2 = getAmountOutV2(rx, ry, dx2, f2)
				if debug {
					logger.Info().Msgf("v3 1->0: dy1: %v dx1: %v v2 0->1 dx2: %v dy2: %v", dy1, dx1, dx2, dy2)
				}
				break
			}
		}
		dy1.Add(dy1, m)
		profit := new(big.Int).Sub(dy2, dy1)
		return dy1, dx2, dy2, profit, profit.Cmp(bigZero) > 0
	}
}

// k1 = (L1*SP1+Q96*N)*f1*f2*SP1*x2
// k2 = y2*Q192*L1+N*f2*Q192*L1
// k3 = f1*SP1*(Q96*y2 + f2*(L1*SP1 + Q96*N))
// k4 = -Q192*L1*N*f2*x2
func calcV3BoundToV2ZeroForOne(rx, ry, l, sp, n, f1, f2 *big.Int, debug bool) (*big.Int, *big.Int) {
	if debug {
		logger.Info().Msgf("rx: %v ry: %v liquidity: %v sp: %v n: %v", rx, ry, l, sp, n)
	}
	l_sp := new(big.Int).Add(mulmul(l, sp), mulmul(Q96, n))
	k1 := mulmul(f1, f2, sp, rx, l_sp)
	k1.Div(k1, e12)

	k2 := mulmul(Q192, l, ry)
	if n.Cmp(bigZero) != 0 {
		k2.Add(k2, div(mulmul(n, f2, Q192, l), e6))
	}

	k31 := new(big.Int).Add(mulmul(Q96, ry, e6), mulmul(l_sp, f2))
	k3 := mulmul(f1, sp, k31)
	k3.Div(k3, e12)

	k4 := new(big.Int)
	if n.Cmp(bigZero) != 0 {
		k4 = new(big.Int).Neg(mulmul(n, rx, f2, l, Q192))
		k4.Div(k4, e6)
	}

	dx1 := quadraticOptimalByK(k1, k2, k3, k4)
	if dx1.Cmp(bigZero) == 0 {
		logger.Warn().Msgf("v3->v2, 0->1 rx: %v ry: %v liquidity: %v sp: %v n: %v", rx, ry, l, sp, n)
		return big.NewInt(0), big.NewInt(0)
	}

	dy1 := getAmountOutV3Tick(true, l, sp, dx1, f1)
	// dy1 := getAmountOutV2(rx, ry, dx1, f1)
	// if dy1.Cmp(n) < 0 {
	// 	logger.Warn().Msgf("dy1 less than n: dx1=%v dy1=%v n=%v", dx1, dy1, n)
	// }
	// dy2 := new(big.Int).Sub(dy1, n)
	// // dx2 = dy2*f2*L2*Q192 / (Sp2*Sp2*L + dy2*Q96*f2*Sp2)
	// dx2 := div(mulmul(dy2, f2, l, Q192), new(big.Int).Add(mulmul(sp, sp, l, e6), mulmul(dy2, Q96, f2, sp)))
	// logger.Info().Msgf("dy2: %v dx2: %v", dy2, dx2)

	return dx1, dy1
}

// k1 = f1*Q96*f2*y2*(L1*Q96+Sp1*N)
// k2 = x2*Sp1*Sp1*L1+N*f2*Sp1*Sp1*L1
// k3 = Q96*f1*(Sp1*x2+L1*Q96*f2+N*f2*Sp1)
// k4 = -N*Sp1*Sp1*L1*f2*y2
func calcV3BoundToV2OneForZero(rx, ry, l, sp, n, f1, f2 *big.Int, debug bool) (*big.Int, *big.Int) {
	if debug {
		logger.Info().Msgf("rx: %v ry: %v liquidity: %v sp: %v n: %v", rx, ry, l, sp, n)
	}
	l_sp := new(big.Int).Add(mulmul(l, Q96), mulmul(sp, n))
	k1 := mulmul(f1, f2, Q96, ry, l_sp)
	k1.Div(k1, e12)

	k2 := mulmul(sp, sp, l, rx)
	if n.Cmp(bigZero) != 0 {
		k2.Add(k2, div(mulmul(n, f2, sp, sp, l), e6))
	}

	k31 := new(big.Int).Add(mulmul(sp, rx, e6), mulmul(l_sp, f2))
	k3 := mulmul(f1, Q96, k31)
	k3.Div(k3, e12)

	k4 := new(big.Int)
	if n.Cmp(bigZero) != 0 {
		k4 = new(big.Int).Neg(mulmul(n, ry, f2, l, sp, sp))
		k4.Div(k4, e6)
	}

	dy1 := quadraticOptimalByK(k1, k2, k3, k4)
	if dy1.Cmp(bigZero) == 0 {
		logger.Warn().Msgf("v3->v2, 1->0 rx: %v ry: %v liquidity: %v sp: %v n: %v", rx, ry, l, sp, n)
		return big.NewInt(0), big.NewInt(0)
	}

	dx1 := getAmountOutV3Tick(false, l, sp, dy1, f1)
	// dy1 := getAmountOutV2(rx, ry, dx1, f1)
	// if dy1.Cmp(n) < 0 {
	// 	logger.Warn().Msgf("dy1 less than n: dx1=%v dy1=%v n=%v", dx1, dy1, n)
	// }
	// dy2 := new(big.Int).Sub(dy1, n)
	// // dx2 = dy2*f2*L2*Q192 / (Sp2*Sp2*L + dy2*Q96*f2*Sp2)
	// dx2 := div(mulmul(dy2, f2, l, Q192), new(big.Int).Add(mulmul(sp, sp, l, e6), mulmul(dy2, Q96, f2, sp)))
	// logger.Info().Msgf("dy2: %v dx2: %v", dy2, dx2)

	return dy1, dx1
}

func CalcAmountInV3ToV3(
	zeroForOne1 bool,
	pl0, pl1 *pool.Pool,
	tls1 []*TickLiquidity,
	tls2 []*TickLiquidity,
	fee1, fee2 uint32,
	debug bool,
) (*big.Int, *big.Int, *big.Int, *big.Int, bool) {
	if len(tls1) == 0 || len(tls2) == 0 {
		// logger.Warn().Msgf("either pool1 %v or pool2 %v has no liquidity: %d %d", pl0.Address, pl1.Address, len(tls1), len(tls2))
		return bigZero, bigZero, bigZero, bigZero, false
	}

	idx1 := 0
	idx2 := 0
	f1 := new(big.Int).Sub(e6, big.NewInt(int64(fee1)))
	f2 := new(big.Int).Sub(e6, big.NewInt(int64(fee2)))

	if zeroForOne1 {
		cumX1 := big.NewInt(0) // dx1
		cumX2 := big.NewInt(0) // dy2
		cumY1 := big.NewInt(0) // sum(dy1)
		cumY2 := big.NewInt(0) // sum(dy2)

		// v3 0->1, v3 1->0
		for {
			if idx1 >= len(tls1) || idx2 >= len(tls2) {
				logger.Warn().Msgf("all ticks used up: len(ticks1)=%d len(ticks2)=%d", len(tls1), len(tls2))
				return bigZero, bigZero, bigZero, bigZero, false
			}
			tl1 := tls1[idx1]
			tl2 := tls2[idx2]
			if debug {
				logger.Info().Msgf("v3 0->1, v3 1->0: idx1=%d idx2=%d tick1: [%v %v] tick2: [%v %v]",
					idx1, idx2, SqrtPriceX96ToPrice(tl1.SPCurr), SqrtPriceX96ToPrice(tl1.SPLimit),
					SqrtPriceX96ToPrice(tl2.SPCurr), SqrtPriceX96ToPrice(tl2.SPLimit))
			}
			dx1 := calcV3BoundZeroForOne(tl1, tl2, cumY1, cumY2, f1, f2, debug)
			if dx1.Cmp(bigZero) <= 0 {
				logger.Warn().Msgf("invalid dx1: %v sell price: %v buy price: %v", dx1, tl1.SPCurr, tl2.SPCurr)
				return bigZero, bigZero, bigZero, bigZero, false
			}
			if dx1.Cmp(tl1.Amount0) >= 0 {
				if debug {
					logger.Info().Msgf("tick1 liquidity used up: price: [%v %v], liquidity: %v",
						SqrtPriceX96ToPrice(tl1.SPCurr), SqrtPriceX96ToPrice(tl1.SPLimit), tl1.Liquidity)
				}
				cumX1.Add(cumX1, tl1.Amount0)
				cumY1.Add(cumY1, tl1.Amount1)
				idx1++
				continue
			}
			dy1 := getAmountOutV3Tick(true, tl1.Liquidity, tl1.SPCurr, dx1, f1)
			dy2 := bAddSub(dy1, cumY1, cumY2)
			if dy2.Cmp(bigZero) <= 0 {
				logger.Warn().Msgf("dy2 less than 0: dy1=%v cumY1=%v cumY2=%v", dy1, cumY1, cumY2)
				return bigZero, bigZero, bigZero, bigZero, false
			}
			if dy2.Cmp(tl2.Amount1) >= 0 {
				if debug {
					logger.Info().Msgf("tick2 liquidity used up: price: [%v %v], liquidity: %v",
						SqrtPriceX96ToPrice(tl2.SPCurr), SqrtPriceX96ToPrice(tl2.SPLimit), tl2.Liquidity)
				}
				cumY2.Add(cumY2, tl2.Amount1)
				cumX2.Add(cumX2, tl2.Amount0)
				idx2++
				continue
			}
			dx2 := getAmountOutV3Tick(false, tl2.Liquidity, tl2.SPCurr, dy2, f2)
			if debug {
				logger.Info().Msgf("swap result: v3(0->1): dx1=%v dy1=%v, v3(1->0): dy2=%v dx2=%v",
					dx1, dy1, dy2, dx2)
			}
			// complete
			dx1.Add(dx1, cumX1)
			dy1.Add(dy1, cumY1)
			dx2.Add(dx2, cumX2)
			dy2.Add(dy2, cumY2)
			profit := new(big.Int).Sub(dx2, dx1)
			return dx1, dy1, dx2, profit, profit.Cmp(bigZero) > 0
		}
	} else {
		cumX1 := big.NewInt(0) // dx1
		cumX2 := big.NewInt(0) // dy2
		cumY1 := big.NewInt(0) // sum(dy1)
		cumY2 := big.NewInt(0) // sum(dy2)

		// v3 1->0, v3 0->1
		for {
			if idx1 >= len(tls1) || idx2 >= len(tls2) {
				logger.Warn().Msgf("all ticks used up: len(ticks1)=%d len(ticks2)=%d", len(tls1), len(tls2))
				return bigZero, bigZero, bigZero, bigZero, false
			}
			tl1 := tls1[idx1]
			tl2 := tls2[idx2]
			if debug {
				logger.Info().Msgf("v3 1->0, v3 0->1: idx1=%d idx2=%d tick1: [%v %v] tick2: [%v %v]",
					idx1, idx2, SqrtPriceX96ToPrice(tl1.SPCurr), SqrtPriceX96ToPrice(tl1.SPLimit),
					SqrtPriceX96ToPrice(tl2.SPCurr), SqrtPriceX96ToPrice(tl2.SPLimit))
			}
			dy1 := calcV3BoundOneForZero(tl1, tl2, cumX1, cumX2, f1, f2, debug)
			if dy1.Cmp(bigZero) <= 0 {
				logger.Warn().Msgf("invalid dy1: %v", dy1)
				return bigZero, bigZero, bigZero, bigZero, false
			}
			if dy1.Cmp(tl1.Amount1) >= 0 {
				if debug {
					logger.Info().Msgf("tick1 liquidity used up: price: [%v %v], liquidity: %v",
						SqrtPriceX96ToPrice(tl1.SPCurr), SqrtPriceX96ToPrice(tl1.SPLimit), tl1.Liquidity)
				}
				cumX1.Add(cumX1, tl1.Amount0)
				cumY1.Add(cumY1, tl1.Amount1)
				idx1++
				continue
			}
			dx1 := getAmountOutV3Tick(false, tl1.Liquidity, tl1.SPCurr, dy1, f1)
			dx2 := bAddSub(dx1, cumX1, cumX2)
			if dx2.Cmp(bigZero) <= 0 {
				logger.Warn().Msgf("dx2 less than 0: dx1=%v cumX1=%v cumX2=%v", dx1, cumX1, cumX2)
				return bigZero, bigZero, bigZero, bigZero, false
			}
			if dx2.Cmp(tl2.Amount0) >= 0 {
				if debug {
					logger.Info().Msgf("tick2 liquidity used up: price: [%v %v], liquidity: %v",
						SqrtPriceX96ToPrice(tl2.SPCurr), SqrtPriceX96ToPrice(tl2.SPLimit), tl2.Liquidity)
				}
				cumY2.Add(cumY2, tl2.Amount1)
				cumX2.Add(cumX2, tl2.Amount0)
				idx2++
				continue
			}
			dy2 := getAmountOutV3Tick(true, tl2.Liquidity, tl2.SPCurr, dx2, f2)
			if debug {
				logger.Info().Msgf("swap result: v3(1->0): dx1=%v dy1=%v, v3(0->1): dy2=%v dx2=%v",
					dx1, dy1, dy2, dx2)
			}
			// complete
			dx1.Add(dx1, cumX1)
			dy1.Add(dy1, cumY1)
			dx2.Add(dx2, cumX2)
			dy2.Add(dy2, cumY2)
			profit := new(big.Int).Sub(dy2, dy1)
			return dy1, dx1, dy2, profit, profit.Cmp(bigZero) > 0
		}
	}
	// logger.Warn().Msg("should not go here")
	// return bigZero, bigZero, bigZero, bigZero, bigZero, false
}

// dx2 = (dx1*f1*SP1*f2*L2*Q192*(L1*SP1+cumY*Q96) + cumY*Q192*L1*f2*L2*Q192) /
// ( (Sp2*L2+cumY*Q96*f2)*Q192*L1*Sp2 + dx1*f1*SP1*Q96*Sp2*(Sp2*L2 + f2*(L1*SP1+cumY*Q96)) )
// k1 = f1*SP1*f2*L2*Q192*(L1*SP1+cumY*Q96)
// k2 = (Sp2*L2+cumY*Q96*f2)*Q192*L1*Sp2
// k3 = f1*SP1*Q96*Sp2*(Sp2*L2 + f2*(L1*SP1+cumY*Q96))
// k4 = -cumY*Q192*L1*f2*L2*Q192
func calcV3BoundZeroForOne(tl1, tl2 *TickLiquidity, cumY1, cumY2, f1, f2 *big.Int, debug bool) *big.Int {
	l1, sp1 := tl1.Liquidity, tl1.SPCurr
	l2, sp2 := tl2.Liquidity, tl2.SPCurr
	cumY := new(big.Int).Sub(cumY1, cumY2)
	l_sp := new(big.Int).Add(mulmul(l1, sp1), mulmul(cumY, Q96))
	k1 := mulmul(f1, sp1, f2, l2, Q192, l_sp)
	k1.Div(k1, e12)

	k20 := new(big.Int).Div(mulmul(cumY, Q96, f2), e6) // cumY*Q96*f2
	k2 := mulmul(Q192, l1, sp2, new(big.Int).Add(mulmul(sp2, l2), k20))

	// f2*(L1*SP1+cumY*Q96)
	k30 := mulmul(f2, new(big.Int).Add(mulmul(l1, sp1), mulmul(cumY, Q96)))
	k3 := mulmul(f1, sp1, Q96, sp2, new(big.Int).Add(mulmul(sp2, l2, e6), k30))
	k3.Div(k3, e12)

	k4 := big.NewInt(0)
	if cumY.Cmp(bigZero) != 0 {
		k4 = mulmul(cumY, Q192, l1, f2, l2, Q192)
		k4.Div(k4, e6)
		k4.Neg(k4)
	}
	dx1 := quadraticOptimalByK(k1, k2, k3, k4)
	if dx1.Cmp(bigZero) <= 0 {
		logger.Warn().Msgf("v3->v3, 0->1 l1: %v sp1: %v l2: %v sp2: %v cumY1: %v cumY2: %v", l1, sp1, l2, sp2, cumY, cumY2)
		return bigZero
	}

	return dx1
}

// dy2 = (dy1*f1*Q96*f2*L2*SP2*SP2*(L1*Q96+cumX*Sp1) + cumX*Sp1*Sp1*L1*f2*L2*SP2*SP2) /
// (Sp1*Sp1*L1*Q96*(Q96*L2+cumX*f2*SP2) + dy1*f1*Q192*(L2*Q96*Sp1 + f2*SP2*(L1*Q96+cumX*Sp1)))
// k1 = f1*Q96*f2*L2*SP2*SP2*(L1*Q96+cumX*Sp1)
// k2 = Sp1*Sp1*L1*Q96*(Q96*L2+cumX*f2*SP2)
// k3 = f1*Q192*(L2*Q96*Sp1 + f2*SP2*(L1*Q96+cumX*Sp1))
// k4 = -cumX*Sp1*Sp1*L1*f2*L2*SP2*SP2
func calcV3BoundOneForZero(tl1, tl2 *TickLiquidity, cumX1, cumX2, f1, f2 *big.Int, debug bool) *big.Int {
	l1, sp1 := tl1.Liquidity, tl1.SPCurr
	l2, sp2 := tl2.Liquidity, tl2.SPCurr
	cumX := new(big.Int).Sub(cumX1, cumX2)
	l_sp := new(big.Int).Add(mulmul(l1, Q96), mulmul(cumX, sp1))
	k1 := mulmul(f1, Q96, f2, l2, sp2, sp2, l_sp)
	k1.Div(k1, e12)

	k20 := new(big.Int).Div(mulmul(cumX, sp2, f2), e6) // cumX*sp2*f2
	k2 := mulmul(sp1, sp1, l1, Q96, new(big.Int).Add(mulmul(Q96, l2), k20))

	// f2*SP2*(L1*Q96+cumX*Sp1)
	k30 := mulmul(f2, sp2, new(big.Int).Add(mulmul(l1, Q96), mulmul(cumX, sp1)))
	k3 := mulmul(f1, Q192, new(big.Int).Add(mulmul(l2, Q96, sp1, e6), k30))
	k3.Div(k3, e12)

	k4 := big.NewInt(0)
	if cumX.Cmp(bigZero) != 0 {
		k4 = mulmul(cumX, sp1, sp1, l1, f2, l2, sp2, sp2)
		k4.Div(k4, e6)
		k4.Neg(k4)
	}
	dy1 := quadraticOptimalByK(k1, k2, k3, k4)
	if dy1.Cmp(bigZero) <= 0 {
		logger.Warn().Msgf("v3->v3, 1->0 l1: %v sp1: %v l2: %v sp2: %v", l1, sp1, l2, sp2)
		return bigZero
	}

	return dy1
}
