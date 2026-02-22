package pool

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/liquid/utils"
	"starbase.ag/liquidity/pkg/logger"
)

// Provider ETH RPC provider.
type Provider struct {
	Rpc    string
	Tps    uint64
	LastTs uint64
	Client *rpc.Client
	used   uint64
	failed uint64
	status bool
	// v3     bool
}

// ProviderPool providers pool.
type ProviderPool struct {
	pools   []*Provider
	poolsV3 []*Provider
	index   int
	indexV3 int

	mu sync.Mutex
}

// NewProvider create provider.
func NewProvider(uri string, tps uint64, v3 bool, tmo uint64) (*Provider, error) {
	// create transport
	if tmo == 0 {
		tmo = 10
	}

	tr, err := transport.NewHTTP(transport.HTTPOptions{
		URL:        uri,
		HTTPClient: &http.Client{Timeout: time.Duration(tmo) * time.Second}, // nolint
	})
	if err != nil {
		// panicWithSentry(err.Error())
		return nil, err
	}

	// Create a JSON-RPC client.
	c, err := rpc.NewClient(rpc.WithTransport(tr))
	if err != nil {
		// panicWithSentry(err.Error())
		return nil, err
	}

	return &Provider{
		Rpc:    uri,
		Tps:    tps,
		LastTs: 0,
		Client: c,
		status: true,
	}, nil
}

// CreateProviderPool create provider pool.
func CreateProviderPool(rpcList []config.ProviderConfig) *ProviderPool {
	pp := &ProviderPool{
		index:   0,
		pools:   []*Provider{},
		poolsV3: []*Provider{},
	}

	v3 := 0

	for _, rpc := range rpcList {
		provider, err := NewProvider(rpc.RPC, rpc.Tps, rpc.V3, rpc.Timeout)
		if err == nil {
			if rpc.V3 {
				v3++

				pp.poolsV3 = append(pp.poolsV3, provider)
			}

			pp.pools = append(pp.pools, provider)
		} else {
			logger.Warn().Str("RPC", rpc.RPC).Err(err).Msg("create provider failed")
		}
	}

	if len(pp.pools) == 0 {
		logger.Fatal().Msg("CreateProviderPool: no provider pool created")
	}

	if v3 == 0 {
		logger.Fatal().Msg("CreateProviderPool: no v3 provider found")
	}

	logger.Info().Msgf("create provider pool: providers: %d v3Providers: %d", len(pp.pools), len(pp.poolsV3))

	return pp
}

// HasIdle check if has idle provider.
func (pp *ProviderPool) HasIdle() bool {
	ts := uint64(time.Now().UnixMilli()) // nolint

	for i := 0; i < len(pp.pools); i++ {
		provider := pp.pools[(pp.index+i)%len(pp.pools)]
		if ts-provider.LastTs < 1000 {
			if provider.used <= provider.Tps {
				return true
			}
		} else {
			return true
		}
	}

	return false
}

// ListSecondProviders list all providers can used in 1 second.
func (pp *ProviderPool) ListSecondProviders() []*Provider {
	ps := []*Provider{}

	for i := 0; i < len(pp.pools); i++ {
		provider := pp.pools[i]
		for j := uint64(0); j < provider.Tps; j++ {
			ps = append(ps, provider)
		}
	}

	return ps
}

/*
// GetNext get next provider.
// func (pp *ProviderPool) GetNext(opts ...any) *Provider {
// 	ts := uint64(time.Now().UnixMilli())

// 	pp.mu.Lock()
// 	defer pp.mu.Unlock()

// 	for i := 0; i < len(pp.pools); i++ {
// 		provider := pp.pools[(pp.index+i+1)%len(pp.pools)]
// 		if !provider.status {
// 			continue
// 		}

// 		if provider.LastTs < ts {
// 			if ts-provider.LastTs < 1000 {
// 				if provider.used < provider.Tps {
// 					pp.index = i
// 					provider.used++

// 					return provider
// 				}
// 			} else {
// 				provider.LastTs = ts
// 				provider.used = 1

// 				return provider
// 			}
// 		}
// 	}

// 	wait := utils.GetBoolVariableParam(opts)
// 	if wait {
// 		time.Sleep(time.Second)
// 		logger.Warn().Msg("no idle provider available, wait 1 second, then use default")
// 	}

	return pp.pools[0]
}
*/

