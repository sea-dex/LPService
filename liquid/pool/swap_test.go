package pool

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// go test -timeout 600s -run ^TestComputeSwapStep periphery/liquid/handlers -count 1 -v.
func TestComputeSwapStep1(t *testing.T) {
	// exact amount in that gets capped at price target in one for zero
	testComputeSwapStep(t,
		encodePriceSqrt(bigOne, bigOne),
		encodePriceSqrt(big.NewInt(101), big.NewInt(100)),
		expandTo18Decimals(2),
		expandTo18Decimals(1),
		600,
		false,
		big.NewInt(9975124224178055),
		toBigIntMust("9925619580021728"),
		toBigIntMust("5988667735148"),
	)
}

func testComputeSwapStep(t *testing.T,
	price, priceTarget, liquidity, amount *big.Int,
	fee int,
	zeroForOne bool,
	expAmountIn, expAmountOut, expFeeAmount *big.Int,
) {
	sqrtQ, amountIn, amountOut, feeAmount := computeSwapStep(price, priceTarget, liquidity, amount, fee)
	assert.Equal(t, expAmountIn, amountIn)
	assert.Equal(t, expAmountOut, amountOut)
	assert.Equal(t, expFeeAmount, feeAmount)
	assert.Equal(t, priceTarget, sqrtQ)

	priceAfterWholeInputAmount := getNextSqrtPriceFromInput(price, liquidity, amount, zeroForOne)
	assert.True(t, sqrtQ.Cmp(priceAfterWholeInputAmount) < 0)
}

func TestComputeSwapStep2(t *testing.T) {
	// exact amount out that gets capped at price target in one for zero
	price := encodePriceSqrt(bigOne, bigOne)
	priceTarget := encodePriceSqrt(big.NewInt(101), big.NewInt(100))
	liquidity := expandTo18Decimals(2)
	amount := expandTo18Decimals(-1)
	fee := 600
	zeroForOne := false
	expAmountIn := big.NewInt(9975124224178055)
	expAmountOut := toBigIntMust("9925619580021728")
	expFeeAmount := toBigIntMust("5988667735148")

	sqrtQ, amountIn, amountOut, feeAmount := computeSwapStep(price, priceTarget, liquidity, amount, fee)
	assert.Equal(t, expAmountIn, amountIn)
	assert.Equal(t, expAmountOut, amountOut)
	assert.Equal(t, expFeeAmount, feeAmount)
	assert.Equal(t, priceTarget, sqrtQ)
	assert.True(t, amountOut.Cmp(big.NewInt(1).Neg(amount)) < 0)

	priceAfterWholeOutputAmount := getNextSqrtPriceFromOutput(
		price,
		liquidity,
		big.NewInt(1).Neg(amount),
		zeroForOne,
	)
	assert.True(t, sqrtQ.Cmp(priceAfterWholeOutputAmount) < 0)
}

func TestComputeSwapStep3(t *testing.T) {
	// exact amount in that is fully spent in one for zero
	price := encodePriceSqrt(bigOne, bigOne)
	priceTarget := encodePriceSqrt(big.NewInt(1000), big.NewInt(100))
	liquidity := expandTo18Decimals(2)
	amount := expandTo18Decimals(1)
	fee := 600
	zeroForOne := false
	expAmountIn := toBigIntMust("999400000000000000")
	expAmountOut := toBigIntMust("666399946655997866")
	expFeeAmount := toBigIntMust("600000000000000")

	sqrtQ, amountIn, amountOut, feeAmount := computeSwapStep(price, priceTarget, liquidity, amount, fee)
	assert.Equal(t, expAmountIn, amountIn)
	assert.Equal(t, expAmountOut, amountOut)
	assert.Equal(t, expFeeAmount, feeAmount)
	assert.True(t, sqrtQ.Cmp(priceTarget) < 0)
	assert.Equal(t, amount, big.NewInt(1).Add(amountIn, feeAmount))

	priceAfterWholeInputAmountLessFee := getNextSqrtPriceFromInput(
		price,
		liquidity,
		big.NewInt(1).Sub(amount, feeAmount),
		zeroForOne,
	)
	assert.True(t, sqrtQ.Cmp(priceAfterWholeInputAmountLessFee) == 0)
}

func TestComputeSwapStep4(t *testing.T) {
	// exact amount out that is fully received in one for zero
	price := encodePriceSqrt(bigOne, bigOne)
	priceTarget := toBigIntMust("792281625142643375935439503360") // encodePriceSqrt(big.NewInt(10000), big.NewInt(100))
	liquidity := expandTo18Decimals(2)
	amount := expandTo18Decimals(-1)
	fee := 600
	zeroForOne := false
	expAmountIn := toBigIntMust("2000000000000000000")
	expAmountOut := big.NewInt(1).Neg(amount)
	expFeeAmount := toBigIntMust("1200720432259356")

	// println("price:", price.String())
	sqrtQ, amountIn, amountOut, feeAmount := computeSwapStep(price, priceTarget, liquidity, amount, fee)
	assert.Equal(t, expAmountIn, amountIn)
	assert.Equal(t, expAmountOut, amountOut)
	assert.Equal(t, expFeeAmount, feeAmount)
	assert.True(t, sqrtQ.Cmp(priceTarget) < 0)

	priceAfterWholeOutputAmount := getNextSqrtPriceFromOutput(
		price,
		liquidity,
		big.NewInt(1).Neg(amount),
		zeroForOne,
	)
	// println("sqrtQ:", sqrtQ.String())
	assert.True(t, sqrtQ.Cmp(priceAfterWholeOutputAmount) == 0)
}

