package handlers

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"starbase.ag/liquidity/contracts/swaprouter"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/liquid/utils"
	"starbase.ag/liquidity/pkg/logger"
)

// CAMMHandler CAMM handler.
type CAMMHandler struct {
	factory string
	vendor  string
}

func CreateCAMMHandler(vendor, factory string) *CAMMHandler {
	return &CAMMHandler{
		factory: strings.ToLower(factory),
		vendor:  vendor,
	}
}

// Factory uniswap v2 factory address.
func (ch *CAMMHandler) Factory() string {
	return ch.factory
}

// Factory uniswap v2 factory address.
func (ch *CAMMHandler) Vendor() string {
	return ch.vendor
}

func (ch *CAMMHandler) ParseEvent(event *types.Log) (result common.ParseResult) {
	addr := strings.ToLower(event.Address.Hex())
	topics := event.Topics
	topic := strings.ToLower(topics[0].Hex())

	if addr == ch.factory {
		if topic == pool.TopicPoolCreated {
			// https://etherscan.io/tx/0xf87d91f3d72a8e912c020c2e316151f3557b1217b44d4f6b6bec126448318530#eventlog
			var created swaprouter.UniswapV3FactoryPoolCreated
			if err := pool.UniswapV3FactoryABI.UnpackIntoInterface(&created, "PoolCreated", event.Data); err != nil {
				logger.Fatal().
					Err(err).
					Str("factory", ch.factory).
					Str("txhahs", event.TxHash.Hex()).
					Msg("Unpack PoolCreated event data failed")
			}

			pool := pool.CreateCAMMPool(
				utils.HashToAddress(topics[1]),
				utils.HashToAddress(topics[2]),
				created.Pool.Hex(),
				ch.factory,
				uint(big.NewInt(0).SetBytes(topics[3].Bytes()).Int64()), // nolint
				int(created.TickSpacing.Int64()),
				event.BlockNumber,
			)

			logger.Info().Msgf("new pool: block=%d vendor=%v pool=%v", event.BlockNumber, ch.vendor, pool.Address)

			return common.ParseResult{
				Status: common.ParsePoolCreated,
				Data:   pool,
			}
		}

		return
	}

	return ch.ParsePoolEvent(event, topic)
}

func (ch *CAMMHandler) ParsePoolEvent(event *types.Log, topic string) (result common.ParseResult) {
	poolAddr := strings.ToLower(event.Address.Hex())
	pl := pool.GetCAMMPool(poolAddr)
	// topics := event.Topics
	// txhash := event.TxHash.Hex()

	switch topic {
	case pool.TopicInitialize:
		var initialize swaprouter.UniswapV3PoolInitialize
		if err := pool.UniswapV3PoolABI.UnpackIntoInterface(&initialize, "Initialize", event.Data); err != nil {
			logger.Fatal().
				Err(err).
				Str("vendor", pl.Vendor).
				Str("factory", pl.Factory).
				Str("txhahs", event.TxHash.Hex()).
				Str("pool", pl.Address).
				Msg("Unpack Initialize event data failed")
		}

		tick := int(initialize.Tick.Int64())
		pl.OnInitialize(initialize.SqrtPriceX96, tick)
		logger.Info().Msgf("%s pool Initialized: pool=%s tick=%v", ch.vendor, poolAddr, tick)

	case pool.TopicMint:
		pl.OnMint(event)

	case pool.TopicSwap:
		pl.OnSwap(event)

	case pool.TopicBurn:
		pl.OnBurn(event)

	case pool.TopicCollect:
		pl.OnCollect(event)

	case pool.TopicPancakeSwap:
		pl.OnPancakeSwap(event)

	default:
		return
	}

	pl.LastBlockUpdated = event.BlockNumber
	result.Status = common.ParsePoolUpdated
	result.Data = pl

	return
}
