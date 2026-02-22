package pool

import (
	"math/big"
	"testing"

	"github.com/defiweb/go-eth/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/pkg/logger"
)

var (
	_token0 = "0x0000000000000000000000000000000000000003"
	_token1 = "0x0000000000000000000000000000000000000004"
	_token2 = "0x0000000000000000000000000000000000000005"
)

func clearTestPools() {
	tokenToPools = map[string]map[string]*Pool{}
}

// 0 -> 2 cross 2 tick.
func TestSwapExactInput1(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)
	// 0 -> 2 cross 2 tick
	amt, sqrtPrices, crossTicks := quoteExactInput([]string{_token0, _token2}, big.NewInt(10000))
	assert.Equal(t, big.NewInt(9871), amt)
	assert.Equal(t, len(sqrtPrices), 1)
	assert.Equal(t, toBigIntMust("78461846509168490764501028180"), sqrtPrices[0])
	assert.True(t, crossTicks[0] == 2)
}

// 0 -> 2 cross 2 tick where after is initialized.
func TestSwapExactInput2(t *testing.T) {
	clearTestPools()

	pool := createPoolWithMultiplePositions(_token0, _token2)

	// 0 -> 2 cross 2 tick where after is initialized
	amt, sqrtPrices, crossTicks := quoteExactInput([]string{_token0, _token2}, big.NewInt(6200))
	assert.Equal(t, amt, big.NewInt(6143))
	assert.Equal(t, len(sqrtPrices), 1)
	assert.Equal(t, sqrtPrices[0], toBigIntMust("78757224507315167622282810783"))
	assert.Equal(t, len(crossTicks), 1)
	// assert.True(t, crossTicks[0] == 1)
	assert.Equal(t, pool.Tick, int(-120))
}

// 0 -> 2 cross 1 tick.
func TestSwapExactInput3(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amt, sqrtPrices, crossTicks := quoteExactInput([]string{_token0, _token2}, big.NewInt(4000))
	assert.Equal(t, amt, big.NewInt(3971))
	assert.Equal(t, len(sqrtPrices), 1)
	assert.Equal(t, sqrtPrices[0], toBigIntMust("78926452400586371254602774705"))
	assert.Equal(t, len(crossTicks), 1)
}

// 0 -> 2 cross 0 tick, starting tick not initialized.
func TestSwapExactInput4(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amt, sqrtPrices, crossTicks := quoteExactInput([]string{_token0, _token2}, big.NewInt(10))
	assert.Equal(t, amt, big.NewInt(8))
	assert.Equal(t, len(sqrtPrices), 1)
	assert.Equal(t, sqrtPrices[0], toBigIntMust("79227483487511329217250071027"))
	assert.Equal(t, len(crossTicks), 1)
	assert.Equal(t, crossTicks[0], uint32(0))
}

// 2 -> 0 cross 2.
func TestSwapExactInput6(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amt, sqrtPrices, crossTicks := quoteExactInput([]string{_token2, _token0}, big.NewInt(10000))
	assert.Equal(t, amt, big.NewInt(9871))
	assert.Equal(t, len(sqrtPrices), 1)
	assert.Equal(t, toBigIntMust("80001962924147897865541384515"), sqrtPrices[0])
	// assert.Equal(t, len(crossTicks), 1)
	assert.Equal(t, crossTicks[0], uint32(2))
}

// 2 -> 0 cross 2 where tick after is initialized.
func TestSwapExactInput7(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amt, sqrtPrices, crossTicks := quoteExactInput([]string{_token2, _token0}, big.NewInt(6250))
	assert.Equal(t, amt, big.NewInt(6190))
	assert.Equal(t, len(sqrtPrices), 1)
	assert.Equal(t, toBigIntMust("79705728824507063507279123685"), sqrtPrices[0])
	assert.Equal(t, len(crossTicks), 1)
	assert.Equal(t, crossTicks[0], uint32(2))
}

