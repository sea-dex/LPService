package pool

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"
	"time"

	defiabi "github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/utils"
)

/*
func TestGetV3PoolInfo(t *testing.T) {
	utils.SkipCI(t)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC:      "https://mainnet.base.org",
		Interval: 500,
	}})

	poolAddrs := []string{
		"0xd0b53D9277642d899DF5C87A3966A349A798F224", // uniswap eth/usdc
		"0xb2cc224c1c9feE385f8ad6a55b4d94E92359DC59", // aero
		"0x72ab388e2e2f6facef59e3c3fa2c4e29011c2d38", // pancake
	}
	poolTypes := []common.PoolType{
		common.PoolTypeCAMM,
		common.PoolTypeAeroCAMM,
		common.PoolTypePancakeCAMM,
	}
	expectPool := []*Pool{
		{
			PoolInfo: PoolInfo{
				Factory:     "0x33128a8fc17869897dce68ed026d694621f6fdfd",
				Token0:      "0x4200000000000000000000000000000000000006",
				Token1:      "0x833589fcd6edb6e08f4c7c32d4f71b54bda02913",
				TickSpacing: 10,
				Fee:         500,
			},
		},
		{
			PoolInfo: PoolInfo{
				Factory:     "0x5e7bb104d84c7cb9b682aac2f3d509f5f406809a",
				Token0:      "0x4200000000000000000000000000000000000006",
				Token1:      "0x833589fcd6edb6e08f4c7c32d4f71b54bda02913",
				TickSpacing: 100,
				Fee:         400,
			},
		},
		{
			PoolInfo: PoolInfo{
				Factory:     "0x0bfbcf9fa4f9c56b0f40a671ad40e0805a091865",
				Token0:      "0x4200000000000000000000000000000000000006",
				Token1:      "0x833589fcd6edb6e08f4c7c32d4f71b54bda02913",
				TickSpacing: 1,
				Fee:         100,
			},
		},
	}

	zero := big.NewInt(0)

	for i, poolAddr := range poolAddrs {
		pool := &Pool{
			PoolInfo: PoolInfo{
				Address: poolAddr,
				Typ:     poolTypes[i],
			},
		}
		err := pp.GetV3PoolLiquidInfo(pool, pool.Typ)
		assert.Nil(t, err)

		exp := expectPool[i]
		// assert.Equal(t, exp.Factory, pool.Factory)
		// assert.Equal(t, exp.Token0, pool.Token0)
		// assert.Equal(t, exp.Token1, pool.Token1)
		assert.Equal(t, exp.TickSpacing, pool.TickSpacing)
		assert.Equal(t, exp.Fee, pool.Fee)
		assert.True(t, pool.Liquidity.Cmp(zero) > 0)
		assert.True(t, pool.SqrtPriceX96.Cmp(zero) > 0)

		t.Log(pool.InitBlock)
	}
}
*/

// func TestABIBug(t *testing.T) {
// 	defiabi.MustParseMethod(`function sendTo(address payable to) external`)
// }

func TestGetV3PoolLiquidity(t *testing.T) {
	utils.SkipCI(t)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://mainnet.base.org",
		Tps: 3,
	}})

	pools := []*Pool{
		{
			PoolInfo: PoolInfo{
				Address: "0xd0b53D9277642d899DF5C87A3966A349A798F224",
				Typ:     common.PoolTypeCAMM,
				Token0:  "0x4200000000000000000000000000000000000006",
				Token1:  "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
			},
		},
		{
			PoolInfo: PoolInfo{
				Address: "0xb2cc224c1c9feE385f8ad6a55b4d94E92359DC59",
				Typ:     common.PoolTypeAeroCAMM,
				Token0:  "0x4200000000000000000000000000000000000006",
				Token1:  "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
			},
		},
		{
			PoolInfo: PoolInfo{
				Address: "0x72ab388e2e2f6facef59e3c3fa2c4e29011c2d38",
				Typ:     common.PoolTypePancakeCAMM,
				Token0:  "0x4200000000000000000000000000000000000006",
				Token1:  "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
			},
		},
	}

	queryAddrs := []string{
		"0xDf7acDFaab84FE57c999aEf080749845C97ca038",
		"0xB369A9B58bE84783F47E66c244F79567E36367B1",
		"0x398FbFe61579090aEcC613a25BdeCffBa8D60313",
	}

	for i := range queryAddrs {
		pool := pools[i]
		err := pp.GetV3PoolLiquidity(pool, queryAddrs[i], 0)
		assert.Nil(t, err)

		// err = pp.GetV3PoolTicksLiquid(pool, queryAddrs[i], 0) // uint(0x01000100)) // uint((65536-3000)<<16|3000))
		// assert.Nil(t, err)
		assert.True(t, pool.SqrtPriceX96.Cmp(MIN_SQRT_RATIO) >= 0)
		assert.True(t, pool.SqrtPriceX96.Cmp(MAX_SQRT_RATIO) <= 0)
		t.Logf("pool tick: %d tickSpacing: %d fee: %d sqrtPrice: %s liquidity: %s ticks: %d",
			pool.Tick, pool.TickSpacing, pool.Fee, pool.SqrtPriceX96, pool.Liquidity, len(pool.Ticks))
	}

	time.Sleep(time.Second)
}

