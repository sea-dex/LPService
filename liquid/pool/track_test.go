package pool

import (
	"testing"

	"starbase.ag/liquidity/config"
)

func TestTrackPool(t *testing.T) {
	poolAddr := ""
	cfg := config.Config{}
	TrackPoolV3(&cfg, poolAddr, 0, 100)
}
