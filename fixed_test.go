package decimal_test

import (
	"math"
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestNewFixed(t *testing.T) {
	tests := []struct {
		name string
		got  decimal.Fixed
		want decimal.Fixed
	}{
		// unsigned integers
		{"uint", decimal.NewFixed(uint(42)), 4200},
		{"uint8", decimal.NewFixed(uint8(255)), 25500},
		{"uint16", decimal.NewFixed(uint16(1000)), 100000},
		{"uint32", decimal.NewFixed(uint32(100)), 10000},
		{"uint64", decimal.NewFixed(uint64(123)), 12300},
		// signed integers
		{"int", decimal.NewFixed(123), 12300},
		{"int_negative", decimal.NewFixed(-5), -500},
		{"int8_min", decimal.NewFixed(int8(-128)), -12800},
		{"int16_min", decimal.NewFixed(int16(-32768)), -3276800},
		{"int16_max", decimal.NewFixed(int16(32767)), 3276700},
		{"int32", decimal.NewFixed(int32(100)), 10000},
		{"int64_negative", decimal.NewFixed(int64(-7)), -700},
		{"zero", decimal.NewFixed(0), 0},
		// float64 (truncated to 2 fractional digits)
		{"float64", decimal.NewFixed(123.456), 12345},
		{"float64_exact", decimal.NewFixed(123.45), 12345},
		{"float64_negative", decimal.NewFixed(-123.456), -12345},
		{"float64_half", decimal.NewFixed(1.5), 150},
		{"float64_fraction_only", decimal.NewFixed(0.99), 99},
		{"float64_small_fraction", decimal.NewFixed(100.01), 10001},
		{"float64_zero", decimal.NewFixed(0.0), 0},
		{"float64_nan", decimal.NewFixed(math.NaN()), 0},
		{"float64_pos_inf", decimal.NewFixed(math.Inf(1)), 0},
		{"float64_neg_inf", decimal.NewFixed(math.Inf(-1)), 0},
		// float64 cast/truncation edge cases: NewFixed(v) == int32(v*100), toward zero
		{"float64_trunc_below", decimal.NewFixed(0.29), 28}, // 0.29 is stored as 0.2899...; truncates to 0.28
		{"float64_trunc_below_large", decimal.NewFixed(19.99), 1998},
		{"float64_trunc_below_115", decimal.NewFixed(1.15), 114},
		{"float64_negative_trunc", decimal.NewFixed(-0.29), -28},
		{"float64_negative_trunc_large", decimal.NewFixed(-19.99), -1998},
		{"float64_subcent_truncated", decimal.NewFixed(0.001), 0},
		{"float64_subcent_truncated_nine", decimal.NewFixed(0.009), 0},
		{"float64_negative_subcent", decimal.NewFixed(-0.009), 0},
		{"float64_exact_half", decimal.NewFixed(0.5), 50},
		{"float64_exact_quarter", decimal.NewFixed(0.25), 25},
		{"float64_third_digit_truncated", decimal.NewFixed(0.125), 12},
		{"float64_max", decimal.NewFixed(21474836.47), 2147483647},
		{"float64_min", decimal.NewFixed(-21474836.48), -2147483648},
		// float32 (truncated to 2 fractional digits)
		{"float32_half", decimal.NewFixed(float32(1.5)), 150},
		{"float32_negative_half", decimal.NewFixed(float32(-1.5)), -150},
		{"float32_fraction_only", decimal.NewFixed(float32(0.99)), 99},
		{"float32_small_fraction", decimal.NewFixed(float32(100.01)), 10001},
		{"float32_nan", decimal.NewFixed(float32(math.NaN())), 0},
		{"float32_pos_inf", decimal.NewFixed(float32(math.Inf(1))), 0},
		{"float32_neg_inf", decimal.NewFixed(float32(math.Inf(-1))), 0},
		// float32 truncation can differ from float64 because float32(x) rounds differently
		{"float32_trunc", decimal.NewFixed(float32(0.29)), 29}, // float32(0.29) rounds above 0.29, unlike float64
		{"float32_trunc_large", decimal.NewFixed(float32(19.99)), 1999},
		{"float32_negative_trunc", decimal.NewFixed(float32(-19.99)), -1999},
		{"float32_subcent_truncated", decimal.NewFixed(float32(0.001)), 0},
		{"float32_mid", decimal.NewFixed(float32(12345.67)), 1234567},
		// Decimal (truncated to 2 fractional digits)
		{"decimal_integer", decimal.NewFixed(decimal.Decimal{Integer: 123}), 12300},
		{"decimal_1_digit", decimal.NewFixed(decimal.Decimal{Integer: 123, Fraction: 4, Digits: 1}), 12340},
		{"decimal_2_digits", decimal.NewFixed(decimal.Decimal{Integer: 123, Fraction: 45, Digits: 2}), 12345},
		{"decimal_3_digits", decimal.NewFixed(decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}), 12345},
		{"decimal_4_digits", decimal.NewFixed(decimal.Decimal{Integer: 123, Fraction: 4567, Digits: 4}), 12345},
		{"decimal_zero", decimal.NewFixed(decimal.Decimal{}), 0},
		{"decimal_negative_integer", decimal.NewFixed(decimal.Decimal{Integer: 5, Negative: true}), -500},
		{"decimal_negative_fraction", decimal.NewFixed(decimal.Decimal{Integer: 123, Fraction: 45, Digits: 2, Negative: true}), -12345},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("NewFixed() = %d, want %d", int32(tt.got), int32(tt.want))
			}
		})
	}
}

