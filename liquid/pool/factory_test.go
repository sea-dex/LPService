package pool

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/utils"
)

func TestGetPairLength(t *testing.T) {
	utils.SkipCI(t)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://mainnet.base.org",
		Tps: 3,
	}})

	factory := []string{
		"0x420DD381b31aEf6683db6B902084cB0FFECe40Da",
		"0x8909Dc15e40173Ff4699343b6eB8132c65e18eC6",
		"0x71524B4f93c58fcbF659783284E38825f0622859",
		"0x3E84D913803b02A4a7f027165E8cA42C14C0FdE7",
		"0x02a84c1b3BBD7401a5f7fa98a384EBC70bB5749E",
		"0x2d5dd5fa7B8a1BFBDbB0916B42280208Ee6DE51e",
		"0xFDa619b6d20975be80A10332cD39b9a4b0FAa8BB",

		/*
			// unknown factory
			"0x70fe4a44ea505cfa3a57b95cf2862d4fd5f0f687",
			"0x0bfbcf9fa4f9c56b0f40a671ad40e0805a091865",
			"0xe6da85feb3b4e0d6aed95c41a125fba859bb9d24",
			"0x591f122d1df761e616c13d265006fcbf4c6d6551",
			"0x7b72c4002ea7c276dd717b96b20f4956c5c904e7",
			"0x3d237ac6d2f425d2e890cc99198818cc1fa48870",
			"0xc5e1116c41aa4525407078123b49bafcb39cfe93",
			"0x04c9f118d21e8b767d2e50c946f0cc9f6c367300",
			"0x510ec01c5ecde6f66febc58ac0e672c8c1b895be",
			"0xb5620f90e803c7f957a9ef351b8db3c746021bea",
			"0xc207628e5e2b59e9c690071e68c7c1c4193b0252",
			"0x079463f811e6eb2e226908e79144cddb59a7fb71",
			"0x0fd83557b2be93617c9c1c1b6fd549401c74558c",
			"0x1b8128c3a1b7d20053d10763ff02466ca7ff99fc",
			"0x4bd16d59a5e1e0db903f724aa9d721a31d7d720d",
			"0x07aced5690e09935b1c0e6e88b772d9440f64718",
			"0xbe720274c24b5ec773559b8c7e28c2503dac7645",
			"0x9e6d21e759a7a288b80eef94e4737d313d31c13f",
			"0x77efb7acee1b0b4f26714c9a61c7de19c053e69f",
			"0xa37359e63d1aa44c0acb2a4605d3b45785c97ee3",
			"0x4c1b8d4ae77a37b94e195cab316391d3c687ebd1",
			"0xdc323d16c451819890805737997f4ede96b95e3e",
			"0x539db2b4fe8016db2594d7cfbeab4d2b730b723e",
			"0x4858c605862a91a34d83c19a9704f837f64fa405",
			"0x2f0d41f94d5d1550b79a83d2fe85c82d68c5a3ca",
			"0x78fa7fa39cf6544dd9768a75d8ad8c45854ae530",
			"0x3e84d913803b02a4a7f027165e8ca42c14c0fde7",
			"0x7b55fa23b303a15ee1261514fdfeafe330e59abb",
			"0xd13aaf098d829aa25eb69cf329a60cef74f2d3bf",
			"0x57592d44eb60011500961ef177bff8d8691d5a8b",
			"0xeddef4273518b137cdbcb3a7fa1c6a688303dfe2",
			"0xa081ce40f079a381b59893b4dc0abf8b1817af70",
			"0xbfae4f07c099798f23f5ac6773532fb637b68ad7",
			"0xbaa207a1673dea8b6890817e5e68e06677471cfb",
			"0xf0384882db4c7be90fcce26d2b9cf0cda80bae22",
			"0x5fd88cd0034a0b2fc77e75ea09b6e512511b0eb9",
			"0x9592cd9b267748cbfbde90ac9f7df3c437a6d51b",
			"0xddf5a3259a88ab79d5530eb3eb14c1c92cd97fcf",
			"0xe4806bdd8e010828324928d25587721f6b58bea2",
			"0xf6c96ac4251905572c7083b1804825850b9bc9e6",
			"0xe396465a85dedb00fa8774162b106833de51ea41",
			"0x51cc508f1f4569073de51fe0ef473e5e4e9bcdc0",
			"0x7bf960b15cbd9976042257be3f6bb2361e107384",
			"0x73cd6389d14522ea15224cd11556c52130b0985a",
			"0x1a62a841e83ecc3d72b0de6002af7a7dbf921cd5",
			"0x576a1301b42942537d38fb147895fe83fb418fd4",
			"0xc7a590291e07b9fe9e64b86c58fd8fc764308c4a",
			"0x45272c3aac84190ed7e3dfcc83315e785376ede7",
			"0xfbb4e52fecc90924c79f980eb24a9794ae4affa4",
			"0x77be21b83cc0592403ed63b2cc012e57155fa527",
			"0xbfd866bd2502e1693bc31735253a2fb9e728f714",
			"0x4ba35f718f8d03e97e4ece6a29724ceecbdf6e8c",
			"0x79b8f15a3beecd5014b748499ec89692665ea368",
		*/
	}
	provider := pp.Get()

	for i, item := range factory {
		l, err := getPairLength(provider, item, common.PoolTypeAMM)
		// assert.Nil(t, err)
		if err == nil {
			t.Logf("factory %s pairs: %d", item, l)
		}

		if i%10 == 0 {
			time.Sleep(time.Second)
		}
	}
}

