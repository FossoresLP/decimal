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

func TestDecimal_MultiplyUint64(t *testing.T) {
	tests := []struct {
		name       string
		decimal    decimal.Decimal
		multiplier uint64
		expected   *decimal.Decimal
	}{
		{"zero", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 0, &decimal.Decimal{Digits: 3}},
		{"negative_to_zero", decimal.Decimal{Integer: 5, Negative: true}, 0, &decimal.Decimal{}},
		{"integer", decimal.Decimal{Integer: 123}, 2, &decimal.Decimal{Integer: 246}},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, 2, &decimal.Decimal{Fraction: 246, Digits: 3}},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 2, &decimal.Decimal{Integer: 246, Fraction: 912, Digits: 3}},
		{"negative", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, 2, &decimal.Decimal{Integer: 246, Fraction: 912, Digits: 3, Negative: true}},
		{"large", decimal.Decimal{Integer: 1234567890123456789, Fraction: 1234567890123456789, Digits: 19}, 2, &decimal.Decimal{Integer: 2469135780246913578, Fraction: 2469135780246913578, Digits: 19}},
		{"fraction_mul_overflow_with_carry", decimal.Decimal{Fraction: 9999999999999999999, Digits: 19}, 2, &decimal.Decimal{Integer: 1, Fraction: 9999999999999999998, Digits: 19}},
		{"fraction_mul_overflow_large_carry", decimal.Decimal{Fraction: 9999999999999999999, Digits: 19}, 7, &decimal.Decimal{Integer: 6, Fraction: 9999999999999999993, Digits: 19}},
		{"fraction_mul_overflow_with_integer_wrap_semantics", decimal.Decimal{Integer: math.MaxUint64, Fraction: 9999999999999999999, Digits: 19}, 2, &decimal.Decimal{Integer: math.MaxUint64, Fraction: 9999999999999999998, Digits: 19}},
		{"large_multiplier", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 1234567890, &decimal.Decimal{Integer: 152414813427, Fraction: 840, Digits: 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := tt.decimal.MultiplyUint64(tt.multiplier)
			if !tt.expected.Equal(product) {
				t.Errorf("Decimal.MultiplyUint64() = %v, want %v", product, tt.expected)
			}
		})
	}
}

func TestDecimal_DivideUint64(t *testing.T) {
	tests := []struct {
		name     string
		decimal  decimal.Decimal
		divisor  uint64
		expected *decimal.Decimal
	}{
		{"one", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 1, &decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}},
		{"integer", decimal.Decimal{Integer: 123}, 2, &decimal.Decimal{Integer: 61}},
		{"negative_to_zero", decimal.Decimal{Fraction: 1, Digits: 1, Negative: true}, 2, &decimal.Decimal{Digits: 1}},
		{"integer_with_fraction_digit", decimal.Decimal{Integer: 123, Digits: 1}, 2, &decimal.Decimal{Integer: 61, Fraction: 5, Digits: 1}},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, 2, &decimal.Decimal{Fraction: 61, Digits: 3}},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 2, &decimal.Decimal{Integer: 61, Fraction: 728, Digits: 3}},
		{"negative", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, 2, &decimal.Decimal{Integer: 61, Fraction: 728, Digits: 3, Negative: true}},
		{"large", decimal.Decimal{Integer: 1234567890, Fraction: 123456789, Digits: 10}, 2, &decimal.Decimal{Integer: 617283945, Fraction: 61728394, Digits: 10}},
		{"high_digits_remainder", decimal.Decimal{Integer: 100, Fraction: 0, Digits: 19}, 7, &decimal.Decimal{Integer: 14, Fraction: 2857142857142857142, Digits: 19}},
		{"max_uint64", decimal.Decimal{Integer: math.MaxUint64, Digits: 19}, 7, &decimal.Decimal{Integer: math.MaxUint64 / 7, Fraction: 1428571428571428571, Digits: 19}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotient := tt.decimal.DivideUint64(tt.divisor)
			if !tt.expected.Equal(quotient) {
				t.Errorf("Decimal.DivideUint64() = %v, want %v", quotient, tt.expected)
			}
		})
	}
}