func BenchmarkNewFixed(b *testing.B) {
	b.Run("uint64", func(b *testing.B) {
		var v uint64 = 123
		for b.Loop() {
			_ = decimal.NewFixed(v)
		}
	})
	b.Run("int", func(b *testing.B) {
		var v int = -123
		for b.Loop() {
			_ = decimal.NewFixed(v)
		}
	})
	b.Run("float32", func(b *testing.B) {
		var v float32 = 123.45
		for b.Loop() {
			_ = decimal.NewFixed(v)
		}
	})
	b.Run("float64", func(b *testing.B) {
		var v float64 = 123.45
		for b.Loop() {
			_ = decimal.NewFixed(v)
		}
	})
	b.Run("decimal", func(b *testing.B) {
		v := decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}
		for b.Loop() {
			_ = decimal.NewFixed(v)
		}
	})
}

func TestFixed_Decimal(t *testing.T) {
	tests := []struct {
		name string
		f    decimal.Fixed
		want decimal.Decimal
	}{
		{"zero", 0, decimal.Decimal{Digits: 2}},
		{"hundredth", 1, decimal.Decimal{Fraction: 1, Digits: 2}},
		{"fraction", 5, decimal.Decimal{Fraction: 5, Digits: 2}},
		{"fraction_full", 99, decimal.Decimal{Fraction: 99, Digits: 2}},
		{"one", 100, decimal.Decimal{Integer: 1, Digits: 2}},
		{"integer", 12300, decimal.Decimal{Integer: 123, Digits: 2}},
		{"digits", 12345, decimal.Decimal{Integer: 123, Fraction: 45, Digits: 2}},
		{"max_int32", 2147483647, decimal.Decimal{Integer: 21474836, Fraction: 47, Digits: 2}},
		{"negative_hundredth", -1, decimal.Decimal{Fraction: 1, Digits: 2, Negative: true}},
		{"negative_fraction", -5, decimal.Decimal{Fraction: 5, Digits: 2, Negative: true}},
		{"negative_integer", -500, decimal.Decimal{Integer: 5, Digits: 2, Negative: true}},
		{"negative_digits", -12345, decimal.Decimal{Integer: 123, Fraction: 45, Digits: 2, Negative: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Decimal(); got != tt.want {
				t.Errorf("Fixed.Decimal() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func BenchmarkFixed_Decimal(b *testing.B) {
	f := decimal.Fixed(-12345)
	for b.Loop() {
		_ = f.Decimal()
	}
}