func TestDecimalPow(t *testing.T) {
	d := decimal.NewFromInt(100)
	sqrt := d.Pow(decimal.NewFromFloat(0.5))
	t.Logf("sqrt: %v", sqrt)

	d = decimal.NewFromBigInt(big.NewInt(100), 0)
	sqrt = d.Pow(decimal.NewFromFloat(0.5))
	t.Logf("sqrt: %v", sqrt)

	d = decimal.NewFromBigInt(big.NewInt(10000), 0)
	d = d.Div(decimal.NewFromBigInt(big.NewInt(100), 0))
	sqrt = d.Pow(decimal.NewFromFloat(0.5))
	t.Logf("sqrt: %v", sqrt)
}

func TestComputeSwapStep5(t *testing.T) {
	// amount out is capped at the desired amount out
	price := toBigIntMust("417332158212080721273783715441582")
	priceTarget := toBigIntMust("1452870262520218020823638996")
	liquidity := toBigIntMust("159344665391607089467575320103")
	amount := big.NewInt(-1)
	fee := 1
	// zeroForOne := false
	expAmountIn := toBigIntMust("1")
	expAmountOut := big.NewInt(1)
	expFeeAmount := toBigIntMust("1")

	// println("price:", price.String())
	sqrtQ, amountIn, amountOut, feeAmount := computeSwapStep(price, priceTarget, liquidity, amount, fee)
	assert.Equal(t, expAmountIn, amountIn)
	assert.Equal(t, expAmountOut, amountOut)
	assert.Equal(t, expFeeAmount, feeAmount)
	// assert.True(t, sqrtQ.Cmp(priceTarget) < 0)
	// priceAfterWholeOutputAmount := getNextSqrtPriceFromOutput(
	// 	price,
	// 	liquidity,
	// 	big.NewInt(1).Neg(amount),
	// 	zeroForOne,
	// )
	// // println("sqrtQ:", sqrtQ.String())
	// assert.True(t, sqrtQ.Cmp(priceAfterWholeOutputAmount) == 0)
	assert.Equal(t, sqrtQ, toBigIntMust("417332158212080721273783715441581"))
}

func TestComputeSwapStep6(t *testing.T) {
	// target price of 1 uses partial input amount
	price := toBigIntMust("2")
	priceTarget := toBigIntMust("1")
	liquidity := toBigIntMust("1")
	amount := toBigIntMust("3915081100057732413702495386755767")
	fee := 1
	// zeroForOne := false
	expAmountIn := toBigIntMust("39614081257132168796771975168")
	expAmountOut := big.NewInt(0)
	expFeeAmount := toBigIntMust("39614120871253040049813")

	// println("price:", price.String())
	sqrtQ, amountIn, amountOut, feeAmount := computeSwapStep(price, priceTarget, liquidity, amount, fee)
	assert.Equal(t, expAmountIn, amountIn)
	assert.True(t, expAmountOut.Cmp(amountOut) == 0)
	assert.Equal(t, expFeeAmount, feeAmount)
	// assert.True(t, sqrtQ.Cmp(priceTarget) < 0)
	// priceAfterWholeOutputAmount := getNextSqrtPriceFromOutput(
	// 	price,
	// 	liquidity,
	// 	big.NewInt(1).Neg(amount),
	// 	zeroForOne,
	// )
	// // println("sqrtQ:", sqrtQ.String())
	// assert.True(t, sqrtQ.Cmp(priceAfterWholeOutputAmount) == 0)
	assert.Equal(t, sqrtQ, toBigIntMust("1"))
}

func TestComputeSwapStep7(t *testing.T) {
	// entire input amount taken as fee
	price := toBigIntMust("2413")
	priceTarget := toBigIntMust("79887613182836312")
	liquidity := toBigIntMust("1985041575832132834610021537970")
	amount := toBigIntMust("10")
	fee := 1872
	// zeroForOne := false
	expAmountIn := toBigIntMust("0")
	expAmountOut := big.NewInt(0)
	expFeeAmount := toBigIntMust("10")

	// println("price:", price.String())
	sqrtQ, amountIn, amountOut, feeAmount := computeSwapStep(price, priceTarget, liquidity, amount, fee)
	assert.True(t, expAmountIn.Cmp(amountIn) == 0)
	assert.True(t, expAmountOut.Cmp(amountOut) == 0)
	assert.Equal(t, expFeeAmount, feeAmount)
	assert.Equal(t, sqrtQ, toBigIntMust("2413"))
}

