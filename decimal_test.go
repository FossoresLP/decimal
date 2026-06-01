package decimal_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		got  decimal.Decimal
		want decimal.Decimal
	}{
		// float32
		{"float32_zero", decimal.New(float32(0)), decimal.Decimal{}},
		{"float32_integer", decimal.New(float32(123)), decimal.Decimal{Integer: 123}},
		{"float32_fraction", decimal.New(float32(0.123)), decimal.Decimal{Fraction: 123, Digits: 3}},
		{"float32_digits", decimal.New(float32(123.123)), decimal.Decimal{Integer: 123, Fraction: 1230011008, Digits: 10}},
		{"float32_negative", decimal.New(float32(-123.123)), decimal.Decimal{Integer: 123, Fraction: 1230011008, Digits: 10, Negative: true}},
		{"float32_half", decimal.New(float32(123.5)), decimal.Decimal{Integer: 123, Fraction: 5, Digits: 1}},
		{"float32_negative_half", decimal.New(float32(-123.5)), decimal.Decimal{Integer: 123, Fraction: 5, Digits: 1, Negative: true}},
		{"float32_quarter", decimal.New(float32(0.25)), decimal.Decimal{Fraction: 25, Digits: 2}},
		{"float32_eighth", decimal.New(float32(0.125)), decimal.Decimal{Fraction: 125, Digits: 3}},
		{"float32_large", decimal.New(float32(1234567)), decimal.Decimal{Integer: 1234567}},
		{"float32_seven_digits", decimal.New(float32(1234567.5)), decimal.Decimal{Integer: 1234567, Fraction: 5, Digits: 1}},
		{"float32_indivisible", decimal.New(float32(1.0 / 3.0)), decimal.Decimal{Fraction: 3333333504, Digits: 10}},
		{"float32_truncate", decimal.New(float32(0.999999)), decimal.Decimal{Fraction: 999998976, Digits: 9}},
		{"float32_small", decimal.New(float32(1e-10)), decimal.Decimal{Fraction: 1, Digits: 10}},
		{"float32_negative_small", decimal.New(float32(-1e-10)), decimal.Decimal{Fraction: 1, Digits: 10, Negative: true}},
		{"float32_underflow", decimal.New(float32(1e-18)), decimal.Decimal{}},
		{"float32_negative_underflow", decimal.New(float32(-1e-20)), decimal.Decimal{}},
		{"float32_subnormal", decimal.New(float32(math.SmallestNonzeroFloat32)), decimal.Decimal{}},
		{"float32_negative_subnormal", decimal.New(float32(-math.SmallestNonzeroFloat32)), decimal.Decimal{}},
		{"float32_nan", decimal.New(float32(math.NaN())), decimal.Decimal{}},
		{"float32_pos_inf", decimal.New(float32(math.Inf(1))), decimal.Decimal{}},
		{"float32_neg_inf", decimal.New(float32(math.Inf(-1))), decimal.Decimal{}},
		{"float32_max", decimal.New(float32(math.MaxFloat32)), decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 10}},
		{"float32_min", decimal.New(float32(-math.MaxFloat32)), decimal.Decimal{Negative: true, Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 10}},
		{"float32_max_int32", decimal.New(float32(math.MaxInt32)), decimal.Decimal{Integer: 2147483648}},
		{"float32_min_int32", decimal.New(float32(math.MinInt32)), decimal.Decimal{Integer: 2147483648, Negative: true}},
		{"float32_max_uint32", decimal.New(float32(math.MaxUint32)), decimal.Decimal{Integer: 4294967296}},
		// float64
		{"float64_zero", decimal.New(0.0), decimal.Decimal{}},
		{"float64_integer", decimal.New(123.0), decimal.Decimal{Integer: 123}},
		{"float64_fraction", decimal.New(0.123), decimal.Decimal{Fraction: 123, Digits: 3}},
		{"float64_digits", decimal.New(123.123), decimal.Decimal{Integer: 123, Fraction: 123000000000004656, Digits: 18}},
		{"float64_negative", decimal.New(-123.123), decimal.Decimal{Integer: 123, Fraction: 123000000000004656, Digits: 18, Negative: true}},
		{"float64_half", decimal.New(123.5), decimal.Decimal{Integer: 123, Fraction: 5, Digits: 1}},
		{"float64_negative_half", decimal.New(-123.5), decimal.Decimal{Integer: 123, Fraction: 5, Digits: 1, Negative: true}},
		{"float64_large", decimal.New(1234567890.1234567890), decimal.Decimal{Integer: 1234567890, Fraction: 123456716537475584, Digits: 18}},
		{"float64_int_precision_loss", decimal.New(1234567890123456789.12345678901234567890), decimal.Decimal{Integer: 1234567890123456768}},
		{"float64_small", decimal.New(0.000000000000000001), decimal.Decimal{Fraction: 1, Digits: 18}},
		{"float64_negative_small", decimal.New(-0.000000000000000001), decimal.Decimal{Fraction: 1, Digits: 18, Negative: true}},
		{"float64_indivisible", decimal.New(3.3333333333), decimal.Decimal{Integer: 3, Fraction: 333333333300000128, Digits: 18}},
		{"float64_truncate_1", decimal.New(0.999999999999999944), decimal.Decimal{Fraction: 999999999999999872, Digits: 18}},
		{"float64_truncate_0.5", decimal.New(0.499999999999999944), decimal.Decimal{Fraction: 499999999999999936, Digits: 18}},
		{"float64_subnormal", decimal.New(5e-324), decimal.Decimal{}},
		{"float64_negative_subnormal", decimal.New(-5e-324), decimal.Decimal{}},
		{"float64_negative_tiny", decimal.New(-1e-19), decimal.Decimal{}},
		{"float64_nan", decimal.New(math.NaN()), decimal.Decimal{}},
		{"float64_pos_inf", decimal.New(math.Inf(1)), decimal.Decimal{}},
		{"float64_neg_inf", decimal.New(math.Inf(-1)), decimal.Decimal{}},
		{"float64_max", decimal.New(math.MaxFloat64), decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}},
		{"float64_min", decimal.New(-math.MaxFloat64), decimal.Decimal{Negative: true, Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}},
		{"float64_max_uint64", decimal.New(float64(math.MaxUint64)), decimal.Decimal{Integer: 0x8000000000000000, Fraction: 0x8000000000000000, Digits: 18}},
		{"float64_max_int64", decimal.New(float64(math.MaxInt64)), decimal.Decimal{Integer: 0x8000000000000000}},
		{"float64_min_int64", decimal.New(float64(math.MinInt64)), decimal.Decimal{Negative: true, Integer: 0x8000000000000000}},
		// signed integers
		{"int", decimal.New(123), decimal.Decimal{Integer: 123}},
		{"int_negative", decimal.New(-123), decimal.Decimal{Integer: 123, Negative: true}},
		// min-value signed integers (must widen before negation)
		{"int8_min", decimal.New(int8(math.MinInt8)), decimal.Decimal{Integer: uint64(math.MaxInt8) + 1, Negative: true}},
		{"int16_min", decimal.New(int16(math.MinInt16)), decimal.Decimal{Integer: uint64(math.MaxInt16) + 1, Negative: true}},
		{"int32_min", decimal.New(int32(math.MinInt32)), decimal.Decimal{Integer: uint64(math.MaxInt32) + 1, Negative: true}},
		{"int64_min", decimal.New(int64(math.MinInt64)), decimal.Decimal{Integer: uint64(math.MaxInt64) + 1, Negative: true}},
		// unsigned integers
		{"uint", decimal.New(uint(42)), decimal.Decimal{Integer: 42}},
		{"uint8", decimal.New(uint8(math.MaxUint8)), decimal.Decimal{Integer: math.MaxUint8}},
		{"uint16", decimal.New(uint16(math.MaxUint16)), decimal.Decimal{Integer: math.MaxUint16}},
		{"uint32", decimal.New(uint32(math.MaxUint32)), decimal.Decimal{Integer: math.MaxUint32}},
		{"uint64", decimal.New(uint64(math.MaxUint64)), decimal.Decimal{Integer: math.MaxUint64}},
		// Fixed
		{"fixed", decimal.New(decimal.Fixed(12345)), decimal.Decimal{Integer: 123, Fraction: 45, Digits: 2}},
		{"fixed_negative", decimal.New(decimal.Fixed(-12345)), decimal.Decimal{Integer: 123, Fraction: 45, Digits: 2, Negative: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("New() = %#v, want %#v", tt.got, tt.want)
			}
		})
	}
}