// 2 -> 0 cross 2 where tick after is initialized.
func TestSwapExactInput8(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amt, sqrtPrices, crossTicks := quoteExactInput([]string{_token2, _token0}, big.NewInt(6250))
	assert.Equal(t, amt, big.NewInt(6190))
	assert.Equal(t, len(sqrtPrices), 1)
	assert.Equal(t, toBigIntMust("79705728824507063507279123685"), sqrtPrices[0])
	assert.Equal(t, len(crossTicks), 1)
	assert.Equal(t, crossTicks[0], uint32(2))
}

// 2 -> 0 cross 0 tick, starting tick not initialized.
func TestSwapExactInput10(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amt, sqrtPrices, crossTicks := quoteExactInput([]string{_token2, _token0}, big.NewInt(103))
	assert.Equal(t, amt, big.NewInt(101))
	assert.Equal(t, len(sqrtPrices), 1)
	assert.Equal(t, toBigIntMust("79235858216754624215638319723"), sqrtPrices[0])
	assert.Equal(t, len(crossTicks), 1)
	assert.Equal(t, crossTicks[0], uint32(0))
}

// 2 -> 1.
func TestSwapExactInput11(t *testing.T) {
	clearTestPools()
	createPoolWithInitLiquidity(_token1, _token2)

	amt, sqrtPrices, crossTicks := quoteExactInput([]string{_token2, _token1}, big.NewInt(10000))
	assert.Equal(t, amt, big.NewInt(9871))
	assert.Equal(t, len(sqrtPrices), 1)
	assert.Equal(t, toBigIntMust("80018067294531553039351583520"), sqrtPrices[0])
	assert.Equal(t, len(crossTicks), 1)
	assert.Equal(t, crossTicks[0], uint32(0))
}

// 0 -> 2 -> 1.
func TestSwapExactInput12(t *testing.T) {
	clearTestPools()
	createPoolWithInitLiquidity(_token0, _token1)
	createPoolWithInitLiquidity(_token1, _token2)
	createPoolWithMultiplePositions(_token0, _token2)
	// _ = pool

	amt, sqrtPrices, crossTicks := quoteExactInput([]string{_token0, _token2, _token1}, big.NewInt(10000))
	assert.Equal(t, amt, big.NewInt(9745))
	assert.Equal(t, len(sqrtPrices), 2)
	assert.Equal(t, toBigIntMust("78461846509168490764501028180"), sqrtPrices[0])
	assert.Equal(t, toBigIntMust("80007846861567212939802016351"), sqrtPrices[1])
	assert.Equal(t, len(crossTicks), 2)
	assert.Equal(t, crossTicks[0], uint32(2))
	assert.Equal(t, crossTicks[1], uint32(0))
}

// 0 -> 2.
func TestQuoteExactInputSingle1(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amountOut,
		sqrtPriceX96After,
		initializedTicksCrossed := quoteExactInputSingle(
		_token0,
		_token2,
		toBigIntMust("10000"),
		encodePriceSqrt(big.NewInt(100), big.NewInt(102)))

	assert.Equal(t, initializedTicksCrossed, uint32(2))
	assert.Equal(t, amountOut, toBigIntMust("9871"))
	assert.Equal(t, sqrtPriceX96After, toBigIntMust("78461846509168490764501028180"))
}

// 2 -> 0.
func TestQuoteExactInputSingle2(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amountOut,
		sqrtPriceX96After,
		initializedTicksCrossed := quoteExactInputSingle(
		_token2,
		_token0,
		toBigIntMust("10000"),
		encodePriceSqrt(big.NewInt(102), big.NewInt(100)))

	assert.Equal(t, initializedTicksCrossed, uint32(2))
	assert.Equal(t, amountOut, toBigIntMust("9871"))
	assert.Equal(t, sqrtPriceX96After, toBigIntMust("80001962924147897865541384515"))
}

