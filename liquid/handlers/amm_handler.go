package handlers

import (
	"github.com/ethereum/go-ethereum/core/types"
	"starbase.ag/liquidity/liquid/common"
)

// AMMHandler AMM handler.
type AMMHandler struct {
	factory string
	vendor  string
}

// Factory uniswap v2 factory address.
func (ah *AMMHandler) Factory() string {
	return ah.factory
}

// Factory uniswap v2 factory address.
func (ah *AMMHandler) Vendor() string {
	return ah.vendor
}

func (ah *AMMHandler) ParseEvent(event *types.Log) (result common.ParseResult) {
	return
}
