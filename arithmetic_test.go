package decimal_test

import (
	"math"
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestDecimal_Add(t *testing.T) {
	tests := []struct {
		name     string
		a        decimal.Decimal
		b        decimal.Decimal
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
			sum := decimal.Add(tt.a, tt.b)
			if !decimal.Equal(tt.expected, sum) {
				t.Errorf("Decimal.Add() = %v, want %v", sum, tt.expected)
			}
		})
	}
}

func TestAdd_WithGenerics(t *testing.T) {
	// Decimal + int
	result := decimal.Add(decimal.Decimal{Integer: 10, Fraction: 5, Digits: 1}, 3)
	expected := decimal.Decimal{Integer: 13, Fraction: 5, Digits: 1}
	if !decimal.Equal(expected, result) {
		t.Errorf("Add(Decimal, int) = %v, want %v", result, expected)
	}

	// int + Decimal
	result = decimal.Add(3, decimal.Decimal{Integer: 10, Fraction: 5, Digits: 1})
	if !decimal.Equal(expected, result) {
		t.Errorf("Add(int, Decimal) = %v, want %v", result, expected)
	}

	// negative int
	result = decimal.Add(decimal.Decimal{Integer: 10}, -3)
	expected = decimal.Decimal{Integer: 7}
	if !decimal.Equal(expected, result) {
		t.Errorf("Add(Decimal, negative int) = %v, want %v", result, expected)
	}

	// float64
	result = decimal.Add(decimal.Decimal{Integer: 1}, 0.5)
	expected = decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}
	if !decimal.Equal(expected, result) {
		t.Errorf("Add(Decimal, float64) = %v, want %v", result, expected)
	}

	// uint64
	result = decimal.Add(uint64(100), decimal.Decimal{Integer: 50})
	expected = decimal.Decimal{Integer: 150}
	if !decimal.Equal(expected, result) {
		t.Errorf("Add(uint64, Decimal) = %v, want %v", result, expected)
	}
}

