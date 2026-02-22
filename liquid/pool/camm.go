package pool

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"starbase.ag/liquidity/contracts/swaprouter"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
)

var (
	tokenToPools = map[string]map[string]*Pool{}
	cammPools    = map[string]*Pool{}
	bigZero      = big.NewInt(0)
	bigOne       = big.NewInt(1)
	// https://mavlevin.com/2023/02/22/Size-Matters-Solidity-Integer-Range-Cheatsheet-From-uint8-To-uint256.html
	// uint128Max, _ = big.NewInt(0).SetString("340282366920938463463374607431768211455", 0).
	uint8Max, _   = big.NewInt(0).SetString("255", 0)
	uint16Max, _  = big.NewInt(0).SetString("65535", 0)
	uint32Max, _  = big.NewInt(0).SetString("4294967295", 0)
	uint64Max, _  = big.NewInt(0).SetString("18446744073709551615", 0)
	uint128Max, _ = big.NewInt(0).SetString("340282366920938463463374607431768211455", 0)
	uint160Max, _ = big.NewInt(0).SetString("1461501637330902918203684832716283019655932542975", 0)
	uint256Max, _ = big.NewInt(0).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 0)
	Q96, _        = big.NewInt(0).SetString("0x1000000000000000000000000", 0)
)

// SyncEvent sync event.
type SyncEvent struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
}

// CreateCAMMPool create CAMM pool, on event PoolCreated.
func CreateCAMMPool(token0, token1, poolAddr, factory string, fee uint, tickSpacing int, block uint64) *Pool {
	token0 = strings.ToLower(token0)
	token1 = strings.ToLower(token1)

	if token0 > token1 {
		token0, token1 = token1, token0
	}

	poolAddr = strings.ToLower(poolAddr)
	pool := &Pool{
		PoolInfo: PoolInfo{
			Address:          poolAddr,
			Token0:           strings.ToLower(token0),
			Token1:           strings.ToLower(token1),
			Factory:          strings.ToLower(factory),
			Fee:              fee,
			TickSpacing:      tickSpacing,
			LastBlockUpdated: block,
			Typ:              common.PoolTypeCAMM,
		},
		Tick: 0,
		// TickList:     []int{},
		Liquidity:    nil,
		SqrtPriceX96: nil,
		Reserve0:     big.NewInt(0),
		Reserve1:     big.NewInt(0),
		Ticks:        map[int]*TickInfo{},
		tickBitmap:   map[int16]*big.Int{},
	}
	cammPools[poolAddr] = pool

	// tokenToPools[token0][token1] = pool
	mp, ok := tokenToPools[token0]
	if !ok {
		mp = map[string]*Pool{
			token1: pool,
		}
		tokenToPools[token0] = mp
	} else {
		mp[token1] = pool
	}

	// tokenToPools[token1][token0] = pool
	mp, ok = tokenToPools[token1]
	if !ok {
		mp = map[string]*Pool{
			token0: pool,
		}
		tokenToPools[token1] = mp
	} else {
		mp[token0] = pool
	}

	return pool
}

// InitializeCAMMPool initialize pool, on event Initialize.
func GetCAMMPool(poolAddr string) *Pool {
	// poolAddr = strings.ToLower(poolAddr)
	pool := cammPools[poolAddr]

	if nil == pool {
		logger.Fatal().Msg("GetCAMMPool: not found pool: " + poolAddr)
	}

	return pool
}

// IPool interface.
func (pool *Pool) PoolAddress() string {
	return pool.Address
}

func (pool *Pool) FactoryAddress() string {
	return pool.Factory
}

func (pool *Pool) Token0Address() string {
	return pool.Token0
}

func (pool *Pool) Token1Address() string {
	return pool.Token1
}

func (pool *Pool) GetPoolType() common.PoolType {
	return common.PoolTypeCAMM
}

func (pool *Pool) OnInitialize(sqrtPriceX96 *big.Int, tick int) {
	pool.Tick = tick
	pool.SqrtPriceX96 = sqrtPriceX96
	pool.Liquidity = big.NewInt(0)
	pool.Initialized = true

	logger.Info().
		Str("vender", pool.Vendor).
		Str("pool", pool.Address).
		Str("sqrtPrice", sqrtPriceX96.String()).
		Int("tick", tick).
		Msg("initialize v3 pool")
}

// OnSwap deal with pool swap event.
func (pool *Pool) OnSwap(event *types.Log) bool {
	// https://etherscan.io/tx/0xe04cb094f051e928ff1673ae3e6b14a72e59800f247be567e8874d522f1c7c24#eventlog
	var swap swaprouter.UniswapV3PoolSwap
	if err := UniswapV3PoolABI.UnpackIntoInterface(&swap, "Swap", event.Data); err != nil {
		logger.Fatal().
			Str("vendor", pool.Vendor).
			Str("factory", pool.Factory).
			Str("txhahs", event.TxHash.Hex()).
			Str("pool", pool.Address).
			Msg("Unpack UniswapV3Pool Swap event data failed")
	}

	tick := int(swap.Tick.Int64())

	correct := pool.VerifySwapResult(event.BlockNumber, event.TxHash.String(), tick, swap.SqrtPriceX96, swap.Amount0, swap.Amount1)
	if !correct {
		return false
	}

	pool.Synced = true
	pool.onSwap(tick,
		swap.SqrtPriceX96, swap.Liquidity, swap.Amount0, swap.Amount1,
		event.BlockNumber,
		event.TxHash.String())

	return true
}