func TestComputeSwapStep8(t *testing.T) {
	// handles intermediate insufficient liquidity in zero for one exact output case
	price := toBigIntMust("20282409603651670423947251286016")
	priceTarget := big.NewInt(1).Mul(price, big.NewInt(11))
	priceTarget.Div(priceTarget, big.NewInt(10))

	liquidity := toBigIntMust("1024")
	amount := toBigIntMust("-4")
	fee := 3000
	// zeroForOne := false
	expAmountIn := toBigIntMust("26215")
	expAmountOut := big.NewInt(0)
	expFeeAmount := toBigIntMust("79")

	// println("price:", price.String())
	sqrtQ, amountIn, amountOut, feeAmount := computeSwapStep(price, priceTarget, liquidity, amount, fee)
	assert.True(t, expAmountIn.Cmp(amountIn) == 0)
	assert.True(t, expAmountOut.Cmp(amountOut) == 0)
	assert.Equal(t, expFeeAmount, feeAmount)
	assert.Equal(t, sqrtQ, priceTarget)
}

func TestComputeSwapStep9(t *testing.T) {
	// handles intermediate insufficient liquidity in zero for one exact output case
	price := toBigIntMust("20282409603651670423947251286016")
	priceTarget := big.NewInt(1).Mul(price, big.NewInt(9))
	priceTarget.Div(priceTarget, big.NewInt(10))

	liquidity := toBigIntMust("1024")
	amount := toBigIntMust("-263000")
	fee := 3000
	// zeroForOne := false
	expAmountIn := toBigIntMust("1")
	expAmountOut := big.NewInt(26214)
	expFeeAmount := toBigIntMust("1")

	// println("price:", price.String())
	sqrtQ, amountIn, amountOut, feeAmount := computeSwapStep(price, priceTarget, liquidity, amount, fee)
	assert.True(t, expAmountIn.Cmp(amountIn) == 0)
	assert.True(t, expAmountOut.Cmp(amountOut) == 0)
	assert.Equal(t, expFeeAmount, feeAmount)
	assert.Equal(t, sqrtQ, priceTarget)
}

// https://dashboard.tenderly.co/tx/base/0x31b55eaacb1110b8bba99c06c1020067fb0ae433d69f92aed2dd09ea52f482a6/debugger?trace=0.7.1.9
func TestComputeSwapStep10(t *testing.T) {
	price := toBigIntMust("670899113489198717055724758074698") // price current
	priceTarget := toBigIntMust("171433008539253222554574933565267")

	liquidity := toBigIntMust("7494463238316011082453")
	amount := toBigIntMust("-500000000000000000000000")
	fee := 10000
	// zeroForOne := false
	expAmountIn := toBigIntMust("7028294329069908")
	expAmountOut := toBigIntMust("500000000000000000000000")
	expFeeAmount := toBigIntMust("70992872010808")
	expSqrtQ := toBigIntMust("665613333841284483888513302454088")

	// println("price:", price.String())
	sqrtQ, amountIn, amountOut, feeAmount := computeSwapStep(price, priceTarget, liquidity, amount, fee)
	assert.True(t, expAmountIn.Cmp(amountIn) == 0)
	assert.True(t, expAmountOut.Cmp(amountOut) == 0)
	assert.Equal(t, expFeeAmount, feeAmount)
	assert.Equal(t, sqrtQ, expSqrtQ)
}

// https://dashboard.tenderly.co/tx/base/0x0d5fe03abfccd9b2e195f80241deae66bdbdd3158d81ff44c04c6ad1a0057609/debugger?trace=0.0.1.1.0.0.11
func TestComputeSwapStep11(t *testing.T) {
	price := toBigIntMust("79422499761647659469960370034") // price current
	priceTarget := toBigIntMust("79414558305817077762184151619")
	// sqrtPriceLimitX96 := toBigIntMust("4295128740")
	// sqrtPriceNextX96 := toBigIntMust("79414558305817077762184151619")

	liquidity := toBigIntMust("997652876304682704914221258")
	amount := toBigIntMust("10000000000000000000")
	fee := 3000
	// zeroForOne := false
	expAmountIn := toBigIntMust("9970000000000000000")
	expAmountOut := toBigIntMust("10018970330728980107")
	expFeeAmount := toBigIntMust("30000000000000000")
	expSqrtQ := toBigIntMust("79422498965995555976032368296")

	// println("price:", price.String())
	sqrtQ, amountIn, amountOut, feeAmount := computeSwapStep(price, priceTarget, liquidity, amount, fee)
	assert.True(t, expAmountIn.Cmp(amountIn) == 0)
	assert.True(t, expAmountOut.Cmp(amountOut) == 0)
	assert.Equal(t, expFeeAmount, feeAmount)
	assert.Equal(t, sqrtQ, expSqrtQ)
}
