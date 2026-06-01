package decimal_test

import (
	"math"
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestDecimal_Float64(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
		want float64
	}{
		{"zero", decimal.Decimal{}, 0},
		{"integer", decimal.Decimal{Integer: 123}, 123},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, 0.123},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, 123.123},
		{"negative", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, -123.123},
		{"indivisible", decimal.Decimal{Integer: 3, Fraction: 3333333333, Digits: 10}, 3.3333333333},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Float64(); got != tt.want {
				t.Errorf("Decimal.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_Float64(b *testing.B) {
	d := decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}
	for b.Loop() {
		_ = d.Float64()
	}
}

func TestFixed_Float64(t *testing.T) {
	tests := []struct {
		name string
		f    decimal.Fixed
		want float64
	}{
		{"zero", 0, 0},
		{"hundredth", 1, 0.01},
		{"fraction", 5, 0.05},
		{"one", 100, 1},
		{"integer", 12300, 123},
		{"digits", 12345, 123.45},
		{"negative_hundredth", -5, -0.05},
		{"negative_integer", -500, -5},
		{"negative_digits", -12345, -123.45},
		{"negative_one_hundredth", -1, -0.01},
		{"two_nines", 99, 0.99},
		{"large_fraction", 9999, 99.99},
		{"negative_large_fraction", -9999, -99.99},
		{"max", math.MaxInt32, 21474836.47},
		{"min", math.MinInt32, -21474836.48},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Float64(); got != tt.want {
				t.Errorf("Fixed.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkFixed_Float64(b *testing.B) {
	f := decimal.Fixed(12345)
	for b.Loop() {
		_ = f.Float64()
	}
}