// OnSwap deal with pancake pool swap event.
func (pool *Pool) OnPancakeSwap(event *types.Log) bool {
	var swap swaprouter.PancakeV3PoolSwap

	if err := PancakeV3PoolABI.UnpackIntoInterface(&swap, "Swap", event.Data); err != nil {
		logger.Fatal().
			Str("vendor", pool.Vendor).
			Str("factory", pool.Factory).
			Str("txhahs", event.TxHash.Hex()).
			Str("pool", pool.Address).
			Msg("Unpack PancakeV3Pool Swap event data failed")
	}

	tick := int(swap.Tick.Int64())

	correct := pool.VerifySwapResult(event.BlockNumber, event.TxHash.String(), tick, swap.SqrtPriceX96, swap.Amount0, swap.Amount1)
	if !correct {
		return false
	}

	pool.Synced = true
	pool.onSwap(tick,
		swap.SqrtPriceX96, swap.Liquidity, swap.Amount0, swap.Amount1,
		event.BlockNumber, event.TxHash.String())

	return true
}

func (pool *Pool) onSwap(tick int,
	sqrtPrice, liquidity, amount0, amount1 *big.Int,
	block uint64,
	txhash string,
) {
	// pool.swap(zeroForOne, )
	pool.SqrtPriceX96 = sqrtPrice
	pool.Liquidity = liquidity
	pool.Tick = tick

	pool.Reserve0.Add(pool.Reserve0, amount0)
	pool.Reserve1.Add(pool.Reserve1, amount1)

	// logger.Info().
	// 	Str("pool", pool.Address).
	// 	Str("vendor", pool.Vendor).
	// 	Str("txhash", txhash).
	// 	Msgf("pool swap: tick=%d liquidity=%v sqrtPrice=%v block=%d",
	// 		pool.Tick, pool.Liquidity, sqrtPrice, block)
}

// OnCollect deal with pool Collect event.
func (pool *Pool) OnCollect(event *types.Log) {
	var col swaprouter.UniswapV3PoolCollect
	if err := UniswapV3PoolABI.UnpackIntoInterface(&col, "Collect", event.Data); err != nil {
		logger.Fatal().
			Str("vendor", pool.Vendor).
			Str("factory", pool.Factory).
			Str("txhahs", event.TxHash.Hex()).
			Str("pool", pool.Address).
			Msg("Unpack Collect event data failed")
	}

	topics := event.Topics
	col.TickLower = big.NewInt(0).SetBytes(topics[2].Bytes())
	col.TickUpper = big.NewInt(0).SetBytes(topics[3].Bytes())

	pool.Reserve0.Sub(pool.Reserve0, col.Amount0)
	pool.Reserve1.Sub(pool.Reserve1, col.Amount1)

	// logger.Info().
	// 	Str("pool", pool.Address).
	// 	Str("vendor", pool.Vendor).
	// 	Str("txhash", event.TxHash.Hex()).
	// 	Msgf("pool Collect: tickRange=[%v, %v] amount0=%v amount1=%v reserve0=%v reserve1=%v block=%d",
	// 		col.TickLower.Int64(), col.TickUpper.Int64(), col.Amount0, col.Amount1, pool.Reserve0, pool.Reserve1, event.BlockNumber)
}

// OnMint deal with pool Mint event.
func (pool *Pool) OnMint(event *types.Log) {
	var mint swaprouter.UniswapV3PoolMint
	if err := UniswapV3PoolABI.UnpackIntoInterface(&mint, "Mint", event.Data); err != nil {
		logger.Fatal().
			Str("vendor", pool.Vendor).
			Str("factory", pool.Factory).
			Str("txhahs", event.TxHash.Hex()).
			Str("pool", pool.Address).
			Msg("Unpack Mint event data failed")
	}

	topics := event.Topics
	mint.TickLower = big.NewInt(0).SetBytes(topics[2].Bytes())
	mint.TickUpper = big.NewInt(0).SetBytes(topics[3].Bytes())
	tickLower := int(mint.TickLower.Int64())
	tickUpper := int(mint.TickUpper.Int64())
	txhash := event.TxHash.Hex()

	pool.modifyPosition(tickLower, tickUpper, mint.Amount, "Mint", txhash)
	pool.Reserve0.Add(pool.Reserve0, mint.Amount0)
	pool.Reserve1.Add(pool.Reserve1, mint.Amount1)

	// logger.Info().
	// 	Str("pool", pool.Address).
	// 	Str("vendor", pool.Vendor).
	// 	Str("txhash", txhash).
	// 	Msgf("pool Mint: tickRange=[%v, %v] liquidity=%v amount0=%v amount1=%v block=%d",
	// 		tickLower, tickUpper, mint.Amount, mint.Amount0, mint.Amount1, event.BlockNumber)
}

