package swapor

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/config"
	starcomm "starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/events"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/utils"
)

func createTestEventHandler() *EventHandler {
	pp := pool.CreateProviderPool([]config.ProviderConfig{
		{
			RPC: "https://mainnet.base.org",
			V3:  true,
			Tps: 10,
		},
	})

	return &EventHandler{pp: pp}
}

func TestFeeEven(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "dev", false)

	es := events.MustNewEventSubscriber("https://mainnet.base.org", "", false, 0, 0, 0, 0)

	logs, err := es.GetLogsFromToReturn(context.Background(), nil, 19432671, 19432671, 1)
	assert.Nil(t, err)

	eh := createTestEventHandler()

	for _, evt := range logs {
		topic := strings.ToLower(evt.Topics[0].String())

		if topic == pool.TopicAeroV2SetFee {
			addr, fee, err := eh.decodeSetCustomFee(&evt, 2)
			assert.Nil(t, err, "txhash: %s", evt.TxHash.String())
			logger.Info().Msgf("pool %s fee change to %d", addr, fee)
		}
	}
}

func TestFeeBlocks(t *testing.T) {
	utils.SkipCI(t)

	// v3 set custom fee 0xb1b52838c29a26f42e00c0657a2f4d5e2f301d2b262e174bd389866949389d6b
	parseFeeBlocks(t, 19948076)
}

func parseFeeBlocks(t *testing.T, block uint64) {
	logger.Init("", "dev", false)

	es := events.MustNewEventSubscriber("https://mainnet.base.org", "", false, 0, 0, 0, 0)
	eh := CreateEventHandler(_feeTestConfig)
	logs, err := es.GetLogsFromToReturn(context.Background(), nil, block, block, 1)
	assert.Nil(t, err)

	logCh := make(chan types.Log, 5000)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for _, item := range logs {
			logCh <- item
		}

		wg.Done()
	}()

	wg.Wait()

	_ = eh.Init("resume", true)
	eh.blockNumber = block - 1
	_, feeEvents, _, _, _ := eh.drainEvents(logCh)
	logger.Info().Msgf("fee events: %d", len(feeEvents))
	eh.handleFeeEvents(feeEvents)
}

var _feeTestConfig = &config.Config{
	Env: "local",
	// Mode          string       `toml:"mode"`
	Log: config.Log{Level: "info"},
	Chain: config.ChainConfig{
		ChainID:        "8453",
		Network:        "mainnet",
		URL:            "https://mainnet.base.org",
		WSS:            "wss://base-rpc.publicnode.com",
		Subscribe:      false,
		StartingHeight: 0,
		Factory: []starcomm.SwapFactory{
			{
				Address: "0x5e7BB104d84c7CB9B682AaC2F3d509f5F406809A",
				Name:    "AerodromeV3",
				Typ:     301,
			},
			{
				Address: "0x33128a8fC17869897dcE68Ed026d694621f6FDfD",
				Name:    "UniswapV3",
				Typ:     300,
			},
			{
				Address: "0xc35DADB65012eC5796536bD9864eD8773aBc74C4",
				Name:    "SushiswapV3",
				Typ:     300,
			},
			{
				Address: "0x38015D05f4fEC8AFe15D7cc0386a126574e8077B",
				Name:    "BaseSwap-Basex",
				Typ:     300,
			},
			{
				Address: "0x0BFbCF9fa4f9C56B0F40a671Ad40E0805A091865",
				Name:    "PancakeswapV3",
				Typ:     302,
			},
			{
				Address: "0x8909Dc15e40173Ff4699343b6eB8132c65e18eC6",
				Name:    "UniswapV2",
				Typ:     200,
			},
			{
				Address: "0x71524B4f93c58fcbF659783284E38825f0622859",
				Name:    "SushiswapV2",
				Typ:     200,
			},
			{
				Address: "0x3E84D913803b02A4a7f027165E8cA42C14C0FdE7",
				Name:    "Alien",
				Typ:     200,
			},
			{
				Address: "0x02a84c1b3BBD7401a5f7fa98a384EBC70bB5749E",
				Name:    "PancakeswapV2",
				Typ:     200,
			},
			{
				Address: "0x2d5dd5fa7B8a1BFBDbB0916B42280208Ee6DE51e",
				Name:    "Alien-Area51",
				Typ:     200,
			},
			{
				Address: "0xFDa619b6d20975be80A10332cD39b9a4b0FAa8BB",
				Name:    "Baseswap",
				Typ:     200,
			},
			{
				Address: "0x420DD381b31aEf6683db6B902084cB0FFECe40Da",
				Name:    "AerodromeV1",
				Typ:     201,
			},
			{
				Address: "0x2d9a3a2bd6400ee28d770c7254ca840c82faf23f",
				Name:    "Infusion",
				Typ:     201,
			},
		},
		PairsQuery: "0x8fb641dfe7173fC58C7Edb5BCC13a7187881b96E",
		PoolQuery: map[string]uint{
			"0xDf7acDFaab84FE57c999aEf080749845C97ca038": 300,
			"0xB369A9B58bE84783F47E66c244F79567E36367B1": 301,
			"0x398FbFe61579090aEcC613a25BdeCffBa8D60313": 302,
		},
		Providers: []config.ProviderConfig{
			{RPC: "https://mainnet.base.org", V3: true, Tps: 10},
		},
	},
	Kafka: config.KafkaConfig{
		Brokers: "localhost:9092",
	},
	DB:            config.DBConfig{},
	Redis:         config.RedisConfig{},
	HTTPServer:    config.ServerConfig{},
	MetricsServer: config.ServerConfig{},
	SentryOpts:    config.SentryOpts{},
}
