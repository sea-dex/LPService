package common

import (
	"github.com/ethereum/go-ethereum/core/types"
)

// ILiquidor.
type ILiquidor interface {
	// Factory factory address
	Factory() string
	// Vendor name, such as UniswapV2, UniswapV3
	Vendor() string
	ParseEvent(event *types.Log) ParseResult
}

// ParseStatus result status.
type ParseStatus int

const (
	ParseIgnore      = ParseStatus(0)
	ParsePoolCreated = ParseStatus(1)
	ParsePoolUpdated = ParseStatus(2)
)

type ParseResult struct {
	Status ParseStatus
	Data   any
}