// func addPoolWithZeroTickInitialized(pool *Pool) {
// 	tickSpacing := 60
// 	getMinTick := func(tickSpacing int) int {
// 		return (-887272 / tickSpacing) * tickSpacing
// 	}
// 	getMaxTick := func(tickSpacing int) int {
// 		return (887272 / tickSpacing) * tickSpacing
// 	}
// 	tickLower := getMinTick(tickSpacing)
// 	tickUpper := getMaxTick(tickSpacing)
// 	addLiquidity(pool, tickLower, tickUpper, big.NewInt(1000000), big.NewInt(1000000))

// 	addLiquidity(pool, 0, 60, big.NewInt(100), big.NewInt(100))
// 	addLiquidity(pool, -120, 0, big.NewInt(100), big.NewInt(100))
// }

func getOrCreatPool(token0, token1 string, sqrtPrice *big.Int) *Pool {
	poolAddr := genPoolAddr(token0, token1)
	factory := "0x0000000000000000000000000000000000000002"

	mp, ok := tokenToPools[token0]
	if !ok {
		return createPool(poolAddr, factory, token0, token1, sqrtPrice)
	}

	pool, ok := mp[token1]
	if !ok {
		return createPool(poolAddr, factory, token0, token1, sqrtPrice)
	}

	return pool
}

func genPoolAddr(token0, token1 string) string {
	msg := []byte{255}
	msg = append(msg, common.HexToAddress(token0).Bytes()...)
	msg = append(msg, common.HexToAddress(token1).Bytes()...)

	hash := crypto.Keccak256(msg)

	return common.BytesToAddress(hash[12:]).Hex()
}

func createPoolWithZeroTickInitialized(token0, token1 string) (pool *Pool) {
	sqrtPrice := encodePriceSqrt(big.NewInt(1), big.NewInt(1))
	pool = getOrCreatPool(token0, token1, sqrtPrice)

	tickSpacing := 60
	getMinTick := func(tickSpacing int) int {
		return (-887272 / tickSpacing) * tickSpacing
	}
	getMaxTick := func(tickSpacing int) int {
		return (887272 / tickSpacing) * tickSpacing
	}
	tickLower := getMinTick(tickSpacing)
	tickUpper := getMaxTick(tickSpacing)
	// addLiquidity(pool, tickLower, tickUpper, big.NewInt(1000000), big.NewInt(1000000))
	// addLiquidity(pool, -60, 60, big.NewInt(100), big.NewInt(100))
	// addLiquidity(pool, -120, 120, big.NewInt(100), big.NewInt(100))

	addLiquidity(pool, tickLower, tickUpper, big.NewInt(1000000), big.NewInt(1000000))

	addLiquidity(pool, 0, 60, big.NewInt(100), big.NewInt(100))
	addLiquidity(pool, -120, 0, big.NewInt(100), big.NewInt(100))

	return pool
}

func createPoolWithMultiplePositions(token0, token1 string) (pool *Pool) {
	sqrtPrice := encodePriceSqrt(big.NewInt(1), big.NewInt(1))
	pool = getOrCreatPool(token0, token1, sqrtPrice)

	tickSpacing := 60
	getMinTick := func(tickSpacing int) int {
		return (-887272 / tickSpacing) * tickSpacing
	}
	getMaxTick := func(tickSpacing int) int {
		return (887272 / tickSpacing) * tickSpacing
	}
	tickLower := getMinTick(tickSpacing)
	tickUpper := getMaxTick(tickSpacing)
	addLiquidity(pool, tickLower, tickUpper, big.NewInt(1000000), big.NewInt(1000000))

	addLiquidity(pool, -60, 60, big.NewInt(100), big.NewInt(100))
	addLiquidity(pool, -120, 120, big.NewInt(100), big.NewInt(100))

	return
}

