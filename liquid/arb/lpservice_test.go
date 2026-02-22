package arb

import (
	"errors"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/utils"
)

func createTestLPService(t *testing.T, load bool) *LPService {
	logger.Init("", "dev", true)

	lpservice := CreateLPService(&config.Config{
		Chain: config.ChainConfig{
			ChainID: "8453",
			URL:     "https://mainnet.base.org",
		},
		Arb: config.ArbConfig{
			Sender:      "0xfEb3509b7099Db900995e964f4586043A3C4BBF1",
			ArbContract: "0xCf7ac26F87a58df9dCBe2eC3E4Cc0697fD2Be419",
			MinPoolETH:  "0.001",
			StablePools: []string{
				"0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
				"0xd9aAEc86B65D86f6A7B5B1b0c42FFA531710b6CA",
				"0x50c5725949A6F0c72E6C4a641F24049A917DB0Cb",
				"0x4621b7A9c75199271F773Ebd9A499dbd165c3191",
			},
			NativeStablePair: "0xd0b53D9277642d899DF5C87A3966A349A798F224",
			BlackListTokens: []string{
				"0x7f12d13b34f5f4f0a9449c16bcd42f0da47af200", // NORMIE
				"0x4e496c0256fb9d4cc7ba2fdf931bc9cbb7731660", // BOGE
			},
		},
	})

	if load {
		err := lpservice.LoadTokenPools()
		assert.Nil(t, err)
	}

	lpservice.InitPools()

	return lpservice
}

func TestCalculatePriceByPool(t *testing.T) {
	utils.SkipCI(t)

	lpservice := createTestLPService(t, true)

	fixtures := []struct {
		poolAddr string
		ethQuote bool
	}{
		{
			poolAddr: "0x72ab388e2e2f6facef59e3c3fa2c4e29011c2d38", // weth/usdc-0.01% pancake v3
			ethQuote: false,
		},
		{
			poolAddr: "0xcDAC0d6c6C59727a65F871236188350531885C43", // weth/usdc-0.3% aerodrome v2
			ethQuote: true,
		},
		{
			poolAddr: "0xb37642e87613d8569fd8ec80888ea6c63684e79e", // vAMM-WETH/KLIMA-1% aerodrome
			ethQuote: true,
		},
	}

	for _, item := range fixtures {
		poolAddr := strings.ToLower(item.poolAddr)

		var tokenAddr string

		pl, ok := lpservice.pools[poolAddr]
		assert.True(t, ok, "pool %s not found", item.poolAddr)

		if !item.ethQuote {
			tokenAddr = WETHAddress
		} else {
			if isNativeOrWrapperNativeToken(pl.Token0) || isStableToken(pl.Token0, lpservice.stableTokens) {
				tokenAddr = pl.Token1
			} else if isNativeOrWrapperNativeToken(pl.Token1) || isStableToken(pl.Token1, lpservice.stableTokens) {
				tokenAddr = pl.Token0
			} else {
				assert.Nil(t, errors.New("neither pool token0 nor token1 is quotable"))
			}
		}

		pl.UpdateBestPrice()

		price := pl.GetPrice0()
		if item.ethQuote {
			logger.Info().Msgf("token %s ETH quoted price: %s", tokenAddr, price)
		} else {
			logger.Info().Msgf("ETH price: %s", price)
		}
	}
}

func TestRedisPool(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "", false)

	rc, err := common.CreatePoolRedisStore(false, "localhost:6379", "", 0)
	assert.Nil(t, err)

	pls, err := pool.LoadAllPools(rc)
	assert.Nil(t, err)

	for _, pl := range pls {
		if pl.Reserve0 == nil || pl.Reserve1 == nil {
			logger.Info().Msgf("pool %v %v reserve is nil", pl.Vendor, pl.Address)
		}
	}
}

func TestArb(t *testing.T) {
	//nolint
	token0 := "0x4200000000000000000000000000000000000006" // weth
	//nolint
	token1 := "0xb79dd08ea68a908a97220c76d19a6aa9cbde4376" // usd+

	lpservice := createTestLPService(t, true)
	lpservice.UpdateArbPairs(false)

	key := token0 + "-" + token1
	arb := lpservice.arbPairs[key]
	assert.NotNil(t, arb)
	t.Logf("arb pairs: %d", len(arb.Pairs))
	ratio := decimal.NewFromFloat(1.00005)
	profitable, param, _ := lpservice.CalcProfitable(arb, arb.Pairs, ratio, true)
	t.Logf("profitable: %v bestAmtIn: %v profit: %v", profitable, param.bestAmtIn, param.bestProfit)
}
