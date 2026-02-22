package handlers

import (
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
)

type Parser struct {
	addressHandlers map[string]common.ILiquidor
}

// NewParser Create new parser.
func NewParser() *Parser {
	return &Parser{
		addressHandlers: map[string]common.ILiquidor{},
	}
}

// AddLiquidor Add liquidor.
func (parser *Parser) AddLiquidor(l common.ILiquidor) {
	parser.addressHandlers[l.Factory()] = l
}

// AddPool Add pool address to filtered.
func (parser *Parser) AddPool(addr string, l common.ILiquidor) {
	parser.addressHandlers[addr] = l
}

// ParseEvents Dispatch events to liquidor, parse events.
func (parser *Parser) ParseEvents(events []types.Log) []common.ParseResult {
	results := []common.ParseResult{}

	for _, event := range events {
		event := event
		addr := strings.ToLower(event.Address.Hex())

		handler, ok := parser.addressHandlers[addr]
		if !ok {
			continue
		}

		// logger.Info().Msgf("event: block=%v address=%v", event.BlockNumber, addr)

		res := handler.ParseEvent(&event)
		if res.Status == common.ParseIgnore {
			continue
		}

		if res.Status == common.ParsePoolCreated {
			pool := res.Data.(*pool.Pool)
			parser.addressHandlers[pool.PoolAddress()] = handler
		}

		results = append(results, res)
	}

	return results
}
