package pool

import (
	"encoding/hex"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"starbase.ag/liquidity/contracts/swaprouter"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
)

const LatestVersion = 1

// IPool Pool interface.
//
//	type IPool interface {
//		PoolAddress() string
//		FactoryAddress() string
//		Token0Address() string
//		Token1Address() string
//		GetPoolType() PoolType
//		GetVendor() string
//	}
type PoolBasicInfo struct {
	Address string `json:"address"` // pool address
	Token0  string `json:"token0"`  // pool token0 address
	Token1  string `json:"token1"`  // pool token1 address
	Factory string `json:"factory"` // pool factory address
}

type PoolInfo struct {
	key       string // token0-token1
	name      string // WETH/USDC-0.3%
	tvl       *big.Int
	price0    *decimal.Decimal // sell token0
	price1    *decimal.Decimal // sell token1
	decimals0 *big.Int
	decimals1 *big.Int

	Address          string          `json:"address"`          // pool address
	Token0           string          `json:"token0"`           // pool token0 address
	Token1           string          `json:"token1"`           // pool token1 address
	Factory          string          `json:"factory"`          // pool factory address
	Vendor           string          `json:"vendor"`           // pool vendor, for example: UniswapV2, UniswapV3
	Typ              common.PoolType `json:"typ"`              // pool type, AMM or CAMM
	Fee              uint            `json:"fee"`              // pool fee bips, 3000 is 0.3%
	Stable           bool            `json:"stable"`           // for Aero v2 pair
	TickSpacing      int             `json:"tickSpacing"`      // uniswap v3 pool tick spacing
	LastBlockUpdated uint64          `json:"lastBlockUpdated"` // the block number of pool changed(Mint, Burn, Swap, Collect, etc)
	InitBlock        uint64          `json:"initBlock"`        // the block get initial liquidity
	// Version          uint            `json:"version"`
	// Known            bool            `json:"known"`            // if the pool is known
}

type SwapEvent struct {
	ZeroForOne bool     `json:"zeroForOne"` // true: sell token0, buy token1; false: buy token0, sell token1
	Amount0    *big.Int `json:"amount0"`
	Amount1    *big.Int `json:"amount1"`
	Decimals0  uint     `json:"decimals0"` // token0 decimals, for example 18
	Decimals1  uint     `json:"decimals1"` // token1 decimals, for example 18
	Txhash     string   `json:"txhash"`

	Address          string          `json:"address"`          // pool address
	Token0           string          `json:"token0"`           // pool token0 address
	Token1           string          `json:"token1"`           // pool token1 address
	Factory          string          `json:"factory"`          // pool factory address
	Vendor           string          `json:"vendor"`           // pool vendor, for example: UniswapV2, UniswapV3
	Typ              common.PoolType `json:"typ"`              // pool type, AMM or CAMM
	Fee              uint            `json:"fee"`              // pool fee bips, 3000 is 0.3%
	Stable           bool            `json:"stable"`           // for Aero v2 pair
	TickSpacing      int             `json:"tickSpacing"`      // uniswap v3 pool tick spacing
	LastBlockUpdated uint64          `json:"lastBlockUpdated"` // the block number of pool changed(Mint, Burn, Swap, Collect, etc)
	InitBlock        uint64          `json:"initBlock"`        // the block get initial liquidity
}

// Pool CAMM(Uniswap v3, Sushi v3, Pancake v3, Aero v3) pool liquidity.
type Pool struct {
	PoolInfo

	Reserve0 *big.Int `json:"reserve0"` // token0 reserve
	Reserve1 *big.Int `json:"reserve1"` // token1 reserve

	// following is Uniswap V3 fields
	// TickList     []int             `json:"tickList,omitempty"`     // uniswap v3 tick list, for iterate
	Tick         int               `json:"tick"`                   // uniswap v3 pool tick
	Liquidity    *big.Int          `json:"liquidity,omitempty"`    // uniswap v3 liquidity
	SqrtPriceX96 *big.Int          `json:"sqrtPriceX96,omitempty"` // uniswap v3 sqrtPrice
	Ticks        map[int]*TickInfo `json:"ticks,omitempty"`        // pool ticks map
	Synced       bool              `json:"synced"`                 // synced meanings swap verify success
	Initialized  bool              `json:"initialized"`

	tickBitmap map[int16]*big.Int //nolint
}

var PoolReloadAll = Pool{
	PoolInfo: PoolInfo{
		Typ: common.PoolTypeReloadAll,
	},
}

// func (pool *Pool) NeedUpdate() bool {
// 	if pool.Typ == common.PoolTypeAeroAMM && pool.Version < LatestVersion {
// 		return true
// 	}

//		return false
//	}

func (pool *Pool) SortedTickList() []int {
	ticks := []int{}
	for tick := range pool.Ticks {
		ticks = append(ticks, tick)
	}

	sort.Ints(ticks)
	return ticks
}

func (pool *Pool) GetTickListAndActiveRange() ([]int, int) {
	ticks := pool.SortedTickList()
	currTickIdx := sort.SearchInts(ticks, pool.Tick)
	if currTickIdx > 0 && currTickIdx < len(ticks) && pool.Tick != ticks[currTickIdx] {
		// if pool.Tick != ticks[currTickIdx-1] {
		currTickIdx--
		// }
	}

	return ticks, currTickIdx
}

