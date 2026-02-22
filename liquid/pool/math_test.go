package pool

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var _e18 = big.NewInt(10).Exp(big.NewInt(10), big.NewInt(18), nil)

func expandTo18Decimals(v int64) *big.Int {
	b := big.NewInt(v)
	b.Mul(b, _e18)
	// println(b.String()) // 1000000000000000000
	return b
}

func TestGetSqrtRatioAtTick(t *testing.T) {
	var v *big.Int

	v = getSqrtRatioAtTick(MIN_TICK)
	assert.Equal(t, v, big.NewInt(4295128739))

	v = getSqrtRatioAtTick(MIN_TICK + 1)
	assert.Equal(t, v, big.NewInt(4295343490))

	v = getSqrtRatioAtTick(MAX_TICK - 1)
	assert.Equal(t, v, toBigIntMust("1461373636630004318706518188784493106690254656249"))

	v = getSqrtRatioAtTick(MIN_TICK)
	// println("v:", v.String())
	// println(decimal.NewFromInt(1).Div(decimal.NewFromInt(2).Pow(decimal.NewFromInt(127))).String())
	_1_rsh127 := encodePriceSqrt(bigOne, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(127), nil))
	// println(_1_rsh127.String())
	assert.True(t, v.Cmp(_1_rsh127) < 0)

	v = getSqrtRatioAtTick(MAX_TICK)
	assert.True(t, v.Cmp(encodePriceSqrt(big.NewInt(1).Exp(big.NewInt(2), big.NewInt(127), nil), bigOne)) > 0)
	assert.Equal(t, v, toBigIntMust("1461446703485210103287273052203988822378723970342"))

	ticks := []int{50, 100, 250, 500, 1000, 2500, 3000, 4000, 5000, 50000, 150000, 250000, 500000, 738203}
	negSqrtPrice := []*big.Int{
		toBigIntMust("79030349367926598376800521322"),
		toBigIntMust("78833030112140176575862854579"),
		toBigIntMust("78244023372248365697264290337"),
		toBigIntMust("77272108795590369356373805297"),
		toBigIntMust("75364347830767020784054125655"),
		toBigIntMust("69919044979842180277688105136"),
		toBigIntMust("68192822843687888778582228483"),
		toBigIntMust("64867181785621769311890333195"),
		toBigIntMust("61703726247759831737814779831"),
		toBigIntMust("6504256538020985011912221507"),
		toBigIntMust("43836292794701720435367485"),
		toBigIntMust("295440463448801648376846"),
		toBigIntMust("1101692437043807371"),
		toBigIntMust("7409801140451"),
	}
	posSqrtPrice := []*big.Int{
		toBigIntMust("79426470787362580746886972461"),
		toBigIntMust("79625275426524748796330556128"),
		toBigIntMust("80224679980005306637834519095"),
		toBigIntMust("81233731461783161732293370115"),
		toBigIntMust("83290069058676223003182343270"),
		toBigIntMust("89776708723587163891445672585"),
		toBigIntMust("92049301871182272007977902845"),
		toBigIntMust("96768528593268422080558758223"),
		toBigIntMust("101729702841318637793976746270"),
		toBigIntMust("965075977353221155028623082916"),
		toBigIntMust("143194173941309278083010301478497"),
		toBigIntMust("21246587762933397357449903968194344"),
		toBigIntMust("5697689776495288729098254600827762987878"),
		toBigIntMust("847134979253254120489401328389043031315994541"),
	}

	for i, tick := range ticks {
		negTick := -tick
		v = getSqrtRatioAtTick(negTick)
		assert.Equal(t, negSqrtPrice[i], v)

		v = getSqrtRatioAtTick(tick)
		assert.Equal(t, posSqrtPrice[i], v)
	}
}

func TestDivRoundingUp(t *testing.T) {
	divEqual(t, "5", "2")
	divEqual(t, "3", "3")
	divEqual(t, "3", "13")
	divEqual(t, "16", "2")
	divEqual(t, "5697689776495288729098254600827762987878", "7409801140451")
	divEqual(t, "847134979253254120489401328389043031315994541", "1101692437043807371")
	divEqual(t, "21246587762933397357449903968194344", "7409801140451")
	divEqual(t, "143194173941309278083010301478497", "1101692437043807371")
}