func TestGetV3PoolLiquidityUniswap(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "", false)

	rpcList := []string{
		"https://base.rpc.subquery.network/public",
		"https://base-mainnet.public.blastapi.io",
		"https://base.llamarpc.com",
		"https://base.meowrpc.com",
		"https://base.blockpi.network/v1/rpc/public",
		"https://base.gateway.tenderly.co",
		"https://1rpc.io/base",
		"https://public.stackup.sh/api/v1/node/base-mainnet",
		"https://base-mainnet.gateway.tatum.io",
		"https://base.api.onfinality.io/public",
		"https://base-pokt.nodies.app",
		"https://endpoints.omniatech.io/v1/base/mainnet/public",
	}

	for _, rpc := range rpcList {
		pp := CreateProviderPool([]config.ProviderConfig{{
			RPC: rpc,
			Tps: 3,
		}})
		pool := &Pool{
			PoolInfo: PoolInfo{
				Address: "0x06959273e9a65433de71f5a452d529544e07ddd0",
				Typ:     common.PoolTypeCAMM,
				Token0:  "0x4200000000000000000000000000000000000006",
				Token1:  "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
			},
		}
		queryAddr := "0xDf7acDFaab84FE57c999aEf080749845C97ca038"

		err := pp.GetV3PoolLiquidity(pool, queryAddr, 0)
		if err != nil {
			t.Logf("call node %v failed: %v", rpc, err)
			continue
		}

		t.Logf("----------- RPC Success: %v", rpc)
		t.Logf("Liquidity: %v", pool.Liquidity)
		t.Logf("SqrtPriceX96: %v", pool.SqrtPriceX96)
		t.Logf("Tick: %v", pool.Tick)
		t.Logf("TickList: %v", len(pool.Ticks))
		t.Logf("Reserve0: %v", pool.Reserve0)
		t.Logf("Reserve1: %v", pool.Reserve1)
	}
}

func TestGetV3PoolLiquidityAero(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "", false)

	rpcList := []string{
		"https://base.rpc.subquery.network/public",
		"https://base-mainnet.public.blastapi.io",
		"https://base.llamarpc.com",
		"https://base.meowrpc.com",
		"https://base.blockpi.network/v1/rpc/public",
		"https://base.gateway.tenderly.co",
		"https://1rpc.io/base",
		"https://public.stackup.sh/api/v1/node/base-mainnet",
		"https://base-mainnet.gateway.tatum.io",
		"https://base.api.onfinality.io/public",
		"https://base-pokt.nodies.app",
		"https://endpoints.omniatech.io/v1/base/mainnet/public",
	}

	for _, rpc := range rpcList {
		pp := CreateProviderPool([]config.ProviderConfig{{
			RPC: rpc,
			Tps: 3,
		}})
		pool := &Pool{
			PoolInfo: PoolInfo{
				Address: "0x861a2922be165a5bd41b1e482b49216b465e1b5f",
				Typ:     common.PoolTypeAeroCAMM,
				Token0:  "0x4200000000000000000000000000000000000006",
				Token1:  "0xc1CBa3fCea344f92D9239c08C0568f6F2F0ee452",
			},
		}
		queryAddr := "0xB369A9B58bE84783F47E66c244F79567E36367B1"

		err := pp.GetV3PoolLiquidity(pool, queryAddr, 0) // 0x200)
		if err != nil {
			t.Logf("call node %v failed: %v", rpc, err)
			continue
		}

		t.Logf("----------- RPC Success: %v", rpc)
		t.Logf("Liquidity: %v", pool.Liquidity)
		t.Logf("SqrtPriceX96: %v", pool.SqrtPriceX96)
		t.Logf("Tick: %v", pool.Tick)
		t.Logf("TickList: %v", len(pool.Ticks))
		t.Logf("Reserve0: %v", pool.Reserve0)
		t.Logf("Reserve1: %v", pool.Reserve1)
	}
}

