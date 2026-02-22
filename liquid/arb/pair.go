package arb

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"sort"

	"github.com/defiweb/go-eth/abi"
	"github.com/ethereum/go-ethereum"
	ethcomm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

type ArbPair struct {
	Pools []string
}

type ArbParams struct {
	ratio      decimal.Decimal
	zeroForOne bool
	// fromV3     bool
	poolLoan   string
	pool0      string
	pool1      string
	tokenIn    string
	tokenOut   string
	poolType0  uint32
	poolType1  uint32
	bestAmtIn  *big.Int
	bestAmtOut *big.Int
	bestProfit *big.Int
}

var (
	methodSwapFromV2 = abi.MustParseMethod("function swapFromV2(bytes calldata data) external")
	methodSwapFromV3 = abi.MustParseMethod("function swapFromV3(bytes calldata data) external")
)

type ArbPairList struct {
	Pairs      []*pool.Pool
	Token0     *common.Token
	Token1     *common.Token
	TokenIn    string
	TokenOut   string
	ZeroForOne bool
}

// DiscoverArbPairList discover arb pairs.
func (lpservice *LPService) DiscoverArbPairList() map[string]*ArbPairList {
	tokens := lpservice.tokens
	blackTokens := lpservice.blacklistTokens
	arbPairsMap := map[string]*ArbPairList{}

	for _, p := range lpservice.pools {
		if !lpservice.poolReservesGreatThan(p, minTokenETHReserve) {
			continue
		}

		if pairs, ok := arbPairsMap[p.GetKey()]; ok {
			pairs.Pairs = append(pairs.Pairs, p)
		} else {
			if blackTokens[p.Token0] || blackTokens[p.Token1] {
				continue
			}

			token0 := tokens[p.Token0]
			token1 := tokens[p.Token1]

			if token0 == nil || token1 == nil {
				// logger.Warn().Msgf("token0 %s or token1 %s is null", p.Token0, p.Token1)
				continue
			}

			zeroForOne := true
			tokenIn := token0.Address
			tokenOut := token1.Address
			quote := lpservice.tokenQuotable(p.Token0, p.Token1)

			if quote != p.Token0 {
				// swap from pool 1: token1 -> token0; token0 -> token1
				zeroForOne = false
				tokenIn = token1.Address
				tokenOut = token0.Address
			}
			arbPairsMap[p.GetKey()] = &ArbPairList{
				Pairs:      []*pool.Pool{p},
				Token0:     token0,
				Token1:     token1,
				TokenIn:    tokenIn,
				TokenOut:   tokenOut,
				ZeroForOne: zeroForOne,
			}
		}
	}

	return arbPairsMap
}

func (arb *ArbPairList) checkPoolPos(pl *pool.Pool) int {
	for i, p := range arb.Pairs {
		if p.Address == pl.Address {
			return i
		}
	}
	logger.Warn().Msgf("not found pool %s in arb pairs", pl.Address)
	n := sort.Search(len(arb.Pairs), func(i int) bool {
		if arb.ZeroForOne {
			// from high to low
			return arb.Pairs[i].GetPrice0().Cmp(pl.GetPrice0()) < 0
		} else {
			return arb.Pairs[i].GetPrice0().Cmp(pl.GetPrice0()) > 0
		}
	})

	arb.Pairs = slices.Insert(arb.Pairs, n, pl)

	return n
}

