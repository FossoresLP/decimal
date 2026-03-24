package decimal_test

import (
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestEqual(t *testing.T) {
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
			if got := decimal.Equal(tt.d1, tt.d2); got != tt.want {
				t.Errorf("Decimal.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkEqual(b *testing.B) {
	for b.Loop() {
		_ = decimal.Equal(decimal.Decimal{}, decimal.Decimal{})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name string
		d1   decimal.Decimal
		d2   decimal.Decimal
		want int
	}{
		// Equal values
		{"equal_zero", decimal.Decimal{}, decimal.Decimal{}, 0},
		{"equal_integer", decimal.Decimal{Integer: 5}, decimal.Decimal{Integer: 5}, 0},
		{"equal_fraction", decimal.Decimal{Fraction: 5, Digits: 1}, decimal.Decimal{Fraction: 5, Digits: 1}, 0},
		{"equal_full", decimal.Decimal{Integer: 3, Fraction: 14, Digits: 2}, decimal.Decimal{Integer: 3, Fraction: 14, Digits: 2}, 0},
		{"equal_negative", decimal.Decimal{Integer: 5, Negative: true}, decimal.Decimal{Integer: 5, Negative: true}, 0},
		{"equal_different_digits", decimal.Decimal{Fraction: 1, Digits: 1}, decimal.Decimal{Fraction: 10, Digits: 2}, 0},

		// Different signs
		{"positive_vs_negative", decimal.Decimal{Integer: 1}, decimal.Decimal{Integer: 1, Negative: true}, 1},
		{"negative_vs_positive", decimal.Decimal{Integer: 1, Negative: true}, decimal.Decimal{Integer: 1}, -1},
		{"zero_vs_negative", decimal.Decimal{}, decimal.Decimal{Integer: 5, Negative: true}, 1},
		{"negative_vs_zero", decimal.Decimal{Integer: 5, Negative: true}, decimal.Decimal{}, -1},

		// Positive comparisons - integer part differs
		{"greater_integer", decimal.Decimal{Integer: 10}, decimal.Decimal{Integer: 5}, 1},
		{"lesser_integer", decimal.Decimal{Integer: 5}, decimal.Decimal{Integer: 10}, -1},

		// Positive comparisons - fraction part differs
		{"greater_fraction", decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}, decimal.Decimal{Integer: 1, Fraction: 3, Digits: 1}, 1},
		{"lesser_fraction", decimal.Decimal{Integer: 1, Fraction: 3, Digits: 1}, decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}, -1},

		// Negative comparisons (sign inverts the result)
		{"neg_greater_abs", decimal.Decimal{Integer: 10, Negative: true}, decimal.Decimal{Integer: 5, Negative: true}, -1},
		{"neg_lesser_abs", decimal.Decimal{Integer: 5, Negative: true}, decimal.Decimal{Integer: 10, Negative: true}, 1},
		{"neg_greater_fraction", decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1, Negative: true}, decimal.Decimal{Integer: 1, Fraction: 3, Digits: 1, Negative: true}, -1},
		{"neg_lesser_fraction", decimal.Decimal{Integer: 1, Fraction: 3, Digits: 1, Negative: true}, decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1, Negative: true}, 1},

		// Edge cases
		{"fraction_only_greater", decimal.Decimal{Fraction: 9, Digits: 1}, decimal.Decimal{Fraction: 1, Digits: 1}, 1},
		{"fraction_only_lesser", decimal.Decimal{Fraction: 1, Digits: 1}, decimal.Decimal{Fraction: 9, Digits: 1}, -1},
		{"integer_vs_fraction", decimal.Decimal{Integer: 1}, decimal.Decimal{Fraction: 9, Digits: 1}, 1},
		{"fraction_vs_integer", decimal.Decimal{Fraction: 9, Digits: 1}, decimal.Decimal{Integer: 1}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decimal.Compare(tt.d1, tt.d2); got != tt.want {
				t.Errorf("Compare(%v, %v) = %v, want %v", tt.d1, tt.d2, got, tt.want)
			}
		})
	}
}

func BenchmarkCompare(b *testing.B) {
	for b.Loop() {
		_ = decimal.Compare(decimal.Decimal{Integer: 123, Fraction: 456, Negative: true}, decimal.Decimal{Integer: 123, Fraction: 789, Negative: true})
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
