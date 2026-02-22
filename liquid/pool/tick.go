package pool

import (
	"fmt"
	"math/big"
	"slices"

	"starbase.ag/liquidity/pkg/logger"
)

const (
	MIN_TICK = int(-887272)
	MAX_TICK = -MIN_TICK
)

// TickInfo CAMM pool tick info.
type TickInfo struct {
	Tick           int      `json:"tick"`
	LiquidityNet   *big.Int `json:"liquidityNet"`
	LiquidityGross *big.Int `json:"liquidityGross"`
}

func (pool *Pool) printInfos() {
	fmt.Printf("Pool %s Info:\n", pool.Address)
	fmt.Printf("  tick: %d\n", pool.Tick)
	fmt.Printf("  fee: %d\n", pool.Fee)
	fmt.Printf("  tickSpacing: %d\n", pool.TickSpacing)
	fmt.Printf("  sqrtPrice: %s\n", pool.SqrtPriceX96.String())
	fmt.Printf("  liquidity: %s\n", pool.Liquidity.String())

	for tick, info := range pool.Ticks {
		fmt.Printf("  Tick %d:\n", tick)
		fmt.Printf("    Liquidity: gross=%s net=%s\n", info.LiquidityGross.String(), info.LiquidityNet.String())
	}
}

// zeroForOne: if true, search to the left; or else search to the right.
func findNextTick(s []int, tick int, zeroForOne bool) (tickNext int, exist bool) {
	idx, ok := slices.BinarySearch(s, tick)
	if ok {
		// found
		if zeroForOne {
			if idx-1 < 0 {
				return MIN_TICK, false
			}

			return s[idx-1], true
		}

		if idx+1 >= len(s) {
			return MAX_TICK, false
		}

		return s[idx+1], true
	}

	// not found
	if zeroForOne {
		if idx == 0 {
			return MIN_TICK, false
		}

		return s[idx-1], true
	}

	if idx == len(s) {
		return MAX_TICK, false
	}

	return s[idx], true
}

func (pool *Pool) crossTicks(tick int) *big.Int {
	info, ok := pool.Ticks[tick]
	if !ok {
		panic("not found tick info: tick=" + fmt.Sprint(tick))
	}

	return big.NewInt(0).Set(info.LiquidityNet)
}

func (pool *Pool) clearTick(tick int) {
	delete(pool.Ticks, tick)
}

func (pool *Pool) flipTick(tick int) {
	if tick%pool.TickSpacing != 0 {
		logger.Fatal().Msgf("invalid tick: tick=%d tickSpacing=%d poolAddr=%s",
			tick, pool.TickSpacing, pool.Address)
	}

	// if len(pool.Ticks) == 0 {
	// pool.TickList = append(pool.TickList, tick)
	// flipTickBitmap(pool.tickBitmap, tick, pool.TickSpacing)

	// return
	// }

	// get index
	// idx, ok := slices.BinarySearch(pool.TickList, tick)
	// if ok {
	// 	// remove the tick from tick list
	// 	pool.TickList = slices.Delete(pool.TickList, idx, idx+1)
	// } else {
	// 	// insert to tick list
	// 	pool.TickList = slices.Insert(pool.TickList, idx, tick)
	// }

	flipTickBitmap(pool.tickBitmap, tick, pool.TickSpacing)
}

// Updates a tick and returns true if the tick was flipped from initialized to uninitialized, or vice versa.
func (pool *Pool) UpdateTicks(
	tick int,
	liquidityDelta *big.Int,
	upper bool,
) (flipped bool) {
	info, ok := pool.Ticks[tick]
	if !ok {
		info = &TickInfo{
			Tick:           tick,
			LiquidityNet:   big.NewInt(0),
			LiquidityGross: big.NewInt(0),
		}
		pool.Ticks[tick] = info
	}

	liquidityGrossBefore := info.LiquidityGross
	liquidityGrossAfter := big.NewInt(0).Add(info.LiquidityGross, liquidityDelta)
	flipped = (liquidityGrossAfter.Cmp(bigZero) == 0) != (liquidityGrossBefore.Cmp(bigZero) == 0)

	// if liquidityGrossBefore.Cmp(bigZero) == 0 {
	// 	if tick <= tickCurrent {
	// 	}
	// }

	info.LiquidityGross = liquidityGrossAfter
	if upper {
		info.LiquidityNet.Sub(info.LiquidityNet, liquidityDelta)
	} else {
		info.LiquidityNet.Add(info.LiquidityNet, liquidityDelta)
	}

	return
}