func (lpservice *LPService) CalcProfitable(arb *ArbPairList, pls []*pool.Pool, ratio decimal.Decimal, debug bool) (bool, *ArbParams, error) {
	if arb.ZeroForOne {
		// from high to low
		sort.Sort(sort.Reverse(PoolPriceSlice(arb.Pairs)))
	} else {
		sort.Sort(PoolPriceSlice(arb.Pairs))
	}

	if !arb.ProfitGreatThan(ratio) {
		return false, &ArbParams{}, nil
	}

	var (
		bestPool0  *pool.Pool
		bestPool1  *pool.Pool
		bestAmtIn  = big.NewInt(0)
		bestAmtOut = big.NewInt(0)
		bestProfit = big.NewInt(0)
	)
	// plast := arb.Pairs[len(arb.Pairs)-1].GetPrice1()

	for i := 0; i < len(pls); i++ {
		curr := pls[i]
		idx := arb.checkPoolPos(curr)

		for k := 0; k < idx; k++ {
			pool0 := arb.Pairs[k]
			if (arb.ZeroForOne && pool0.GetPrice0().Mul(curr.GetPrice1()).Cmp(ratio) < 0) ||
				((!arb.ZeroForOne) && pool0.GetPrice1().Mul(curr.GetPrice0()).Cmp(ratio) < 0) {
				break
			}
			amtIn, amtOut, profit := lpservice.GetOptimalArbAmount(pool0, curr, arb.ZeroForOne, ratio)
			if debug && profit.Cmp(bigZero) > 0 {
				logger.Info().Msgf("pair profitable: pool0=%v pool1=%v amtIn=%v amtOut=%v profit=%v",
					pool0.Address, curr.Address, amtIn, amtOut, profit)
			}
			if profit.Cmp(bestProfit) > 0 {
				bestAmtIn = amtIn
				bestAmtOut = amtOut
				bestProfit = profit
				bestPool0 = pool0
				bestPool1 = curr
			}
		}

		for j := len(arb.Pairs) - 1; j > idx; j-- {
			pool1 := arb.Pairs[j]
			if (arb.ZeroForOne && curr.GetPrice0().Mul(pool1.GetPrice1()).Cmp(ratio) < 0) ||
				(!arb.ZeroForOne && curr.GetPrice1().Mul(pool1.GetPrice0()).Cmp(ratio) < 0) {
				break
			}

			amtIn, amtOut, profit := lpservice.GetOptimalArbAmount(curr, pool1, arb.ZeroForOne, ratio)
			if debug && profit.Cmp(bigZero) > 0 {
				logger.Info().Msgf("pair profitable: pool0=%v pool1=%v amtIn=%v amtOut=%v profit=%v",
					curr.Address, pool1.Address, amtIn, amtOut, profit)
			}
			if profit.Cmp(bestProfit) > 0 {
				bestAmtIn = amtIn
				bestAmtOut = amtOut
				bestProfit = profit
				bestPool0 = curr
				bestPool1 = pool1
			}
		}
	}

	if bestPool0 != nil {
		return true, &ArbParams{
			ratio:      ratio,
			zeroForOne: arb.ZeroForOne,
			pool0:      bestPool0.Address,
			pool1:      bestPool1.Address,
			tokenIn:    arb.TokenIn,
			tokenOut:   arb.TokenOut,
			poolType0:  uint32(bestPool0.Typ), // nolint
			poolType1:  uint32(bestPool0.Typ), // nolint
			bestAmtIn:  bestAmtIn,
			bestAmtOut: bestAmtOut,
			bestProfit: bestProfit,
		}, nil
	}
	return false, nil, nil
	/*
		var (
			p0    = decimal.Zero
			p1    = decimal.Zero
		)

		for _, p := range arb.Pairs {
			if lpservice.poolReservesGreatThan(p, minTokenETHReserve) {
				if p.GetPrice0().Cmp(p0) > 0 {
					p0 = p.GetPrice0()
					pool0 = p
				}

				if p.GetPrice1().Cmp(p1) > 0 {
					p1 = p.GetPrice1()
					pool1 = p
				}
			}
		}

		// p0 := arb.Pairs[0]
		// p1 := arb.Pairs[len(arb.Pairs)-1]
		result := p0.Mul(p1)
		if result.Cmp(ratio) < 0 {
			return false, nil, nil
		}

		return lpservice.optimizeMaxAmountIn(pool0, pool1, result, 0)
	*/
}

