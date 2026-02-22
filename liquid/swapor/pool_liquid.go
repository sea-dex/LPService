package swapor

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

func (eh *EventHandler) fetchNewPoolAndTokens(poolAddrs []string) {
	if len(poolAddrs) == 0 {
		return
	}

	pools, _, err := eh.pp.GetPoolsBasicInfo(poolAddrs, 3)
	if err != nil {
		return
	}

	unknown := []*pool.PoolBasicInfo{}
	aerov2 := []*pool.Pool{}
	newTokens := map[string]bool{}

	for _, pl := range pools {
		fac, ok := eh.factory[pl.Factory]
		newTokens[pl.Token0] = true
		newTokens[pl.Token1] = true

		if !ok {
			// this is an unknown pool
			pbi := &pool.PoolBasicInfo{
				Factory: pl.Factory,
				Address: pl.Address,
				Token0:  pl.Token0,
				Token1:  pl.Token1,
			}
			unknown = append(unknown, pbi)
		} else {
			fac.Pools++
			pl.Vendor = fac.Name
			pl.Typ = fac.Typ
			pl.Fee = uint(fac.Fee) // nolint
			// update to latest version
			// pl.Version = pool.LatestVersion

			if pl.Typ.IsAMMVariety() {
				pl.Initialized = true
				if pl.Typ == common.PoolTypeAeroAMM {
					aerov2 = append(aerov2, pl)
				}
			}

			eh.pools[pl.Address] = pl
		}
	}

	if len(aerov2) > 0 {
		eh.getAeroV2PoolsInfo(aerov2)
	}

	eh.foundNewUnknownPools(unknown)
	eh.foundNewTokens(newTokens)
}

// get pool stable and pool fee.
func (eh *EventHandler) getAeroV2PoolsInfo(pls []*pool.Pool) {
	if err := eh.pp.GetAeroPoolV2Stable(pls, 3); err != nil {
		logger.Warn().Err(err).Int("pools", len(pls)).Msg("get aero pool v2 pool stable failed")
	}

	if _, err := eh.pp.GetAeroPoolV2Fee(pls, 3); err != nil {
		logger.Warn().Err(err).Int("pools", len(pls)).Msg("get aero pool v2 pool fee failed")
	}
}

// 1. first, get pool from eh.pools, if exist, return it
// 2. second, get pool basic info: factory, token0, token1
// 3. get pool liquidity.
func (eh *EventHandler) getPool(addr string) (*pool.Pool, error) {
	pl, ok := eh.pools[addr]
	if ok {
		var err error

		if pl.Typ.IsCAMMVariety() && !pl.Initialized {
			err = eh.getPoolLiquidity(pl)
			if err != nil {
				logger.Warn().Err(err).
					Str("pool", pl.Address).
					Str("poolType", pl.Typ.String()).
					Msg("get v3 pool liquidity failed")
			}
		}

		return pl, err
	}

	if _, ok = eh.unknownPools[addr]; ok {
		return nil, ErrUnknownPool
	}

	logger.Warn().Str("pool", addr).Msg("not found pool")

	var err error

	pl, err = eh.pp.GetPoolBasicInfo(addr, 3)
	if err != nil {
		logger.Warn().Err(err).Str("pool", addr).Msg("get pool basic info failed")
		return nil, err
	}

	fac, ok := eh.factory[pl.Factory]
	if !ok {
		// this is an unknown pool
		pbi := &pool.PoolBasicInfo{
			Factory: pl.Factory,
			Address: pl.Address,
			Token0:  pl.Token0,
			Token1:  pl.Token1,
		}
		eh.foundNewUnknownPool(pbi)

		return nil, ErrUnknownPool
	}

	pl.Vendor = fac.Name
	pl.Typ = fac.Typ
	pl.Fee = uint(fac.Fee) // nolint

	err = eh.getPoolLiquidity(pl)
	if err != nil {
		logger.Warn().Err(err).
			Str("pool", pl.Address).
			Str("poolType", pl.Typ.String()).
			Msg("get pool liquidity failed")
	}

	return pl, err
}