func TestAdd_Commutativity(t *testing.T) {
	a := decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}
	b := decimal.Decimal{Integer: 789, Fraction: 12, Digits: 2}
	ab := decimal.Add(a, b)
	ba := decimal.Add(b, a)
	if !decimal.Equal(ab, ba) {
		t.Errorf("Add not commutative: Add(a,b) = %v, Add(b,a) = %v", ab, ba)
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a        decimal.Decimal
		b        decimal.Decimal
		expected decimal.Decimal
	}{
		{"zero_minus_zero", decimal.Decimal{}, decimal.Decimal{}, decimal.Decimal{}},
		{"integer_minus_integer", decimal.Decimal{Integer: 7}, decimal.Decimal{Integer: 3}, decimal.Decimal{Integer: 4}},
		{"result_negative", decimal.Decimal{Integer: 3}, decimal.Decimal{Integer: 7}, decimal.Decimal{Integer: 4, Negative: true}},
		{"fraction_minus_fraction", decimal.Decimal{Fraction: 75, Digits: 2}, decimal.Decimal{Fraction: 25, Digits: 2}, decimal.Decimal{Fraction: 50, Digits: 2}},
		{"mixed_minus_mixed", decimal.Decimal{Integer: 5, Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 2, Fraction: 3, Digits: 1}, decimal.Decimal{Integer: 3, Fraction: 2, Digits: 1}},
		{"borrow_from_integer", decimal.Decimal{Integer: 5, Fraction: 1, Digits: 1}, decimal.Decimal{Integer: 2, Fraction: 7, Digits: 1}, decimal.Decimal{Integer: 2, Fraction: 4, Digits: 1}},
		{"cancel_to_zero", decimal.Decimal{Integer: 5, Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 5, Fraction: 5, Digits: 1}, decimal.Decimal{Digits: 1}},
		{"negative_minus_positive", decimal.Decimal{Integer: 3, Negative: true}, decimal.Decimal{Integer: 4}, decimal.Decimal{Integer: 7, Negative: true}},
		{"negative_minus_negative", decimal.Decimal{Integer: 3, Negative: true}, decimal.Decimal{Integer: 7, Negative: true}, decimal.Decimal{Integer: 4}},
		{"positive_minus_negative", decimal.Decimal{Integer: 3}, decimal.Decimal{Integer: 4, Negative: true}, decimal.Decimal{Integer: 7}},
		{"different_digits", decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}, decimal.Decimal{Fraction: 25, Digits: 2}, decimal.Decimal{Integer: 1, Fraction: 25, Digits: 2}},
		{"zero_minus_value", decimal.Decimal{}, decimal.Decimal{Integer: 5, Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 5, Fraction: 5, Digits: 1, Negative: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decimal.Subtract(tt.a, tt.b)
			if !decimal.Equal(tt.expected, result) {
				t.Errorf("Subtract() = %v (%#v), want %v (%#v)", result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestSubtract_WithGenerics(t *testing.T) {
	result := decimal.Subtract(decimal.Decimal{Integer: 10}, 3)
	expected := decimal.Decimal{Integer: 7}
	if !decimal.Equal(expected, result) {
		t.Errorf("Subtract(Decimal, int) = %v, want %v", result, expected)
	}

	result = decimal.Subtract(10, decimal.Decimal{Integer: 3})
	if !decimal.Equal(expected, result) {
		t.Errorf("Subtract(int, Decimal) = %v, want %v", result, expected)
	}
}

func TestSubtract_AddInverse(t *testing.T) {
	a := decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}
	b := decimal.Decimal{Integer: 789, Fraction: 12, Digits: 2}
	// a - b + b == a
	result := decimal.Add(decimal.Subtract(a, b), b)
	if !decimal.Equal(a.ToDigits(result.Digits), result) {
		t.Errorf("Subtract/Add inverse failed: got %v, want %v", result, a)
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a        decimal.Decimal
		b        decimal.Decimal
		expected decimal.Decimal
	}{
		{"zero_times_zero", decimal.Decimal{}, decimal.Decimal{}, decimal.Decimal{}},
		{"zero_times_value", decimal.Decimal{}, decimal.Decimal{Integer: 5}, decimal.Decimal{}},
		{"value_times_zero", decimal.Decimal{Integer: 5}, decimal.Decimal{}, decimal.Decimal{}},
		{"zero_times_negative", decimal.Decimal{}, decimal.Decimal{Integer: 5, Negative: true}, decimal.Decimal{}},
		{"integer_times_integer", decimal.Decimal{Integer: 3}, decimal.Decimal{Integer: 4}, decimal.Decimal{Integer: 12}},
		{"integer_times_fraction", decimal.Decimal{Integer: 3}, decimal.Decimal{Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}},
		{"fraction_times_integer", decimal.Decimal{Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 3}, decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}},
		{"fraction_times_fraction", decimal.Decimal{Fraction: 5, Digits: 1}, decimal.Decimal{Fraction: 5, Digits: 1}, decimal.Decimal{Fraction: 25, Digits: 2}},
		{"mixed_times_mixed", decimal.Decimal{Integer: 2, Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 3, Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 8, Fraction: 75, Digits: 2}},
		{"identity_one", decimal.Decimal{Integer: 42, Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 1}, decimal.Decimal{Integer: 42, Fraction: 5, Digits: 1}},
		{"negative_times_positive", decimal.Decimal{Integer: 3, Negative: true}, decimal.Decimal{Integer: 4}, decimal.Decimal{Integer: 12, Negative: true}},
		{"positive_times_negative", decimal.Decimal{Integer: 3}, decimal.Decimal{Integer: 4, Negative: true}, decimal.Decimal{Integer: 12, Negative: true}},
		{"negative_times_negative", decimal.Decimal{Integer: 3, Negative: true}, decimal.Decimal{Integer: 4, Negative: true}, decimal.Decimal{Integer: 12}},
		{"fraction_carry_into_integer", decimal.Decimal{Fraction: 9, Digits: 1}, decimal.Decimal{Fraction: 9, Digits: 1}, decimal.Decimal{Fraction: 81, Digits: 2}},
		{"large_fraction_carry", decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 2, Fraction: 25, Digits: 2}},
		{"precision_capped_at_19", decimal.Decimal{Fraction: 1234567890, Digits: 10}, decimal.Decimal{Fraction: 1234567890, Digits: 10}, decimal.Decimal{Fraction: 152415787501905210, Digits: 19}},
		{"ten_times_ten", decimal.Decimal{Integer: 10}, decimal.Decimal{Integer: 10}, decimal.Decimal{Integer: 100}},
		{"large_integers", decimal.Decimal{Integer: 1000000}, decimal.Decimal{Integer: 1000000}, decimal.Decimal{Integer: 1000000000000}},
		{"integer_overflow_wraps", decimal.Decimal{Integer: math.MaxUint64}, decimal.Decimal{Integer: 2}, decimal.Decimal{Integer: math.MaxUint64 - 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decimal.Multiply(tt.a, tt.b)
			if !decimal.Equal(tt.expected, result) {
				t.Errorf("Multiply() = %v (%#v), want %v (%#v)", result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestMultiply_Int(t *testing.T) {
	tests := []struct {
		name     string
		a        decimal.Decimal
		b        uint64
		expected decimal.Decimal
	}{
		{"zero", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 0, decimal.Decimal{Digits: 3}},
		{"negative_to_zero", decimal.Decimal{Integer: 5, Negative: true}, 0, decimal.Decimal{}},
		{"integer", decimal.Decimal{Integer: 123}, 2, decimal.Decimal{Integer: 246}},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, 2, decimal.Decimal{Fraction: 246, Digits: 3}},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 2, decimal.Decimal{Integer: 246, Fraction: 912, Digits: 3}},
		{"negative", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, 2, decimal.Decimal{Integer: 246, Fraction: 912, Digits: 3, Negative: true}},
		{"large", decimal.Decimal{Integer: 1234567890123456789, Fraction: 1234567890123456789, Digits: 19}, 2, decimal.Decimal{Integer: 2469135780246913578, Fraction: 2469135780246913578, Digits: 19}},
		{"fraction_mul_overflow_with_carry", decimal.Decimal{Fraction: 9999999999999999999, Digits: 19}, 2, decimal.Decimal{Integer: 1, Fraction: 9999999999999999998, Digits: 19}},
		{"fraction_mul_overflow_large_carry", decimal.Decimal{Fraction: 9999999999999999999, Digits: 19}, 7, decimal.Decimal{Integer: 6, Fraction: 9999999999999999993, Digits: 19}},
		{"fraction_mul_overflow_with_integer_wrap_semantics", decimal.Decimal{Integer: math.MaxUint64, Fraction: 9999999999999999999, Digits: 19}, 2, decimal.Decimal{Integer: math.MaxUint64, Fraction: 9999999999999999998, Digits: 19}},
		{"large_multiplier", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 1234567890, decimal.Decimal{Integer: 152414813427, Fraction: 840, Digits: 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := decimal.Multiply(tt.a, tt.b)
			if !decimal.Equal(tt.expected, product) {
				t.Errorf("Decimal.MultiplyUint64() = %v, want %v", product, tt.expected)
			}
		})
	}
}

func TestMultiply_WithGenerics(t *testing.T) {
	// Decimal * int
	result := decimal.Multiply(decimal.Decimal{Integer: 5, Fraction: 5, Digits: 1}, 3)
	expected := decimal.Decimal{Integer: 16, Fraction: 5, Digits: 1}
	if !decimal.Equal(expected, result) {
		t.Errorf("Multiply(Decimal, int) = %v, want %v", result, expected)
	}

	// int * Decimal
	result = decimal.Multiply(3, decimal.Decimal{Integer: 5, Fraction: 5, Digits: 1})
	if !decimal.Equal(expected, result) {
		t.Errorf("Multiply(int, Decimal) = %v, want %v", result, expected)
	}

	// float64
	result = decimal.Multiply(decimal.Decimal{Integer: 4}, 2.5)
	expected = decimal.Decimal{Integer: 10}
	if !decimal.Equal(expected, result) {
		t.Errorf("Multiply(Decimal, float64) = %v, want %v", result, expected)
	}
}

func TestMultiply_Commutativity(t *testing.T) {
	a := decimal.Decimal{Integer: 12, Fraction: 345, Digits: 3}
	b := decimal.Decimal{Integer: 67, Fraction: 89, Digits: 2}
	ab := decimal.Multiply(a, b)
	ba := decimal.Multiply(b, a)
	if !decimal.Equal(ab, ba) {
		t.Errorf("Multiply not commutative: Multiply(a,b) = %v, Multiply(b,a) = %v", ab, ba)
	}
}

func TestMultiply_Distributive(t *testing.T) {
	a := decimal.Decimal{Integer: 3, Fraction: 5, Digits: 1}
	b := decimal.Decimal{Integer: 2}
	c := decimal.Decimal{Integer: 4}
	// a * (b + c) == a*b + a*c
	lhs := decimal.Multiply(a, decimal.Add(b, c))
	rhs := decimal.Add(decimal.Multiply(a, b), decimal.Multiply(a, c))
	if !decimal.Equal(lhs, rhs) {
		t.Errorf("distributive law failed: %v != %v", lhs, rhs)
	}
}

func TestConvert_SpecialFloats(t *testing.T) {
	// Inf and NaN should convert to zero when used with arithmetic functions
	zero := decimal.Decimal{}

	five := decimal.Decimal{Integer: 5}
	result := decimal.Add(five, math.Inf(1))
	if !decimal.Equal(five, result) {
		t.Errorf("Add(5, +Inf) = %v, want 5", result)
	}

	result = decimal.Multiply(decimal.Decimal{Integer: 5}, math.NaN())
	if !decimal.Equal(zero, result) {
		t.Errorf("Multiply(5, NaN) = %v, want 0", result)
	}

	result = decimal.Multiply(decimal.Decimal{Integer: 5}, math.Inf(-1))
	if !decimal.Equal(zero, result) {
		t.Errorf("Multiply(5, -Inf) = %v, want 0", result)
	}
}

func BenchmarkAdd(b *testing.B) {
	d1 := decimal.Decimal{Integer: 123, Fraction: 456789, Digits: 6}
	d2 := decimal.Decimal{Integer: 987, Fraction: 654321, Digits: 6}
	for b.Loop() {
		_ = decimal.Add(d1, d2)
	}
}

func BenchmarkSubtract(b *testing.B) {
	d1 := decimal.Decimal{Integer: 987, Fraction: 654321, Digits: 6}
	d2 := decimal.Decimal{Integer: 123, Fraction: 456789, Digits: 6}
	for b.Loop() {
		_ = decimal.Subtract(d1, d2)
	}
}

func BenchmarkMultiply(b *testing.B) {
	d1 := decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}
	d2 := decimal.Decimal{Integer: 789, Fraction: 12, Digits: 2}
	for b.Loop() {
		_ = decimal.Multiply(d1, d2)
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name     string
		a        decimal.Decimal
		b        decimal.Decimal
		expected decimal.Decimal
	}{
		{"integer_exact", decimal.Decimal{Integer: 10}, decimal.Decimal{Integer: 2}, decimal.Decimal{Integer: 5, Digits: 19}},
		{"integer_with_remainder", decimal.Decimal{Integer: 10}, decimal.Decimal{Integer: 3}, decimal.Decimal{Integer: 3, Fraction: 3333333333333333333, Digits: 19}},
		{"one_over_three", decimal.Decimal{Integer: 1}, decimal.Decimal{Integer: 3}, decimal.Decimal{Fraction: 3333333333333333333, Digits: 19}},
		{"fraction_by_integer", decimal.Decimal{Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 2}, decimal.Decimal{Fraction: 2500000000000000000, Digits: 19}},
		{"integer_by_fraction", decimal.Decimal{Integer: 1}, decimal.Decimal{Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 2, Digits: 19}},
		{"fraction_by_fraction", decimal.Decimal{Fraction: 1, Digits: 1}, decimal.Decimal{Fraction: 4, Digits: 1}, decimal.Decimal{Fraction: 2500000000000000000, Digits: 19}},
		{"negative_dividend", decimal.Decimal{Integer: 10, Negative: true}, decimal.Decimal{Integer: 3}, decimal.Decimal{Integer: 3, Fraction: 3333333333333333333, Digits: 19, Negative: true}},
		{"negative_divisor", decimal.Decimal{Integer: 10}, decimal.Decimal{Integer: 3, Negative: true}, decimal.Decimal{Integer: 3, Fraction: 3333333333333333333, Digits: 19, Negative: true}},
		{"both_negative", decimal.Decimal{Integer: 10, Negative: true}, decimal.Decimal{Integer: 3, Negative: true}, decimal.Decimal{Integer: 3, Fraction: 3333333333333333333, Digits: 19}},
		{"zero_dividend", decimal.Decimal{}, decimal.Decimal{Integer: 5}, decimal.Decimal{}},
		{"identity", decimal.Decimal{Integer: 7, Fraction: 25, Digits: 2}, decimal.Decimal{Integer: 1}, decimal.Decimal{Integer: 7, Fraction: 2500000000000000000, Digits: 19}},
		{"self_division", decimal.Decimal{Integer: 7, Fraction: 25, Digits: 2}, decimal.Decimal{Integer: 7, Fraction: 25, Digits: 2}, decimal.Decimal{Integer: 1, Digits: 19}},
		{"100_div_7", decimal.Decimal{Integer: 100}, decimal.Decimal{Integer: 7}, decimal.Decimal{Integer: 14, Fraction: 2857142857142857142, Digits: 19}},
		{"one_over_seven", decimal.Decimal{Integer: 1}, decimal.Decimal{Integer: 7}, decimal.Decimal{Fraction: 1428571428571428571, Digits: 19}},
		{"large_by_small_fraction", decimal.Decimal{Integer: 1000000}, decimal.Decimal{Fraction: 1, Digits: 1}, decimal.Decimal{Integer: 10000000, Digits: 19}},
		{"small_fraction_by_large", decimal.Decimal{Fraction: 1, Digits: 19}, decimal.Decimal{Integer: 1000000000}, decimal.Decimal{Fraction: 0, Digits: 19}},
		{"negative_zero_result", decimal.Decimal{Fraction: 1, Digits: 19, Negative: true}, decimal.Decimal{Integer: 1000000000000000000}, decimal.Decimal{Fraction: 0, Digits: 19}},
		{"integer_quotient_wraps", decimal.Decimal{Integer: math.MaxUint64}, decimal.Decimal{Fraction: 1, Digits: 1}, decimal.Decimal{Integer: 18446744073709551606}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decimal.Divide(tt.a, tt.b)
			if !decimal.Equal(tt.expected, result) {
				t.Errorf("Divide() = %v (%#v), want %v (%#v)", result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestDivide_DivByZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Divide by zero did not panic")
		}
	}()
	_ = decimal.Divide(decimal.Decimal{Integer: 1}, decimal.Decimal{})
}

func TestDivide_MultiplyRoundTrip(t *testing.T) {
	// a * b / b should approximately equal a (within precision limits)
	a := decimal.Decimal{Integer: 123, Fraction: 456789, Digits: 6}
	b := decimal.Decimal{Integer: 7}
	product := decimal.Multiply(a, b)
	result := decimal.Divide(product, b)
	// Compare at 6 digits of precision
	if !decimal.Equal(a.ToDigits(6), result.ToDigits(6)) {
		t.Errorf("round trip: got %v, want %v", result.ToDigits(6), a.ToDigits(6))
	}
}

func TestDivide_Int(t *testing.T) {
	tests := []struct {
		name     string
		dividend decimal.Decimal
		divisor  int
		expected decimal.Decimal
	}{
		{"one", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 1, decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}},
		{"integer", decimal.Decimal{Integer: 123}, 2, decimal.Decimal{Integer: 61, Fraction: 5, Digits: 1}},
		{"integer_with_fraction_digit", decimal.Decimal{Integer: 123, Digits: 1}, 2, decimal.Decimal{Integer: 61, Fraction: 5, Digits: 1}},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, 2, decimal.Decimal{Fraction: 615, Digits: 4}},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 2, decimal.Decimal{Integer: 61, Fraction: 728, Digits: 3}},
		{"negative", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, 2, decimal.Decimal{Integer: 61, Fraction: 728, Digits: 3, Negative: true}},
		{"large", decimal.Decimal{Integer: 1234567890, Fraction: 123456789, Digits: 10}, 2, decimal.Decimal{Integer: 617283945, Fraction: 617283945, Digits: 11}},
		{"high_digits_remainder", decimal.Decimal{Integer: 100, Fraction: 0, Digits: 19}, 7, decimal.Decimal{Integer: 14, Fraction: 2857142857142857142, Digits: 19}},
		{"max_uint64", decimal.Decimal{Integer: math.MaxUint64, Digits: 19}, 7, decimal.Decimal{Integer: math.MaxUint64 / 7, Fraction: 1428571428571428571, Digits: 19}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotient := decimal.Divide(tt.dividend, tt.divisor)
			if !decimal.Equal(tt.expected, quotient) {
				t.Errorf("Divide(Decimal, int) = %v, want %v", quotient, tt.expected)
			}
		})
	}
}

func TestDivide_WithGenerics(t *testing.T) {
	// Test that generic type parameters work
	result := decimal.Divide(decimal.Decimal{Integer: 10}, 3)
	expected := decimal.Decimal{Integer: 3, Fraction: 3333333333333333333, Digits: 19}
	if !decimal.Equal(expected, result) {
		t.Errorf("Divide(Decimal, int) = %v, want %v", result, expected)
	}

	result = decimal.Divide(10, decimal.Decimal{Integer: 3})
	if !decimal.Equal(expected, result) {
		t.Errorf("Divide(int, Decimal) = %v, want %v", result, expected)
	}
}

func TestAdd_OverflowWraps(t *testing.T) {
	result := decimal.Add(decimal.Decimal{Integer: math.MaxUint64}, decimal.Decimal{Integer: 1})
	expected := decimal.Decimal{}
	if !decimal.Equal(expected, result) {
		t.Errorf("Add(MaxUint64, 1) = %v, want %v", result, expected)
	}
}

func TestNegate(t *testing.T) {
	tests := []struct {
		name     string
		input    decimal.Decimal
		expected decimal.Decimal
	}{
		{"positive_integer", decimal.Decimal{Integer: 5}, decimal.Decimal{Integer: 5, Negative: true}},
		{"negative_integer", decimal.Decimal{Integer: 5, Negative: true}, decimal.Decimal{Integer: 5}},
		{"positive_fraction", decimal.Decimal{Fraction: 25, Digits: 2}, decimal.Decimal{Fraction: 25, Digits: 2, Negative: true}},
		{"negative_fraction", decimal.Decimal{Fraction: 25, Digits: 2, Negative: true}, decimal.Decimal{Fraction: 25, Digits: 2}},
		{"positive_mixed", decimal.Decimal{Integer: 3, Fraction: 14, Digits: 2}, decimal.Decimal{Integer: 3, Fraction: 14, Digits: 2, Negative: true}},
		{"negative_mixed", decimal.Decimal{Integer: 3, Fraction: 14, Digits: 2, Negative: true}, decimal.Decimal{Integer: 3, Fraction: 14, Digits: 2}},
		{"zero", decimal.Decimal{}, decimal.Decimal{}},
		{"zero_with_digits", decimal.Decimal{Digits: 5}, decimal.Decimal{Digits: 5}},
		{"large_value", decimal.Decimal{Integer: math.MaxUint64, Fraction: 9999999999999999999, Digits: 19}, decimal.Decimal{Integer: math.MaxUint64, Fraction: 9999999999999999999, Digits: 19, Negative: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decimal.Negate(tt.input)
			if !decimal.Equal(tt.expected, result) {
				t.Errorf("Negate() = %v (%#v), want %v (%#v)", result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestNegate_WithGenerics(t *testing.T) {
	// int
	result := decimal.Negate(5)
	expected := decimal.Decimal{Integer: 5, Negative: true}
	if !decimal.Equal(expected, result) {
		t.Errorf("Negate(int) = %v, want %v", result, expected)
	}

	// negative int
	result = decimal.Negate(-5)
	expected = decimal.Decimal{Integer: 5}
	if !decimal.Equal(expected, result) {
		t.Errorf("Negate(negative int) = %v, want %v", result, expected)
	}

	// float64
	result = decimal.Negate(3.14)
	expected = decimal.Decimal{Integer: 3, Fraction: 140000000000000128, Digits: 18, Negative: true}
	if !decimal.Equal(expected, result) {
		t.Errorf("Negate(float64) = %v, want %v", result, expected)
	}

	// uint64 zero
	result = decimal.Negate(uint64(0))
	expected = decimal.Decimal{}
	if !decimal.Equal(expected, result) {
		t.Errorf("Negate(uint64(0)) = %v, want %v", result, expected)
	}
}

func TestNegate_DoubleNegate(t *testing.T) {
	original := decimal.Decimal{Integer: 42, Fraction: 123, Digits: 3}
	result := decimal.Negate(decimal.Negate(original))
	if !decimal.Equal(original, result) {
		t.Errorf("double negate: got %v, want %v", result, original)
	}
}

func TestNegate_AddInverse(t *testing.T) {
	v := decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	result := decimal.Add(v, decimal.Negate(v))
	zero := decimal.Decimal{Digits: 1}
	if !decimal.Equal(zero, result) {
		t.Errorf("v + Negate(v) = %v, want zero", result)
	}
}

func TestAbsolute(t *testing.T) {
	tests := []struct {
		name     string
		input    decimal.Decimal
		expected decimal.Decimal
	}{
		{"positive_integer", decimal.Decimal{Integer: 5}, decimal.Decimal{Integer: 5}},
		{"negative_integer", decimal.Decimal{Integer: 5, Negative: true}, decimal.Decimal{Integer: 5}},
		{"positive_fraction", decimal.Decimal{Fraction: 25, Digits: 2}, decimal.Decimal{Fraction: 25, Digits: 2}},
		{"negative_fraction", decimal.Decimal{Fraction: 25, Digits: 2, Negative: true}, decimal.Decimal{Fraction: 25, Digits: 2}},
		{"positive_mixed", decimal.Decimal{Integer: 3, Fraction: 14, Digits: 2}, decimal.Decimal{Integer: 3, Fraction: 14, Digits: 2}},
		{"negative_mixed", decimal.Decimal{Integer: 3, Fraction: 14, Digits: 2, Negative: true}, decimal.Decimal{Integer: 3, Fraction: 14, Digits: 2}},
		{"zero", decimal.Decimal{}, decimal.Decimal{}},
		{"zero_with_digits", decimal.Decimal{Digits: 5}, decimal.Decimal{Digits: 5}},
		{"negative_zero", decimal.Decimal{Negative: true}, decimal.Decimal{}},
		{"large_value", decimal.Decimal{Integer: math.MaxUint64, Fraction: 9999999999999999999, Digits: 19, Negative: true}, decimal.Decimal{Integer: math.MaxUint64, Fraction: 9999999999999999999, Digits: 19}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decimal.Absolute(tt.input)
			if !decimal.Equal(tt.expected, result) {
				t.Errorf("Absolute() = %v (%#v), want %v (%#v)", result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestAbsolute_WithGenerics(t *testing.T) {
	// negative int
	result := decimal.Absolute(-5)
	expected := decimal.Decimal{Integer: 5}
	if !decimal.Equal(expected, result) {
		t.Errorf("Absolute(negative int) = %v, want %v", result, expected)
	}

	// positive int
	result = decimal.Absolute(5)
	expected = decimal.Decimal{Integer: 5}
	if !decimal.Equal(expected, result) {
		t.Errorf("Absolute(positive int) = %v, want %v", result, expected)
	}

	// negative float64
	result = decimal.Absolute(-3.14)
	expected = decimal.Decimal{Integer: 3, Fraction: 140000000000000128, Digits: 18}
	if !decimal.Equal(expected, result) {
		t.Errorf("Absolute(negative float64) = %v, want %v", result, expected)
	}

	// uint64 zero
	result = decimal.Absolute(uint64(0))
	expected = decimal.Decimal{}
	if !decimal.Equal(expected, result) {
		t.Errorf("Absolute(uint64(0)) = %v, want %v", result, expected)
	}
}

func TestAbsolute_Idempotent(t *testing.T) {
	v := decimal.Decimal{Integer: 42, Fraction: 123, Digits: 3, Negative: true}
	result := decimal.Absolute(decimal.Absolute(v))
	expected := decimal.Decimal{Integer: 42, Fraction: 123, Digits: 3}
	if !decimal.Equal(expected, result) {
		t.Errorf("double Absolute: got %v, want %v", result, expected)
	}
}

func TestAbsolute_NegateRelationship(t *testing.T) {
	v := decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1, Negative: true}
	absResult := decimal.Absolute(v)
	negResult := decimal.Negate(v)
	if !decimal.Equal(absResult, negResult) {
		t.Errorf("Absolute(negative) should equal Negate(negative): abs=%v, neg=%v", absResult, negResult)
	}
}

func BenchmarkDivide(b *testing.B) {
	d1 := decimal.Decimal{Integer: 355, Fraction: 113, Digits: 3}
	d2 := decimal.Decimal{Integer: 7, Fraction: 22, Digits: 2}
	for b.Loop() {
		_ = decimal.Divide(d1, d2)
	}
}
