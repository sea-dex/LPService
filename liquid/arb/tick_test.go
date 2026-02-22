package arb

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/utils"
)

// pool: 0x02f2ccd82a640ba1ef458e9c0992cd5c85acbfa5 tick: 229965 ticks: [206880 -887220]
// pool: 0x56297a0c63351bc240cd493fbf6f0b4a4ed027cb tick: 90215 ticks: [87000 -46000]
func TestBuildTickLiquidity(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "", false)

	rc, err := common.CreatePoolRedisStore(false, "", "", 0)
	assert.Nil(t, err)

	pls, err := pool.LoadAllPools(rc)
	assert.Nil(t, err)

	for _, pl := range pls {
		if len(pl.Ticks) == 3 {
			ticks := slices.Collect(maps.Keys(pl.Ticks))
			sort.Ints(ticks)
			t.Logf("pool: %v tick: %d ticks: %v", pl.Address, pl.Tick, ticks)

			tls0 := BuildTickLiquidity(pl, true, uint32(len(pl.Ticks))) // nolint
			// assert.Equal(t, 1, len(tls0), "left")

			tls1 := BuildTickLiquidity(pl, false, uint32(len(pl.Ticks))) // nolint
			assert.True(t, len(tls0)+len(tls1) >= len(pl.Ticks)-1 && len(tls0)+len(tls1) <= len(pl.Ticks)+1,
				fmt.Sprintf("left: %v right: %v", len(tls0), len(tls1)))

			// break
		}
	}
}

func TestPoolTickLiquidity(t *testing.T) {
	utils.SkipCI(t)
	logger.Init("", "", false)

	s := "{\"address\":\"0xe60bca07cc2d35f20baf5b429957ace1a415e72c\",\"token0\":\"0x4200000000000000000000000000000000000006\",\"token1\":\"0xce674591474c45d4f1e51c4a19611ba705105c4a\",\"factory\":\"0x33128a8fc17869897dce68ed026d694621f6fdfd\",\"vendor\":\"UniswapV3\",\"typ\":300,\"fee\":500,\"stable\":false,\"tickSpacing\":10,\"lastBlockUpdated\":18352083,\"initBlock\":18352078,\"reserve0\":191625052335023,\"reserve1\":9498389726608878964101073444,\"tick\":225479,\"liquidity\":120687150347878834802357,\"sqrtPriceX96\":6235460549979035245950043414982686,\"ticks\":{\"-887270\":{\"tick\":-887270,\"liquidityNet\":120687150347878834802357,\"liquidityGross\":120687150347878834802357},\"225480\":{\"tick\":225480,\"liquidityNet\":-120672692048887975557766,\"liquidityGross\":120672692048887975557766},\"887270\":{\"tick\":887270,\"liquidityNet\":-14458298990859244591,\"liquidityGross\":14458298990859244591}},\"synced\":true,\"initialized\":true}"

	var pl pool.Pool
	err := json.Unmarshal([]byte(s), &pl)
	assert.Nil(t, err)

	tls0 := BuildTickLiquidity(&pl, true, uint32(len(pl.Ticks))) // nolint
	// assert.Equal(t, 1, len(tls0), "left")
	t.Logf("left tick ranges: %v", len(tls0))
	t.Logf("left tick[0]: %v", tls0[0])

	tls1 := BuildTickLiquidity(&pl, false, uint32(len(pl.Ticks))) // nolint
	t.Logf("right tick ranges: %v", len(tls1))
	t.Logf("right tick[0]: %v", tls1[0])

	assert.True(t, len(tls0)+len(tls1) >= len(pl.Ticks)-1 && len(tls0)+len(tls1) <= len(pl.Ticks)+1,
		fmt.Sprintf("left: %v right: %v", len(tls0), len(tls1)))
}

func TestSort(t *testing.T) {
	iv := []int{1, 3, 5, 8}

	assert.Equal(t, sort.SearchInts(iv, 0), 0)
	assert.Equal(t, sort.SearchInts(iv, 1), 0)
	assert.Equal(t, sort.SearchInts(iv, 2), 1)
	assert.Equal(t, sort.SearchInts(iv, 3), 1)
	assert.Equal(t, sort.SearchInts(iv, 5), 2)
	assert.Equal(t, sort.SearchInts(iv, 6), 3)
	assert.Equal(t, sort.SearchInts(iv, 8), 3)
	assert.Equal(t, sort.SearchInts(iv, 9), 4)
}

