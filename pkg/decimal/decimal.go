package decimal

import (
	"github.com/shopspring/decimal"
)

func CompByFloat64AndString(a float64, b string) int {
	a1 := decimal.NewFromFloat(a)
	b1, _ := decimal.NewFromString(b)

	return a1.Cmp(b1)
}

func Float64ToString(a float64, round ...int32) string {
	if len(round) > 0 {
		return decimal.NewFromFloat(a).Round(round[0]).String()
	}

	return decimal.NewFromFloat(a).String()
}

func Float64Round(a float64, round ...int32) float64 {
	f, _ := decimal.NewFromFloat(a).RoundFloor(round[0]).Float64()

	return f
}

func StringToStringRound(a string, round ...int32) string {
	dd, _ := decimal.NewFromString(a)
	s := dd.RoundFloor(round[0]).String()

	return s
}

func Float64MulFloat64ToFloat64(a, b float64, round ...int32) float64 {
	a1 := decimal.NewFromFloat(a)
	b1 := decimal.NewFromFloat(b)

	if len(round) > 0 {
		f, _ := a1.Mul(b1).Round(round[0]).Float64()

		return f
	}

	f, _ := a1.Mul(b1).Float64()

	return f
}

func StringMulStringToString(a, b string, round ...int32) string {
	a1, _ := decimal.NewFromString(a)
	b1, _ := decimal.NewFromString(b)

	if len(round) > 0 {
		f := a1.Mul(b1).Round(round[0]).String()

		return f
	}

	f := a1.Mul(b1).String()

	return f
}

func Float64AddFloat64ToFloat64(a, b float64, round ...int32) float64 {
	a1 := decimal.NewFromFloat(a)
	b1 := decimal.NewFromFloat(b)
	f, _ := a1.Add(b1).Float64()

	if len(round) > 0 {
		f, _ := a1.Add(b1).Round(round[0]).Float64()

		return f
	}

	return f
}

func Float64AddFloat64ToFloat642(a, b float64, round ...int32) float64 {
	a1 := decimal.NewFromFloat(a)
	b1 := decimal.NewFromFloat(b)
	f, _ := a1.Add(b1).Float64()

	if len(round) > 0 {
		f, _ := a1.Add(b1).RoundFloor(round[0]).Float64()

		return f
	}

	return f
}

func Float64SubFloat64ToFloat64(a, b float64, round ...int32) float64 {
	a1 := decimal.NewFromFloat(a)
	b1 := decimal.NewFromFloat(b)

	if len(round) > 0 {
		f, _ := a1.Sub(b1).Round(round[0]).Float64()

		return f
	}

	f, _ := a1.Sub(b1).Float64()

	return f
}

func Float64DivFloat64ToFloat64(a, b float64, round ...int32) float64 {
	a1 := decimal.NewFromFloat(a)
	b1 := decimal.NewFromFloat(b)

	if len(round) > 0 {
		f, _ := a1.Div(b1).Round(round[0]).Float64()

		return f
	}

	f, _ := a1.Div(b1).Float64()

	return f
}

func Float64DivFloat64ToString(a, b float64, round ...int32) string {
	a1 := decimal.NewFromFloat(a)
	b1 := decimal.NewFromFloat(b)

	if len(round) > 0 {
		return a1.Div(b1).Round(round[0]).String()
	}

	return a1.Div(b1).String()
}

func StringCmpString(a, b decimal.Decimal) (bool, error) {
	result := a.Cmp(b)

	switch result {
	case 1:
		return true, nil
	case -1, 0:
		return false, nil
	default:
		return false, nil
	}
}

func EqualsZero(a decimal.Decimal) bool {
	return a.Cmp(decimal.NewFromInt(0)) == 0
}
