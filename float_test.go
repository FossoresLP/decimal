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

func TestDecimal_FromFloat64(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
		f    float64
		want decimal.Decimal
	}{
		{"zero", decimal.Decimal{}, 0, decimal.Decimal{}},
		{"integer", decimal.Decimal{}, 123, decimal.Decimal{Integer: 123}},
		{"fraction", decimal.Decimal{}, 0.123, decimal.Decimal{Fraction: 123, Digits: 3}},
		{"digits", decimal.Decimal{}, 123.123, decimal.Decimal{Integer: 123, Fraction: 123000000000004656, Digits: 18}},
		{"negative", decimal.Decimal{}, -123.123, decimal.Decimal{Integer: 123, Fraction: 123000000000004656, Digits: 18, Negative: true}},
		{"large", decimal.Decimal{}, 1234567890.1234567890, decimal.Decimal{Integer: 1234567890, Fraction: 123456716537475584, Digits: 18}},
		{"huge", decimal.Decimal{}, 1234567890123456789.12345678901234567890, decimal.Decimal{Integer: 1234567890123456768}},
		{"small", decimal.Decimal{}, 0.000000000000000001, decimal.Decimal{Fraction: 1, Digits: 18}},
		{"negative_small", decimal.Decimal{}, -0.000000000000000001, decimal.Decimal{Fraction: 1, Digits: 18, Negative: true}},
		{"indivisible", decimal.Decimal{}, 3.3333333333, decimal.Decimal{Integer: 3, Fraction: 333333333300000128, Digits: 18}},
		{"truncate_1", decimal.Decimal{}, 0.999999999999999944, decimal.Decimal{Fraction: 999999999999999872, Digits: 18}},
		{"truncate_0.5", decimal.Decimal{}, 0.499999999999999944, decimal.Decimal{Fraction: 499999999999999936, Digits: 18}},
		{"subnormal", decimal.Decimal{}, 5e-324, decimal.Decimal{}},
		{"nan", decimal.Decimal{}, math.NaN(), decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}},
		{"positive_infinity", decimal.Decimal{}, math.Inf(1), decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}},
		{"negative_infinity", decimal.Decimal{}, math.Inf(-1), decimal.Decimal{Negative: true, Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}},
		{"max_float64", decimal.Decimal{}, math.MaxFloat64, decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}},
		{"min_float64", decimal.Decimal{}, -math.MaxFloat64, decimal.Decimal{Negative: true, Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}},
		{"max_uint64", decimal.Decimal{}, math.MaxUint64, decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}},
		{"max_int64", decimal.Decimal{}, math.MaxInt64, decimal.Decimal{Integer: 0x8000000000000000}},
		{"min_int64", decimal.Decimal{}, math.MinInt64, decimal.Decimal{Negative: true, Integer: 0x8000000000000000}},
		{"reuse", decimal.Decimal{Negative: true, Integer: 123, Fraction: 123, Digits: 3}, 456.456, decimal.Decimal{Integer: 456, Fraction: 45600000000001728, Digits: 17}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.FromFloat64(tt.f)
			if !tt.d.Equal(tt.want) {
				t.Errorf("Decimal.FromFloat64() = %#v, want %#v", tt.d, tt.want)
			}
		})
	}
}