// OnBurn deal with pool Burn event.
func (pool *Pool) OnBurn(event *types.Log) {
	// tickLower, tickUpper int, amount *big.Int) {
	// https://etherscan.io/tx/0x11cfb0e1c780f8e704c1d4e8938cfc9d3cc969d1766f861c3d7b9bf75e57fef2#eventlog
	var burn swaprouter.UniswapV3PoolBurn
	if err := UniswapV3PoolABI.UnpackIntoInterface(&burn, "Burn", event.Data); err != nil {
		logger.Fatal().
			Str("vendor", pool.Vendor).
			Str("factory", pool.Factory).
			Str("txhahs", event.TxHash.Hex()).
			Str("pool", pool.Address).
			Msg("Unpack Burn event data failed")
	}

	topics := event.Topics
	tickLower := int(big.NewInt(0).SetBytes(topics[2].Bytes()).Int64())
	tickUpper := int(big.NewInt(0).SetBytes(topics[3].Bytes()).Int64())
	txhash := event.TxHash.Hex()

	pool.modifyPosition(tickLower, tickUpper, big.NewInt(0).Neg(burn.Amount), "Burn", txhash)
	// burn have no token transfer
	// pool.Reserve0.Sub(pool.Reserve0, burn.Amount0)
	// pool.Reserve1.Sub(pool.Reserve1, burn.Amount1)

	// logger.Info().
	// 	Str("pool", pool.Address).
	// 	Str("Vendor", pool.Vendor).
	// 	Str("txhash", txhash).
	// 	Msgf("pool Burn: tickRange=[%v, %v] liquidity=%v block=%d amount0=%v amount1=%v reserve0=%v reserve1=%v",
	// 		tickLower, tickUpper, burn.Amount, event.BlockNumber, burn.Amount0, burn.Amount1, pool.Reserve0, pool.Reserve1)
}

func (pool *Pool) modifyPosition(tickLower, tickUpper int, amount *big.Int, evtName, txhash string) {
	pool.checkStatus()
	pool.checkTicks(tickLower, tickUpper)

	if amount.Cmp(bigZero) == 0 {
		// logger.Warn().
		// 	Str("pool", pool.Address).
		// 	Str("vendor", pool.Vendor).
		// 	Str("event", evtName).
		// 	Str("txhash", txhash).
		// 	Msg("liquidityDelta is 0, do nothing.")

		return
	}

	pool.updatePosition(tickLower, tickUpper, amount)

	// we only care active liquidity
	if pool.Tick >= tickLower && pool.Tick < tickUpper {
		pool.Liquidity.Add(pool.Liquidity, amount)

		if pool.Liquidity.Cmp(bigZero) < 0 {
			logger.Fatal().Str("pool", pool.Address).Msgf("pool liquidity less than zero: amount=%s liquidity=%s",
				amount.String(), pool.Liquidity.String())
		}
	}
}

// Gets and updates pool tick liquidity with the given liquidity delta.
func (pool *Pool) updatePosition(tickLower, tickUpper int, amount *big.Int) {
	flippedLower := pool.UpdateTicks(tickLower, amount, false)
	flippedUpper := pool.UpdateTicks(tickUpper, amount, true)

	if flippedLower {
		pool.flipTick(tickLower)
	}

	if flippedUpper {
		pool.flipTick(tickUpper)
	}

	if amount.Cmp(bigZero) < 0 {
		if flippedLower {
			pool.clearTick(tickLower)
		}

		if flippedUpper {
			pool.clearTick(tickUpper)
		}
	}
}

// checkTicks validate tickLower, tickUpper.
func (pool *Pool) checkTicks(tickLower, tickUpper int) {
	if tickLower >= tickUpper {
		logger.Fatal().Msgf("tickUpper should great than tickLower, poolAddr=%s", pool.Address)
	}

	if tickLower < MIN_TICK {
		logger.Fatal().Msgf("tickLower should NOT less than MIN_TICK, poolAddr=%s", pool.Address)
	}

	if tickUpper > MAX_TICK {
		logger.Fatal().Msgf("tickUpper should NOT great than MIN_TICK, poolAddr=%s", pool.Address)
	}
}

// make sure pool was initialized.
func (pool *Pool) checkStatus() {
	if !pool.Initialized {
		logger.Fatal().Msgf("pool should be initialized, poolAddr=%s", pool.Address)
	}
}

// GetAmountOut give exact amountOut, calculate min amountIn.
func (pool *Pool) GetAmountIn(amountOut *big.Int) *big.Int {
	return big.NewInt(0)
}