func TestGetV3Pools(t *testing.T) {
	utils.SkipCI(t)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://mainnet.base.org",
		Tps: 3,
	}})

	pairQuery := "0x8fb641dfe7173fC58C7Edb5BCC13a7187881b96E" //"0x21fF22F27a0D345F3D3faAB6A35Ee8294569715E"
	factory := "0x33128a8fC17869897dcE68Ed026d694621f6FDfD"
	positionManager := "0x03a520b32C04BF3bEEf7BEb72E919cf822Ed34f1"
	pools, err := pp.GetV3Pools(pairQuery, factory, positionManager, "uniswapv3", 10000, true)
	assert.Nil(t, err)
	t.Logf("pools: %d", len(pools))
	t.Logf("last pool: %s", pools[len(pools)-1])

	poolMap := map[string]bool{}
	for _, p := range pools {
		poolMap[p] = true
	}

	t.Logf("pool 0xE18ABE492b3001151034EBCC5c2Dd236d64af553 exist: %v",
		poolMap[strings.ToLower("0xE18ABE492b3001151034EBCC5c2Dd236d64af553")])
}

func TestGetPancakeV3Pools(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "", false)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://mainnet.base.org",
		Tps: 3,
	}})

	pairQuery := "0x8fb641dfe7173fC58C7Edb5BCC13a7187881b96E" //"0x21fF22F27a0D345F3D3faAB6A35Ee8294569715E"
	factory := "0x41ff9AA7e16B8B1a8a8dc4f0eFacd93D02d071c9"
	positionManager := "0x46A15B0b27311cedF172AB29E4f4766fbE7F4364"
	pools, err := pp.GetV3Pools(pairQuery, factory, positionManager, "pancake", 10000, true)
	assert.Nil(t, err)
	t.Logf("pools: %d", len(pools))
	t.Logf("last pool: %s", pools[len(pools)-1])

	poolMap := map[string]bool{}
	for _, p := range pools {
		poolMap[p] = true
	}

	t.Logf("pool 0xE18ABE492b3001151034EBCC5c2Dd236d64af553 exist: %v",
		poolMap[strings.ToLower("0xE18ABE492b3001151034EBCC5c2Dd236d64af553")])
}

func TestGetAeroV3Pools(t *testing.T) {
	utils.SkipCI(t)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://mainnet.base.org",
		Tps: 3,
	}})

	pairQuery := "0x8fb641dfe7173fC58C7Edb5BCC13a7187881b96E" //"0x21fF22F27a0D345F3D3faAB6A35Ee8294569715E"
	factory := "0x5e7BB104d84c7CB9B682AaC2F3d509f5F406809A"   // aero v3
	pools, err := pp.GetAeroV3Pools(pairQuery, factory, 10000, false)
	assert.Nil(t, err)
	t.Logf("pools: %d", len(pools))
}

