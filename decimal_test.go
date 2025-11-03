package decimal_test

import (
	"reflect"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d1.Equal(&tt.d2); got != tt.want {
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
		{"integer", decimal.Decimal{Integer: 123}, 2, &decimal.Decimal{Integer: 246}},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, 2, &decimal.Decimal{Fraction: 246, Digits: 3}},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 2, &decimal.Decimal{Integer: 246, Fraction: 912, Digits: 3}},
		{"negative", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, 2, &decimal.Decimal{Integer: 246, Fraction: 912, Digits: 3, Negative: true}},
		{"large", decimal.Decimal{Integer: 1234567890123456789, Fraction: 1234567890123456789, Digits: 19}, 2, &decimal.Decimal{Integer: 2469135780246913578, Fraction: 2469135780246913578, Digits: 19}},
		{"large_multiplier", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 1234567890, &decimal.Decimal{Integer: 152414813427, Fraction: 840, Digits: 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.decimal.MultiplyUint64(tt.multiplier), tt.expected) {
				t.Errorf("Decimal.MultiplyUint64() = %v, want %v", tt.decimal, tt.expected)
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
		{"integer_with_fraction_digit", decimal.Decimal{Integer: 123, Digits: 1}, 2, &decimal.Decimal{Integer: 61, Fraction: 5, Digits: 1}},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, 2, &decimal.Decimal{Fraction: 61, Digits: 3}},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 2, &decimal.Decimal{Integer: 61, Fraction: 728, Digits: 3}},
		{"negative", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, 2, &decimal.Decimal{Integer: 61, Fraction: 728, Digits: 3, Negative: true}},
		{"large", decimal.Decimal{Integer: 1234567890, Fraction: 123456789, Digits: 10}, 2, &decimal.Decimal{Integer: 617283945, Fraction: 61728394, Digits: 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.decimal.DivideUint64(tt.divisor), tt.expected) {
				t.Errorf("Decimal.DivideUint64() = %v, want %v", tt.decimal, tt.expected)
			}
		})
	}
}

func TestDecimal_Add(t *testing.T) {
	tests := []struct {
		name     string
		d1       decimal.Decimal
		d2       decimal.Decimal
		expected *decimal.Decimal
	}{
		{"zero", decimal.Decimal{}, decimal.Decimal{}, &decimal.Decimal{}},
		{"integer", decimal.Decimal{Integer: 123}, decimal.Decimal{Integer: 456}, &decimal.Decimal{Integer: 579}},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, decimal.Decimal{Fraction: 456, Digits: 3}, &decimal.Decimal{Fraction: 579, Digits: 3}},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 3}, &decimal.Decimal{Integer: 579, Fraction: 579, Digits: 3}},
		{"negative_d1", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 3}, &decimal.Decimal{Integer: 333, Fraction: 333, Digits: 3}},
		{"negative_d2", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 3, Negative: true}, &decimal.Decimal{Integer: 333, Fraction: 333, Digits: 3, Negative: true}},
		{"negative_both", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 3, Negative: true}, &decimal.Decimal{Integer: 579, Fraction: 579, Digits: 3, Negative: true}},
		{"large", decimal.Decimal{Integer: 1234567890, Fraction: 123456789, Digits: 10}, decimal.Decimal{Integer: 9876543210, Fraction: 987654321, Digits: 10}, &decimal.Decimal{Integer: 11111111100, Fraction: 1111111110, Digits: 10}},
		{"different_digits_d1_less", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 4}, &decimal.Decimal{Integer: 579, Fraction: 1686, Digits: 4}},
		{"different_digits_d2_less", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 4}, decimal.Decimal{Integer: 456, Fraction: 456, Digits: 3}, &decimal.Decimal{Integer: 579, Fraction: 4683, Digits: 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.d1.Add(&tt.d2).Equal(tt.expected) {
				t.Errorf("Decimal.Add() = %v, want %v", tt.d1, tt.expected)
			}
		})
	}
}

func TestNew(t *testing.T) {
	d := decimal.New(123.123).ToDigits(3)
	if d.Integer != 123 || d.Fraction != 123 || d.Digits != 3 || d.Negative {
		t.Errorf("New() = %v, want %v", d, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3})
	}
	d = decimal.New(-123.123).ToDigits(3)
	if d.Integer != 123 || d.Fraction != 123 || d.Digits != 3 || !d.Negative {
		t.Errorf("New() = %v, want %v", d, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true})
	}
	d = decimal.New(123)
	if d.Integer != 123 || d.Fraction != 0 || d.Digits != 0 || d.Negative {
		t.Errorf("New() = %v, want %v", d, decimal.Decimal{Integer: 123})
	}
	d = decimal.New(-123)
	if d.Integer != 123 || d.Fraction != 0 || d.Digits != 0 || !d.Negative {
		t.Errorf("New() = %v, want %v", d, decimal.Decimal{Integer: 123, Negative: true})
	}
}

func TestDecimal_Clone(t *testing.T) {
	d := decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}
	c := d.Clone()
	if !reflect.DeepEqual(c, &d) {
		t.Errorf("Decimal.Clone() = %v, want %v", c, d)
	}
	if c == &d {
		t.Errorf("Decimal.Clone() did not return a copy")
	}
}

func TestDecimal_ToDigits(t *testing.T) {
	tests := []struct {
		name   string
		d      decimal.Decimal
		digits uint8
		want   *decimal.Decimal
	}{
		{"zero", decimal.Decimal{}, 3, &decimal.Decimal{Digits: 3}},
		{"integer", decimal.Decimal{Integer: 123}, 3, &decimal.Decimal{Integer: 123, Digits: 3}},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, 6, &decimal.Decimal{Fraction: 123000, Digits: 6}},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, 6, &decimal.Decimal{Integer: 123, Fraction: 123000, Digits: 6}},
		{"negative", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, 6, &decimal.Decimal{Integer: 123, Fraction: 123000, Digits: 6, Negative: true}},
		{"less_low", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, 2, &decimal.Decimal{Integer: 123, Fraction: 12, Digits: 2}},
		{"less_high", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}, 2, &decimal.Decimal{Integer: 123, Fraction: 45, Digits: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.ToDigits(tt.digits); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decimal.ToDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_Zero(t *testing.T) {
	if !reflect.DeepEqual(decimal.Zero(), &decimal.Decimal{}) {
		t.Errorf("Decimal.Zero() = %v, want %v", decimal.Zero(), &decimal.Decimal{})
	}
}
