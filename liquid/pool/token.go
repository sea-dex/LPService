package pool

import (
	"strings"

	defiabi "github.com/defiweb/go-eth/abi"
	defitypes "github.com/defiweb/go-eth/types"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
)

var (
	nameMethodABI     = defiabi.MustParseMethod("function name() external view returns (string memory)")
	symbolMethodABI   = defiabi.MustParseMethod("function symbol() external view returns (string memory)")
	decimalsMethodABI = defiabi.MustParseMethod("function decimals() external view returns (uint8)")
)

func (pp *ProviderPool) LoadTokens(addresses []string) ([]*common.Token, []string, error) {
	params := []Call3{}

	for _, addr := range addresses {
		tokenAddr := defitypes.MustAddressFromHex(addr)
		// token name
		params = append(params, Call3{
			Target:       tokenAddr,
			CallData:     nameMethodABI.MustEncodeArgs(),
			AllowFailure: true,
		},
			// token symbol
			Call3{
				Target:       tokenAddr,
				CallData:     symbolMethodABI.MustEncodeArgs(),
				AllowFailure: true,
			},
			// token decimals
			Call3{
				Target:       tokenAddr,
				CallData:     decimalsMethodABI.MustEncodeArgs(),
				AllowFailure: true,
			})
	}

	results, err := pp.Multicall(params, 3, "LoadTokens", false)
	if err != nil {
		logger.Warn().Err(err).Msgf("load tokens failed: tokens=%d", len(addresses))
		return nil, nil, err
	}

	tokens := []*common.Token{}
	failedTokens := []string{}

	for i := 0; i < len(results)/3; i++ {
		var token common.Token

		token.Address = addresses[i]
		if results[i*3].Success && len(results[i*3].ReturnData) > 0 {
			nameMethodABI.MustDecodeValues(results[i*3].ReturnData, &token.Name)
		} else {
			logger.Warn().Str("token", token.Address).Msg("token name invalid")
		}

		if results[i*3+1].Success && len(results[i*3+1].ReturnData) > 0 {
			symbolMethodABI.MustDecodeValues(results[i*3+1].ReturnData, &token.Symbol)
		} else {
			logger.Warn().Str("token", token.Address).Msg("token symbol invalid")
		}

		if results[i*3+2].Success && len(results[i*3+2].ReturnData) > 0 {
			decimalsMethodABI.MustDecodeValues(results[i*3+2].ReturnData, &token.Decimals)
			tokens = append(tokens, &token)
		} else {
			logger.Warn().Str("token", token.Address).Msg("token decimals invalid")

			failedTokens = append(failedTokens, addresses[i])
		}
	}

	return tokens, failedTokens, nil
}

func (pp *ProviderPool) LoadTokensMap(newTokens map[string]bool) ([]*common.Token, error) {
	maxStep := 500

	logger.Info().Msgf("prepare load tokens: %d", len(newTokens))

	tokenList := []string{}
	tokens := []*common.Token{}

	for addr := range newTokens {
		addr = strings.TrimSpace(addr)
		if addr == "" || addr == "0x" {
			logger.Warn().Msgf("empty token address: %v", addr)
			continue
		}

		tokenList = append(tokenList, addr)
		if len(tokenList) >= maxStep {
			items, _, err := pp.LoadTokens(tokenList)
			if err == nil {
				tokens = append(tokens, items...)
			} else {
				logger.Warn().Strs("tokens", tokenList).Msg("get token info from chain failed")
			}

			tokenList = []string{}
		}
	}

	if len(tokenList) > 0 {
		items, _, err := pp.LoadTokens(tokenList)
		if err != nil {
			logger.Error().Err(err).Msgf("get token info failed, tokenList: %d", len(tokenList))
		} else {
			tokens = append(tokens, items...)
		}
	}

	return tokens, nil
}
