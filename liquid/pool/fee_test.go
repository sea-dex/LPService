package pool

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/utils"
)

func TestProviderPool_GetAeroPoolV2Fee(t *testing.T) {
	utils.SkipCI(t)

	logger.Init("", "dev", true)

	rc := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	m, err := rc.HGetAll(context.Background(), "poolLiquidity").Result()
	assert.Nil(t, err)

	pls := []*Pool{}

	for _, val := range m {
		var p Pool
		err = json.Unmarshal([]byte(val), &p)
		assert.Nil(t, err)

		if p.Typ == common.PoolTypeAeroAMM {
			pls = append(pls, &p)
		}
	}

	logger.Info().Msgf("aero pools: %d", len(pls))

	type args struct {
		pls      []*Pool
		maxRetry uint
	}

	tt := struct {
		name    string
		args    args
		want    []PoolFeeChanged
		wantErr bool
	}{
		name: "get aerov2 pools fee",
		args: args{
			pls:      pls,
			maxRetry: 3,
		},
		wantErr: false,
	}
	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://base-rpc.publicnode.com",
		Tps: 3,
		V3:  true,
	}})

	t.Run(tt.name, func(t *testing.T) {
		got, err := pp.GetAeroPoolV2Fee(tt.args.pls, tt.args.maxRetry)
		if (err != nil) != tt.wantErr {
			t.Errorf("ProviderPool.GetAeroPoolV2Fee() error = %v, wantErr %v", err, tt.wantErr)
			return
		}

		if len(got) > 0 {
			for _, item := range got {
				logger.Info().Msgf("pool fee changed: %v %v", item.Pool, item.Fee)
			}
		}
	})
}

func TestProviderPool_GetInfusionSwapFee(t *testing.T) {
	tests := []struct {
		name    string
		factory *common.SwapFactory
	}{
		// TODO: Add test cases.
		{
			name: "infusion swap",
			factory: &common.SwapFactory{
				Name:    "Infusion Swap",
				Address: "0x2d9a3a2bd6400ee28d770c7254ca840c82faf23f",
				Typ:     common.PoolTypeInfusionAMM,
			},
		},
	}
	pp := CreateProviderPool([]config.ProviderConfig{{
		RPC: "https://base-rpc.publicnode.com",
		Tps: 3,
		V3:  true,
	}})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pp.GetInfusionSwapFee(tt.factory)
		})
	}
}