func BenchmarkNew(b *testing.B) {
	b.Run("uint", func(b *testing.B) {
		var v uint = 123
		for b.Loop() {
			_ = decimal.New(v)
		}
	})
	b.Run("uint8", func(b *testing.B) {
		var v uint8 = 123
		for b.Loop() {
			_ = decimal.New(v)
		}
	})
	b.Run("uint16", func(b *testing.B) {
		var v uint16 = 123
		for b.Loop() {
			_ = decimal.New(v)
		}
	})
	b.Run("uint32", func(b *testing.B) {
		var v uint32 = 123
		for b.Loop() {
			_ = decimal.New(v)
		}
	})
	b.Run("uint64", func(b *testing.B) {
		var v uint64 = 123
		for b.Loop() {
			_ = decimal.New(v)
		}
	})
	b.Run("int", func(b *testing.B) {
		var v int = -123
		for b.Loop() {
			_ = decimal.New(v)
		}
	})
	b.Run("int8", func(b *testing.B) {
		var v int8 = -123
		for b.Loop() {
			_ = decimal.New(v)
		}
	})
	b.Run("int16", func(b *testing.B) {
		var v int16 = -123
		for b.Loop() {
			_ = decimal.New(v)
		}
	})
	b.Run("int32", func(b *testing.B) {
		var v int32 = -123
		for b.Loop() {
			_ = decimal.New(v)
		}
	})
	b.Run("int64", func(b *testing.B) {
		var v int64 = -123
		for b.Loop() {
			_ = decimal.New(v)
		}
	})
	b.Run("float32", func(b *testing.B) {
		var v float32 = -123.456
		for b.Loop() {
			_ = decimal.New(v)
		}
	})
	b.Run("float64", func(b *testing.B) {
		var v float64 = -123.456
		for b.Loop() {
			_ = decimal.New(v)
		}
	})
}

