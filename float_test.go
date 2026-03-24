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

func TestDecimal_NewFloat(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		want decimal.Decimal
	}{
		{"zero", 0, decimal.Decimal{}},
		{"integer", 123, decimal.Decimal{Integer: 123}},
		{"fraction", 0.123, decimal.Decimal{Fraction: 123, Digits: 3}},
		{"digits", 123.123, decimal.Decimal{Integer: 123, Fraction: 123000000000004656, Digits: 18}},
		{"negative", -123.123, decimal.Decimal{Integer: 123, Fraction: 123000000000004656, Digits: 18, Negative: true}},
		{"large", 1234567890.1234567890, decimal.Decimal{Integer: 1234567890, Fraction: 123456716537475584, Digits: 18}},
		{"huge", 1234567890123456789.12345678901234567890, decimal.Decimal{Integer: 1234567890123456768}},
		{"small", 0.000000000000000001, decimal.Decimal{Fraction: 1, Digits: 18}},
		{"negative_small", -0.000000000000000001, decimal.Decimal{Fraction: 1, Digits: 18, Negative: true}},
		{"indivisible", 3.3333333333, decimal.Decimal{Integer: 3, Fraction: 333333333300000128, Digits: 18}},
		{"truncate_1", 0.999999999999999944, decimal.Decimal{Fraction: 999999999999999872, Digits: 18}},
		{"truncate_0.5", 0.499999999999999944, decimal.Decimal{Fraction: 499999999999999936, Digits: 18}},
		{"subnormal", 5e-324, decimal.Decimal{}},
		{"nan", math.NaN(), decimal.Decimal{}},
		{"positive_infinity", math.Inf(1), decimal.Decimal{}},
		{"negative_infinity", math.Inf(-1), decimal.Decimal{}},
		{"max_float64", math.MaxFloat64, decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}},
		{"min_float64", -math.MaxFloat64, decimal.Decimal{Negative: true, Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}},
		{"max_uint64", math.MaxUint64, decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}},
		{"max_int64", math.MaxInt64, decimal.Decimal{Integer: 0x8000000000000000}},
		{"min_int64", math.MinInt64, decimal.Decimal{Negative: true, Integer: 0x8000000000000000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := decimal.New(tt.f); !decimal.Equal(tt.want, got) {
				t.Errorf("Decimal.FromFloat64() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func FuzzDecimal_NewFloat(f *testing.F) {
	testcases := []float64{123.456, 0.123, 0.000000000000000001, 1234567890123456789.12345678901234567890, 3.3333333333}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, f float64) {
		if f > math.MaxUint64 {
			return
		}
		d := decimal.New(f)
		if math.Abs(d.Float64()-f) > 1e-15 {
			t.Errorf("Decimal.FromFloat64() = %.18f, want %.18f", d.Float64(), f)
		}
	})
}