func (eh *EventHandler) foundNewUnknownPools(pls []*pool.PoolBasicInfo) {
	pipe := eh.store.Pipeline()
	ctx := context.Background()

	for _, pl := range pls {
		eh.unknownPools[pl.Address] = pl
		// store to redis
		buf, err := json.Marshal(pl)
		if err != nil {
			logger.Error().Err(err).Str("pool", pl.Address).Msg("marshal pool basic info failed")
			continue
		}

		pipe.HSet(ctx, unknownPoolInfoKey, pl.Address, buf)
		pipe.HSet(ctx, unknownFactory, pl.Factory, pl.Factory)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("add pool factory to redis unknownPoolInfo hashmap failed")
	}
}

func (eh *EventHandler) foundNewTokens(newTokens map[string]bool) {
	tokens, err := eh.pp.LoadTokensMap(newTokens)
	if err != nil {
		return
	}

	if len(tokens) == 0 {
		return
	}

	pipe := eh.store.Pipeline()
	ctx := context.Background()

	for _, token := range tokens {
		eh.tokens[token.Address] = token
		// store to redis
		buf, err := json.Marshal(token)
		if err != nil {
			logger.Error().Err(err).Str("token", token.Address).Msg("marshal token failed")
			continue
		}

		pipe.HSet(ctx, common.TokenInfoKey, token.Address, buf)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("set token to redis tokenInfoKey hashmap failed")
	} else {
		logger.Info().Msgf("set %d new tokens to redis", len(tokens))
	}
}

// check token not found in pools.
func (eh *EventHandler) checkPoolTokens() {
	newTokens := map[string]bool{}

	for _, pool := range eh.pools {
		if _, ok := eh.tokens[pool.Token0]; !ok {
			newTokens[pool.Token0] = true
		}

		if _, ok := eh.tokens[pool.Token1]; !ok {
			newTokens[pool.Token1] = true
		}
	}

	if len(newTokens) > 0 {
		logger.Info().Msgf("found %d new tokens in pools", len(newTokens))
		eh.foundNewTokens(newTokens)
	}
}

func (eh *EventHandler) foundNewUnknownPool(pl *pool.PoolBasicInfo) {
	eh.unknownPools[pl.Address] = pl
	// store to redis
	buf, err := json.Marshal(pl)
	if err != nil {
		logger.Error().Err(err).
			Str("pool", pl.Address).
			Msg("marshal pool basic info failed")

		return
	}

	ctx := context.Background()

	err = eh.store.HSet(ctx, unknownPoolInfoKey, pl.Address, buf).Err()
	if err != nil {
		logger.Error().Err(err).
			Str("pool", pl.Address).
			Msg("add pool to redis unknownPoolInfo hashmap failed")
	}

	if pl.Factory != "" {
		err = eh.store.HSet(ctx, unknownFactory, pl.Factory, pl.Factory).Err()
		if err != nil {
			logger.Error().Err(err).
				Str("pool", pl.Address).
				Str("factory", pl.Factory).
				Msg("add factory to redis unknownPoolInfo hashmap failed")
		}
	}
}

func (eh *EventHandler) doV2PoolEvents(pl *pool.Pool, events []types.Log) {
	for _, evt := range events {
		topic := strings.ToLower(evt.Topics[0].Hex())

		switch topic {
		case pool.TopicSYNC:
			fallthrough
		case pool.TopicSYNCAero:
			// logger.Info().
			// 	Str("vendor", pl.Vendor).
			// 	Str("pool", pl.Address).
			// 	Str("reserve0", pl.Reserve0.String()).
			// 	Str("reserve1", pl.Reserve1.String()).
			// 	Uint64("block", evt.BlockNumber).
			// 	Msg("v2 pool sync event")
			pl.OnSync(&evt)

		case pool.TopicAeroSwapV2:
			fallthrough
		case pool.TopicSwapV2:
			// do nothing

		default:
			logger.Fatal().Msgf("unknown v2 pool event: %s txhash=%v", topic, evt.TxHash.String())
		}
	}
}

