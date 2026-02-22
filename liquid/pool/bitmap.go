package pool

import "math/big"

// tick bitmap

func position(tick int) (wordPos int16, bitPos uint8) {
	wordPos = int16(tick >> 8) // nolint
	bitPos = uint8(tick % 256) // nolint

	return
}

func flipTickBitmap(tickBitmap map[int16]*big.Int, tick, tickSpacing int) {
	if tick%tickSpacing != 0 {
		panic("invalid tick")
	}

	wordPos, bitPos := position(tick / tickSpacing)
	mask := big.NewInt(0).Lsh(big.NewInt(1), uint(bitPos))

	v, ok := tickBitmap[wordPos]
	if ok {
		tickBitmap[wordPos] = v.Xor(v, mask)
	} else {
		tickBitmap[wordPos] = big.NewInt(0).Xor(big.NewInt(0), mask)
	}
}

// / Returns the next initialized tick contained in the same word (or adjacent word) as the tick that is either
// / to the left (less than or equal to) or right (greater than) of the given tick.
func nextInitializedTickWithinOneWord(
	tickBitmap map[int16]*big.Int,
	tick, tickSpacing int,
	lte bool,
) (next int, initialized bool) {
	compressed := tick / tickSpacing
	if tick < 0 && tick%tickSpacing != 0 {
		compressed--
	}

	if lte {
		wordPos, bitPos := position(compressed)
		// uint256 mask = (1 << bitPos) - 1 + (1 << bitPos);
		mask := big.NewInt(0).Lsh(big.NewInt(1), uint(bitPos))
		mask.Sub(mask, bigOne)
		mask.Add(mask, big.NewInt(1).Lsh(big.NewInt(1), uint(bitPos)))

		masked := big.NewInt(0)
		if v, ok := tickBitmap[wordPos]; ok {
			// println("mask:", mask.Text(16))
			// println("v:", v.Text(16))
			masked.And(v, mask)
		}

		initialized = masked.Cmp(bigZero) != 0
		// println("nextInitializedTickWithinOneWord lte=true", initialized, bitPos, masked.Text(2))
		if initialized {
			msb := mostSignificantBit(masked)
			// println("msb:", msb)
			next = (compressed - int(bitPos-msb)) * tickSpacing
		} else {
			next = (compressed - int(bitPos)) * tickSpacing
		}
	} else {
		wordPos, bitPos := position(compressed + 1)
		// uint256 mask = ~((1 << bitPos) - 1);
		mask := big.NewInt(0).Lsh(big.NewInt(1), uint(bitPos))
		mask.Sub(mask, bigOne)
		mask.Not(mask)

		masked := big.NewInt(0)
		if v, ok := tickBitmap[wordPos]; ok {
			masked.And(v, mask)
		}

		initialized = masked.Cmp(bigZero) != 0
		if initialized {
			next = (compressed + 1 + int(leastSignificantBit(masked)-bitPos)) * tickSpacing
		} else {
			next = (compressed + 1 + int(255-bitPos)) * tickSpacing
		}
	}

	return
}

func mostSignificantBit(x *big.Int) (r uint8) {
	if x.Cmp(bigZero) <= 0 {
		panic("x should great than 0")
	}

	cmpRsh := func(n *big.Int, bits uint8) {
		if x.Cmp(n) >= 0 {
			x.Rsh(x, uint(bits))
			r += bits
		}
	}

	cmpRsh(toBigIntMust("0x100000000000000000000000000000000"), 128)
	cmpRsh(toBigIntMust("0x10000000000000000"), 64)
	cmpRsh(toBigIntMust("0x100000000"), 32)
	cmpRsh(toBigIntMust("0x10000"), 16)
	cmpRsh(toBigIntMust("0x100"), 8)
	cmpRsh(toBigIntMust("0x10"), 4)
	cmpRsh(toBigIntMust("0x4"), 2)
	cmpRsh(toBigIntMust("0x2"), 1)

	return
}

func leastSignificantBit(x *big.Int) (r uint8) {
	if x.Cmp(bigZero) <= 0 {
		panic("x should great than 0")
	}

	r = 255
	andRsh := func(v *big.Int, bits uint8) {
		if big.NewInt(0).And(x, v).Cmp(bigZero) > 0 {
			r -= bits
		} else {
			x.Rsh(x, uint(bits))
		}
	}

	// type(uint128).max
	andRsh(uint128Max, 128)
	andRsh(uint64Max, 64)
	andRsh(uint32Max, 32)
	andRsh(uint16Max, 16)
	andRsh(uint8Max, 8)
	andRsh(big.NewInt(15), 4)
	andRsh(big.NewInt(3), 2)
	andRsh(big.NewInt(1), 1)

	return
}
