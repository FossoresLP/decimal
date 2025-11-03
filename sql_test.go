package decimal_test

import (
	"database/sql/driver"
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestDecimal_Scan(t *testing.T) {
	tests := []struct {
		name    string
		d       decimal.Decimal
		value   any
		wantErr bool
	}{
		{"nil", decimal.Decimal{}, nil, false},
		{"string", decimal.Decimal{}, "123.123", false},
		{"bytes", decimal.Decimal{}, []byte("123.123"), false},
		{"float64", decimal.Decimal{}, 123.123, false},
		{"int64", decimal.Decimal{}, int64(123), false},
		{"int64_negative", decimal.Decimal{}, int64(-123), false},
		{"uint64", decimal.Decimal{}, uint64(123), false},
		{"invalid", decimal.Decimal{}, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.Scan(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Decimal.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkDecimal_Scan(b *testing.B) {
	d := decimal.Decimal{}
	for b.Loop() {
		_ = d.Scan("123.123")
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
