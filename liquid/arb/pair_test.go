package arb

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/defiweb/go-eth/hexutil"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/contracts/arbitrage"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/utils"
)

func TestPairList(t *testing.T) {
	utils.SkipCI(t)

	lpservice := createTestLPService(t, true)
	// lpservice.LoadTokenPools()
	t.Logf("pools: %d tokens: %d", len(lpservice.pools), len(lpservice.tokens))

	arbPairs := lpservice.DiscoverArbPairList()
	t.Logf("arb pairs: %d", len(arbPairs))

	minProfit := d_one.Add(d_00001) // .Add(d_one)
	for _, pair := range arbPairs {
		profitable, param, err := lpservice.CalcProfitable(pair, pair.Pairs, minProfit, true)
		if profitable && err == nil {
			gas, err := lpservice.Estimate(param)
			if err != nil {
				continue
			}

			logger.Info().Msgf("gas: %v error: %v", gas, err)

			logger.Info().Msgf("arb pair profitable: %s %v", lpservice.getPoolName(pair.Pairs[0]), param.ratio.StringFixed(5))

			pool0 := lpservice.pools[param.pool0]
			pool1 := lpservice.pools[param.pool1]
			logger.Info().Msgf("zeroForOne: %v pool0: %v pool1: %v token0: %v token1: %v",
				param.zeroForOne, pool0.Address, pool1.Address, pool0.Token0, pool0.Token1)
			logger.Info().Msgf("pool0 Reserves: %v %v pool1 Reserves: %v %v",
				pool0.Reserve0, pool0.Reserve1, pool1.Reserve0, pool1.Reserve1)

			logger.Info().Msgf("bestAmt: %v bestOut: %v profit: %v", param.bestAmtIn, param.bestAmtOut, param.bestProfit)

			for _, p := range []*pool.Pool{pool0, pool1} {
				logger.Info().Msgf("pool: %s %s %s price0: %s price1: %s TVL: %v",
					p.Address, p.Vendor, p.GetName(), p.GetPrice0(), p.GetPrice1(), readableETHAmount(p.GetTVL()))
			}
		}
	}

	/*
		for _, p := range lpservice.pools {
			if p.Reserve0.Cmp(bigZero) < 0 || p.Reserve1.Cmp(bigZero) < 0 {
				logger.Info().Msgf("pool %v %v reserves invalid: %v %v", p.Address, p.Vendor, p.Reserve0, p.Reserve1)
			}
			if p.Typ.IsCAMMVariety() {
				if p.Liquidity.Cmp(bigZero) <= 0 &&
					(lpservice.isStableToken(p.Token0) || isNativeOrWrapperNativeToken(p.Token0) ||
						lpservice.isStableToken(p.Token0) || isNativeOrWrapperNativeToken(p.Token0)) &&
					(p.Reserve0.Cmp(bigZero) > 0 || p.Reserve1.Cmp(bigZero) > 0) {
					logger.Info().Msgf("pool %v %v has reservers: %v %v", p.Address, p.Vendor, p.Reserve0, p.Reserve1)
				}
			}
		}
	*/
	_ = minProfit
}

// func TestCalculateArbPairs(t *testing.T) {
// 	lpservice := CreateLPService(&config.Config{})
// }

func TestMaxProfit(t *testing.T) {
	utils.SkipCI(t)

	lpservice := createTestLPService(t, true)
	// lpservice.LoadTokenPools()
	t.Logf("pools: %d tokens: %d", len(lpservice.pools), len(lpservice.tokens))

	pool0 := lpservice.pools["0x14e0d45c7b0d82e226990e9ddf260e06bb9cd78a"]
	pool1 := lpservice.pools["0xe2807f2058d6c618d20267c71fa0316c9ab6cb0c"]
	token0 := lpservice.tokens[pool0.Token0]
	token1 := lpservice.tokens[pool1.Token1]

	pair := &ArbPairList{
		Pairs:  []*pool.Pool{pool0, pool1},
		Token0: token0,
		Token1: token1,
	}

	_, param, _ := lpservice.CalcProfitable(pair, []*pool.Pool{pool0, pool1}, d_one.Add(d_one), true)
	logger.Info().Msgf("ratio: %v", param.ratio)
	logger.Info().Msgf("pool0 reserves: %s %s %v %v prices: %v %v",
		pool0.Vendor, lpservice.getPoolName(pool0), pool0.Reserve0, pool0.Reserve1, pool0.GetPrice0(), pool0.GetPrice1())
	logger.Info().Msgf("pool1 reserves: %s %s %v %v prices: %v %v",
		pool1.Vendor, lpservice.getPoolName(pool1), pool1.Reserve0, pool1.Reserve1, pool1.GetPrice0(), pool1.GetPrice1())

	// amountIn, bestAmt, bestProfit := lpservice.optimizeMaxAmountIn(pool0, pool1, ratio, 0)
	logger.Info().Msgf("bestAmtIn: %v bestProfit: %v",
		param.bestAmtIn, param.bestProfit)
}