func TestDecimal_Add(t *testing.T) {
	tests := []struct {
		name     string
		d1       decimal.Decimal
		d2       decimal.Decimal
		expected decimal.Decimal
	}{
		{"zero", decimal.Decimal{}, decimal.Decimal{}, decimal.Decimal{}},
		{"integer", decimal.Decimal{Integer: 123}, decimal.Decimal{Integer: 456}, decimal.Decimal{Integer: 579}},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, decimal.Decimal{Fraction: 456, Digits: 3}, decimal.Decimal{Fraction: 579, Digits: 3}},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 3}, decimal.Decimal{Integer: 579, Fraction: 579, Digits: 3}},
		{"negative_d1", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 3}, decimal.Decimal{Integer: 333, Fraction: 333, Digits: 3}},
		{"negative_d2", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 3, Negative: true}, decimal.Decimal{Integer: 333, Fraction: 333, Digits: 3, Negative: true}},
		{"negative_both", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 3, Negative: true}, decimal.Decimal{Integer: 579, Fraction: 579, Digits: 3, Negative: true}},
		{"mixed_sign_integer_d1_larger", decimal.Decimal{Integer: 7}, decimal.Decimal{Integer: 5, Negative: true}, decimal.Decimal{Integer: 2}},
		{"mixed_sign_integer_d2_larger", decimal.Decimal{Integer: 5}, decimal.Decimal{Integer: 7, Negative: true}, decimal.Decimal{Integer: 2, Negative: true}},
		{"mixed_sign_integer_cancel", decimal.Decimal{Integer: 7}, decimal.Decimal{Integer: 7, Negative: true}, decimal.Decimal{}},
		{"large", decimal.Decimal{Integer: 1234567890, Fraction: 123456789, Digits: 10}, decimal.Decimal{Integer: 9876543210, Fraction: 987654321, Digits: 10}, decimal.Decimal{Integer: 11111111100, Fraction: 1111111110, Digits: 10}},
		{"different_digits_d1_less", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 4}, decimal.Decimal{Integer: 579, Fraction: 1686, Digits: 4}},
		{"different_digits_d2_less", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 4}, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 3}, decimal.Decimal{Integer: 579, Fraction: 4683, Digits: 4}},
		{"mixed_sign_smaller_int_larger_frac", decimal.Decimal{Integer: 5, Fraction: 7, Digits: 1, Negative: true}, decimal.Decimal{Integer: 6, Fraction: 3, Digits: 1}, decimal.Decimal{Integer: 0, Fraction: 6, Digits: 1}},
		{"mixed_sign_borrow", decimal.Decimal{Integer: 1, Fraction: 9, Digits: 1}, decimal.Decimal{Negative: true, Integer: 2, Fraction: 1, Digits: 1}, decimal.Decimal{Negative: true, Integer: 0, Fraction: 2, Digits: 1}},
		{"mixed_sign_exact_cancellation", decimal.Decimal{Integer: 1, Fraction: 2, Digits: 1}, decimal.Decimal{Negative: true, Integer: 1, Fraction: 2, Digits: 1}, decimal.Decimal{Digits: 1}},
		{"mixed_sign_equal_int_d1_frac_larger", decimal.Decimal{Integer: 3, Fraction: 7, Digits: 1}, decimal.Decimal{Negative: true, Integer: 3, Fraction: 2, Digits: 1}, decimal.Decimal{Integer: 0, Fraction: 5, Digits: 1}},
		{"mixed_sign_equal_int_d2_frac_larger", decimal.Decimal{Integer: 3, Fraction: 2, Digits: 1}, decimal.Decimal{Negative: true, Integer: 3, Fraction: 7, Digits: 1}, decimal.Decimal{Negative: true, Integer: 0, Fraction: 5, Digits: 1}},
		{"mixed_sign_d2_larger_int_borrow", decimal.Decimal{Integer: 1, Fraction: 7, Digits: 1}, decimal.Decimal{Negative: true, Integer: 3, Fraction: 2, Digits: 1}, decimal.Decimal{Negative: true, Integer: 1, Fraction: 5, Digits: 1}},
		{"mixed_sign_d1_larger_int_borrow", decimal.Decimal{Integer: 5, Fraction: 1, Digits: 1}, decimal.Decimal{Negative: true, Integer: 3, Fraction: 7, Digits: 1}, decimal.Decimal{Integer: 1, Fraction: 4, Digits: 1}},
		{"mixed_sign_digits_19_borrow_d1_larger", decimal.Decimal{Integer: 2, Fraction: 9999999999999999998, Digits: 19}, decimal.Decimal{Negative: true, Integer: 1, Fraction: 9999999999999999999, Digits: 19}, decimal.Decimal{Integer: 0, Fraction: 9999999999999999999, Digits: 19}},
		{"mixed_sign_digits_19_borrow_d2_larger", decimal.Decimal{Negative: true, Integer: 1, Fraction: 9999999999999999999, Digits: 19}, decimal.Decimal{Integer: 2, Fraction: 9999999999999999998, Digits: 19}, decimal.Decimal{Integer: 0, Fraction: 9999999999999999999, Digits: 19}},
		{"fraction_overflow", decimal.Decimal{Fraction: 9999999999999999999, Digits: 19}, decimal.Decimal{Fraction: 9999999999999999999, Digits: 19}, decimal.Decimal{Integer: 1, Fraction: 9999999999999999998, Digits: 19}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum := tt.d1.Add(tt.d2)
			if !tt.expected.Equal(sum) {
				t.Errorf("Decimal.Add() = %v, want %v", sum, tt.expected)
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

func BenchmarkDecimal_NewFloat32(b *testing.B) {
	for b.Loop() {
		_ = decimal.New(float32(123.456))
	}
}

func BenchmarkDecimal_NewFloat64(b *testing.B) {
	for b.Loop() {
		_ = decimal.New(123.456)
	}
}

func BenchmarkDecimal_NewInt(b *testing.B) {
	for b.Loop() {
		_ = decimal.New(123456)
	}
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

func TestDecimal_DivideUint64_DivByZero(t *testing.T) {
	// Division by zero panics, matching native uint64 behavior.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("DivideUint64(0) did not panic — expected panic consistent with uint64 division by zero")
		}
	}()
	d := decimal.Decimal{Integer: 1}
	_ = d.DivideUint64(0)
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
