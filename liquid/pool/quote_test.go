package pool

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/utils"
)

func TestGetAmountOut(t *testing.T) {
	utils.SkipCI(t)

	tests := []struct {
		name       string
		pool       Pool
		zeroForOne bool
		amountIn   *big.Int
		expected   *big.Int
	}{
		{
			name: "AMM Variety - Non-Stable",
			pool: Pool{
				PoolInfo: PoolInfo{
					Typ: common.PoolTypeAMM,
					Fee: 3000, // 0.3%
				},
				Reserve0: big.NewInt(1000000),
				Reserve1: big.NewInt(2000000),
			},
			zeroForOne: true,
			amountIn:   big.NewInt(100000),
			expected:   big.NewInt(181322), // Calculated manually
		},
		// Add more test cases for different pool types and scenarios
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pool.GetAmountOut(tt.zeroForOne, tt.amountIn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetStableOut(t *testing.T) {
	utils.SkipCI(t)

	tests := []struct {
		name       string
		pool       Pool
		zeroForOne bool
		amountIn   *big.Int
		expected   *big.Int
	}{
		{
			name: "Stable Pool - Zero for One",
			pool: Pool{
				PoolInfo: PoolInfo{
					Typ:       common.PoolTypeAeroAMM,
					Stable:    true,
					decimals0: big.NewInt(18),
					decimals1: big.NewInt(18),
					Fee:       100,
				},
				Reserve0: big.NewInt(1000000000000000000),
				Reserve1: big.NewInt(1000000000000000000),
			},
			zeroForOne: true,
			amountIn:   big.NewInt(10000000),
			expected:   big.NewInt(9998999), // Approximate value, may need adjustment
		},
		// Add more test cases for different scenarios
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pool.getStableOut(tt.zeroForOne, tt.amountIn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAmountOutV2(t *testing.T) {
	utils.SkipCI(t)

	tests := []struct {
		name       string
		pool       Pool
		zeroForOne bool
		amountIn   *big.Int
		expected   *big.Int
	}{
		{
			name: "Uniswap V2 - Zero for One",
			pool: Pool{
				Reserve0: big.NewInt(1000000),
				Reserve1: big.NewInt(2000000),
				PoolInfo: PoolInfo{
					Fee: 3000, // 0.3%
				},
			},
			zeroForOne: true,
			amountIn:   big.NewInt(100000),
			expected:   big.NewInt(181322), // Calculated manually
		},
		// Add more test cases for different scenarios
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pool.GetAmountOutV2(tt.zeroForOne, tt.amountIn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUpdateBestPrice(t *testing.T) {
	utils.SkipCI(t)

	tests := []struct {
		name           string
		pool           Pool
		expectedPrice0 decimal.Decimal
		expectedPrice1 decimal.Decimal
	}{
		{
			name: "AMM Variety - Non-Stable",
			pool: Pool{
				PoolInfo: PoolInfo{
					Typ: common.PoolTypeAMM,
					Fee: 3000, // 0.3%
				},
				Reserve0: big.NewInt(1000000),
				Reserve1: big.NewInt(2000000),
			},
			expectedPrice0: decimal.NewFromFloat(1.994), // Approximate value, may need adjustment
			expectedPrice1: decimal.NewFromFloat(0.4985),
		},
		// Add more test cases for different pool types and scenarios
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pool.UpdateBestPrice()
			assert.True(t, tt.expectedPrice0.Equal(*tt.pool.price0))
			assert.True(t, tt.expectedPrice1.Equal(*tt.pool.price1))
		})
	}
}

func TestStablePoolPrice(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "dev", false)

	rc := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "",
	})

	address := []string{
		"0x6d0b9c9e92a3de30081563c3657b5258b3ffa38b", // usdz/usdc
		"0x468de3406a083d9cd976dae505293126a1f00b29",
		"0xe6307ca2e717668176280d2333d70b68600a73da",
	}

	logger.Info().Msgf("big.Int 3/2=%v", new(big.Int).Div(big.NewInt(3), big.NewInt(2)))
	logger.Info().Msgf("big.Int 5/3=%v", new(big.Int).Div(big.NewInt(5), big.NewInt(3)))

	for _, addr := range address {
		p, err := getPoolFromRedis(rc, addr)
		assert.Nil(t, err)

		p.UpdateBestPrice()
		logger.Info().Msgf("pool price %v %v", p.price0, p.price1)
	}
}

func getPoolFromRedis(rc *redis.Client, addr string) (*Pool, error) {
	addr = strings.ToLower(addr)

	val, err := rc.HGet(context.Background(), "poolLiquidity", addr).Result()
	if err != nil {
		return nil, err
	}

	var p Pool

	err = json.Unmarshal([]byte(val), &p)
	if err != nil {
		return nil, err
	}

	tok0, err := getTokenFromRedis(rc, p.Token0)
	if err != nil {
		return nil, err
	}

	tok1, err := getTokenFromRedis(rc, p.Token1)
	if err != nil {
		return nil, err
	}

	p.decimals0 = new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(tok0.Decimals)), nil) // nolint
	p.decimals1 = new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(tok1.Decimals)), nil) // nolint

	return &p, err
}

func getTokenFromRedis(rc *redis.Client, addr string) (*common.Token, error) {
	addr = strings.ToLower(addr)

	val, err := rc.HGet(context.Background(), "tokenInfo", addr).Result()
	if err != nil {
		return nil, err
	}

	var tok common.Token
	err = json.Unmarshal([]byte(val), &tok)

	return &tok, err
}

func TestPrice01(t *testing.T) {
	logger.Init("", "", false)
	utils.SkipCI(t)

	s2 := "{\"address\":\"0x806eab3b2f63343da07fe3c462a0b38a8bec5fd9\",\"token0\":\"0x4200000000000000000000000000000000000006\",\"token1\":\"0xb79dd08ea68a908a97220c76d19a6aa9cbde4376\",\"factory\":\"0x38015d05f4fec8afe15d7cc0386a126574e8077b\",\"vendor\":\"BaseSwap-Basex\",\"typ\":300,\"fee\":450,\"stable\":false,\"tickSpacing\":10,\"lastBlockUpdated\":19351112,\"initBlock\":18403283,\"reserve0\":978364187031964367,\"reserve1\":16567173,\"tick\":-887272,\"liquidity\":0,\"sqrtPriceX96\":4295128740,\"ticks\":{\"-194330\":{\"tick\":-194330,\"liquidityNet\":-3259878533457,\"liquidityGross\":3259878533457},\"-194440\":{\"tick\":-194440,\"liquidityNet\":-1211780679066012,\"liquidityGross\":1211780679066012},\"-195380\":{\"tick\":-195380,\"liquidityNet\":1211780679066012,\"liquidityGross\":1211780679066012},\"-195590\":{\"tick\":-195590,\"liquidityNet\":3259878533457,\"liquidityGross\":3259878533457}},\"synced\":true,\"initialized\":true}"

	var pl Pool
	err := json.Unmarshal([]byte(s2), &pl)
	assert.Nil(t, err)

	pl.Reload()
	pl.UpdateBestPrice()
	t.Logf("price0: %v price1: %v", pl.price0, pl.price1)
}