func TestEstimate(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "dev", false)

	amt, _ := new(big.Int).SetString("32310385150800423", 10)
	lpservice := createTestLPService(t, false)

	param := &ArbParams{
		zeroForOne: true,
		poolLoan:   "0x20E068D76f9E90b90604500B84c7e19dCB923e7e", // WETH/wstETH-0.01
		pool0:      "0x14e0d45c7b0d82e226990e9ddf260e06bb9cd78a",
		pool1:      "0xe2807f2058d6c618d20267c71fa0316c9ab6cb0c",
		tokenIn:    "0x4200000000000000000000000000000000000006",
		tokenOut:   "0xBA5eDF631828EBbe81B850F476FA5936e3C15783",
		poolType0:  300,
		poolType1:  200,
		bestAmtIn:  amt,
	}
	gas, err := lpservice.Estimate(param)

	logger.Info().Msgf("gas: %v error: %v", gas, err)
}

func TestABI(t *testing.T) {
	a, err := arbitrage.ArbitrageContractMetaData.GetAbi()
	assert.Nil(t, err)

	data, err := hexutil.HexToBytes("0x0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000bf5f6359745298ee65c52a32be5fe9fafb1b434d00000000000000000000000004eab58dd2cc60580c6b866983c87a561079d79000000000000000000000000042000000000000000000000000000000000000060000000000000000000000004e496c0256fb9d4cc7ba2fdf931bc9cbb773166000000000000000000000000000000000000000000000000000000000000000c800000000000000000000000000000000000000000000000000000000000000c80000000000000000000000000000000000000000000000000043a694138dfb1c")
	assert.Nil(t, err)

	// Otherwise pack up the parameters and invoke the contract
	input, err := a.Pack("swapFromV2", data)
	assert.Nil(t, err)

	fmt.Println(hexutil.BytesToHex(input))
}

func TestSignature(t *testing.T) {
	t.Logf("swapFromV2: %v", methodSwapFromV2.FourBytes())
	t.Logf("swapFromV3: %v", methodSwapFromV3.FourBytes())
}

func TestArbCheckPoolPos(t *testing.T) {
	arb := &ArbPairList{
		ZeroForOne: true,
		Pairs:      []*pool.Pool{},
	}

	pos := 0
	pos++
	pl1 := createPair(3, pos)
	n := arb.checkPoolPos(pl1)
	t.Logf("pl1 pos: %d  pairs: %d", n, len(arb.Pairs))
	printArbPairs(arb)

	pos++
	pl2 := createPair(10, pos)
	n = arb.checkPoolPos(pl2)
	t.Logf("pl2 pos: %d  pairs: %d", n, len(arb.Pairs))
	printArbPairs(arb)

	pos++
	pl3 := createPair(5, pos)
	n = arb.checkPoolPos(pl3)
	t.Logf("pl3 pos: %d  pairs: %d", n, len(arb.Pairs))
	printArbPairs(arb)

	pos++
	pl4 := createPair(15, pos)
	n = arb.checkPoolPos(pl4)
	t.Logf("pl4 pos: %d  pairs: %d", n, len(arb.Pairs))
	printArbPairs(arb)

	pos++
	pl5 := createPair(2, pos)
	n = arb.checkPoolPos(pl5)
	t.Logf("pl5 pos: %d  pairs: %d", n, len(arb.Pairs))
	printArbPairs(arb)
}

func TestArbCheckPoolPos10(t *testing.T) {
	arb := &ArbPairList{
		ZeroForOne: false,
		Pairs:      []*pool.Pool{},
	}

	pos := 0
	pos++
	pl1 := createPair(3, pos)
	n := arb.checkPoolPos(pl1)
	t.Logf("pl1 pos: %d  pairs: %d", n, len(arb.Pairs))
	printArbPairs(arb)

	pos++
	pl2 := createPair(10, pos)
	n = arb.checkPoolPos(pl2)
	t.Logf("pl2 pos: %d  pairs: %d", n, len(arb.Pairs))
	printArbPairs(arb)

	pos++
	pl3 := createPair(5, pos)
	n = arb.checkPoolPos(pl3)
	t.Logf("pl3 pos: %d  pairs: %d", n, len(arb.Pairs))
	printArbPairs(arb)

	pos++
	pl4 := createPair(15, pos)
	n = arb.checkPoolPos(pl4)
	t.Logf("pl4 pos: %d  pairs: %d", n, len(arb.Pairs))
	printArbPairs(arb)

	pos++
	pl5 := createPair(2, pos)
	n = arb.checkPoolPos(pl5)
	t.Logf("pl5 pos: %d  pairs: %d", n, len(arb.Pairs))
	printArbPairs(arb)
}

func printArbPairs(arb *ArbPairList) {
	for _, pl := range arb.Pairs {
		fmt.Printf("[%v %v] ", pl.GetPrice0(), pl.Address)
	}
	fmt.Println("\n--------")
}

func createPair(price0 float64, idx int) *pool.Pool {
	p := &pool.Pool{
		PoolInfo: pool.PoolInfo{Address: fmt.Sprint(idx)},
	}
	p0 := decimal.NewFromFloat(price0)
	p1 := decimal.NewFromFloat(1 / price0)
	p.SetPrice0(&p0)
	p.SetPrice1(&p1)
	return p
}