func TestGetV2Pools(t *testing.T) {
	utils.SkipCI(t)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://mainnet.base.org",
		Tps: 3,
	}})

	factory := &common.SwapFactory{
		Address: "0x8909Dc15e40173Ff4699343b6eB8132c65e18eC6",
		Name:    "UniswapV2",
		Typ:     common.PoolTypeAMM,
	}
	pairQuery := "0x8fb641dfe7173fC58C7Edb5BCC13a7187881b96E" //"0x21fF22F27a0D345F3D3faAB6A35Ee8294569715E"

	pools, err := pp.GetV2PoolInfos(pairQuery, factory, 20000, true)
	assert.Nil(t, err)

	_ = pools

	time.Sleep(time.Second)
}

func TestGetAeroV2PoolsParallel(t *testing.T) {
	utils.SkipCI(t)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://base-rpc.publicnode.com",
		Tps: 3,
	}})

	factory := &common.SwapFactory{
		Address: "0x420DD381b31aEf6683db6B902084cB0FFECe40Da",
		Name:    "AeroV2",
		Typ:     common.PoolTypeAeroAMM,
	}
	pairQuery := "0x8fb641dfe7173fC58C7Edb5BCC13a7187881b96E" //"0x21fF22F27a0D345F3D3faAB6A35Ee8294569715E"

	pools, err := pp.GetV2PoolInfos(pairQuery, factory, 20000, true)
	assert.Nil(t, err)

	_ = pools
	// t.Logf("pool: %v", pools[0])
	time.Sleep(time.Second)
}

func TestGetAeroV2Pools(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "dev", true)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://base-rpc.publicnode.com",
		Tps: 3,
		V3:  true,
	}})

	factory := &common.SwapFactory{
		Address: "0x420DD381b31aEf6683db6B902084cB0FFECe40Da",
		Name:    "AeroV2",
		Typ:     common.PoolTypeAeroAMM,
	}
	pairQuery := "0x8fb641dfe7173fC58C7Edb5BCC13a7187881b96E" //"0x21fF22F27a0D345F3D3faAB6A35Ee8294569715E"

	pools, err := pp.GetV2PoolInfos(pairQuery, factory, 1000, false)
	assert.Nil(t, err)

	_ = pools
	t.Logf("aero v2 pool: %d", len(pools))

	// pool: 0xf8c6644e3641e23072bfce62aa15d9feefd0824c stable: false fee: 10000
	// pool: 0x92b22ce639ae7d5bd9bbb0fe3b134cf8b52d5c5d stable: false fee: 3000
	// pool: 0xa57645a4188747926f25dc086b10dff604c8d930 stable: true fee: 500
	for _, pl := range pools {
		t.Logf("pool: %s stable: %v fee: %v\n", pl.Address, pl.Stable, pl.Fee)
	}

	time.Sleep(time.Second)
}

func TestGetInfusionPools(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "dev", true)

	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://base-rpc.publicnode.com",
		Tps: 3,
		V3:  true,
	}})

	factory := &common.SwapFactory{
		Address: "0x2d9a3a2bd6400ee28d770c7254ca840c82faf23f",
		Name:    "Infusion",
		Typ:     common.PoolTypeInfusionAMM,
	}
	pairQuery := "0x8fb641dfe7173fC58C7Edb5BCC13a7187881b96E" //"0x21fF22F27a0D345F3D3faAB6A35Ee8294569715E"

	pools, err := pp.GetV2PoolInfos(pairQuery, factory, 1000, false)
	assert.Nil(t, err)

	_ = pools
	t.Logf("Infusion v2 pool: %d", len(pools))

	for _, pl := range pools {
		t.Logf("pool: %s stable: %v fee: %v\n", pl.Address, pl.Stable, pl.Fee)
	}

	time.Sleep(time.Second)
}

func TestWait(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Wait()
	t.Log("wait end 1")
	wg.Wait()
	t.Log("wait end 2")
}
