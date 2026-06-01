package decimal_test

import (
	"database/sql/driver"
	"math"
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestDecimal_Scan(t *testing.T) {
	sentinel := decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	negative := decimal.Decimal{Negative: true, Integer: 9}
	tests := []struct {
		name    string
		value   any
		initial decimal.Decimal
		want    decimal.Decimal
		wantErr bool
	}{
		{"nil", nil, sentinel, decimal.Decimal{}, false},
		{"string", "123.123", sentinel, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"string_overflow", "18446744073709551616", sentinel, sentinel, true},
		{"bytes", []byte("123.123"), sentinel, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"bytes_maxuint64", []byte("18446744073709551615"), sentinel, decimal.Decimal{Integer: 18446744073709551615}, false},
		{"bytes_maxuint64_plus_one", []byte("18446744073709551616"), sentinel, sentinel, true},
		{"bytes_bad", []byte("bad"), sentinel, sentinel, true},
		{"float64", 123.123, sentinel, decimal.New(123.123), false},
		{"float64_negative_subnormal", -5e-324, sentinel, decimal.Decimal{}, false},
		{"int64", int64(123), sentinel, decimal.Decimal{Integer: 123}, false},
		{"int64_negative", int64(-123), sentinel, decimal.Decimal{Integer: 123, Negative: true}, false},
		{"int64_min", int64(math.MinInt64), sentinel, decimal.Decimal{Integer: uint64(math.MaxInt64) + 1, Negative: true}, false},
		{"int64_clears_negative", int64(5), negative, decimal.Decimal{Integer: 5}, false},
		{"uint64", uint64(123), sentinel, decimal.Decimal{Integer: 123}, false},
		{"uint64_clears_negative", uint64(5), negative, decimal.Decimal{Integer: 5}, false},
		{"invalid", true, sentinel, sentinel, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.initial
			if err := d.Scan(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Decimal.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if d != tt.want {
				t.Errorf("Decimal.Scan() = %#v, want %#v", d, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_Scan(b *testing.B) {
	d := decimal.Decimal{}
	var value any = []byte("123.123")
	for b.Loop() {
		_ = d.Scan(value)
	}
}

func TestDecimal_Value(t *testing.T) {
	tests := []struct {
		name    string
		d       decimal.Decimal
		want    driver.Value
		wantErr bool
	}{
		{"zero", decimal.Decimal{}, "0", false},
		{"integer", decimal.Decimal{Integer: 123}, "123", false},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, "0.123", false},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, "123.123", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Decimal.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decimal.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_Value(b *testing.B) {
	d := decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}
	for b.Loop() {
		_, _ = d.Value()
	}
}

func TestFixed_Scan(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		initial decimal.Fixed
		want    decimal.Fixed
		wantErr bool
	}{
		{"nil", nil, 12345, 0, false},
		{"string", "123.45", 0, 12345, false},
		{"string_negative", "-123.45", 0, -12345, false},
		{"bytes", []byte("123.45"), 0, 12345, false},
		{"bytes_truncated", []byte("123.456"), 12345, 12345, true},
		{"string_invalid", "abc", 12345, 12345, true},
		{"float64", 123.45, 0, 12345, false},
		{"float64_negative", -123.45, 0, -12345, false},
		{"float64_zero", 0.0, 12345, 0, false},
		{"float64_negative_fraction", -0.05, 12345, -5, false},
		{"float64_large_magnitude", 21474800.01, 0, 2147480001, false},
		{"float64_max", 21474836.47, 0, 2147483647, false},
		{"float64_min", -21474836.48, 0, -2147483648, false},
		{"float64_million", 1000000.07, 0, 100000007, false},
		{"float64_subcent", 1.005, 12345, 12345, true},
		{"float64_too_many_digits", 12345.678, 12345, 12345, true},
		{"float64_subcent_small", 0.001, 12345, 12345, true},
		{"float64_overflow_max", 21474836.48, 12345, 12345, true},
		{"float64_overflow", 21474837.0, 12345, 12345, true},
		{"float64_overflow_min", -21474836.49, 12345, 12345, true},
		{"float64_nan", math.NaN(), 12345, 12345, true},
		{"float64_pos_inf", math.Inf(1), 12345, 12345, true},
		{"float64_neg_inf", math.Inf(-1), 12345, 12345, true},
		{"int64", int64(123), 0, 12300, false},
		{"int64_negative", int64(-123), 0, -12300, false},
		{"uint64", uint64(123), 0, 12300, false},
		{"invalid", true, 12345, 12345, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.initial
			err := f.Scan(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fixed.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if f != tt.want {
				t.Errorf("Fixed.Scan() = %d, want %d", int32(f), int32(tt.want))
			}
		})
	}
}

func BenchmarkFixed_Scan(b *testing.B) {
	var f decimal.Fixed
	var value any = []byte("123.45")
	for b.Loop() {
		_ = f.Scan(value)
	}
}

func TestFixed_Value(t *testing.T) {
	tests := []struct {
		name    string
		f       decimal.Fixed
		want    driver.Value
		wantErr bool
	}{
		{"zero", 0, "0.00", false},
		{"fraction", 5, "0.05", false},
		{"integer", 12300, "123.00", false},
		{"digits", 12345, "123.45", false},
		{"negative_digits", -12345, "-123.45", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Fixed.Value() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Fixed.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkFixed_Value(b *testing.B) {
	f := decimal.Fixed(12345)
	for b.Loop() {
		_, _ = f.Value()
	}
}
