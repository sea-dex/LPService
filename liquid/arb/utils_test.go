package arb

import "testing"

func TestTickAtPrice(t *testing.T) {
	ticks := []int32{20000, 78244, 110826}
	for _, tick := range ticks {
		t.Logf("price at %d is : %v", tick, PriceAtTick(tick))
	}

	prices := []float64{2500.0, 65000.0}
	for _, price := range prices {
		t.Logf("tick at %v is : %v", price, TickAtPrice(price))
	}
}
