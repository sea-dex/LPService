package swapor

import (
	"encoding/json"
	"fmt"
	"strings"

	// "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	// "github.com/confluentinc/confluent-kafka-go/v2/kafka".
	"github.com/ethereum/go-ethereum/core/types"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

const (
	aeroFeeMultiplier = 100
)

// handleFeeEvents parse set custom fee events.
func (eh *EventHandler) handleFeeEvents(feeEvents []types.Log) {
	for _, evt := range feeEvents {
		topic := strings.ToLower(evt.Topics[0].String())

		switch topic {
		case pool.TopicAeroV2SetFee:
			eh.handleSetCustomFee(&evt, 2)
		case pool.TopicAeroV3SetFee:
			eh.handleSetCustomFee(&evt, 3)
		default:
			logger.Error().Str("topic", topic).Msg("handleFeeEvents: invalid event topic")
		}
	}
}

func (eh *EventHandler) handleSetCustomFee(feeEvent *types.Log, version int) {
	addr, fee, err := eh.decodeSetCustomFee(feeEvent, version)
	if err != nil {
		logger.Error().Err(err).
			Uint64("block", feeEvent.BlockNumber).
			Str("txhash", feeEvent.TxHash.String()).
			Str("address", feeEvent.Address.String()).
			Msgf("decode SetCustomFee V%d failed", version)

		return
	}

	pl := eh.pools[addr]
	if pl == nil {
		logger.Warn().Str("pool", addr).
			Uint64("block", feeEvent.BlockNumber).
			Str("txhash", feeEvent.TxHash.String()).
			Str("address", feeEvent.Address.String()).
			Msg("handleSetCustomFee: not found pool")

		return
	}

	updatePoolFee(pl, version, fee, feeEvent.BlockNumber, feeEvent.TxHash.String())
}

func (eh *EventHandler) decodeSetCustomFee(feeEvent *types.Log, version int) (string, int, error) {
	switch version {
	case 2:
		return eh.pp.DecodeSetCustomFeeV2(feeEvent)
	case 3:
		return eh.pp.DecodeSetCustomFeeV3(feeEvent)
	default:
		return "", 0, fmt.Errorf("invalid version %d for handleSetCustomFee", version)
	}
}

func updatePoolFee(pl *pool.Pool, version, fee int, blocknumber uint64, txhash string) {
	if fee == pool.ZERO_FEE_INDICATOR {
		logger.Info().Str("pool", pl.Address).
			Uint64("block", blocknumber).
			Str("txhash", txhash).
			Msg("set pool fee to 0")

		pl.Fee = 0
	} else {
		var newFee uint
		if version == 2 {
			newFee = uint(fee * aeroFeeMultiplier) // nolint
		} else if version == 3 {
			newFee = uint(fee) // nolint
		} else {
			panic("invalid version")
		}

		logger.Info().Str("pool", pl.Address).
			Str("vendor", pl.Vendor).
			Uint("oldFee", pl.Fee).
			Uint("newFee", newFee).
			Uint64("blocknumber", blocknumber).
			Str("txhash", txhash).
			Msg("pool fee updated")

		pl.Fee = newFee
	}
}

// check aerov2 pools fee
// check infusion factory fee.
func (eh *EventHandler) checkPoolsFee() {
	feeChangedPools := eh.checkAeroPoolsFee()
	infusionFeeChanges := eh.checkInfusionPoolsFee()

	changedPools, changedFees := mergePools(feeChangedPools, infusionFeeChanges)

	logger.Info().Msgf("pool fee changed: aerov2=%d infusion=%d", len(feeChangedPools), len(infusionFeeChanges))

	if err := eh.FlushToRedis(changedPools); err != nil {
		logger.Warn().Err(err).Msg("Failed to flush changed pools to Redis")
	}

	// msgs := []kafka.Message{}
	msgs := [][]byte{}

	for _, item := range changedFees {
		buf, err := json.Marshal(item)
		if err != nil {
			logger.Error().Err(err).Str("pool", item.Pool).Uint("fee", item.Fee).Msg("marshal pool fee changed failed")
			continue
		}

		// msgs = append(msgs, kafka.Message{
		// 	Value: buf,
		// })
		msgs = append(msgs, buf)
	}

	_ = msgs
	// if err := eh.feeProducer.WriteMessages(context.Background(), msgs...); err != nil {
	// 	logger.Error().Err(err).Msg("producer pools fees to kafka failed")
	// }
	// remove kafka dependency
	// eh.producer.ProduceMessage(config.PoolFeeTopic, msgs)
}

func (eh *EventHandler) checkAeroPoolsFee() []*pool.Pool {
	aerov2 := eh.getAeroV2Pools()
	pls := []*pool.Pool{}

	if len(aerov2) == 0 {
		return pls
	}

	items, err := eh.pp.GetAeroPoolV2Fee(aerov2, 3)
	if err != nil {
		logger.Warn().Err(err).Int("pools", len(aerov2)).Msg("GetAeroPoolV2Fee failed")
		return pls
	}

	for _, item := range items {
		pls = append(pls, eh.pools[item.Pool])
	}

	logger.Info().Msgf("AeroAMM pool fee changed: %d", len(pls))

	return pls
}

func mergePools(v1, v2 []*pool.Pool) (map[string]*pool.Pool, []pool.PoolFeeChanged) {
	pls := map[string]*pool.Pool{}
	changed := []pool.PoolFeeChanged{}

	for _, item := range v1 {
		pls[item.Address] = item
		changed = append(changed, pool.PoolFeeChanged{
			Pool: item.Address,
			Fee:  item.Fee,
		})
	}

	for _, item := range v2 {
		pls[item.Address] = item
		changed = append(changed, pool.PoolFeeChanged{
			Pool: item.Address,
			Fee:  item.Fee,
		})
	}

	return pls, changed
}

func (eh *EventHandler) checkInfusionPoolsFee() []*pool.Pool {
	pls := []*pool.Pool{}
	changedFactory := map[string]bool{}

	for _, fac := range eh.factory {
		if fac.Typ == common.PoolTypeInfusionAMM {
			oldStableFee := fac.StableFee
			oldVolatileFee := fac.Fee

			eh.pp.GetInfusionSwapFee(fac)

			if fac.StableFee == oldStableFee && fac.Fee == oldVolatileFee {
				continue
			} else {
				changedFactory[fac.Address] = true
			}
		}
	}

	if len(changedFactory) > 0 {
		for _, pl := range eh.pools {
			if changedFactory[pl.Factory] {
				fac := eh.factory[pl.Factory]
				if pl.Stable {
					pl.Fee = uint(fac.StableFee) // nolint
				} else {
					pl.Fee = uint(fac.Fee) // nolint
				}

				pls = append(pls, pl)
			}
		}
	}

	logger.Info().Msgf("type InfusionAMM pool fee changed: %d", len(pls))

	return pls
}

func (eh *EventHandler) getAeroV2Pools() []*pool.Pool {
	aerov2 := []*pool.Pool{}

	for _, pl := range eh.pools {
		if pl.Typ == common.PoolTypeAeroAMM {
			aerov2 = append(aerov2, pl)
		}
	}

	return aerov2
}