func createPoolWithInitLiquidity(token0, token1 string) (pool *Pool) {
	sqrtPrice := encodePriceSqrt(big.NewInt(1), big.NewInt(1))
	pool = getOrCreatPool(token0, token1, sqrtPrice)

	tickSpacing := 60
	getMinTick := func(tickSpacing int) int {
		return (-887272 / tickSpacing) * tickSpacing
	}
	getMaxTick := func(tickSpacing int) int {
		return (887272 / tickSpacing) * tickSpacing
	}
	tickLower := getMinTick(tickSpacing)
	tickUpper := getMaxTick(tickSpacing)
	addLiquidity(pool, tickLower, tickUpper, big.NewInt(1000000), big.NewInt(1000000))

	// addLiquidity(pool, -60, 60, big.NewInt(100), big.NewInt(100))
	// addLiquidity(pool, -120, 120, big.NewInt(100), big.NewInt(100))

	return
}

func addLiquidity(pool *Pool, tickLower, tickUpper int, amount0, amount1 *big.Int) {
	liquid := getLiquidityForAmounts(tickLower, tickUpper, pool.SqrtPriceX96, amount0, amount1)
	// evt := &types.Log{}
	pool.modifyPosition(tickLower, tickUpper, liquid, "Mint", "")
}

func createPool(poolAddr, factory, tokenA, tokenB string, sqrtPrice *big.Int) *Pool {
	if tokenA > tokenB {
		tokenA, tokenB = tokenB, tokenA
	}

	tick := getTickAtSqrtRatio(sqrtPrice)
	pool := CreateCAMMPool(tokenA, tokenB, poolAddr, factory, 3000, 60, 0)
	pool.OnInitialize(sqrtPrice, tick)

	return pool
}

// 0 -> 2 cross 2 tick.
func TestQuoteExactOutput1(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amountIn, sqrtPriceX96AfterList, crossedTicks := quoteExactOutput([]string{_token2, _token0}, toBigIntMust("15000"))
	assert.Equal(t, len(crossedTicks), 1)
	assert.Equal(t, crossedTicks[0], uint32(2))
	assert.Equal(t, toBigIntMust("15273"), amountIn)
	assert.Equal(t, len(sqrtPriceX96AfterList), 1)
	assert.Equal(t, sqrtPriceX96AfterList[0], toBigIntMust("78055527257643669242286029831"))
}

// 0 -> 2 cross 2 where tick after is initialized.
func TestQuoteExactOutput2(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amountIn, sqrtPriceX96AfterList, crossedTicks := quoteExactOutput([]string{_token2, _token0}, toBigIntMust("6143"))
	assert.Equal(t, len(crossedTicks), 1)
	assert.Equal(t, crossedTicks[0], uint32(1))
	assert.Equal(t, toBigIntMust("6200"), amountIn)
	assert.Equal(t, len(sqrtPriceX96AfterList), 1)
	assert.Equal(t, sqrtPriceX96AfterList[0], toBigIntMust("78757225449310403327341205211"))
}

// 0 -> 2 cross 1 tick.
func TestQuoteExactOutput3(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amountIn, sqrtPriceX96AfterList, crossedTicks := quoteExactOutput([]string{_token2, _token0}, toBigIntMust("4000"))
	assert.Equal(t, len(crossedTicks), 1)
	assert.Equal(t, crossedTicks[0], uint32(1))
	assert.Equal(t, toBigIntMust("4029"), amountIn)
	assert.Equal(t, len(sqrtPriceX96AfterList), 1)
	assert.Equal(t, sqrtPriceX96AfterList[0], toBigIntMust("78924219757724709840818372098"))
}

// 0 -> 2 cross 0 tick starting tick not initialized.
func TestQuoteExactOutput5(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amountIn, sqrtPriceX96AfterList, crossedTicks := quoteExactOutput([]string{_token2, _token0}, toBigIntMust("10"))
	assert.Equal(t, len(crossedTicks), 1)
	assert.Equal(t, crossedTicks[0], uint32(0))
	assert.Equal(t, toBigIntMust("12"), amountIn)
	assert.Equal(t, len(sqrtPriceX96AfterList), 1)
	assert.Equal(t, sqrtPriceX96AfterList[0], toBigIntMust("79227408033628034983534698435"))
}

