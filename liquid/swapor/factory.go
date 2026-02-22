package swapor

import (
	"github.com/ethereum/go-ethereum/core/types"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

// factory events.

func (eh *EventHandler) doFactoryEvents(fac *common.SwapFactory, event *types.Log) (result common.ParseResult) {
	if fac.Typ == common.PoolTypeInfusionAMM {
		pl, err := pool.DecodeInfusionPairCreatedEvent(fac, event)
		if err == nil {
			result.Status = common.ParsePoolCreated
			result.Data = pl

			logger.Info().
				Str("pool", pl.Address).
				Str("vendor", pl.Vendor).
				Str("Typ", pl.Typ.String()).
				Bool("Stable", pl.Stable).
				Uint("Fee", pl.Fee).
				Msg("found new pool created")

			return
		}
	}

	return
}