// Get get an v3 idle provider from pool.
func (pp *ProviderPool) GetV3(opts ...any) *Provider {
	ts := uint64(time.Now().UnixMilli()) // nolint

	var bestV3 *Provider

	pp.mu.Lock()
	defer pp.mu.Unlock()

	for i := 0; i < len(pp.poolsV3); i++ {
		provider := pp.poolsV3[(pp.indexV3+i)%len(pp.poolsV3)]
		if bestV3 == nil {
			bestV3 = provider
		} else {
			if provider.Tps > bestV3.Tps {
				bestV3 = provider
			}
		}

		if provider.LastTs < ts {
			if ts-provider.LastTs < 1000 {
				if provider.used < provider.Tps {
					pp.indexV3 = i + 1
					provider.used++

					return provider
				}
			} else {
				provider.LastTs = ts
				provider.used = 1
				pp.indexV3 = i + 1

				return provider
			}
		}
	}

	wait := utils.GetBoolVariableParam(opts)
	if wait {
		time.Sleep(time.Second)
		logger.Warn().Msg("no idle provider available, wait 1 second, then use best v3")
	}

	if bestV3 == nil {
		logger.Fatal().Msg(fmt.Sprintf("not found v3 provider: index=%d", pp.indexV3))
	}

	return bestV3
}

// Get get an idle provider from pool.
func (pp *ProviderPool) Get(opts ...any) *Provider {
	ts := uint64(time.Now().UnixMilli()) // nolint

	pp.mu.Lock()
	defer pp.mu.Unlock()

	for i := 0; i < len(pp.pools); i++ {
		provider := pp.pools[(pp.index+i)%len(pp.pools)]

		if provider.LastTs < ts {
			if ts-provider.LastTs < 1000 {
				if provider.used < provider.Tps {
					pp.index = i + 1
					// provider.LastTs = ts
					provider.used++

					return provider
				}
			} else {
				provider.LastTs = ts
				provider.used = 1
				pp.index = i + 1

				return provider
			}
		}
	}

	logger.Warn().Msg("no available provider currently, use provider 0")

	wait := utils.GetBoolVariableParam(opts)
	if wait {
		time.Sleep(time.Second)
		logger.Warn().Msg("no idle provider available, wait 1 second, then use default")
	}

	return pp.pools[0]
}

// Call call contract.
func (provider *Provider) Call(ctx context.Context, call *types.Call) ([]byte, *types.Call, error) {
	// provider.LastTs = uint64(time.Now().UnixMilli())
	buf, call, err := provider.Client.Call(ctx, call, types.LatestBlockNumber)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			provider.failed++
			// setback
			if provider.failed >= 3 {
				logger.Warn().Str("rpc", provider.Rpc).Msg("provider failed too many times, setback 10 minutes")
				provider.LastTs = uint64(time.Now().UnixMilli()) + 600*1000 // nolint
				provider.failed = 0
			}
		}
	} else {
		provider.failed = 0
	}

	return buf, call, err
}

func (provider *Provider) SetBack(n uint64) {
	provider.failed++
	if provider.failed >= 3 {
		// continues failed, set back 300 seconds
		logger.Warn().Str("provider", provider.Rpc).Msg("provider continues failed too many times, setback 5 minutes")
		provider.LastTs = uint64(time.Now().UnixMilli()) + 300*1000 // nolint
		provider.failed = 0

		return
	}

	logger.Info().Str("provider", provider.Rpc).Msgf("set provider back %d seconds", n/1000)
	provider.LastTs += uint64(time.Now().UnixMilli()) + n // nolint
}

func (pp *ProviderPool) getProvider(v3 bool) *Provider {
	if v3 {
		return pp.GetV3()
	}

	return pp.Get()
}

func (pp *ProviderPool) Multicall(params []Call3, maxRetry uint, name string, v3 bool) (results []Result, err error) {
	calldata := multicallABI.Methods["aggregate3"].MustEncodeArgs(params)
	call := types.NewCall().
		SetTo(multicallAddr).
		SetInput(calldata)

	var b []byte

	for i := uint(0); i < maxRetry; i++ {
		provider := pp.getProvider(v3)

		b, _, err = provider.Call(context.Background(), call)
		if err != nil {
			if seconds := isProviderError(err); seconds > 0 {
				provider.SetBack(uint64(seconds) * 1000)
			}

			if i >= maxRetry {
				return nil, err
			} else {
				logger.Warn().Err(err).Str("RPC", provider.Rpc).Msgf("multicall %s failed, retry=%d", name, i)
			}
		} else {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	// err = multicallABI.Methods["aggregate3"].DecodeValues(b, &results)
	multicallABI.Methods["aggregate3"].MustDecodeValues(b, &results)

	return results, err
}

func isProviderError(err error) uint {
	es := strings.ToLower(err.Error())

	// execution reverted, not provider error.
	if strings.Contains(es, "reverted") {
		return 0
	}

	if strings.Contains(err.Error(), "429") {
		// rate limit, set back 10 seconds
		return 10
	}

	if strings.Contains(es, "context deadline exceeded") {
		return 10
	}

	if strings.Contains(es, "service unavailable") {
		return 30
	}

	if strings.Contains(es, "unexpected end of json input") {
		return 60
	}

	if strings.Contains(es, "not found") {
		return 60
	}

	if strings.Contains(es, "403 forbidden") {
		return 60
	}

	if strings.Contains(es, ": eof") {
		return 60
	}

	return 0
}