func (eh *EventHandler) doV3PoolEvents(pl *pool.Pool, events []types.Log) {
	// poolAddr := pl.Address
	for _, event := range events {
		if pl.Synced {
			if event.BlockNumber < pl.InitBlock {
				// logger.Info().Str("pool", pl.Address).
				// 	Str("txhash", event.TxHash.String()).
				// 	Msgf("synced v3 pool event skip: block=%d initBlock=%d", event.BlockNumber, pl.InitBlock)

				continue
			}
		} else {
			if event.BlockNumber <= pl.InitBlock {
				// logger.Info().Str("pool", pl.Address).
				// 	Str("txhash", event.TxHash.String()).
				// 	Msgf("v3 pool event skip: block=%d initBlock=%d", event.BlockNumber, pl.InitBlock)
				// only
				continue
			}
		}

		topics := event.Topics

		topic := strings.ToLower(topics[0].Hex())
		switch topic {
		// case TopicInitialize:
		case pool.TopicMint:
			pl.OnMint(&event)

		case pool.TopicSwap:
			if !pl.OnSwap(&event) {
				// reload pool liqudity
				_ = eh.getPoolLiquidity(pl)
			}

		case pool.TopicBurn:
			pl.OnBurn(&event)

		case pool.TopicCollect:
			pl.OnCollect(&event)

		case pool.TopicPancakeSwap:
			if !pl.OnPancakeSwap(&event) {
				// reload pool liqudity
				_ = eh.getPoolLiquidity(pl)
			}

		default:
			logger.Error().Str("pool", pl.Address).Str("txhash", event.TxHash.String()).Msg("unknown v3 pool event topic")
		}
	}
}

// getPoolLiquidity get pool liquidity, addr is always lowercase.
func (eh *EventHandler) getPoolLiquidity(pl *pool.Pool) (err error) {
	switch pl.Typ {
	case common.PoolTypeAMM:
		fallthrough
	case common.PoolTypeInfusionAMM:
		fallthrough
	case common.PoolTypeAeroAMM:
		if err = eh.pp.GetV2PoolLiquid(pl); err == nil {
			if pl.Fee != 0 {
				pl.Fee = uint(eh.factory[pl.Factory].Fee) // nolint
			}

			logger.Info().
				Str("pool", pl.Address).
				Str("vendor", pl.Vendor).
				Uint("fee", pl.Fee).
				Str("reserve0", pl.Reserve0.String()).
				Str("reserve1", pl.Reserve1.String()).
				Uint64("block", pl.InitBlock).
				Msg("load initial v2 pool liquidity onchain")
		}

	case common.PoolTypeCAMM:
		fallthrough
	case common.PoolTypeAeroCAMM:
		fallthrough
	case common.PoolTypePancakeCAMM:
		if err = eh.pp.GetV3PoolLiquidity(pl, eh.v3poolQuery[pl.Typ], 0); err == nil {
			logger.Info().
				Str("pool", pl.Address).
				Str("vendor", pl.Vendor).
				Int("tick", pl.Tick).
				Uint("fee", pl.Fee).
				Int("tickSpacing", pl.TickSpacing).
				Int("tickList", len(pl.Ticks)).
				Str("liquidity", pl.Liquidity.String()).
				Str("sqrtPriceX96", pl.SqrtPriceX96.String()).
				Uint64("block", pl.InitBlock).
				Msg("load initial v3 pool liquidity onchain")
		}

	default:
		logger.Fatal().
			Str("vendor", pl.Vendor).
			Str("pool", pl.Address).
			Uint("poolType", uint(pl.Typ)).
			Msg("unknown pool type")
	}

	return
}
