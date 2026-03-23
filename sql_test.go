package decimal_test

import (
	"database/sql/driver"
	"math"
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestDecimal_Scan(t *testing.T) {
	tests := []struct {
		name    string
		d       decimal.Decimal
		value   any
		want    decimal.Decimal
		wantErr bool
	}{
		{"nil", decimal.Decimal{}, nil, decimal.Decimal{}, false},
		{"string", decimal.Decimal{}, "123.123", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"bytes", decimal.Decimal{}, []byte("123.123"), decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"bytes_maxuint64", decimal.Decimal{}, []byte("18446744073709551615"), decimal.Decimal{Integer: 18446744073709551615}, false},
		{"bytes_maxuint64_plus_one", decimal.Decimal{}, []byte("18446744073709551616"), decimal.Decimal{}, true},
		{"float64", decimal.Decimal{}, 123.123, decimal.New(123.123), false},
		{"float64_negative_subnormal", decimal.Decimal{}, -5e-324, decimal.Decimal{}, false},
		{"int64", decimal.Decimal{}, int64(123), decimal.Decimal{Integer: 123}, false},
		{"int64_negative", decimal.Decimal{}, int64(-123), decimal.Decimal{Integer: 123, Negative: true}, false},
		{"int64_min", decimal.Decimal{}, int64(math.MinInt64), decimal.Decimal{Integer: uint64(math.MaxInt64) + 1, Negative: true}, false},
		{"uint64", decimal.Decimal{}, uint64(123), decimal.Decimal{Integer: 123}, false},
		{"invalid", decimal.Decimal{}, true, decimal.Decimal{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.Scan(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Decimal.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.d != tt.want {
				t.Errorf("Decimal.Scan() = %#v, want %#v", tt.d, tt.want)
			}
		})
	}
}

func TestDecimal_Scan_NilResetsReceiver(t *testing.T) {
	d := decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	if err := d.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if d != (decimal.Decimal{}) {
		t.Errorf("Scan(nil) did not reset receiver: got %+v, want zero value", d)
	}
}

func TestDecimal_Scan_ResetsNegative(t *testing.T) {
	d := decimal.Decimal{Negative: true, Integer: 9}
	if err := d.Scan(int64(5)); err != nil {
		t.Fatal(err)
	}
	if d.Negative {
		t.Error("Negative should be false after scanning positive int64")
	}
	if d.Integer != 5 {
		t.Errorf("Integer = %d, want 5", d.Integer)
	}

	d = decimal.Decimal{Negative: true, Integer: 9}
	if err := d.Scan(uint64(5)); err != nil {
		t.Fatal(err)
	}
	if d.Negative {
		t.Error("Negative should be false after scanning uint64")
	}
	if d.Integer != 5 {
		t.Errorf("Integer = %d, want 5", d.Integer)
	}
}

func TestDecimal_Scan_ErrorKeepsReceiver(t *testing.T) {
	d := decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	if err := d.Scan([]byte("bad")); err == nil {
		t.Fatal("Scan([]byte(\"bad\")) error = nil, want non-nil")
	}
	if d != (decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}) {
		t.Errorf("Scan([]byte(\"bad\")) changed receiver: got %#v", d)
	}

	d = decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	if err := d.Scan("18446744073709551616"); err == nil {
		t.Fatal("Scan(overflow string) error = nil, want non-nil")
	}
	if d != (decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}) {
		t.Errorf("Scan(overflow string) changed receiver: got %#v", d)
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