func TestGetV3PoolLiquidityPancake(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "", false)

	rpcList := []string{
		"https://mainnet.base.org",
		"https://base.rpc.subquery.network/public",
		// "https://base-mainnet.public.blastapi.io",
		// "https://base.llamarpc.com",
		// "https://base.meowrpc.com",
		// "https://base.blockpi.network/v1/rpc/public",
		// "https://base.gateway.tenderly.co",
		// "https://1rpc.io/base",
		// "https://public.stackup.sh/api/v1/node/base-mainnet",
		// "https://base-mainnet.gateway.tatum.io",
		// "https://base.api.onfinality.io/public",
		// "https://base-pokt.nodies.app",
		// "https://endpoints.omniatech.io/v1/base/mainnet/public",
	}

	for _, rpc := range rpcList {
		pp := CreateProviderPool([]config.ProviderConfig{{
			RPC: rpc,
			Tps: 3,
			V3:  true,
		}})
		pool := &Pool{
			PoolInfo: PoolInfo{
				Address: "0xf6c0a374a483101e04ef5f7ac9bd15d9142bac95",
				// Address: "0x72ab388e2e2f6facef59e3c3fa2c4e29011c2d38",
				Typ:    common.PoolTypePancakeCAMM,
				Token0: "0x4200000000000000000000000000000000000006",
				Token1: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
			},
		}
		queryAddr := "0x398FbFe61579090aEcC613a25BdeCffBa8D60313"

		err := pp.GetV3PoolLiquidity(pool, queryAddr, 0)
		if err != nil {
			t.Logf("call node %v failed: %v", rpc, err)
			continue
		}

		t.Logf("----------- RPC Success: %v", rpc)
		t.Logf("Liquidity: %v", pool.Liquidity)
		t.Logf("SqrtPriceX96: %v", pool.SqrtPriceX96)
		t.Logf("Tick: %v", pool.Tick)
		t.Logf("TickList: %v", len(pool.Ticks))
		t.Logf("Reserve0: %v", pool.Reserve0)
		t.Logf("Reserve1: %v", pool.Reserve1)
	}
}

func TestTickLens(t *testing.T) {
	utils.SkipCI(t)

	// Uniswap V3 pool ETH/USDC on base mainnet
	poolAddr := "0xd0b53D9277642d899DF5C87A3966A349A798F224"
	poolTickLensABI := defiabi.MustParseJSON([]byte(`[{"inputs":[{"internalType":"address","name":"pool","type":"address"},{"internalType":"int16","name":"tickBitmapIndex","type":"int16"}],"name":"getPopulatedTicksInWord","outputs":[{"components":[{"internalType":"int24","name":"tick","type":"int24"},{"internalType":"int128","name":"liquidityNet","type":"int128"},{"internalType":"uint128","name":"liquidityGross","type":"uint128"}],"internalType":"struct ITickLens.PopulatedTick[]","name":"populatedTicks","type":"tuple[]"}],"stateMutability":"view","type":"function"}]`))

	calldata := poolTickLensABI.Methods["getPopulatedTicksInWord"].MustEncodeArgs(poolAddr, -100)
	queryAddr := "0x0CdeE061c75D43c82520eD998C23ac2991c9ac6d" // Uniswap V3 TickLens contract address on base mainnet

	uri := "https://mainnet.base.org"
	tr, err := transport.NewHTTP(transport.HTTPOptions{URL: uri})
	assert.Nil(t, err)

	// Create a JSON-RPC client.
	client, err := rpc.NewClient(rpc.WithTransport(tr))
	assert.Nil(t, err)

	// Prepare a call.
	call := types.NewCall().
		SetTo(types.MustAddressFromHex(queryAddr)).
		SetInput(calldata)

	// 351fb478 000000000000000000000000d0b53d9277642d899df5c87a3966a349a798f224 000000000000000000000000000000000000000000000000000000000000ff9c
	println(hex.EncodeToString(calldata))

	// Call method.
	b, _, err := client.Call(context.Background(), call, types.LatestBlockNumber)
	if err != nil {
		panic(err)
	}

	_ = b
}