func TestDecimal_FromFloat64Fixed(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
		f    float64
		want decimal.Decimal
		prec uint8
	}{
		{"zero", decimal.Decimal{}, 0, decimal.Decimal{Digits: 18}, 18},
		{"integer", decimal.Decimal{}, 123, decimal.Decimal{Integer: 123, Digits: 18}, 18},
		{"fraction", decimal.Decimal{}, 0.123, decimal.Decimal{Fraction: 123, Digits: 3}, 3},
		{"fraction_full", decimal.Decimal{}, 0.123, decimal.Decimal{Fraction: 123000000000000000, Digits: 18}, 18},
		{"digits", decimal.Decimal{}, 123.123, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, 3},
		{"digits_full", decimal.Decimal{}, 123.123, decimal.Decimal{Integer: 123, Fraction: 123000000000004656, Digits: 18}, 18},
		{"negative", decimal.Decimal{}, -123.123, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, 3},
		{"negative_full", decimal.Decimal{}, -123.123, decimal.Decimal{Integer: 123, Fraction: 123000000000004656, Digits: 18, Negative: true}, 18},
		{"large", decimal.Decimal{}, 1234567890.1234567890, decimal.Decimal{Integer: 1234567890, Fraction: 123456716537475584, Digits: 18}, 18},
		{"huge", decimal.Decimal{}, 1234567890123456789.12345678901234567890, decimal.Decimal{Integer: 1234567890123456768, Digits: 18}, 18},
		{"small", decimal.Decimal{}, 0.000000000000000001, decimal.Decimal{Fraction: 1, Digits: 18}, 18},
		{"negative_small", decimal.Decimal{}, -0.000000000000000001, decimal.Decimal{Fraction: 1, Digits: 18, Negative: true}, 18},
		{"indivisible", decimal.Decimal{}, 3.3333333333, decimal.Decimal{Integer: 3, Fraction: 333333333300000128, Digits: 18}, 18},
		{"truncate_1", decimal.Decimal{}, 0.999999999999999944, decimal.Decimal{Fraction: 999999999999, Digits: 12}, 12},
		{"truncate_0.5", decimal.Decimal{}, 0.499999999999999944, decimal.Decimal{Fraction: 499999999999, Digits: 12}, 12},
		{"subnormal", decimal.Decimal{}, 5e-324, decimal.Decimal{Digits: 18}, 18},
		{"nan", decimal.Decimal{}, math.NaN(), decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}, 18},
		{"positive_infinity", decimal.Decimal{}, math.Inf(1), decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}, 18},
		{"negative_infinity", decimal.Decimal{}, math.Inf(-1), decimal.Decimal{Negative: true, Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}, 18},
		{"max_float64", decimal.Decimal{}, math.MaxFloat64, decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}, 18},
		{"min_float64", decimal.Decimal{}, -math.MaxFloat64, decimal.Decimal{Negative: true, Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}, 18},
		{"max_uint64", decimal.Decimal{}, math.MaxUint64, decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}, 18},
		{"max_int64", decimal.Decimal{}, math.MaxInt64, decimal.Decimal{Integer: 0x8000000000000000, Digits: 18}, 18},
		{"min_int64", decimal.Decimal{}, math.MinInt64, decimal.Decimal{Negative: true, Integer: 0x8000000000000000, Digits: 18}, 18},
		{"reuse", decimal.Decimal{Negative: true, Integer: 123, Fraction: 123, Digits: 3}, 456.456, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 3}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.FromFloat64Fixed(tt.f, tt.prec)
			if !tt.d.Equal(tt.want) {
				t.Errorf("Decimal.FromFloat64Fixed() = %#v, want %#v", tt.d, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_FromFloat64(b *testing.B) {
	d := decimal.Decimal{}
	benches := []struct {
		name string
		f    float64
	}{
		{"zero", 0},
		{"integer", 123},
		{"fraction", 0.123},
		{"digits", 123.123},
		{"negative", -123.123},
		{"large", 1234567890.1234567890},
		{"huge", 1234567890123456789.12345678901234567890},
		{"small", 0.000000000000000001},
		{"negative_small", -0.000000000000000001},
		{"indivisible", 3.3333333333},
	}
	for _, bench := range benches {
		b.Run(bench.name, func(b *testing.B) {
			for b.Loop() {
				d.FromFloat64(bench.f)
			}
			b.ReportMetric(float64(d.Digits), "digits")
		})
	}
}

func BenchmarkDecimal_FromFloat64Fixed(b *testing.B) {
	d := decimal.Decimal{}
	benches := []struct {
		name string
		f    float64
	}{
		{"zero", 0},
		{"integer", 123},
		{"fraction", 0.123},
		{"digits", 123.123},
		{"negative", -123.123},
		{"large", 1234567890.1234567890},
		{"huge", 1234567890123456789.12345678901234567890},
		{"small", 0.000000000000000001},
		{"negative_small", -0.000000000000000001},
		{"indivisible", 3.3333333333},
	}
	for _, bench := range benches {
		b.Run(bench.name, func(b *testing.B) {
			for b.Loop() {
				d.FromFloat64Fixed(bench.f, 18)
			}
		})
	}
}

func FuzzDecimal_FromFloat64(f *testing.F) {
	testcases := []float64{123.456, 0.123, 0.000000000000000001, 1234567890123456789.12345678901234567890, 3.3333333333}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, f float64) {
		if f > math.MaxUint64 {
			return
		}
		da := decimal.Decimal{}
		da.FromFloat64(f)
		if math.Abs(da.Float64()-f) > 1e-15 {
			t.Errorf("Decimal.FromFloat64() = %.18f, want %.18f", da.Float64(), f)
		}
		df := decimal.Decimal{}
		df.FromFloat64Fixed(f, 18)
		if math.Abs(df.Float64()-f) > 1e-15 {
			t.Errorf("Decimal.FromFloat64Fixed() = %.18f, want %.18f", df.Float64(), f)
		}
	})
}
