package arb

import (
	"math/big"

	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

type TickLiquidity struct {
	Amount0   *big.Int // max amountIn
	Amount1   *big.Int // max amountOut
	TickCurr  int
	TickNext  int
	SPCurr    *big.Int
	SPLimit   *big.Int
	Liquidity *big.Int
}

func BuildTickLiquidity(p *pool.Pool, toLeft bool, n uint32) []*TickLiquidity {
	ticks, currTickIdx := p.GetTickListAndActiveRange()

	tls := []*TickLiquidity{}
	if len(ticks) == 0 {
		return tls
	}

	var (
		amt0     *big.Int
		amt1     *big.Int
		tickCurr int
		tickNext int
		spCurr   *big.Int
		spLimit  *big.Int
	)

	// println("ticks: ", ticks, " idx:", currTickIdx, " current tick:", p.Tick)

	idx := currTickIdx
	liquidity := p.Liquidity
	if toLeft {
		if idx >= len(ticks) {
			idx = len(ticks) - 2
			// liquidity should be 0
			tick := p.Ticks[ticks[idx+1]]
			liquidity = new(big.Int).Sub(liquidity, tick.LiquidityNet)
		}

		// prices go down
		for i := 0; i < int(n) && idx >= 0; {
			if liquidity.Cmp(bigZero) > 0 {
				// tickLower := ticks[idx]
				if idx == currTickIdx {
					tickCurr = p.Tick
					spCurr = p.SqrtPriceX96
				} else {
					tickCurr = ticks[idx+1]
					spCurr = SqrtPriceX96AtTick(int32(tickCurr)) // nolint
				}

				tickNext = ticks[idx]
				spLimit = SqrtPriceX96AtTick(int32(tickNext))                  // nolint
				amt0 = pool.GetAmount0Delta(spLimit, spCurr, liquidity, false) // token0
				amt1 = pool.GetAmount1Delta(spLimit, spCurr, liquidity, false) // token1

				tls = append(tls, &TickLiquidity{
					Amount0:   amt0,
					Amount1:   amt1,
					TickCurr:  tickCurr,
					TickNext:  tickNext,
					SPCurr:    spCurr,
					SPLimit:   spLimit,
					Liquidity: new(big.Int).Set(liquidity),
				})
				i++
			}

			tick := p.Ticks[ticks[idx]]
			liquidity = new(big.Int).Sub(liquidity, tick.LiquidityNet)

			idx--
		}
	} else {
		for i := 0; i < int(n) && idx < len(ticks)-1; {
			tick := p.Ticks[ticks[idx+1]]
			if liquidity.Cmp(bigZero) > 0 {
				// tickUpper := ticks[idx+1]
				if idx == currTickIdx {
					tickCurr = p.Tick
					spCurr = p.SqrtPriceX96
				} else {
					tickCurr = ticks[idx]
					spCurr = SqrtPriceX96AtTick(int32(tickCurr)) // nolint
				}
				tickNext = ticks[idx+1]
				spLimit = SqrtPriceX96AtTick(int32(tickNext))                  // nolint
				amt0 = pool.GetAmount0Delta(spCurr, spLimit, liquidity, false) // token y
				amt1 = pool.GetAmount1Delta(spCurr, spLimit, liquidity, false) // token x

				tls = append(tls, &TickLiquidity{
					Amount0:   amt0,
					Amount1:   amt1,
					TickCurr:  tickCurr,
					TickNext:  tickNext,
					SPCurr:    spCurr,
					SPLimit:   spLimit,
					Liquidity: new(big.Int).Set(liquidity),
				})
				i++
			}

			if i == 0 && p.Tick < ticks[0] {
				tick := p.Ticks[ticks[0]]
				liquidity = new(big.Int).Add(liquidity, tick.LiquidityNet)
				idx++
			} else {
				liquidity = new(big.Int).Add(liquidity, tick.LiquidityNet)
				idx++
			}
		}
	}

	return tls
}

func printBoundedTick(tl *TickLiquidity) {
	priceCurr := SqrtPriceX96ToPrice(tl.SPCurr)
	priceLimit := SqrtPriceX96ToPrice(tl.SPLimit)
	logger.Info().Msgf("tick: spCurr: %v %v spLimit: %v %v", priceCurr, tl.SPCurr, priceLimit, tl.SPLimit)
	logger.Info().Msgf("Liquidity: %v amount0: %v amount1: %v", tl.Liquidity, tl.Amount0, tl.Amount1)
}

func printBoundedTicks(tls []*TickLiquidity) {
	logger.Info().Msgf("ticks: %d", len(tls))
	for _, tl := range tls {
		printBoundedTick(tl)
	}
}
