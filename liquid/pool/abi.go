package pool

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"starbase.ag/liquidity/contracts/swaprouter"
)

var (
	UniswapV3FactoryABI abi.ABI
	UniswapV3PoolABI    abi.ABI
	PancakeV3PoolABI    abi.ABI
	UniswapV2PairABI    abi.ABI
	AeroV2PoolABI       abi.ABI
)

func loadABI(abiJSON string, name string) abi.ABI {
	loabi, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		panic(fmt.Sprintf("Invalid %s abi: %v", name, err.Error()))
	}

	return loabi
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func InitABIs() {
	// uniswap v3 factory
	UniswapV3FactoryABI = loadABI(swaprouter.UniswapV3FactoryABI, "UniswapV3Factory")

	// uniswap v3 pool
	UniswapV3PoolABI = loadABI(swaprouter.UniswapV3PoolABI, "UniswapV3Pool")

	// pancake v3 pool
	PancakeV3PoolABI = loadABI(swaprouter.PancakeV3PoolABI, "PancakeV3Pool")

	// uniswap v2 pair
	UniswapV2PairABI = loadABI(swaprouter.UniswapV2PairABI, "UniswapV2Pair")

	// aerodrom v1/v2 pool
	AeroV2PoolABI = loadABI(swaprouter.AeroV2PoolABI, "AeroV2PoolABI")
}