// 2 -> 0 cross 2 ticks.
func TestQuoteExactOutput6(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amountIn, sqrtPriceX96AfterList, crossedTicks := quoteExactOutput([]string{_token0, _token2}, toBigIntMust("15000"))
	assert.Equal(t, len(crossedTicks), 1)
	assert.Equal(t, crossedTicks[0], uint32(2))
	assert.Equal(t, toBigIntMust("15273"), amountIn)
	assert.Equal(t, len(sqrtPriceX96AfterList), 1)
	assert.Equal(t, sqrtPriceX96AfterList[0], toBigIntMust("80418414376567919517220409857"))
}

// 2 -> 0 cross 2 where tick after is initialized.
func TestQuoteExactOutput7(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amountIn, sqrtPriceX96AfterList, crossedTicks := quoteExactOutput([]string{_token0, _token2}, toBigIntMust("6223"))
	assert.Equal(t, len(crossedTicks), 1)
	assert.Equal(t, crossedTicks[0], uint32(2))
	assert.Equal(t, toBigIntMust("6283"), amountIn)
	assert.Equal(t, len(sqrtPriceX96AfterList), 1)
	assert.Equal(t, sqrtPriceX96AfterList[0], toBigIntMust("79708304437530892332449657932"))
}

// 2 -> 0 cross 1 tick.
func TestQuoteExactOutput8(t *testing.T) {
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)

	amountIn, sqrtPriceX96AfterList, crossedTicks := quoteExactOutput([]string{_token0, _token2}, toBigIntMust("6000"))
	assert.Equal(t, len(crossedTicks), 1)
	assert.Equal(t, crossedTicks[0], uint32(1))
	assert.Equal(t, toBigIntMust("6055"), amountIn)
	assert.Equal(t, len(sqrtPriceX96AfterList), 1)
	assert.Equal(t, sqrtPriceX96AfterList[0], toBigIntMust("79690640184021170956740081887"))
}

// 2 -> 1.
func TestQuoteExactOutput9(t *testing.T) {
	clearTestPools()
	createPoolWithInitLiquidity(_token0, _token1)
	createPoolWithInitLiquidity(_token1, _token2)
	createPoolWithMultiplePositions(_token0, _token2)

	amountIn, sqrtPriceX96AfterList, crossedTicks := quoteExactOutput([]string{_token1, _token2}, toBigIntMust("9871"))
	assert.Equal(t, len(crossedTicks), 1)
	assert.Equal(t, crossedTicks[0], uint32(0))
	assert.Equal(t, toBigIntMust("10000"), amountIn)
	assert.Equal(t, len(sqrtPriceX96AfterList), 1)
	assert.Equal(t, sqrtPriceX96AfterList[0], toBigIntMust("80018020393569259756601362385"))
}

// 0 -> 2 -> 1.
func TestQuoteExactOutput10(t *testing.T) {
	clearTestPools()
	createPoolWithInitLiquidity(_token0, _token1)
	createPoolWithInitLiquidity(_token1, _token2)
	createPoolWithMultiplePositions(_token0, _token2)

	amountIn, sqrtPriceX96AfterList, crossedTicks := quoteExactOutput([]string{_token1, _token2}, toBigIntMust("9871"))
	assert.Equal(t, len(crossedTicks), 1)
	assert.Equal(t, crossedTicks[0], uint32(0))
	assert.Equal(t, toBigIntMust("10000"), amountIn)
	assert.Equal(t, len(sqrtPriceX96AfterList), 1)
	assert.Equal(t, sqrtPriceX96AfterList[0], toBigIntMust("80018020393569259756601362385"))
}

