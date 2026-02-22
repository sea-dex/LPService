package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindNextTick(t *testing.T) {
	tickNext, exist := findNextTick([]int{}, 1, true)
	assert.False(t, exist)
	assert.Equal(t, tickNext, MIN_TICK)

	tickNext, exist = findNextTick([]int{}, 1, false)
	assert.False(t, exist)
	assert.Equal(t, tickNext, MAX_TICK)

	tickNext, exist = findNextTick([]int{2}, 1, true)
	assert.False(t, exist)
	assert.Equal(t, tickNext, MIN_TICK)

	tickNext, exist = findNextTick([]int{2}, 1, false)
	assert.True(t, exist)
	assert.Equal(t, tickNext, 2)

	tickNext, exist = findNextTick([]int{2}, 2, true)
	assert.False(t, exist)
	assert.Equal(t, tickNext, MIN_TICK)

	tickNext, exist = findNextTick([]int{2}, 2, false)
	assert.False(t, exist)
	assert.Equal(t, tickNext, MAX_TICK)
}

func TestFindNextTick2(t *testing.T) {
	tickNext, exist := findNextTick([]int{-1, 1}, 1, true)
	assert.True(t, exist)
	assert.Equal(t, tickNext, -1)

	tickNext, exist = findNextTick([]int{1, 2}, 1, false)
	assert.True(t, exist)
	assert.Equal(t, tickNext, 2)

	tickNext, exist = findNextTick([]int{2, 4, 6}, 1, true)
	assert.False(t, exist)
	assert.Equal(t, tickNext, MIN_TICK)

	tickNext, exist = findNextTick([]int{2, 4, 6}, 3, true)
	assert.True(t, exist)
	assert.Equal(t, tickNext, 2)

	tickNext, exist = findNextTick([]int{2, 4, 5, 6}, 5, true)
	assert.True(t, exist)
	assert.Equal(t, tickNext, 4)

	tickNext, exist = findNextTick([]int{2, 4, 5, 6}, 7, true)
	assert.True(t, exist)
	assert.Equal(t, tickNext, 6)

	tickNext, exist = findNextTick([]int{2, 4, 6}, 1, false)
	assert.True(t, exist)
	assert.Equal(t, tickNext, 2)

	tickNext, exist = findNextTick([]int{2, 4, 6}, 3, false)
	assert.True(t, exist)
	assert.Equal(t, tickNext, 4)

	tickNext, exist = findNextTick([]int{2, 4, 5, 6}, 5, false)
	assert.True(t, exist)
	assert.Equal(t, tickNext, 6)

	tickNext, exist = findNextTick([]int{2, 4, 5, 6}, 7, false)
	assert.False(t, exist)
	assert.Equal(t, tickNext, MAX_TICK)
}
