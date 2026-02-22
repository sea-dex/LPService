package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/utils"
)

func TestLoadTokens(t *testing.T) {
	utils.SkipCI(t)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://mainnet.base.org",
		Tps: 3,
	}})

	tokens, _, err := pp.LoadTokens([]string{
		"0x4200000000000000000000000000000000000006",
		"0x0bD4887f7D41B35CD75DFF9FfeE2856106f86670",
		"0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
	})
	assert.Nil(t, err)

	expTokens := []common.Token{
		{
			Name:     "Wrapped Ether",
			Symbol:   "WETH",
			Decimals: 18,
		}, {
			Name:     "FRIEND",
			Symbol:   "FRIEND",
			Decimals: 18,
		}, {
			Name:     "USD Coin",
			Symbol:   "USDC",
			Decimals: 6,
		},
	}
	for i, token := range tokens {
		assert.Equal(t, expTokens[i].Name, token.Name)
		assert.Equal(t, expTokens[i].Symbol, token.Symbol)
		assert.Equal(t, expTokens[i].Decimals, token.Decimals)
	}
}

func TestLoadMaxTokens(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "", true)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://base-rpc.publicnode.com",
		Tps: 3,
		V3:  true,
	}})

	addr := "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"
	tokenList := []string{}

	for i := 500; i <= 1000; i += 100 {
		for j := 0; j < i; j++ {
			tokenList = append(tokenList, addr)
		}

		_, _, err := pp.LoadTokens(tokenList)
		if err != nil {
			t.Logf("Failed: i = %d error: %v", i, err.Error())
		} else {
			t.Logf("Success: i = %d", i)
		}

		tokenList = []string{}
	}
}

func TestLoadToken(t *testing.T) {
	utils.SkipCI(t)

	tokens := []string{
		"0x92a062bb26a0cbc704bf7dcd4c833d4e1beeb83d",
		"0x4a89e24ad5f82823989022b5c8f3df6fb24fb3e8",
		"0xe741fb0cd252885a9d12646e57c29f6fd70f904c",
		"0xfe861577822f49f52df5111ba0a07da76beed35f",
		"0xfcaf32b29bfaef79c1419ea0ad927b1a0c247e7b",
		"0xbc73bfbe4c047cbb244f3f12a2fd212f569cccb7",
		"0x3038715eae479f9399711ac6d07394a266e58ab5",
		"0x5f61daa1705dace426cc9a2456d70d674e6bda0e",
		"0xc430e171f35e699409014d610fcfda72f03859b3",
		"0x84cc2585f575ef7c06f6c23771f09f324fb483ff",
		"0xffa3fc373c11c0168c7fcc31a638e2b4a35110a2",
		"0x2f39c1ea64b8435687cb708ca571df468ade031c",
		"0xef4fc624ea1a2acfd806240ada70d6802a81eaf3",
		"0x879d0393477380d3ec06aa05d6b99620bbd4d886",
		"0xf70c46c30eb35fd8d7f24a4eb71e065bb58732cc",
		"0x9065c9f6337fb742709a85402bad53c78236db36",
		"0xb7e01ed9dd6365cd84904a354213f614e01ac8ad",
		"0xf0c77d30788be7844533166dd3992d6dc938184a",
		"0x93201d45122b039628cdde471b5a5d87bb78a83a",
		"0x8a7ac6c3353f0a5b9302096e96432cbb53335ce6",
		// some token not found
		"0x059f1f9d3c071db1ac1899c6f84cfdc83be854fe",
		"0xb0dddb11bdf9ab55e31f9dcdff9152415871b744",
		"0x380921b91c4bd9acc0cbfce321bc93fce676de68",
	}

	logger.Init("", "dev", false)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://base-rpc.publicnode.com",
		Tps: 3,
		V3:  true,
	}})
	// for _, tok := range tokens {
	results, failed, err := pp.LoadTokens(tokens)
	if err != nil {
		logger.Error().Msgf("get token info failed: %d %v", len(failed), err)
	} else {
		logger.Info().Msgf("failed: %d", len(failed))

		for _, res := range results {
			logger.Info().Msgf("token info: %s %s %s %d", res.Address, res.Name, res.Symbol, res.Decimals)
		}
	}
}
