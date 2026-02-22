package arb

import (
	"math"
	"math/big"

	"github.com/shopspring/decimal"
	"starbase.ag/liquidity/liquid/pool"
)

var (
	base1_0001    = float64(1.0001)
	logBase1_0001 = math.Log(base1_0001)
	d0_5          = decimal.NewFromFloat(0.5)
)

func mulmul(args ...*big.Int) *big.Int {
	a := big.NewInt(1)
	for _, arg := range args {
		a.Mul(a, arg)
	}

	return a
}

func div(x *big.Int, dividors ...*big.Int) *big.Int {
	y := new(big.Int).Set(x)
	for _, dividor := range dividors {
		y.Div(y, dividor)
	}

	return y
}

func TickAtPrice(price float64) int32 {
	return int32(math.Log(price) / logBase1_0001)
}

func TickNorm(tick int32, tickSpacing int) int32 {
	return tick / int32(tickSpacing) * int32(tickSpacing) // nolint
}

func PriceAtTick(tick int32) float64 {
	return math.Exp(float64(tick) * logBase1_0001)
}

func SqrtPriceX96AtTick(tick int32) *big.Int {
	price := PriceAtTick(tick)
	sqrtPrice := decimal.NewFromFloat(price).Pow(d0_5)

	return sqrtPrice.Mul(q96).BigInt()
}

func SqrtPriceX96ToPrice(sp *big.Int) float64 {
	p := decimal.NewFromBigInt(sp, 0).Div(q96)
	f, _ := p.Mul(p).Float64()
	return f
}

// r = e6 - fee
// dy = dy*r*ry/(rx+r*dy)
// uint amountInWithFee = amountIn * 997;
// uint numerator = amountInWithFee * reserveOut;
// uint denominator = reserveIn * 1000 + amountInWithFee;
// amountOut = numerator / denominator;
func getAmountOutV2(rx, ry, dx, r *big.Int) *big.Int {
	amountInWithFee := mulmul(dx, r)
	numerator := mulmul(amountInWithFee, ry)
	denominator := new(big.Int).Add(mulmul(rx, e6), amountInWithFee)

	return div(numerator, denominator)
}

// 0->1: dy2 = dx1*f2*L2*SP2*SP2 / (Q192*L2 + dx1*Q96*f2*SP2)
// 1->0: dx2 = dy2*f2*L2*Q192 / (Sp2*Sp2*L + dy2*Q96*f2*Sp2)
func getAmountOutV3Tick(zeroForOne bool, l, sp, amtIn *big.Int, r *big.Int) *big.Int {
	if zeroForOne {
		return div(mulmul(amtIn, r, l, sp, sp), new(big.Int).Add(mulmul(Q192, l, e6), mulmul(amtIn, Q96, r, sp)))
	} else {
		return div(mulmul(amtIn, r, l, Q192), new(big.Int).Add(mulmul(sp, sp, l, e6), mulmul(amtIn, Q96, r, sp)))
	}
}

func getAmountOutV3(pl *pool.Pool, amtIn *big.Int, zeroForOne bool) *big.Int {
	amt0, amt1, _, _ := pl.MockSwap(zeroForOne, amtIn, nil)
	if zeroForOne {
		return new(big.Int).Neg(amt1)
	} else {
		return new(big.Int).Neg(amt0)
	}
}

// a+b-c
func bAddSub(a, b, c *big.Int) *big.Int {
	return new(big.Int).Sub(new(big.Int).Add(a, b), c)
}
