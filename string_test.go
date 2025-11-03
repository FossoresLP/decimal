package decimal_test

import (
	"reflect"
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestDecimal_NewFromString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    *decimal.Decimal
		wantErr bool
	}{
		{"zero", "0.0", &decimal.Decimal{Digits: 1}, false},
		{"integer", "123.0", &decimal.Decimal{Integer: 123, Digits: 1}, false},
		{"fraction", "0.123", &decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"digits", "123.123", &decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"negative", "-123.123", &decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decimal.NewFromString(tt.s)
			if err != nil && !tt.wantErr {
				t.Errorf("Decimal.NewString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decimal.NewString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_NewFromStringFuzzy(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    *decimal.Decimal
		wantErr bool
	}{
		{"zero", "0.0", &decimal.Decimal{Digits: 1}, false},
		{"integer", "123.0", &decimal.Decimal{Integer: 123, Digits: 1}, false},
		{"fraction", "0.123", &decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"digits", "123.123", &decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"negative", "-123.123", &decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, false},
		{"prefix", "prefix123.123", &decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"suffix", "123.123suffix", &decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decimal.NewFromStringFuzzy(tt.s)
			if err != nil && !tt.wantErr {
				t.Errorf("Decimal.NewString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decimal.NewString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_String(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
		want string
	}{
		{"zero", decimal.Decimal{}, "0"},
		{"integer", decimal.Decimal{Integer: 123}, "123"},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, "0.123"},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, "123.123"},
		{"fraction_zero", decimal.Decimal{Integer: 123, Digits: 3}, "123.000"},
		{"fraction_leading_zero", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 4}, "123.0123"},
		{"fraction_trailing_zero", decimal.Decimal{Integer: 123, Fraction: 120, Digits: 3}, "123.120"},
		{"negative", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, "-123.123"},
		{"negative_small", decimal.Decimal{Integer: 0, Fraction: 123, Digits: 3, Negative: true}, "-0.123"},
		{"negative_integer", decimal.Decimal{Integer: 123, Digits: 0, Negative: true}, "-123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("Decimal.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_String(b *testing.B) {
	d := decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}
	for b.Loop() {
		_ = d.String()
	}
}

func BenchmarkDecimal_StringPad(b *testing.B) {
	d := decimal.Decimal{Integer: 123, Fraction: 123, Digits: 6}
	for b.Loop() {
		_ = d.String()
	}
}

func BenchmarkDecimal_StringPadLong(b *testing.B) {
	d := decimal.Decimal{Integer: 12345, Fraction: 67890, Digits: 9}
	for b.Loop() {
		_ = d.String()
	}
}

func TestDecimal_FromString(t *testing.T) {
	tests := []struct {
		name    string
		d       *decimal.Decimal
		s       string
		wantErr bool
	}{
		{"zero", &decimal.Decimal{Digits: 1}, "0.0", false},
		{"integer", &decimal.Decimal{Integer: 123, Digits: 1}, "123.0", false},
		{"fraction", &decimal.Decimal{Fraction: 123, Digits: 3}, "0.123", false},
		{"digits", &decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, "123.123", false},
		{"invalid", nil, "123.123.123", true},
		{"leading_dot", &decimal.Decimal{Fraction: 123, Digits: 3}, ".123", false},
		{"trailing_dot", &decimal.Decimal{Integer: 123}, "123.", false},
		{"empty", &decimal.Decimal{}, "", false},
		{"negative", &decimal.Decimal{Negative: true, Integer: 123, Fraction: 123, Digits: 3}, "-123.123", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := decimal.Decimal{}
			err := d.FromString(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decimal.FromString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.d != nil && !reflect.DeepEqual(d, *tt.d) {
				t.Errorf("Decimal.FromString() = %v, want %v", d, *tt.d)
			}
		})
	}
}

func BenchmarkDecimal_FromString(b *testing.B) {
	d := decimal.Decimal{}
	for b.Loop() {
		_ = d.FromString("123.123")
	}
}

func BenchmarkDecimal_FromStringInt(b *testing.B) {
	d := decimal.Decimal{}
	for b.Loop() {
		_ = d.FromString("12345")
	}
}

func BenchmarkDecimal_FromStringLong(b *testing.B) {
	d := decimal.Decimal{}
	for b.Loop() {
		_ = d.FromString("12345.67890")
	}
}

func TestDecimal_FromStringFuzzy(t *testing.T) {
	tests := []struct {
		name    string
		d       *decimal.Decimal
		s       string
		wantErr bool
	}{
		{"zero", &decimal.Decimal{Digits: 1}, "0.0", false},
		{"integer", &decimal.Decimal{Integer: 123, Digits: 1}, "123.0", false},
		{"fraction", &decimal.Decimal{Fraction: 123, Digits: 3}, "0.123", false},
		{"digits", &decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, "123.123", false},
		{"leading_dot", &decimal.Decimal{Fraction: 123, Digits: 3}, ".123", false},
		{"trailing_dot", &decimal.Decimal{Integer: 123}, "123.", false},
		{"empty", &decimal.Decimal{}, "", false},
		{"negative", &decimal.Decimal{Negative: true, Integer: 123, Fraction: 123, Digits: 3}, "-123.123", false},
		{"prefix", &decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, "prefix123.123", false},
		{"suffix", &decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, "123.123suffix", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := decimal.Decimal{}
			err := d.FromStringFuzzy(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decimal.FromString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.d != nil && !reflect.DeepEqual(d, *tt.d) {
				t.Errorf("Decimal.FromString() = %v, want %v", d, *tt.d)
			}
		})
	}
}

func BenchmarkDecimal_FromStringFuzzy(b *testing.B) {
	d := decimal.Decimal{}
	for b.Loop() {
		_ = d.FromStringFuzzy("123.123")
	}
}

func BenchmarkDecimal_FromStringPrefix(b *testing.B) {
	d := decimal.Decimal{}
	for b.Loop() {
		_ = d.FromStringFuzzy("Angle 12.3")
	}
}

func BenchmarkDecimal_FromStringSuffix(b *testing.B) {
	d := decimal.Decimal{}
	for b.Loop() {
		_ = d.FromStringFuzzy("1234.5 lm")
	}
}