func TestDecimal_ToDigits(t *testing.T) {
	tests := []struct {
		name   string
		d      decimal.Decimal
		digits uint8
		want   decimal.Decimal
	}{
		{"zero", decimal.Decimal{}, 3, decimal.Decimal{Digits: 3}},
		{"integer", decimal.Decimal{Integer: 123}, 3, decimal.Decimal{Integer: 123, Digits: 3}},
		{"negative_to_zero", decimal.Decimal{Fraction: 1, Digits: 1, Negative: true}, 0, decimal.Decimal{}},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, 6, decimal.Decimal{Fraction: 123000, Digits: 6}},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, 6, decimal.Decimal{Integer: 123, Fraction: 123000, Digits: 6}},
		{"negative", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, 6, decimal.Decimal{Integer: 123, Fraction: 123000, Digits: 6, Negative: true}},
		{"less_low", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, 2, decimal.Decimal{Integer: 123, Fraction: 12, Digits: 2}},
		{"less_high", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 2, decimal.Decimal{Integer: 123, Fraction: 45, Digits: 2}},
		{"clamp_out_of_range", decimal.Decimal{Integer: 1}, 48, decimal.Decimal{Integer: 1, Digits: 19}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.d.ToDigits(tt.digits)
			if !decimal.Equal(tt.want, result) {
				t.Errorf("Decimal.ToDigits() = %v, want %v", result, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_ToDigits(b *testing.B) {
	var d decimal.Decimal
	for b.Loop() {
		_ = d.ToDigits(12)
	}
}

func TestDecimal_Round(t *testing.T) {
	tests := []struct {
		name   string
		d      decimal.Decimal
		digits uint8
		want   decimal.Decimal
	}{
		// No-op: same digits
		{"same_digits", decimal.Decimal{Integer: 1, Fraction: 456, Digits: 3}, 3, decimal.Decimal{Integer: 1, Fraction: 456, Digits: 3}},
		// Expand: fewer digits to more
		{"expand", decimal.Decimal{Fraction: 12, Digits: 2}, 5, decimal.Decimal{Fraction: 12000, Digits: 5}},
		{"expand_zero", decimal.Decimal{}, 3, decimal.Decimal{Digits: 3}},
		{"expand_integer", decimal.Decimal{Integer: 42}, 2, decimal.Decimal{Integer: 42, Digits: 2}},
		// Round down (remainder < 5)
		{"round_down", decimal.Decimal{Integer: 1, Fraction: 1234, Digits: 4}, 2, decimal.Decimal{Integer: 1, Fraction: 12, Digits: 2}},
		{"round_down_zero_rem", decimal.Decimal{Fraction: 100, Digits: 3}, 1, decimal.Decimal{Fraction: 1, Digits: 1}},
		// Round up (remainder >= 5)
		{"round_up", decimal.Decimal{Integer: 1, Fraction: 456, Digits: 3}, 2, decimal.Decimal{Integer: 1, Fraction: 46, Digits: 2}},
		{"round_up_5", decimal.Decimal{Fraction: 15, Digits: 2}, 1, decimal.Decimal{Fraction: 2, Digits: 1}},
		{"round_up_carry_fraction_only", decimal.Decimal{Fraction: 99, Digits: 2}, 1, decimal.Decimal{Integer: 1, Digits: 1}},
		{"round_up_carry_integer", decimal.Decimal{Integer: 1, Fraction: 999, Digits: 3}, 2, decimal.Decimal{Integer: 2, Digits: 2}},
		{"round_up_carry_negative", decimal.Decimal{Integer: 1, Fraction: 999, Digits: 3, Negative: true}, 2, decimal.Decimal{Integer: 2, Digits: 2, Negative: true}},
		// Round to zero digits
		{"round_to_zero_down", decimal.Decimal{Integer: 5, Fraction: 4, Digits: 1}, 0, decimal.Decimal{Integer: 5}},
		{"round_to_zero_up", decimal.Decimal{Integer: 5, Fraction: 5, Digits: 1}, 0, decimal.Decimal{Integer: 6}},
		// Negative values
		{"negative_round_down", decimal.Decimal{Integer: 1, Fraction: 1234, Digits: 4, Negative: true}, 2, decimal.Decimal{Integer: 1, Fraction: 12, Digits: 2, Negative: true}},
		{"negative_round_up", decimal.Decimal{Integer: 1, Fraction: 456, Digits: 3, Negative: true}, 2, decimal.Decimal{Integer: 1, Fraction: 46, Digits: 2, Negative: true}},
		// Negative to zero clears sign
		{"negative_to_zero", decimal.Decimal{Fraction: 4, Digits: 1, Negative: true}, 0, decimal.Decimal{}},
		// Clamp digits > 19
		{"clamp_out_of_range", decimal.Decimal{Integer: 1}, 48, decimal.Decimal{Integer: 1, Digits: 19}},
		// High precision
		{"high_precision_round", decimal.Decimal{Integer: 3, Fraction: 1415926535897932384, Digits: 19}, 4, decimal.Decimal{Integer: 3, Fraction: 1416, Digits: 4}},
		{"high_precision_round_carry", decimal.Decimal{Fraction: 9999999999999999999, Digits: 19}, 18, decimal.Decimal{Integer: 1, Digits: 18}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.d.Round(tt.digits)
			if !decimal.Equal(tt.want, result) {
				t.Errorf("Decimal.Round() = %v, want %v", result, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_Round(b *testing.B) {
	d := decimal.Decimal{Integer: 1, Fraction: 99999, Digits: 5}
	for b.Loop() {
		_ = d.Round(2)
	}
}

func TestDecimal_Zero(t *testing.T) {
	zero := decimal.Decimal{}
	if !decimal.Equal(zero, decimal.Zero()) {
		t.Errorf("Decimal.Zero() = %v, want %v", decimal.Zero(), zero)
	}
}

func TestDecimal_Truncate(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
		want decimal.Decimal
	}{
		{"zero", decimal.Decimal{Digits: 3}, decimal.Decimal{Digits: 0}},
		{"integer", decimal.Decimal{Integer: 123, Digits: 3}, decimal.Decimal{Integer: 123, Digits: 0}},
		{"fraction", decimal.Decimal{Fraction: 123000, Digits: 6}, decimal.Decimal{Fraction: 123, Digits: 3}},
		{"18_zeros", decimal.Decimal{Integer: 123, Fraction: 1000000000000000000, Digits: 19}, decimal.Decimal{Integer: 123, Fraction: 1, Digits: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.d.Truncate()
			if !decimal.Equal(tt.want, result) {
				t.Errorf("Decimal.Truncate() = %v, want %v", result, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_Truncate(b *testing.B) {
	zeros := 0
	for frac := uint64(1); frac < 10000000000000000000; frac *= 10 {
		zeros++
		d := decimal.Decimal{Negative: true, Digits: 18, Integer: 1, Fraction: frac}
		b.Run(fmt.Sprintf("zeros=%d", zeros), func(b *testing.B) {
			for b.Loop() {
				_ = d.Truncate()
			}
		})
	}
}