//lint:ignore U1000 Ignore unused function temporarily for debugging.
func (lpservice *LPService) optimizeMaxAmountIn(
	pool0, pool1 *pool.Pool,
	ratio decimal.Decimal,
	maxIterations int,
) (bool, *ArbParams, error) {
	// optimizeMaxAmountIn calculates the best amount to trade for maximum profit.
	// Define constants
	// epsilon := decimal.NewFromFloat(1e-6)
	if maxIterations <= 0 {
		maxIterations = 100
	}

	zeroForOne := true
	quote := lpservice.tokenQuotable(pool0.Token0, pool0.Token1)

	if quote != pool0.Token0 {
		// swap from pool 1: token1 -> token0; token0 -> token1
		zeroForOne = false
	}

	var reserveIn *big.Int

	if zeroForOne {
		if pool0.GetPrice0().Cmp(pool1.GetPrice0()) < 0 {
			pool0, pool1 = pool1, pool0
		}

		reserveIn = pool0.Reserve0
	} else {
		if pool0.GetPrice1().Cmp(pool1.GetPrice1()) < 0 {
			pool0, pool1 = pool1, pool0
		}

		reserveIn = pool0.Reserve1
	}

	half := new(big.Int).Div(reserveIn, big.NewInt(2))
	amountIn := ratio.Sub(d_one).Div(d_two).Mul(decimal.NewFromBigInt(reserveIn, 0)).BigInt()

	if amountIn.Cmp(half) > 0 {
		amountIn = half
	}

	step := new(big.Int).Set(amountIn)
	bestAmt := big.NewInt(0)
	bestOut := big.NewInt(0)
	bestProfit := big.NewInt(0)
	lowest := big.NewInt(100)
	highest := new(big.Int).Sub(reserveIn, big.NewInt(100))

	// logger.Info().Msgf("zeroForOne: %v pool0: %v pool1: %v ratio: %v",
	// 	zeroForOne, pool0.Address, pool1.Address, ratio)

	prevDir := true

	for i := 0; i < maxIterations; i++ {
		output0, priceA := pool0.GetAmountOutAndPrice(zeroForOne, amountIn, zeroForOne)
		if output0.Cmp(bigZero) == 0 {
			break
		}

		output1, priceB := pool1.GetAmountOutAndPrice(!zeroForOne, output0, zeroForOne)
		profit := new(big.Int).Sub(output1, amountIn)

		if step.Cmp(lowest) <= 0 {
			break
		}

		if profit.Cmp(bestProfit) > 0 {
			bestProfit = profit
			bestAmt = amountIn
			bestOut = output0
		}

		if priceA.Cmp(decimal.Zero) == 0 || priceB.Cmp(decimal.Zero) == 0 {
			logger.Warn().Msgf("price zero: i=%d pool0: %s %s %v pool1: %s %s %v priceA: %v priceB: %v, amountIn=%v output0=%v output1=%v",
				i, pool0.Vendor, lpservice.getPoolName(pool0), pool0.Address,
				pool1.Vendor, lpservice.getPoolName(pool1), pool1.Address,
				priceA, priceB, amountIn, output0, output1)

			return false, &ArbParams{ratio: ratio}, fmt.Errorf("price zero")
		}

		// logger.Info().Msgf("i=%d: amountIn: %v step=%v output0: %v output1: %v priceA: %v priceB: %v profit: %v",
		// i, amountIn, step, output0, output1, priceA, priceB, profit)
		// if zeroForOne {
		if priceA.Cmp(priceB) > 0 {
			amt := new(big.Int).Add(amountIn, step)
			if amt.Cmp(highest) >= 0 || !prevDir {
				step = new(big.Int).Div(step, bigTwo)
				amt = new(big.Int).Add(amountIn, step)
			}

			amountIn = amt
			prevDir = true
		} else {
			amt := new(big.Int).Sub(amountIn, step)
			if amt.Cmp(lowest) <= 0 || prevDir {
				step = new(big.Int).Div(step, bigTwo)
				amt = new(big.Int).Sub(amountIn, step)
			}

			amountIn = amt
			prevDir = false
		}
	}

	if bestProfit.Cmp(bigZero) <= 0 {
		return false, nil, fmt.Errorf("no profits")
	}

	if (zeroForOne && bestOut.Cmp(pool0.Reserve1) > 0) ||
		((!zeroForOne) && bestOut.Cmp(pool0.Reserve0) > 0) {
		return false, nil, fmt.Errorf("pool 0 has not enought reserve1")
	}

	var (
		tokenIn  string
		tokenOut string
		poolLoan string
		fromV3   bool = pool0.Typ.IsCAMMVariety()
	)

	if zeroForOne {
		tokenIn = pool0.Token0
		tokenOut = pool0.Token1
	} else {
		tokenIn = pool0.Token1
		tokenOut = pool0.Token0
	}

	if fromV3 {
		poolLoan = lpservice.getPoolFlash(tokenIn, pool0.Address)
	}

	// Return the best amount found
	param := &ArbParams{
		ratio:      ratio,
		zeroForOne: zeroForOne,
		// fromV3:     fromV3,
		poolLoan:  poolLoan,
		pool0:     pool0.Address,
		pool1:     pool1.Address,
		tokenIn:   tokenIn,
		tokenOut:  tokenOut,
		poolType0: getPoolSwapType(pool0),
		poolType1: getPoolSwapType(pool1),
		// amountIn:   amountIn,
		bestAmtIn:  bestAmt,
		bestAmtOut: bestOut,
		bestProfit: bestProfit,
	}

	return true, param, nil
}

