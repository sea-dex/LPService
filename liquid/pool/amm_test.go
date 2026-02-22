package pool

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/liquid/events"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/utils"
)

// 0xc8fcea6430fcfa4cd7cf1f3d1dc83c07b7241e6a 20790143
// 0x148d39c055a73f10d24e8db90976a4567191d009 reserve is nil"} 20789511
// 0x02937f5a9488366241cb625b9209cbb5c0f51777 reserve is nil"}
// 0xff936abb01b87585f8cde34d775c16c21567c4cf reserve is nil"}
// 0xc6208cf6dfa80de770f5ae965aa5e1879ac2b1e0 reserve is nil"}
// 0xf2153719b6d855b2f11ec36d601a3a1c95c1ef3a reserve is nil"} 20788736
// 0x3af0632aadfd69fac67db28af063c5045e0364e7 reserve is nil"}
// 0x2191ec8b9c273e480835b3c3ed0bcde0c272c38f reserve is nil"}
// 0x3210d91f6028e4fe98f935a2c02d3d6dbf2f30df reserve is nil"}
// 0xa3a306f9cdb7d071da6bb551b2704282265bccac reserve is nil"}
// 0x25c39790b322289594d578e4fee203554d2ae525 reserve is nil"}

func TestSync(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "dev", false)
	InitABIs()

	es := events.MustNewEventSubscriber("https://mainnet.base.org", "", false, 0, 0, 0, 0)

	fixtures := []struct {
		address string
		block   uint64
	}{
		{address: "0x148d39c055a73f10d24e8db90976a4567191d009", block: 20789511},
	}

	for _, item := range fixtures {
		logs, err := es.GetLogsFromToReturn(context.Background(), nil, item.block, item.block, 1)
		assert.Nil(t, err)

		for _, evt := range logs {
			// logger.Info().Msgf("txhash: %v address: %v", evt.TxHash.String(), evt.Address)
			if strings.ToLower(evt.Address.String()) == item.address {
				logger.Info().Msgf("address evt: %v %v", item.address, evt.Topics[0].String())

				topic := evt.Topics[0].String()
				if topic == TopicSYNC || topic == TopicSYNCAero {
					pl := &Pool{}
					pl.OnSync(&evt)

					logger.Info().Msgf("sync event: pool=%v reserve0=%v reserve1=%v",
						item.address, pl.Reserve0, pl.Reserve1)
				}
			}
		}
	}
}