// ResetLiquidity reset pool liquidity.
func (pool *Pool) ResetLiquidity() {
	pool.Reserve0 = big.NewInt(0)
	pool.Reserve1 = big.NewInt(0)
	pool.Tick = 0
	// pool.TickList = []int{}
	pool.Liquidity = big.NewInt(0)
	pool.SqrtPriceX96 = big.NewInt(0)
	pool.Ticks = map[int]*TickInfo{}
	pool.tickBitmap = map[int16]*big.Int{}
	pool.Initialized = false
	pool.Synced = false
}

// pool address + tickSpacing.
func (pool *Pool) SingleTickParams() []byte {
	s := pool.Address

	ts := strconv.FormatUint(uint64(pool.TickSpacing), 16) // nolint
	for i := 0; i < 6-len(ts); i++ {
		s += "0"
	}

	s += ts

	buf, _ := hex.DecodeString(s)

	return buf
}

func (pool *Pool) SetInitialized(v bool) {
	pool.Initialized = v
}

func (pool *Pool) Reload() {
	if pool.Reserve0 == nil {
		pool.Reserve0 = big.NewInt(0)
	}

	if pool.Reserve1 == nil {
		pool.Reserve1 = big.NewInt(0)
	}

	if !pool.Typ.IsCAMMVariety() {
		pool.key = pool.Token0 + "-" + pool.Token1
		return
	}

	pool.tickBitmap = map[int16]*big.Int{}
	for tick := range pool.Ticks {
		flipTickBitmap(pool.tickBitmap, tick, pool.TickSpacing)
	}

	if pool.Ticks == nil {
		if pool.Liquidity.Cmp(bigZero) != 0 {
			logger.Error().Str("pool", pool.Address).Msgf("pool ticks is nil, liquidity: %v", pool.Liquidity)
		}

		pool.Ticks = map[int]*TickInfo{}
	}

	pool.key = pool.Token0 + "-" + pool.Token1
	// if pool.TickList == nil {
	// 	pool.TickList = []int{}
	// 	logger.Warn().Str("pool", pool.Address).Msg("pool tickList is nil")
	// }

	if pool.SqrtPriceX96 == nil || pool.Liquidity == nil {
		logger.Fatal().Str("pool", pool.Address).Msg("pool SqrtPriceX96 or Liquidity is nil")
	}
}

func (pool *Pool) GetKey() string {
	if pool.key == "" {
		pool.UpdateKey()
	}

	return pool.key
}

func (pool *Pool) UpdateKey() {
	pool.key = pool.Token0 + "-" + pool.Token1
}

func (pool *Pool) GetName() string {
	return pool.name
}

func (pool *Pool) UpdateName(name string) {
	pool.name = name
}

func (pool *Pool) GetTVL() *big.Int {
	return pool.tvl
}

func (pool *Pool) GetPrice0() decimal.Decimal {
	return *pool.price0
}

func (pool *Pool) GetPrice1() decimal.Decimal {
	return *pool.price1
}

func (pool *Pool) GetDecimals0() *big.Int {
	return pool.decimals0
}

func (pool *Pool) GetDecimals1() *big.Int {
	return pool.decimals1
}

func (pool *Pool) SetTVL(tvl *big.Int) {
	pool.tvl = tvl
}

func (pool *Pool) SetPrices(p0, p1 *decimal.Decimal) {
	pool.price0 = p0
	pool.price1 = p1
}

func (pool *Pool) SetPrice0(price *decimal.Decimal) {
	pool.price0 = price
}

func (pool *Pool) SetPrice1(price *decimal.Decimal) {
	pool.price1 = price
}

func (pool *Pool) SetDecimals0(decimals0 *big.Int) {
	pool.decimals0 = decimals0
}

func (pool *Pool) SetDecimals1(decimals1 *big.Int) {
	pool.decimals1 = decimals1
}

func (pool *Pool) OnV3Event(event *types.Log) {
	topic := strings.ToLower(event.Topics[0].Hex())
	switch topic {
	case TopicInitialize:
		var initialize swaprouter.UniswapV3PoolInitialize
		if err := UniswapV3PoolABI.UnpackIntoInterface(&initialize, "Initialize", event.Data); err != nil {
			logger.Fatal().
				Err(err).
				Str("vendor", pool.Vendor).
				Str("factory", pool.Factory).
				Str("txhahs", event.TxHash.Hex()).
				Str("pool", pool.Address).
				Msg("Unpack Initialize event data failed")
		}

		tick := int(initialize.Tick.Int64())
		pool.OnInitialize(initialize.SqrtPriceX96, tick)

	case TopicMint:
		pool.OnMint(event)

	case TopicSwap:
		pool.OnSwap(event)

	case TopicBurn:
		pool.OnBurn(event)

	case TopicCollect:
		pool.OnCollect(event)

	case TopicPancakeSwap:
		pool.OnPancakeSwap(event)
	}
}