func encodePriceSqrt(reserve1, reserve0 *big.Int) *big.Int {
	d := decimal.NewFromBigInt(reserve1, 0)
	// println("reserve1:", d.String())
	d = d.DivRound(decimal.NewFromBigInt(reserve0, 0), 40)
	if d.Cmp(decimal.NewFromInt(0)) == 0 {
		// println("quotient too small to 0"
		d = decimal.NewFromBigInt(reserve1, 0)
		d = d.DivRound(decimal.NewFromBigInt(reserve0, 0), 99)
	}
	// println("reserve1/reserv0:", d.String())
	// println("sqrt reserve1/reserv0:", d.Pow(decimal.NewFromFloat(0.5)).String())
	d = d.Pow(decimal.NewFromFloat(0.5)).Mul(decimal.NewFromInt(2).Pow(decimal.NewFromInt(96)))
	// println("sqrt d:", d.String())

	return d.BigInt()
}

func divEqual(t *testing.T, x, y string) {
	v1 := divRoundingUp(toBigIntMust(x), toBigIntMust(y))
	dx := decimal.RequireFromString(x)
	dy := decimal.RequireFromString(y)
	dz := dx.Div(dy)

	if dx.Mod(dy).Cmp(decimal.Zero) != 0 {
		dz = dz.Add(decimal.NewFromInt(1))
	}

	assert.Equal(t, v1, dz.BigInt())
}

func TestMulDiv(t *testing.T) {
	q128 := big.NewInt(1).Exp(big.NewInt(2), big.NewInt(128), nil)

	_50 := big.NewInt(50)
	_50.Mul(_50, q128)
	_50.Div(_50, big.NewInt(100))

	_150 := big.NewInt(150)
	_150.Mul(_150, q128)
	_150.Div(_150, big.NewInt(100))
	assert.Equal(t, big.NewInt(1).Div(q128, big.NewInt(3)), mulDiv(q128, _50, _150))

	assert.Equal(t,
		decimal.NewFromInt(4375).Mul(decimal.NewFromBigInt(q128, 0)).Div(decimal.NewFromInt(1000)).BigInt(),
		mulDiv(q128, big.NewInt(1).Mul(big.NewInt(35), q128), big.NewInt(1).Mul(big.NewInt(8), q128)))

	assert.Equal(t,
		big.NewInt(1).Div(q128, big.NewInt(3)),
		mulDiv(q128, big.NewInt(1).Mul(big.NewInt(1000), q128), big.NewInt(1).Mul(big.NewInt(3000), q128)))

	assert.Panics(t, func() { mulDiv(uint256Max, uint256Max, big.NewInt(1).Sub(uint256Max, bigOne)) })
	assert.Equal(t, uint256Max, mulDiv(uint256Max, uint256Max, uint256Max))
}

func TestMulDivRoundingUp(t *testing.T) {
	q128 := big.NewInt(1).Exp(big.NewInt(2), big.NewInt(128), nil)

	assert.Equal(t, uint256Max, mulDivRoundingUp(uint256Max, uint256Max, uint256Max))

	assert.Equal(t,
		big.NewInt(1).Add(bigOne, big.NewInt(1).Div(q128, big.NewInt(3))),
		mulDivRoundingUp(q128, big.NewInt(1).Mul(big.NewInt(1000), q128), big.NewInt(1).Mul(big.NewInt(3000), q128)))
}

func TestGetAmount0Delta(t *testing.T) {
	// returns 0 if liquidity is 0
	// t.Log("encodePriceSqrt 1 1:", encodePriceSqrt(bigOne, bigOne))
	// t.Log("encodePriceSqrt 2 1:", encodePriceSqrt(big.NewInt(2), bigOne))
	// t.Logf("bigOne: %v", bigOne.String())
	amount0 := GetAmount0Delta(encodePriceSqrt(bigOne, bigOne),
		encodePriceSqrt(big.NewInt(2), bigOne),
		big.NewInt(0),
		true)
	assert.Equal(t, bigZero, amount0)

	// returns 0 if prices are equal
	amount0 = GetAmount0Delta(encodePriceSqrt(bigOne, bigOne),
		encodePriceSqrt(bigOne, bigOne),
		big.NewInt(0),
		true)
	assert.Equal(t, bigZero, amount0)

	// returns 0.1 amount1 for price of 1 to 1.21
	amount0 = GetAmount0Delta(encodePriceSqrt(bigOne, bigOne),
		encodePriceSqrt(big.NewInt(121), big.NewInt(100)),
		expandTo18Decimals(1),
		true)
	assert.Equal(t, toBigIntMust("90909090909090910"), amount0)

	amount0RoundedDown := GetAmount0Delta(
		encodePriceSqrt(bigOne, bigOne),
		encodePriceSqrt(big.NewInt(121), big.NewInt(100)),
		expandTo18Decimals(1),
		false,
	)
	assert.Equal(t, big.NewInt(1).Sub(amount0, bigOne), amount0RoundedDown)

	// works for prices that overflow
	amount0Up := GetAmount0Delta(
		encodePriceSqrt(big.NewInt(0).Lsh(big.NewInt(2), 90), bigOne),
		encodePriceSqrt(big.NewInt(0).Lsh(big.NewInt(2), 96), bigOne),
		expandTo18Decimals(1),
		true,
	)
	amount0Down := GetAmount0Delta(
		encodePriceSqrt(big.NewInt(0).Lsh(big.NewInt(2), 90), bigOne),
		encodePriceSqrt(big.NewInt(0).Lsh(big.NewInt(2), 96), bigOne),
		expandTo18Decimals(1),
		false,
	)
	assert.Equal(t, amount0Up, big.NewInt(1).Add(amount0Down, bigOne))
}

