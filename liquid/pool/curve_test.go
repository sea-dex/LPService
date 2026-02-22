package pool

import (
	"math/big"
	"reflect"
	"strings"
	"testing"

	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
)

func Test_getAeroV2AmountOut(t *testing.T) {
	type args struct {
		pl       *Pool
		amountIn *big.Int
		tokenIn  string
	}

	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			name: "stable-USDC+/USD+",
			args: args{
				pl: &Pool{
					PoolInfo: PoolInfo{
						decimals0: toBigIntMust("1000000"),
						decimals1: toBigIntMust("1000000"),
						Address:   strings.ToLower("0xE96c788E66a97Cf455f46C5b27786191fD3bC50B"),
						Typ:       common.PoolTypeAeroAMM,
						Stable:    true,
						Fee:       500,
						Token0:    strings.ToLower("0x85483696Cc9970Ad9EdD786b2C5ef735F38D156f"),
						Token1:    strings.ToLower("0xB79DD08EA68A908A97220C76d19A6aA9cBDE4376"),
					},
					Reserve0: toBigIntMust("8565573955477"),
					Reserve1: toBigIntMust("8738221410012"),
				},
				amountIn: toBigIntMust("1000000"),
				tokenIn:  strings.ToLower("0x85483696Cc9970Ad9EdD786b2C5ef735F38D156f"),
			},
			want: toBigIntMust("999501"),
		},
		{
			name: "stable-DOLA/USD+",
			args: args{
				pl: &Pool{
					PoolInfo: PoolInfo{
						decimals0: e18,
						decimals1: toBigIntMust("1000000"),
						Address:   strings.ToLower("0x8E9154AC849e839d60299E85156bcb589De2693A"),
						Typ:       common.PoolTypeAeroAMM,
						Stable:    true,
						Fee:       3000,
						Token0:    strings.ToLower("0x4621b7A9c75199271F773Ebd9A499dbd165c3191"),
						Token1:    strings.ToLower("0xB79DD08EA68A908A97220C76d19A6aA9cBDE4376"),
					},
					Reserve0: toBigIntMust("4536284569644280666488402"),
					Reserve1: toBigIntMust("3341772539654"),
				},
				amountIn: toBigIntMust("1000000"),
				tokenIn:  strings.ToLower("0xB79DD08EA68A908A97220C76d19A6aA9cBDE4376"),
			},
			want: toBigIntMust("1003975170347218274"),
		},
		{
			name: "vmm-WETH/USDC",
			args: args{
				pl: &Pool{
					PoolInfo: PoolInfo{
						decimals0: e18,
						decimals1: toBigIntMust("1000000"),
						Address:   strings.ToLower("0xcDAC0d6c6C59727a65F871236188350531885C43"),
						Typ:       common.PoolTypeAeroAMM,
						Stable:    false,
						Fee:       3000,
						Token0:    strings.ToLower("0x4200000000000000000000000000000000000006"),
						Token1:    strings.ToLower("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913 "),
					},
					Reserve0: toBigIntMust("14004027074565688150700"),
					Reserve1: toBigIntMust("33648446844282"),
				},
				amountIn: toBigIntMust("1000000000000000"), // 10**15
				tokenIn:  strings.ToLower("0x4200000000000000000000000000000000000006"),
			},
			want: toBigIntMust("2395560"),
		},
	}

	logger.Init("", "dev", true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAeroV2AmountOut(tt.args.pl, tt.args.amountIn, tt.args.tokenIn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAeroV2AmountOut() = %v, want %v", got, tt.want)
			}
		})
	}
}
