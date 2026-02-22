package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopic(t *testing.T) {
	assert.Equal(t, "0xae468ce586f9a87660fdffc1448cee942042c16ae2f02046b134b5224f31936b", TopicAeroV2SetFee)
	assert.Equal(t, "0xd444e1b10a2a0c61e10ee9f0167820955df343074f16b69614952caef34de21d", TopicAeroV3SetFee)
}