func TestGetAmount1Delta(t *testing.T) {
	// returns 0 if liquidity is 0
	amount1 := GetAmount1Delta(encodePriceSqrt(bigOne, bigOne),
		encodePriceSqrt(big.NewInt(2), bigOne), bigZero, true)

	assert.True(t, bigZero.Cmp(amount1) == 0)

	// returns 0 if prices are equal
	amount1 = GetAmount1Delta(encodePriceSqrt(bigOne, bigOne),
		encodePriceSqrt(bigOne, bigOne), bigZero, true)

	assert.True(t, bigZero.Cmp(amount1) == 0)

	// returns 0.1 amount1 for price of 1 to 1.21
	amount1 = GetAmount1Delta(
		encodePriceSqrt(bigOne, bigOne),
		encodePriceSqrt(big.NewInt(121), big.NewInt(100)),
		expandTo18Decimals(1),
		true,
	)

	assert.True(t, amount1.Cmp(toBigIntMust("100000000000000000")) == 0)

	amount1RoundedDown := GetAmount1Delta(
		encodePriceSqrt(bigOne, bigOne),
		encodePriceSqrt(big.NewInt(121), big.NewInt(100)),
		expandTo18Decimals(1),
		false,
	)
	assert.Equal(t, amount1RoundedDown, big.NewInt(1).Sub(amount1, bigOne))
}

func TestGetTickAtSqrtRatio(t *testing.T) {
	tick := getTickAtSqrtRatio(MIN_SQRT_RATIO)
	assert.Equal(t, MIN_TICK, tick)

	tick = getTickAtSqrtRatio(toBigIntMust("4295343490"))
	assert.Equal(t, MIN_TICK+1, tick)

	tick = getTickAtSqrtRatio(toBigIntMust("1461373636630004318706518188784493106690254656249"))
	assert.Equal(t, MAX_TICK-1, tick)

	tick = getTickAtSqrtRatio(toBigIntMust("1461446703485210103287273052203988822378723970341"))
	assert.Equal(t, MAX_TICK-1, tick)

	prices := []*big.Int{
		MIN_SQRT_RATIO,
		toBigIntMust("79228162514264337593543"),
		toBigIntMust("79228162514264337593543950"),
		toBigIntMust("9903520314283042199192993792"),
		toBigIntMust("28011385487393069959365969113"),
		toBigIntMust("56022770974786139918731938227"),
		toBigIntMust("79228162514264337593543950336"),
		toBigIntMust("112045541949572279837463876454"),
		toBigIntMust("224091083899144559674927752909"),
		toBigIntMust("633825300114114700748351602688"),
		toBigIntMust("79228162514264337593543950336000"),
		toBigIntMust("79228162514264337593543950336000000"),
		toBigIntMust("1461446703485210103287273052203988822378723970341"),
	}
	ticks := []int{
		-887272,
		-276325,
		-138163,
		-41591,
		-20796,
		-6932,
		0,
		6931,
		20795,
		41590,
		138162,
		276324,
		887271,
	}

	for i, price := range prices {
		tick = getTickAtSqrtRatio(price)

		assert.Equal(t, ticks[i], tick)
	}
}

func TestComputeSwapStep(t *testing.T) {
	sqrtRatioCurrentX96 := encodePriceSqrt(big.NewInt(1), big.NewInt(1))
	assert.Equal(t, sqrtRatioCurrentX96, toBigIntMust("79228162514264337593543950336"))

	liquidity := big.NewInt(2083487)
	sqrtQ, amtIn, amtOut, feeAmt := computeSwapStep(sqrtRatioCurrentX96, big.NewInt(1), liquidity, big.NewInt(10), 3000)
	t.Logf("sqrtQ: %v", sqrtQ.String())   // 79227820275324292410333451302
	t.Logf("amtIn: %v", amtIn.String())   // 9
	t.Logf("amtOut: %v", amtOut.String()) // 8
	t.Logf("feeAmt: %v", feeAmt.String()) // 1
}