func TestTickLens2(t *testing.T) {
	tickabi := defiabi.MustParseSignatures(
		`struct PopulatedTick {
        int24 tick;
        int128 liquidityNet;
        uint128 liquidityGross;
    }`,
		`function getPopulatedTicksInWord(address pool, int16 tickBitmapIndex)
        view
        returns (PopulatedTick[] memory populatedTicks)`,
	)
	_ = tickabi
	// Uniswap V3 pool ETH/USDC on base mainnet
	poolAddr := "0xd0b53D9277642d899DF5C87A3966A349A798F224"
	_ = poolAddr
	// calldata := tickabi.Methods["getPopulatedTicksInWord"].MustEncodeArgs(poolAddr, -100)
	calldata, err := hex.DecodeString("351fb478000000000000000000000000d0b53d9277642d899df5c87a3966a349a798f224ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9c")
	assert.Nil(t, err)
	// Uniswap V3 TickLens contract address on base mainnet
	queryAddr := "0x0CdeE061c75D43c82520eD998C23ac2991c9ac6d"

	uri := "https://mainnet.base.org"
	tr, err := transport.NewHTTP(transport.HTTPOptions{URL: uri})
	assert.Nil(t, err)

	// Create a JSON-RPC client.
	client, err := rpc.NewClient(rpc.WithTransport(tr))
	assert.Nil(t, err)

	// Prepare a call.
	call := types.NewCall().
		SetTo(types.MustAddressFromHex(queryAddr)).
		SetInput(calldata)

	// 351fb478
	// 000000000000000000000000d0b53d9277642d899df5c87a3966a349a798f224
	// 000000000000000000000000000000000000000000000000000000000000ff9c
	println(hex.EncodeToString(calldata))

	// Call method.
	b, _, err := client.Call(context.Background(), call, types.LatestBlockNumber)
	if err != nil {
		panic(err)
	}

	_ = b
}

func TestGetPoolBasicInfo(t *testing.T) {
	utils.SkipCI(t)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://mainnet.base.org",
		Tps: 3,
		V3:  true,
	}})

	// poolAddr := "0x4cfd5ba4b8e0475d9a3cfa863e0e18ccf9d3eb25"
	// failed
	poolAddr := "0xa414bB51b948c193322f3465F342788484c91750"
	pool, err := pp.GetPoolBasicInfo(poolAddr, 3)
	assert.Nil(t, err)

	t.Logf("pool: %v", pool)
}

func TestGetPoolsBasicInfo(t *testing.T) {
	utils.SkipCI(t)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://mainnet.base.org",
		Tps: 3,
		V3:  true,
	}})

	poolAddrs := []string{
		"0x4cfd5ba4b8e0475d9a3cfa863e0e18ccf9d3eb25",
		"0xa414bB51b948c193322f3465F342788484c91750",
	}
	pools := map[string]*Pool{}
	failed, err := pp.getPoolsBasicInfo(poolAddrs, 3, pools)
	assert.Nil(t, err)

	t.Logf("failed: %d", len(failed))
}

func TestGetV2PoolLiquid(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "dev", true)
	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://mainnet.base.org",
		Tps: 3,
		V3:  true,
	}})

	pls := []*Pool{
		{
			PoolInfo: PoolInfo{
				Address: "0x3ca95d2a2452c3646415fe6e7420f6999e06152e",
			},
		},
		{
			PoolInfo: PoolInfo{
				Address: "0x6c81c67aa4de39821b9d5bcbc620dde2b8d18bb5", // aero
			},
		},
	}
	for _, pl := range pls {
		err := pp.GetV2PoolLiquid(pl)
		assert.Nil(t, err)

		t.Logf("pool: %v %v", pl.Reserve0.String(), pl.Reserve1.String())
	}
}

func TestGetAeroPoolV2StableFee(t *testing.T) {
	utils.SkipCI(t)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://base-rpc.publicnode.com",
		Tps: 3,
		V3:  true,
	}})
	poolAdds := []string{
		"0x33cd8eaf16ff4df29eedb83fe1963876843cc7ec", // true, 5
		"0x6c81c67aa4de39821b9d5bcbc620dde2b8d18bb5", // false, 30
		"0xe9c25c83401f1763b390122bb644eaf041a40919", // false, 30
		"0x56ae4a50252c505c2663346660edb711648221dc", // false, 30
		"0x2e2a6758bd5a4d4311fafe50f673536cf2995350", // false, 30
		"0x8163c6d12bdc1e08af9720f1ef4f496b596127a9", // false, 30
		"0xb543a23ebf95b35f7e472076c34705821e3b0817", // true, 5
		"0xb37642e87613d8569fd8ec80888ea6c63684e79e", // false, 100
		"0x276cc711842379d99038da4d31ffbea8e8c3f77e", // true, 5
		"0x4f7f293e7b7b6d1caa1ca48784c5986f3de5e19c", // true, 1
	}
	pools := []*Pool{}

	for _, addr := range poolAdds {
		pl := &Pool{}
		pl.Address = addr
		pl.Factory = strings.ToLower("0x420DD381b31aEf6683db6B902084cB0FFECe40Da")
		pools = append(pools, pl)
	}

	err := pp.GetAeroPoolV2Stable(pools, 3)
	assert.Nil(t, err)
	_, err = pp.GetAeroPoolV2Fee(pools, 3)
	assert.Nil(t, err)

	for _, pl := range pools {
		t.Logf("pool %s: stable=%v fee=%v\n", pl.Address, pl.Stable, pl.Fee)
	}
}