func getPoolSwapType(pl *pool.Pool) uint32 {
	switch pl.Typ {
	case common.PoolTypeAMM: // 200
		if pl.Vendor == "PancakeswapV2" {
			return 202 // pancake
		}

		if pl.Vendor == "Baseswap" {
			return 201
		}

		return 200

	case common.PoolTypeAeroAMM: //  = PoolType(201)
		return 201

	case common.PoolTypeInfusionAMM: // = PoolType(202)
		return 201

	case common.PoolTypeCAMM: // = PoolType(300)
		return 300
	case common.PoolTypeAeroCAMM: // = PoolType(301)
		return 301
	case common.PoolTypePancakeCAMM: //
		return 302
	default:
		panic("unknonw pool type: " + fmt.Sprintf("%d", uint(pl.Typ)))
	}
}

func (lpservice *LPService) getPoolFlash(tokenIn string, fromPool string) string {
	poolAddr := lpservice.flashPools[tokenIn]
	if poolAddr != "" && poolAddr != fromPool {
		return poolAddr
	}

	poolList := []*pool.Pool{}

	for _, p := range lpservice.poolsGreatThanMinETH {
		if p.Typ.IsCAMMVariety() && (p.Token0 == tokenIn || p.Token1 == tokenIn) && p.Address != fromPool {
			poolList = append(poolList, p)
		}
	}

	if len(poolList) == 0 {
		return ""
	}

	sort.Sort(sort.Reverse(PoolTVLSlice(poolList)))

	for _, p := range poolList {
		if p.Fee <= 100 {
			return p.Address
		}
	}

	return poolList[0].Address
}

func (lpservice *LPService) Estimate(param *ArbParams) (uint64, error) {
	var (
		input []byte
		err   error
	)

	if param.poolLoan != "" {
		// bool,
		// address,
		// address,
		// address,
		// address,
		// address,
		// uint32,
		// uint32,
		// uint256
		typ := abi.MustParseType("(bool,address,address,address,address,address,uint32,uint32,uint256)")
		input, err = abi.EncodeValues(typ,
			param.zeroForOne,
			param.poolLoan,
			param.pool0,
			param.pool1,
			param.tokenIn,
			param.tokenOut,
			param.poolType0,
			param.poolType1,
			param.bestAmtIn,
		)
		input = methodSwapFromV3.MustEncodeArgs(input)
	} else {
		typ := abi.MustParseType("(bool,address,address,address,address,uint32,uint32,uint256)")
		input, err = abi.EncodeValues(typ,
			param.zeroForOne,
			param.pool0,
			param.pool1,
			param.tokenIn,
			param.tokenOut,
			param.poolType0,
			param.poolType1,
			param.bestAmtIn,
		)
		input = methodSwapFromV2.MustEncodeArgs(input)
	}

	// println("method with input: " + hexutil.BytesToHex(input))
	// println("input with method: " + hexutil.BytesToHex(methodSwapFromV2.MustEncodeArgs(input)))
	if err != nil {
		return 0, err
	}

	client, err := ethclient.Dial(lpservice.cfg.Chain.URL)
	if err != nil {
		return 0, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return 0, err
	}

	// logger.Info().Msgf("gas price: %v", gasPrice)

	to := ethcomm.HexToAddress(lpservice.cfg.Arb.ArbContract)

	gas, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     ethcomm.HexToAddress(lpservice.cfg.Arb.Sender),
		To:       &to,
		GasPrice: gasPrice,
		Data:     input,
	})
	if err != nil {
		return 0, err
	}

	return gas, nil
}

func (arb *ArbPairList) ProfitGreatThan(ratio decimal.Decimal) bool {
	p0 := arb.Pairs[0].GetPrice0()
	p1 := arb.Pairs[len(arb.Pairs)-1].GetPrice1()

	return p0.Mul(p1).Cmp(ratio) > 0
}