func TestTickRight(t *testing.T) {
	utils.SkipCI(t)
	logger.Init("", "", false)

	s := "{\"address\":\"0xcdd367446122ba5afbc0eacc675ce9f5030f94a1\",\"token0\":\"0x4200000000000000000000000000000000000006\",\"token1\":\"0xb79dd08ea68a908a97220c76d19a6aa9cbde4376\",\"factory\":\"0x38015d05f4fec8afe15d7cc0386a126574e8077b\",\"vendor\":\"BaseSwap-Basex\",\"typ\":300,\"fee\":80,\"stable\":false,\"tickSpacing\":1,\"lastBlockUpdated\":21708761,\"initBlock\":18149571,\"reserve0\":34062533905565,\"reserve1\":224024624,\"tick\":-198734,\"liquidity\":15640503107075,\"sqrtPriceX96\":3834096742146751123167332,\"ticks\":{\"-198733\":{\"tick\":-198733,\"liquidityNet\":-15640503107075,\"liquidityGross\":15640503107075},\"-198747\":{\"tick\":-198747,\"liquidityNet\":-22575939867,\"liquidityGross\":22575939867},\"-198756\":{\"tick\":-198756,\"liquidityNet\":-188435075534582,\"liquidityGross\":188435075534582},\"-198757\":{\"tick\":-198757,\"liquidityNet\":15640503107075,\"liquidityGross\":15640503107075},\"-198767\":{\"tick\":-198767,\"liquidityNet\":22575939867,\"liquidityGross\":22575939867},\"-198776\":{\"tick\":-198776,\"liquidityNet\":188435075534582,\"liquidityGross\":188435075534582}},\"synced\":true,\"initialized\":true}"
	var pl pool.Pool
	err := json.Unmarshal([]byte(s), &pl)
	assert.Nil(t, err)

	t.Logf("current tick: %v", pl.Tick)
	tls0 := BuildTickLiquidity(&pl, true, uint32(len(pl.Ticks))) // nolint
	// assert.Equal(t, 1, len(tls0), "left")
	t.Logf("left tick ranges: %v", len(tls0))
	t.Logf("left tick[0]: %v", tls0[0])

	tls1 := BuildTickLiquidity(&pl, false, uint32(len(pl.Ticks))) // nolint
	t.Logf("right tick ranges: %v", len(tls1))
	t.Logf("right tick[0]: %v", tls1[0])
}

func TestCalcV2ToV3(t *testing.T) {
	utils.SkipCI(t)
	logger.Init("", "", false)
	s1 := "{\"address\":\"0x8cf50d1581b801168a4f9ff1ae087e78d03bed5b\",\"token0\":\"0x4200000000000000000000000000000000000006\",\"token1\":\"0xb79dd08ea68a908a97220c76d19a6aa9cbde4376\",\"factory\":\"0x3e84d913803b02a4a7f027165e8ca42c14c0fde7\",\"vendor\":\"Alien\",\"typ\":200,\"fee\":1600,\"stable\":false,\"tickSpacing\":0,\"lastBlockUpdated\":21748084,\"initBlock\":0,\"reserve0\":42100056939977,\"reserve1\":106684,\"tick\":0,\"synced\":false,\"initialized\":true}"
	s2 := "{\"address\":\"0x806eab3b2f63343da07fe3c462a0b38a8bec5fd9\",\"token0\":\"0x4200000000000000000000000000000000000006\",\"token1\":\"0xb79dd08ea68a908a97220c76d19a6aa9cbde4376\",\"factory\":\"0x38015d05f4fec8afe15d7cc0386a126574e8077b\",\"vendor\":\"BaseSwap-Basex\",\"typ\":300,\"fee\":450,\"stable\":false,\"tickSpacing\":10,\"lastBlockUpdated\":19351112,\"initBlock\":18403283,\"reserve0\":978364187031964367,\"reserve1\":16567173,\"tick\":-887272,\"liquidity\":0,\"sqrtPriceX96\":4295128740,\"ticks\":{\"-194330\":{\"tick\":-194330,\"liquidityNet\":-3259878533457,\"liquidityGross\":3259878533457},\"-194440\":{\"tick\":-194440,\"liquidityNet\":-1211780679066012,\"liquidityGross\":1211780679066012},\"-195380\":{\"tick\":-195380,\"liquidityNet\":1211780679066012,\"liquidityGross\":1211780679066012},\"-195590\":{\"tick\":-195590,\"liquidityNet\":3259878533457,\"liquidityGross\":3259878533457}},\"synced\":true,\"initialized\":true}"

	var pl1, pl2 pool.Pool
	err := json.Unmarshal([]byte(s1), &pl1)
	assert.Nil(t, err)

	err = json.Unmarshal([]byte(s2), &pl2)
	assert.Nil(t, err)

	ticks := BuildTickLiquidity(&pl2, false, 5)
	printBoundedTicks(ticks)
	//nolint
	amtIn, amtMid, amtOut, profit, profitable := CalcAmountInV2ToV3(true, &pl1, &pl2, pl1.Reserve0, pl1.Reserve1, ticks, uint32(pl1.Fee), uint32(pl2.Fee), true)
	t.Logf("amtIn=%v, amtMid=%v, amtOut=%v, profit=%v, profitable=%v", amtIn, amtMid, amtOut, profit, profitable)
}
