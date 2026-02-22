package pool

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMostSignificantBit(t *testing.T) {
	assert.Equal(t, uint8(0), mostSignificantBit(big.NewInt(1)))
	assert.Equal(t, uint8(1), mostSignificantBit(big.NewInt(2)))

	for i := 0; i < 255; i++ {
		assert.Equal(t, uint8(i), mostSignificantBit(big.NewInt(0).Lsh(big.NewInt(1), uint(i)))) // nolint
	}

	assert.Equal(t, uint8(255), mostSignificantBit(big.NewInt(0).Set(uint256Max))) // nolint
}

func TestLeastSignificantBit(t *testing.T) {
	assert.Equal(t, uint8(0), leastSignificantBit(big.NewInt(1)))
	assert.Equal(t, uint8(1), leastSignificantBit(big.NewInt(2)))

	for i := 0; i < 255; i++ {
		assert.Equal(t, uint8(i), leastSignificantBit(big.NewInt(0).Lsh(big.NewInt(1), uint(i)))) // nolint
	}

	assert.Equal(t, uint8(0), leastSignificantBit(big.NewInt(0).Set(uint256Max)))
}

func initTicks(ticks []int) map[int16]*big.Int {
	tickBitmap := map[int16]*big.Int{}
	for _, tick := range ticks {
		flipTickBitmap(tickBitmap, tick, 1)
	}

	return tickBitmap
}

func isInitialized(bm map[int16]*big.Int, tick int) bool {
	next, initialized := nextInitializedTickWithinOneWord(bm, tick, 1, true)
	if next == tick {
		return initialized
	}
	// println("not found tick:", tick)
	return false
}

func TestIsInitialized(t *testing.T) {
	tickBitmap := initTicks([]int{})

	assert.Equal(t, false, isInitialized(tickBitmap, 1))

	// is flipped by #flipTick
	flipTickBitmap(tickBitmap, 1, 1)
	assert.Equal(t, true, isInitialized(tickBitmap, 1), "is flipped by #flipTick")

	//
	flipTickBitmap(tickBitmap, 1, 1)
	assert.Equal(t, false, isInitialized(tickBitmap, 1), "is flipped back by #flipTick")
	//
	flipTickBitmap(tickBitmap, 2, 1)
	assert.Equal(t, false, isInitialized(tickBitmap, 1), "is not changed by another flip to a different tick")
	assert.Equal(t, true, isInitialized(tickBitmap, 2), "is not changed by another flip to a different tick")

	flipTickBitmap(tickBitmap, 1+256, 1)
	assert.Equal(t, false, isInitialized(tickBitmap, 1), "is not changed by another flip to a different tick on another word")
	assert.Equal(t, true, isInitialized(tickBitmap, 257), "is not changed by another flip to a different tick on another word")
}

func TestFlipTick1(t *testing.T) {
	tickBitmap := initTicks([]int{})

	flipTickBitmap(tickBitmap, -230, 1)
	assert.Equal(t, true, isInitialized(tickBitmap, -230))
	assert.Equal(t, false, isInitialized(tickBitmap, -231))
	assert.Equal(t, false, isInitialized(tickBitmap, -229))
	assert.Equal(t, false, isInitialized(tickBitmap, -230+256))
	assert.Equal(t, false, isInitialized(tickBitmap, -230-256))

	flipTickBitmap(tickBitmap, -230, 1)
	assert.Equal(t, false, isInitialized(tickBitmap, -230))
	assert.Equal(t, false, isInitialized(tickBitmap, -231))
	assert.Equal(t, false, isInitialized(tickBitmap, -229))
	assert.Equal(t, false, isInitialized(tickBitmap, -230+256))
	assert.Equal(t, false, isInitialized(tickBitmap, -230-256))
}

// reverts only itself.
func TestFlipTick2(t *testing.T) {
	tickBitmap := initTicks([]int{})
	flipTickBitmap(tickBitmap, -230, 1)
	flipTickBitmap(tickBitmap, -259, 1)
	flipTickBitmap(tickBitmap, -229, 1)
	flipTickBitmap(tickBitmap, 500, 1)
	flipTickBitmap(tickBitmap, -259, 1)
	flipTickBitmap(tickBitmap, -229, 1)
	flipTickBitmap(tickBitmap, -259, 1)

	assert.Equal(t, true, isInitialized(tickBitmap, 500))
	assert.Equal(t, true, isInitialized(tickBitmap, -230))
	assert.Equal(t, true, isInitialized(tickBitmap, -259))
	assert.Equal(t, false, isInitialized(tickBitmap, -229))
}

// lte = false.
func TestNextInitializedTickWithinOneWord1(t *testing.T) {
	tickBitmap := initTicks([]int{-200, -55, -4, 70, 78, 84, 139, 240, 535})

	check := func(tick int, expNext int, expInit bool) {
		next, initialized := nextInitializedTickWithinOneWord(tickBitmap, tick, 1, false)
		assert.Equal(t, expNext, next)
		assert.Equal(t, expInit, initialized)
	}

	check(78, 84, true)
	check(-55, -4, true)
	check(77, 78, true)
	check(-56, -55, true)

	check(255, 511, false)
	check(-257, -200, true)

	// flip 340
	flipTickBitmap(tickBitmap, 340, 1)
	check(328, 340, true)
	flipTickBitmap(tickBitmap, 340, 1)

	// does not exceed boundary
	check(508, 511, false)

	// skips entire word
	check(255, 511, false)

	// skips half word
	check(383, 511, false)
}

// lte = true.
func TestNextInitializedTickWithinOneWord2(t *testing.T) {
	tickBitmap := initTicks([]int{-200, -55, -4, 70, 78, 84, 139, 240, 535})
	check := func(tick int, expNext int, expInit bool) {
		next, initialized := nextInitializedTickWithinOneWord(tickBitmap, tick, 1, true)
		assert.Equal(t, expNext, next)
		assert.Equal(t, expInit, initialized)
	}

	// returns same tick if initialized
	check(78, 78, true)

	// returns tick directly to the left of input tick if not initialized
	check(79, 78, true)

	// will not exceed the word boundary
	check(258, 256, false)

	// at the word boundary
	check(256, 256, false)

	// word boundary less 1 (next initialized tick in next word
	check(72, 70, true)

	// word boundary
	check(-257, -512, false)

	// entire empty word
	check(1023, 768, false)

	// halfway through empty word
	check(900, 768, false)

	// boundary is initialized
	flipTickBitmap(tickBitmap, 329, 1)
	check(456, 329, true)
}
