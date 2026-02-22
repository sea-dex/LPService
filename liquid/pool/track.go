package pool

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/liquid/events"
)

func TrackPoolV3(cfg *config.Config, poolAddr string, start, end uint64) {
	pl := preparePool(cfg, poolAddr)
	evts := getEvents(cfg, poolAddr, start, end)

	for _, event := range evts {
		pl.OnV3Event(&event)
	}
}

func getEvents(cfg *config.Config, poolAddr string, start, end uint64) []types.Log {
	es, err := events.NewEventSubscirber(cfg.Chain.URL, "", false, 0, 0, 0, 0)
	if err != nil {
		panic(err.Error())
	}

	step := uint32(100)
	ctx := context.Background()
	filters := []common.Address{common.HexToAddress(poolAddr)}

	evts, err := es.GetLogsFromToReturn(ctx, filters, start, end, step)
	if err != nil {
		panic(err.Error())
	}

	return evts
}

func preparePool(cfg *config.Config, poolAddr string) *Pool {
	pp := CreateProviderPool(cfg.Chain.Providers)

	pool, err := pp.GetPoolBasicInfo(poolAddr, 3)
	if err != nil {
		panic(err.Error())
	}

	if pool.Typ.IsAMMVariety() {
		err = pp.GetV2PoolLiquid(pool)
	} else if pool.Typ.IsCAMMVariety() {
		query := ""

		for k, v := range cfg.Chain.PoolQuery {
			if v == uint(pool.Typ) {
				query = k
			}
		}

		err = pp.GetV3PoolLiquidity(pool, query, 0)
	}

	if err != nil {
		panic(err.Error())
	}

	pool.LastBlockUpdated = 0
	pool.tickBitmap = map[int16]*big.Int{}
	pool.Ticks = map[int]*TickInfo{}

	return pool
}
