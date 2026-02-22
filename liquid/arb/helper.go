package arb

import (
	"math/big"

	"github.com/shopspring/decimal"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

const (
	WETHAddress        = "0x4200000000000000000000000000000000000006"
	NativeTokenAddress = "0x0000000000000000000000000000000000000000"
)

type PoolTVLSlice []*pool.Pool

type PoolPriceSlice []*pool.Pool

var (
	d_10_exp_18 = decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))
	d_one       = decimal.NewFromInt(1)
	d_two       = decimal.NewFromInt(2)
	d_00001     = decimal.NewFromFloat(0.0001)
	//nolint
	d_0001 = decimal.NewFromFloat(0.001)
)

func (ps PoolTVLSlice) Len() int {
	return len(ps)
}

func (ps PoolTVLSlice) Less(i, j int) bool {
	return ps[i].GetTVL().Cmp(ps[j].GetTVL()) < 0
}

func (ps PoolTVLSlice) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps PoolPriceSlice) Len() int {
	return len(ps)
}

func (ps PoolPriceSlice) Less(i, j int) bool {
	return ps[i].GetPrice0().Cmp(ps[j].GetPrice0()) < 0
}

func (ps PoolPriceSlice) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func isNativeToken(token string) bool {
	return token == NativeTokenAddress
}

func isWrapperNativeToken(token string) bool {
	return token == WETHAddress
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func isNativeTokens(token0, token1 string) bool {
	return token0 == NativeTokenAddress || token1 == NativeTokenAddress
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func isWrapperNativeTokens(token0, token1 string) bool {
	return token0 == WETHAddress || token1 == NativeTokenAddress
}

func isNativeOrWrapperNativeToken(token string) bool {
	return isNativeToken(token) || isWrapperNativeToken(token)
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func isNativeOrWrapperNativeTokens(token0, token1 string) bool {
	return isNativeOrWrapperNativeToken(token0) || isNativeOrWrapperNativeToken(token1)
}

func isStableToken(token string, stables map[string]*common.Token) bool {
	_, ok := stables[token]
	return ok
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func isStableTokens(token0, token1 string, stables map[string]*common.Token) bool {
	_, ok := stables[token0]
	if ok {
		return true
	}

	_, ok = stables[token1]

	return ok
}

// float to big.Int: v*(10**18).
func convertToBigInt(v, def, name string) *big.Int {
	d, err := decimal.NewFromString(v)
	if err != nil {
		logger.Warn().Err(err).Str("value", v).Msgf("invalid %s: %v, set to default %s", name, v, def)
		d = decimal.RequireFromString(def)
	}

	return d.Mul(d_10_exp_18).BigInt()
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func readableETHAmount(amt *big.Int) string {
	return decimal.NewFromBigInt(amt, 0).Div(d_10_exp_18).StringFixed(3)
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func toETHAmountWithPrecision(amt *big.Int, token *common.Token, price decimal.Decimal, reverse bool) *big.Int {
	precision := decimal.NewFromInt(10).Pow(decimal.NewFromInt(18 - int64(token.Decimals))) // nolint
	if reverse {
		return decimal.NewFromBigInt(amt, 0).Mul(precision).Div(price).BigInt()
	}

	return decimal.NewFromBigInt(amt, 0).Mul(precision).Mul(price).BigInt()
}

func toETHAmount(amt *big.Int, price decimal.Decimal, reverse bool) *big.Int {
	if reverse {
		return decimal.NewFromBigInt(amt, 0).Div(price).BigInt()
	}

	return decimal.NewFromBigInt(amt, 0).Mul(price).BigInt()
}

func setDecimalPrecision(precision int) {
	decimal.DivisionPrecision = precision
}