// 0 -> 1.
func TestQuoteExactOutputSingle1(t *testing.T) {
	clearTestPools()
	createPoolWithInitLiquidity(_token0, _token1)
	createPoolWithInitLiquidity(_token1, _token2)
	createPoolWithMultiplePositions(_token0, _token2)

	amountIn, sqrtPriceX96After, crossedTicks := quoteExactOutputSingle(_token0, _token1,
		toBigIntMust("340282366920938463463374607431768211455"), // uint128
		encodePriceSqrt(big.NewInt(100), big.NewInt(102)))

	assert.Equal(t, crossedTicks, uint32(0))
	assert.Equal(t, toBigIntMust("9981"), amountIn)
	assert.Equal(t, sqrtPriceX96After, toBigIntMust("78447570448055484695608110440"))
}

// 1 -> 0.
func TestQuoteExactOutputSingle2(t *testing.T) {
	clearTestPools()
	createPoolWithInitLiquidity(_token0, _token1)
	createPoolWithInitLiquidity(_token1, _token2)
	createPoolWithMultiplePositions(_token0, _token2)

	amountIn, sqrtPriceX96After, crossedTicks := quoteExactOutputSingle(_token1, _token0,
		toBigIntMust("340282366920938463463374607431768211455"), // uint128
		encodePriceSqrt(big.NewInt(102), big.NewInt(100)))

	assert.Equal(t, crossedTicks, uint32(0))
	assert.Equal(t, toBigIntMust("9981"), amountIn)
	assert.Equal(t, sqrtPriceX96After, toBigIntMust("80016521857016594389520272648"))
}

func TestQuoteExactOutputSingle3(t *testing.T) {
	// utils.SkipCI(t)
	clearTestPools()
	createPoolWithInitLiquidity(_token0, _token1)
	createPoolWithInitLiquidity(_token1, _token2)
	createPoolWithMultiplePositions(_token0, _token2)
	pool := createPoolWithZeroTickInitialized(_token0, _token2)
	_ = pool
	// println("pool liquid:", pool.Liquidity.String())
	// println("pool SqrtPriceX96:", pool.SqrtPriceX96.String())
	// println("pool tick:", pool.Tick)

	amountIn, sqrtPriceX96After, crossedTicks := quoteExactOutputSingle(_token0, _token2,
		toBigIntMust("100"), // uint128
		big.NewInt(0))
	// encodePriceSqrt(big.NewInt(102), big.NewInt(100)))

	assert.Equal(t, crossedTicks, uint32(1))
	assert.Equal(t, toBigIntMust("102"), amountIn)
	//                            79224359842948689147568835350
	assert.Equal(t, toBigIntMust("79224329176051641448521403903"), sqrtPriceX96After)
}

// failed!!!
// 0 -> 2 cross 0 tick, starting tick initialized.
func TestQuoteExactOutput4(t *testing.T) {
	// utils.SkipCI(t)
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)
	pool := createPoolWithZeroTickInitialized(_token0, _token2)
	_ = pool
	// println("pool liquid:", pool.Liquidity.String())
	// println("pool SqrtPriceX96:", pool.SqrtPriceX96.String())
	// println("pool tick:", pool.Tick)

	amountIn, sqrtPriceX96AfterList, crossedTicks := quoteExactOutput([]string{_token2, _token0}, toBigIntMust("100"))
	assert.Equal(t, len(crossedTicks), 1)
	assert.Equal(t, crossedTicks[0], uint32(1))
	assert.Equal(t, toBigIntMust("102"), amountIn)
	assert.Equal(t, len(sqrtPriceX96AfterList), 1)
	//                            79224359842948689147568835350
	assert.Equal(t, toBigIntMust("79224329176051641448521403903"), sqrtPriceX96AfterList[0])
}

