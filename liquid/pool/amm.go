package pool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"starbase.ag/liquidity/contracts/swaprouter"
	"starbase.ag/liquidity/pkg/logger"
)

// OnSync on sync event.
func (pool *Pool) OnSync(evt *types.Log) {
	// if pl.Typ ==
	var sync swaprouter.UniswapV2PairSync

	err := UniswapV2PairABI.UnpackIntoInterface(&sync, "Sync", evt.Data)
	if err != nil {
		logger.Fatal().
			Err(err).
			Str("vendor", pool.Vendor).
			Str("factory", pool.Factory).
			Str("txhash", evt.TxHash.Hex()).
			Str("pool", pool.Address).
			Msg("Unpack Sync event data failed")
	}

	if sync.Reserve0 == nil || sync.Reserve1 == nil {
		logger.Fatal().
			Str("vendor", pool.Vendor).
			Str("factory", pool.Factory).
			Str("txhash", evt.TxHash.Hex()).
			Str("pool", pool.Address).
			Msg("event Sync Reserve is nil")
	}

	pool.Reserve0 = sync.Reserve0
	pool.Reserve1 = sync.Reserve1
}

func (pool *Pool) ParseSwapEvent(topic string, evt *types.Log) (zeroForOne bool, amount0 *big.Int, amount1 *big.Int) {
	switch topic {
	case TopicAeroSwapV2:
		var aeroSwap swaprouter.AeroV2PoolSwap

		err := AeroV2PoolABI.UnpackIntoInterface(&aeroSwap, "Swap", evt.Data)
		if err != nil {
			logger.Fatal().
				Str("vendor", pool.Vendor).
				Str("factory", pool.Factory).
				Str("txhahs", evt.TxHash.Hex()).
				Str("pool", pool.Address).
				Msg("Unpack AeroV2 Swap event data failed")
		}

		zeroForOne, amount0, amount1 = parseV2Swap(
			aeroSwap.Amount0In,
			aeroSwap.Amount1In,
			aeroSwap.Amount0Out,
			aeroSwap.Amount1Out)

	case TopicSwap:
		var v3Swap swaprouter.UniswapV3PoolSwap

		err := UniswapV3PoolABI.UnpackIntoInterface(&v3Swap, "Swap", evt.Data)
		if err != nil {
			logger.Fatal().
				Str("vendor", pool.Vendor).
				Str("factory", pool.Factory).
				Str("txhahs", evt.TxHash.Hex()).
				Str("pool", pool.Address).
				Msg("Unpack UniswapV3 Swap event data failed")
		}

		if v3Swap.Amount0.Cmp(bigZero) > 0 {
			zeroForOne = true
			amount0 = new(big.Int).Set(v3Swap.Amount0)
			amount1 = new(big.Int).Abs(v3Swap.Amount1)
		} else {
			zeroForOne = false
			amount0 = new(big.Int).Abs(v3Swap.Amount0)
			amount1 = new(big.Int).Set(v3Swap.Amount1)
		}

	case TopicSwapV2:
		var uniswapV2Swap swaprouter.UniswapV2PairSwap

		err := UniswapV2PairABI.UnpackIntoInterface(&uniswapV2Swap, "Swap", evt.Data)
		if err != nil {
			logger.Fatal().
				Str("vendor", pool.Vendor).
				Str("factory", pool.Factory).
				Str("txhahs", evt.TxHash.Hex()).
				Str("pool", pool.Address).
				Msg("Unpack UniswapV2 Swap event data failed")
		}

		zeroForOne, amount0, amount1 = parseV2Swap(
			uniswapV2Swap.Amount0In,
			uniswapV2Swap.Amount1In,
			uniswapV2Swap.Amount0Out,
			uniswapV2Swap.Amount1Out)

	default:
		panic("not swap event, should NOT reach here")
	}

	return
}

func parseV2Swap(amt0In, amt1In, amt0Out, amt1Out *big.Int) (bool, *big.Int, *big.Int) {
	if amt0In.Cmp(bigZero) > 0 {
		return true, new(big.Int).Set(amt0In), new(big.Int).Set(amt1Out)
	} else {
		//
		return false, new(big.Int).Set(amt0Out), new(big.Int).Set(amt1In)
	}
}