func (pool *Pool) MockSwap(
	zeroForOne bool,
	amountSpecified,
	sqrtPriceLimitX96 *big.Int,
) (amount0, amount1, sqrtPriceAfter *big.Int, tickAfter int) {
	if sqrtPriceLimitX96 == nil {
		if zeroForOne {
			sqrtPriceLimitX96 = new(big.Int).Add(MIN_SQRT_RATIO, big.NewInt(1))
		} else {
			sqrtPriceLimitX96 = new(big.Int).Sub(MAX_SQRT_RATIO, big.NewInt(1))
		}
	}
	return pool.swap(zeroForOne, amountSpecified, sqrtPriceLimitX96, true)
}

func (pool *Pool) swap(zeroForOne bool,
	amountSpecified,
	sqrtPriceLimitX96 *big.Int,
	args ...interface{},
) (amount0, amount1, sqrtPriceAfter *big.Int, tickAfter int) {
	if amountSpecified.Cmp(bigZero) == 0 {
		panic("param amountSpecified should NOT be zero")
	}

	amountSpecifiedRemaining := big.NewInt(0).Set(amountSpecified)
	amountCalculated := big.NewInt(0)
	tick := pool.Tick
	liquidity := big.NewInt(0).Set(pool.Liquidity)
	sqrtPriceX96 := big.NewInt(0).Set(pool.SqrtPriceX96)

	var (
		amountIn  *big.Int
		amountOut *big.Int
		feeAmount *big.Int
	)

	exactInput := true
	if amountSpecified.Cmp(bigZero) < 0 {
		exactInput = false
	}

	for amountSpecifiedRemaining.Cmp(bigZero) != 0 && sqrtPriceX96.Cmp(sqrtPriceLimitX96) != 0 {
		sqrtPriceStartX96 := big.NewInt(0).Set(sqrtPriceX96)
		tickNext, initialized := nextInitializedTickWithinOneWord(pool.tickBitmap, tick, pool.TickSpacing, zeroForOne)
		sqrtPriceNextX96 := getSqrtRatioAtTick(tickNext)

		// println("pool: ", pool.Address, "zeroForOne:", zeroForOne, "exactInput:", exactInput)
		// println("amountSpecifiedRemaining: ", amountSpecifiedRemaining.String())
		// println("sqrtPriceX96: ", sqrtPriceX96.String())
		// println("sqrtPriceLimitX96: ", sqrtPriceLimitX96.String())
		// println("sqrtPriceStartX96: ", sqrtPriceStartX96.String())
		// println("tick: ", tick)
		// println("tick:", tick, "tickNext:", tickNext, " sqrtPriceNextX96:", sqrtPriceNextX96.String())
		// println("initialized: ", initialized)
		// println("liquidity: ", liquidity.String())

		var priceLimitX96 *big.Int

		if zeroForOne {
			if sqrtPriceNextX96.Cmp(sqrtPriceLimitX96) < 0 {
				priceLimitX96 = big.NewInt(0).Set(sqrtPriceLimitX96)
			} else {
				priceLimitX96 = big.NewInt(0).Set(sqrtPriceNextX96)
			}
		} else {
			if sqrtPriceNextX96.Cmp(sqrtPriceLimitX96) > 0 {
				priceLimitX96 = big.NewInt(0).Set(sqrtPriceLimitX96)
			} else {
				priceLimitX96 = big.NewInt(0).Set(sqrtPriceNextX96)
			}
		}

		sqrtPriceX96, amountIn, amountOut, feeAmount = computeSwapStep(sqrtPriceX96,
			priceLimitX96,
			liquidity,
			amountSpecifiedRemaining,
			int(pool.Fee)) // nolint

		// logger.Info().Msgf("after computeSwapStep: sqrtPriceX96: %v L: %v amtSpec: %v amtIn: %v amtOut: %v",
		// sqrtPriceX96.String(), liquidity, amountSpecifiedRemaining, amountIn, amountOut)

		if exactInput {
			amountSpecifiedRemaining.Sub(amountSpecifiedRemaining, big.NewInt(0).Add(amountIn, feeAmount))
			amountCalculated.Sub(amountCalculated, amountOut)
		} else {
			amountSpecifiedRemaining.Add(amountSpecifiedRemaining, amountOut)
			amountCalculated.Add(amountCalculated, big.NewInt(0).Add(amountIn, feeAmount))
		}

		// println("--------")
		// println("sqrtPriceX96: ", sqrtPriceX96.String())
		// println("amountIn: ", amountIn.String())
		// println("amountOut: ", amountOut.String())
		// println("feeAmount: ", feeAmount.String())
		// println("amountSpecifiedRemaining: ", amountSpecifiedRemaining.String())
		// println("amountCalculated: ", amountCalculated.String())
		// println("--------")

		// if liquidity.Cmp(bigZero) > 0 {
		// 	// update fee tracker
		// }

		// shift tick if we reached the next price
		if sqrtPriceX96.Cmp(sqrtPriceNextX96) == 0 {
			if initialized {
				liquidityNet := pool.crossTicks(tickNext)
				if zeroForOne {
					liquidity.Sub(liquidity, liquidityNet)
				} else {
					liquidity.Add(liquidity, liquidityNet)
				}
			}

			if zeroForOne {
				tick = tickNext - 1
			} else {
				tick = tickNext
			}
		} else if sqrtPriceX96.Cmp(sqrtPriceStartX96) != 0 {
			tick = getTickAtSqrtRatio(sqrtPriceX96)
		}
	}

	dryrun := false

	if len(args) > 0 {
		val, ok := args[0].(bool)
		if ok {
			dryrun = val
		}
	}

	if !dryrun {
		if tick != pool.Tick {
			pool.Tick = tick
		}

		pool.SqrtPriceX96 = sqrtPriceX96
		pool.Liquidity = liquidity
	}

	if zeroForOne == exactInput {
		amount0, amount1 = big.NewInt(0).Sub(amountSpecified, amountSpecifiedRemaining), big.NewInt(0).Set(amountCalculated)
	} else {
		amount0, amount1 = big.NewInt(0).Set(amountCalculated), big.NewInt(0).Sub(amountSpecified, amountSpecifiedRemaining)
	}

	sqrtPriceAfter = big.NewInt(0).Set(sqrtPriceX96)
	tickAfter = tick

	return
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func quoteExactInputSingle(tokenIn, tokenOut string, amountIn, sqrtPriceLimitX96 *big.Int) (amountOut, sqrtPriceX96After *big.Int, ticksCrossed uint32) {
	pool := tokenToPools[tokenIn][tokenOut]

	zeroForOne := tokenIn < tokenOut
	if sqrtPriceLimitX96.Cmp(bigZero) == 0 {
		if zeroForOne {
			sqrtPriceLimitX96 = big.NewInt(0).Add(MIN_SQRT_RATIO, bigOne)
		} else {
			sqrtPriceLimitX96 = big.NewInt(0).Sub(MAX_SQRT_RATIO, bigOne)
		}
	}

	tickBefore := pool.Tick
	amount0, amount1, sqrtPriceX96After, tickAfter := pool.swap(zeroForOne, amountIn, sqrtPriceLimitX96)

	var (
		isExactInput   bool
		amountToPay    *big.Int
		amountReceived *big.Int
	)

	if amount0.Cmp(bigZero) > 0 {
		isExactInput = tokenIn < tokenOut
		amountToPay = big.NewInt(0).Set(amount0)
		amountReceived = big.NewInt(0).Neg(amount1)
	} else {
		isExactInput = tokenOut < tokenIn
		amountToPay = big.NewInt(0).Set(amount1)
		amountReceived = big.NewInt(0).Neg(amount0)
	}

	if isExactInput {
		amountOut = amountReceived
	} else {
		amountOut = amountToPay
	}

	ticksCrossed = countInitializedTicksCrossed(pool, tickBefore, tickAfter)

	return
}

func quoteExactOutputSingle(tokenIn, tokenOut string, amount, sqrtPriceLimitX96 *big.Int) (amountIn, sqrtPriceX96After *big.Int, ticksCrossed uint32) {
	pool := tokenToPools[tokenIn][tokenOut]

	zeroForOne := tokenIn < tokenOut
	if sqrtPriceLimitX96.Cmp(bigZero) == 0 {
		if zeroForOne {
			sqrtPriceLimitX96 = big.NewInt(0).Add(MIN_SQRT_RATIO, bigOne)
		} else {
			sqrtPriceLimitX96 = big.NewInt(0).Sub(MAX_SQRT_RATIO, bigOne)
		}
	}

	tickBefore := pool.Tick
	amount0, amount1, sqrtPriceX96After, tickAfter := pool.swap(zeroForOne, big.NewInt(0).Neg(amount), sqrtPriceLimitX96)
	// tickAfter := pool.Tick
	// sqrtPriceX96After = big.NewInt(0).Set(pool.SqrtPriceX96)

	// println("amount0:", amount0.String())
	// println("amount1:", amount1.String())

	var (
		isExactInput   bool
		amountToPay    *big.Int
		amountReceived *big.Int
	)

	if amount0.Cmp(bigZero) > 0 {
		isExactInput = tokenIn > tokenOut
		amountToPay = big.NewInt(0).Set(amount0)
		amountReceived = big.NewInt(0).Neg(amount1)
	} else {
		isExactInput = tokenOut > tokenIn
		amountToPay = big.NewInt(0).Set(amount1)
		amountReceived = big.NewInt(0).Neg(amount0)
	}

	if isExactInput {
		amountIn = amountReceived
	} else {
		amountIn = amountToPay
	}

	ticksCrossed = countInitializedTicksCrossed(pool, tickBefore, tickAfter)

	return
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func quoteExactInput(tokens []string, amountIn *big.Int) (*big.Int, []*big.Int, []uint32) {
	sqrtPriceX96AfterList := make([]*big.Int, len(tokens)-1)
	initializedTicksCrossedList := make([]uint32, len(tokens)-1)

	amt := big.NewInt(0).Set(amountIn)
	for i := 0; i < len(tokens)-1; i++ {
		amountOut, _sqrtPriceX96After, ticksCrossed := quoteExactInputSingle(tokens[i], tokens[i+1], amt, big.NewInt(0))
		amt = big.NewInt(0).Set(amountOut)
		sqrtPriceX96AfterList[i] = big.NewInt(0).Set(_sqrtPriceX96After)
		initializedTicksCrossedList[i] = ticksCrossed
	}

	return amt, sqrtPriceX96AfterList, initializedTicksCrossedList
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func quoteExactOutput(tokens []string, amountOut *big.Int) (*big.Int, []*big.Int, []uint32) {
	sqrtPriceX96AfterList := make([]*big.Int, len(tokens)-1)
	initializedTicksCrossedList := make([]uint32, len(tokens)-1)

	amt := big.NewInt(0).Set(amountOut)
	for i := 0; i < len(tokens)-1; i++ {
		amountOut, _sqrtPriceX96After, ticksCrossed := quoteExactOutputSingle(tokens[i+1], tokens[i], amt, big.NewInt(0))
		amt = big.NewInt(0).Set(amountOut)
		sqrtPriceX96AfterList[i] = big.NewInt(0).Set(_sqrtPriceX96After)
		initializedTicksCrossedList[i] = ticksCrossed
	}

	return amt, sqrtPriceX96AfterList, initializedTicksCrossedList
}

func countInitializedTicksCrossed(pool *Pool, tickBefore, tickAfter int) uint32 {
	// fmt.Println("tick before/after:", tickBefore, tickAfter)
	var (
		tickBeforeInitialized   bool
		tickAfterInitialized    bool
		initializedTicksCrossed uint32
		wordPosLower            int16
		wordPosHigher           int16
		bitPosLower             uint8
		bitPosHigher            uint8
	)

	tickSpacing := pool.TickSpacing
	// Get the key and offset in the tick bitmap of the active tick before and after the swap.
	wordPos := int16((tickBefore / tickSpacing) >> 8) // nolint
	bitPos := uint8((tickBefore / tickSpacing) % 256) // nolint

	wordPosAfter := int16((tickAfter / tickSpacing) >> 8) // nolint
	bitPosAfter := uint8((tickAfter / tickSpacing) % 256) // nolint
	// In the case where tickAfter is initialized, we only want to count it if we are swapping downwards.
	// If the initializable tick after the swap is initialized, our original tickAfter is a
	// multiple of tick spacing, and we are swapping downwards we know that tickAfter is initialized
	// and we shouldn't count it.
	// ((self.tickBitmap(wordPosAfter) & (1 << bitPosAfter)) > 0)
	word := pool.tickBitmap[wordPosAfter]
	if word == nil {
		word = big.NewInt(0)
	}

	tmp1 := big.NewInt(0).And(word, big.NewInt(0).Lsh(bigOne, uint(bitPosAfter)))
	tickAfterInitialized = (tmp1.Cmp(bigZero) > 0) &&
		((tickAfter % tickSpacing) == 0) &&
		(tickBefore > tickAfter)

	// In the case where tickBefore is initialized, we only want to count it if we are swapping upwards.
	// Use the same logic as above to decide whether we should count tickBefore or not.
	// ((pool.tickBitmap[wordPos] & (1 << bitPos)) > 0)
	word = pool.tickBitmap[wordPos]
	if word == nil {
		word = big.NewInt(0)
	}

	tmp2 := big.NewInt(0).And(word, big.NewInt(0).Lsh(bigOne, uint(bitPos)))
	// fmt.Printf("tmp2: %v word=%v bitPos=%v\n", tmp2, word, bitPos)
	tickBeforeInitialized = (tmp2.Cmp(bigZero) > 0) &&
		((tickBefore % tickSpacing) == 0) &&
		(tickBefore < tickAfter)

	if wordPos < wordPosAfter || (wordPos == wordPosAfter && bitPos <= bitPosAfter) {
		wordPosLower = wordPos
		bitPosLower = bitPos
		wordPosHigher = wordPosAfter
		bitPosHigher = bitPosAfter
	} else {
		wordPosLower = wordPosAfter
		bitPosLower = bitPosAfter
		wordPosHigher = wordPos
		bitPosHigher = bitPos
	}
	// Count the number of initialized ticks crossed by iterating through the tick bitmap.
	// Our first mask should include the lower tick and everything to its left.
	mask := big.NewInt(0).Set(uint256Max) // type(uint256).max << bitPosLower;
	mask.Lsh(mask, uint(bitPosLower))

	for wordPosLower <= wordPosHigher {
		// If we're on the final tick bitmap page, ensure we only count up to our
		// ending tick.
		if wordPosLower == wordPosHigher {
			// mask = mask & (type(uint256).max >> (255 - bitPosHigher));
			mask.And(mask, big.NewInt(0).Rsh(big.NewInt(0).Set(uint256Max), uint(255-bitPosHigher)))
		}

		word = pool.tickBitmap[wordPosLower]
		if word == nil {
			word = big.NewInt(0)
		}

		masked := big.NewInt(0).And(word, mask)
		initializedTicksCrossed += uint32(countOneBits(masked))
		wordPosLower++
		// Reset our mask so we consider all bits on the next iteration.
		mask = big.NewInt(0).Set(uint256Max)
	}

	// fmt.Printf("initializedTicksCrossed: %d tickAfterInitialized: %v tickBeforeInitialized: %v\n",
	// initializedTicksCrossed, tickAfterInitialized, tickBeforeInitialized)
	if tickAfterInitialized {
		initializedTicksCrossed -= 1
	}

	if tickBeforeInitialized {
		initializedTicksCrossed -= 1
	}
	/*
		idxBefore, ok1 := slices.BinarySearch(pool.TickList, tickBefore)
		if ok1 && tickBefore < tickAfter {
			tickBeforeInitialized = true
		}

		idxAfter, ok2 := slices.BinarySearch(pool.TickList, tickAfter)
		if ok2 && tickBefore > tickAfter {
			tickAfterInitialized = true
		}

		// println("tick before/after:", tickBefore, tickAfter, idxBefore, idxAfter)

		if tickAfter > tickBefore {
			// to right
			count = uint32(idxAfter - idxBefore)
			if ok2 {
				count += 1
			}
		} else {
			// to left
			count = uint32(idxBefore - idxAfter)
			if ok1 {
				count += 1
			}
		}

		if tickAfterInitialized {
			count -= 1
		}

		if tickBeforeInitialized {
			count -= 1
		}
	*/
	return initializedTicksCrossed
}

func countOneBits(x *big.Int) (bits uint16) {
	for x.Cmp(bigZero) != 0 {
		bits++

		x.And(x, big.NewInt(0).Sub(x, bigOne))
	}

	return
}

var (
	// ignoredFailedTxHash = map[string]bool{
	// 	// 0xfb83ff48d18242620ac10a859bb8baa411f49de43745a958fff0ddd98f97f531
	// 	// 0xb1c1e628b57c86d4c829082232ca7f64daa163fdb2fae1cbd8a83875b261868d
	// 	"0xef978f1e64f8332047a0a628068a4aa01792f37ea841b6de06e065e65c02d06d": true,
	// 	// 0xd57662d3ff9c8aa247786f92432f5272bbf632971f0bf8a77d2ec4c9bf1f48ef
	// }.
	swapVerifySuccess = uint64(0) //nolint
	swapVerifyFailed  = uint64(0) //nolint
)

func (pool *Pool) VerifySwapResult(blocknumber uint64, txhash string,
	tick int,
	sqrtPrice *big.Int,
	amount0 *big.Int,
	amount1 *big.Int,
) bool {
	// swap *swaprouter.UniswapV3PoolSwap) {
	// tick := int(swap.Tick.Int64())
	// sqrtPrice := swap.SqrtPriceX96
	// amount0 := swap.Amount0
	// amount1 := swap.Amount1
	// logger.Info().Msgf("pool sqrtPrice: %v after swap sqrtPrice: %v", pool.SqrtPriceX96, sqrtPrice)

	err := pool.VerifySwapResultByExact(txhash, tick, sqrtPrice, amount0, amount1, true, nil)
	if err != nil {
		// logger.Warn().Err(err).Msg("verify swap result with exactIn=true failed")
		err = pool.VerifySwapResultByExact(txhash, tick, sqrtPrice, amount0, amount1, false, nil)
		if err != nil {
			// logger.Warn().Err(err).Msg("verify swap result with exactIn=false failed")
			// use swap.SqrtPriceX96
			err = pool.VerifySwapResultByExact(txhash, tick, sqrtPrice, amount0, amount1, true, sqrtPrice)
			if err != nil {
				// logger.Warn().Err(err).Msg("verify swap result with param exactIn=true And sqrtPriceLimitX96 failed")
				// if err != nil {
				// 	logger.Warn().Err(err).Msg("verify swap result with param exactIn=false And sqrtPriceLimitX96 failed")
				// }
				err = pool.VerifySwapResultByExact(txhash, tick, sqrtPrice, amount0, amount1, false, sqrtPrice)
			}
		}
	}

	if err == nil {
		swapVerifySuccess++
	} else {
		swapVerifyFailed++
	}

	// logger.Info().Msgf("================== swap verify: success=%d failed=%d ==================",
	// 	swapVerifySuccess, swapVerifyFailed)

	if err != nil {
		pool.printInfos()
		logger.Error().
			Str("pool", pool.Address).
			Str("vendor", pool.Vendor).
			Uint64("blocknumber", blocknumber).
			Msgf("swap result NOT equal, param: amount0=%s amount1=%s block=%d txhash=%s",
				amount0.String(), amount1.String(), blocknumber, txhash)

		// if config.IsLocal() {
		// 	panic("swap verify failed")
		// }

		return false
	}

	return true
}

func (pool *Pool) VerifySwapResultByExact(txhash string, // swap *swaprouter.UniswapV3PoolSwap,
	tick int,
	sqrtPriceX96 *big.Int,
	amount0 *big.Int,
	amount1 *big.Int,
	exactIn bool,
	sqrtPriceLimit *big.Int,
) error {
	var sqrtPriceLimitX96 *big.Int

	zeroForOne := true

	if sqrtPriceX96.Cmp(pool.SqrtPriceX96) > 0 {
		zeroForOne = false
		sqrtPriceLimitX96 = big.NewInt(0).Set(MAX_SQRT_RATIO)
		// sqrtPriceLimitX96.Add(sqrtPriceLimitX96, big.NewInt(1))
		sqrtPriceLimitX96.Sub(sqrtPriceLimitX96, bigOne)
	} else {
		sqrtPriceLimitX96 = big.NewInt(0).Set(MIN_SQRT_RATIO)
		// sqrtPriceLimitX96.Sub(sqrtPriceLimitX96, big.NewInt(1))
		sqrtPriceLimitX96.Add(sqrtPriceLimitX96, bigOne)
	}

	if sqrtPriceLimit != nil {
		sqrtPriceLimitX96 = sqrtPriceLimit
	}

	// amount0 := swap.Amount0
	// amount1 := swap.Amount1

	var amountSpecified *big.Int
	if zeroForOne == exactIn {
		amountSpecified = big.NewInt(0).Set(amount0)
	} else {
		amountSpecified = big.NewInt(0).Set(amount1)
	}

	if amountSpecified.Cmp(bigZero) == 0 {
		logger.Warn().Msgf("amount specified is zero: txhash=%v", txhash)
		return nil
	}

	compareFn := func() error {
		amt0, amt1, sqrtPriceAfter, tickAfter := pool.swap(zeroForOne, amountSpecified, sqrtPriceLimitX96, true)
		if amount0.Cmp(amt0) != 0 {
			err := fmt.Errorf("pool swap result amount0 diff: zeroForOne=%v exactIn=%v amountSpecified=%s sqrtPriceLimitX96=%s expectAmount0=%s realAmount0=%s",
				zeroForOne, exactIn, amountSpecified.String(), sqrtPriceLimitX96.String(), amount0.String(), amt0.String())
			// logger.Warn().Msg(err.Error())

			return err
		}

		if amount1.Cmp(amt1) != 0 {
			err := fmt.Errorf("pool swap result amount1 diff: zeroForOne=%v exactIn=%v amountSpecified=%s sqrtPriceLimitX96=%s expectAmount1=%s realAmount1=%s",
				zeroForOne, exactIn, amountSpecified.String(), sqrtPriceLimitX96.String(), amount1.String(), amt1.String())
			// logger.Warn().Msg(err.Error())

			return err
		}

		if sqrtPriceX96.Cmp(sqrtPriceAfter) != 0 {
			maxPrice := big.NewInt(0).Sub(MAX_SQRT_RATIO, bigOne)
			minPrice := big.NewInt(0).Add(MIN_SQRT_RATIO, bigOne)

			if sqrtPriceX96.Cmp(minPrice) == 0 || sqrtPriceX96.Cmp(maxPrice) == 0 {
				// logger.Warn().Msgf("!!!! swap sqrtPrice reach min/max sqrtPriceX96: zeroForOne=%v sqrtPrice=%s",
				// zeroForOne, sqrtPriceX96.String())
				return nil
			}

			err := fmt.Errorf("pool swap result sqrtPrice diff: zeroForOne=%v exactIn=%v amountSpecified=%s sqrtPriceLimitX96=%s expectSqrtPrice=%s realSqrtPrice=%s",
				zeroForOne, exactIn, amountSpecified.String(), sqrtPriceLimitX96.String(), sqrtPriceX96.String(), sqrtPriceAfter.String())
			// logger.Warn().Msg(err.Error())

			return err
		}

		if tick != tickAfter {
			err := fmt.Errorf("pool swap result tick diff: zeroForOne=%v exactIn=%v amountSpecified=%s expectTick=%d realTick=%d",
				zeroForOne, exactIn, amountSpecified.String(), tick, tickAfter)
			// logger.Warn().Msg(err.Error())

			return err
		}

		return nil
	}

	err := compareFn()
	if err == nil {
		return nil
	}

	if sqrtPriceLimit != nil {
		if amountSpecified.Cmp(bigZero) > 0 {
			amountSpecified.Add(amountSpecified, bigOne)
		} else {
			amountSpecified.Sub(amountSpecified, bigOne)
		}

		err = compareFn()
	}

	return err
}