// 0 -> 2 cross 0 tick, starting tick initialized todo failed!!!
func TestSwapExactInput5(t *testing.T) {
	// utils.SkipCI(t)
	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)
	pool := createPoolWithZeroTickInitialized(_token0, _token2)
	// println("pool liquid:", pool.Liquidity.String()) // 2083487
	// 79228162514264337593543950336
	// println("pool SqrtPriceX96:", pool.SqrtPriceX96.String())
	pool.printInfos()

	amt, sqrtPrices, crossTicks := quoteExactInput([]string{_token0, _token2}, big.NewInt(10))
	assert.Equal(t, amt, big.NewInt(8))
	assert.Equal(t, len(sqrtPrices), 1)
	//                           79227820275324292410333451302
	assert.Equal(t, toBigIntMust("79227817515327498931091950511"), sqrtPrices[0])
	// assert.Equal(t, len(crossTicks), 1)
	// println(pool.SqrtPriceX96.String())
	assert.Equal(t, crossTicks[0], uint32(1))
}

// 2 -> 0 cross 0 tick, starting tick initialized.
func TestSwapExactInput9(t *testing.T) {
	// utils.SkipCI(t)
	logger.Init("", "", true)

	clearTestPools()
	createPoolWithMultiplePositions(_token0, _token2)
	pool := createPoolWithZeroTickInitialized(_token0, _token2)
	_ = pool

	amt, sqrtPrices, crossTicks := quoteExactInput([]string{_token2, _token0}, big.NewInt(200))
	assert.Equal(t, amt, big.NewInt(198))
	assert.Equal(t, len(sqrtPrices), 1)
	assert.Equal(t, toBigIntMust("79235729830182478001034429156"), sqrtPrices[0])
	assert.Equal(t, len(crossTicks), 1)
	assert.Equal(t, uint32(0), crossTicks[0])
}

func TestSwapVerify(t *testing.T) {
	logger.Init("", "", false)

	pool := &Pool{
		PoolInfo: PoolInfo{
			Address:     "0x883e4ae0a817f2901500971b353b5dd89aa52184",
			TickSpacing: 10,
			Fee:         450,
		},
		Tick:       -196252,
		tickBitmap: map[int16]*big.Int{},
		Ticks: map[int]*TickInfo{
			-196270: {
				Tick:         -196270,
				LiquidityNet: toBigIntMust("65717988679110"),
			},
			-196250: {
				Tick:         -196250,
				LiquidityNet: toBigIntMust("229787366529207"),
			},
			-196240: {
				Tick:         -196240,
				LiquidityNet: toBigIntMust("1830021344017647"),
			},
			-196230: {
				Tick:         -196230,
				LiquidityNet: toBigIntMust("-295505355208317"),
			},
		},
	}
	pool.Liquidity, _ = big.NewInt(0).SetString("11557312308968112", 0)
	pool.SqrtPriceX96, _ = big.NewInt(0).SetString("4340592695139423692312076", 0)
	pool.flipTick(-196270)
	pool.flipTick(-196250)
	pool.flipTick(-196240)
	pool.flipTick(-196230)

	blockNumber := uint64(17951326)
	_ = blockNumber

	txhash := "0x46202d7ab799f312420c9c7ff760adbd106c52f4dd766fea7b90f8160d4ef3cf"
	tick := -196250
	sqrtPrice, _ := big.NewInt(0).SetString("4340980894691140664160317", 0)
	amount0, _ := big.NewInt(0).SetString("-18959857029948014", 0)
	amount1, _ := big.NewInt(0).SetString("56938816", 0)
	sqrtPriceLimitX96 := toBigIntMust("4340980894702089090613367")
	_ = sqrtPriceLimitX96

	// err := pool.VerifySwapResultByExact(txhash, tick, sqrtPrice, amount0, amount1, true, sqrtPriceLimitX96)
	// assert.Nil(t, err)
	pool.VerifySwapResult(blockNumber, txhash, tick, sqrtPrice, amount0, amount1)
}
