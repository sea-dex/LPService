package pool

import (
	"fmt"

	defitypes "github.com/defiweb/go-eth/types"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
)

type PoolFeeChanged struct {
	Pool string `json:"pool"`
	Fee  uint   `json:"fee"`
}

// stable fee
// volatile fee.
func (pp *ProviderPool) GetInfusionSwapFee(factory *common.SwapFactory) {
	factoryAddr := defitypes.MustAddressFromHex(factory.Address)
	params := []Call3{
		{
			Target:   factoryAddr,
			CallData: getInfusionFeeMethodABI.MustEncodeArgs(true),
		}, {
			Target:   factoryAddr,
			CallData: getInfusionFeeMethodABI.MustEncodeArgs(false),
		},
	}

	results, err := pp.Multicall(params, 3, "getInfusionFee", false)
	if err != nil {
		logger.Warn().
			Err(err).
			Str("factory", factory.Address).
			Msgf("getInfusionSwapFee failed")

		return
	}

	var (
		stableFee   uint
		volatileFee uint
	)

	// Decode the results.
	getInfusionFeeMethodABI.MustDecodeValues(results[0].ReturnData, &stableFee)
	getInfusionFeeMethodABI.MustDecodeValues(results[1].ReturnData, &volatileFee)

	factory.StableFee = int(stableFee * 100) // nolint convet to 1000000 denominator
	factory.Fee = int(volatileFee * 100)     // nolint convert to 1000000 denominator

	logger.Info().Str("factory", factory.Address).
		Str("factory", factory.Name).
		Msgf("get factory fee: stableFee=%d volatileFee=%d", factory.StableFee, factory.Fee)
}

// GetAeroPoolV2Stable get pool is stable pool.
func (pp *ProviderPool) GetAeroPoolV2Fee(pls []*Pool, maxRetry uint) ([]PoolFeeChanged, error) {
	step := 1000
	feeChanges := []PoolFeeChanged{}

	for i := 0; i < len(pls); i += step {
		end := i + step
		if end > len(pls) {
			end = len(pls)
		}

		changes, err := pp.getAeroPoolV2Fee(pls[i:end], i, maxRetry)
		if err == nil {
			feeChanges = append(feeChanges, changes...)
		}

		i = end
	}

	return feeChanges, nil
}

func (pp *ProviderPool) getAeroPoolV2Fee(pls []*Pool, start int, maxRetry uint) ([]PoolFeeChanged, error) {
	if len(pls) == 0 {
		return []PoolFeeChanged{}, nil
	}

	callParams := []Call3{}
	for _, pl := range pls {
		callParams = append(callParams, Call3{
			Target:   defitypes.MustAddressFromHex(pl.Factory),
			CallData: aeroGetFeeABI.MustEncodeArgs(pl.Address, pl.Stable),
		})
	}

	feeChanges := []PoolFeeChanged{}

	results, err := pp.Multicall(callParams, maxRetry, fmt.Sprintf("getAeroPoolV2Fee-%d-%d", start, len(pls)), false)
	if err != nil {
		logger.Warn().
			Err(err).
			Int("pools", len(callParams)).
			Msgf("getAeroPoolV2Fee failed")

		return feeChanges, err
	}

	for i, result := range results {
		var fee uint

		aeroGetFeeABI.MustDecodeValues(result.ReturnData, &fee)
		fee *= 100
		pl := pls[i]

		if fee != pl.Fee {
			// make fee to bps(1000,000)
			pl.Fee = fee
			feeChanges = append(feeChanges, PoolFeeChanged{
				Pool: pl.Address,
				Fee:  fee,
			})

			logger.Info().
				Str("pool", pl.Address).
				Str("vendor", pl.Vendor).
				Msgf("pool fee changes to %d", fee)
		}
	}

	return feeChanges, nil
}
