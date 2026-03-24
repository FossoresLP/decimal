package decimal_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestDecimal_Equal(t *testing.T) {
	tests := []struct {
		name string
		d1   decimal.Decimal
		d2   decimal.Decimal
		want bool
	}{
		{"zero", decimal.Decimal{}, decimal.Decimal{}, true},
		{"integer", decimal.Decimal{Integer: 123}, decimal.Decimal{Integer: 123}, true},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, decimal.Decimal{Fraction: 123, Digits: 3}, true},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, true},
		{"negative", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, true},
		{"different", decimal.Decimal{Integer: 123}, decimal.Decimal{Integer: 456}, false},
		{"different_fraction", decimal.Decimal{Fraction: 123, Digits: 3}, decimal.Decimal{Fraction: 456, Digits: 3}, false},
		{"different_digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 4}, false},
		{"different_negative", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"value_equal_different_digits", decimal.Decimal{Fraction: 1, Digits: 1}, decimal.Decimal{Fraction: 10, Digits: 2}, true},
		{"value_equal_different_digits_3", decimal.Decimal{Fraction: 1, Digits: 1}, decimal.Decimal{Fraction: 100, Digits: 3}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d1.Equal(tt.d2); got != tt.want {
				t.Errorf("Decimal.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		got  decimal.Decimal
		want decimal.Decimal
	}{
		// float64
		{"float64", decimal.New(123.123).ToDigits(3), decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}},
		{"float64_negative", decimal.New(-123.123).ToDigits(3), decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}},
		{"float64_half", decimal.New(123.5), decimal.Decimal{Integer: 123, Fraction: 5, Digits: 1}},
		{"float64_negative_half", decimal.New(-123.5), decimal.Decimal{Integer: 123, Fraction: 5, Digits: 1, Negative: true}},
		// float32
		{"float32", decimal.New(float32(123.123)).ToDigits(3), decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}},
		{"float32_negative", decimal.New(float32(-123.123)).ToDigits(3), decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}},
		{"float32_half", decimal.New(float32(123.5)), decimal.Decimal{Integer: 123, Fraction: 5, Digits: 1}},
		{"float32_negative_half", decimal.New(float32(-123.5)), decimal.Decimal{Integer: 123, Fraction: 5, Digits: 1, Negative: true}},
		// int
		{"int", decimal.New(123), decimal.Decimal{Integer: 123}},
		{"int_negative", decimal.New(-123), decimal.Decimal{Integer: 123, Negative: true}},
		// min-value signed integers (must widen before negation)
		{"int8_min", decimal.New(int8(-128)), decimal.Decimal{Integer: 128, Negative: true}},
		{"int16_min", decimal.New(int16(-32768)), decimal.Decimal{Integer: 32768, Negative: true}},
		{"int32_min", decimal.New(int32(-2147483648)), decimal.Decimal{Integer: 2147483648, Negative: true}},
		{"int64_min", decimal.New(int64(math.MinInt64)), decimal.Decimal{Integer: uint64(math.MaxInt64) + 1, Negative: true}},
		// unsigned integer types
		{"uint", decimal.New(uint(42)), decimal.Decimal{Integer: 42}},
		{"uint8", decimal.New(uint8(255)), decimal.Decimal{Integer: 255}},
		{"uint16", decimal.New(uint16(65535)), decimal.Decimal{Integer: 65535}},
		{"uint32", decimal.New(uint32(4294967295)), decimal.Decimal{Integer: 4294967295}},
		{"uint64", decimal.New(uint64(math.MaxUint64)), decimal.Decimal{Integer: math.MaxUint64}},
		// float32 special values
		{"float32_nan", decimal.New(float32(math.NaN())), decimal.Decimal{}},
		{"float32_pos_inf", decimal.New(float32(math.Inf(1))), decimal.Decimal{}},
		{"float32_neg_inf", decimal.New(float32(math.Inf(-1))), decimal.Decimal{}},
		{"float64_negative_subnormal", decimal.New(-5e-324), decimal.Decimal{}},
		{"float64_negative_tiny", decimal.New(-1e-19), decimal.Decimal{}},
		{"float32_negative_tiny", decimal.New(float32(-1e-20)), decimal.Decimal{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.want.Equal(tt.got) {
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
			if !tt.want.Equal(result) {
				t.Errorf("Decimal.ToDigits() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDecimal_IsZero(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
		want bool
	}{
		{"zero_value", decimal.Decimal{}, true},
		{"zero_with_digits", decimal.Decimal{Digits: 3}, true},
		{"zero_negative", decimal.Decimal{Negative: true}, true},
		{"integer_only", decimal.Decimal{Integer: 1}, false},
		{"fraction_only", decimal.Decimal{Fraction: 1, Digits: 1}, false},
		{"both", decimal.Decimal{Integer: 1, Fraction: 1, Digits: 1}, false},
		{"negative_integer", decimal.Decimal{Integer: 1, Negative: true}, false},
		{"negative_fraction", decimal.Decimal{Fraction: 5, Digits: 1, Negative: true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.IsZero(); got != tt.want {
				t.Errorf("Decimal.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_Zero(t *testing.T) {
	zero := decimal.Decimal{}
	if !zero.Equal(decimal.Zero()) {
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
			if !tt.want.Equal(result) {
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
