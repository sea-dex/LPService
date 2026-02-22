package arb

import (
	"testing"

	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/utils"
)

func TestGas(t *testing.T) {
	utils.SkipCI(t)
	logger.Init("", "", false)
	FetchGasPrice("base")
}
