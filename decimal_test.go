package decimal_test

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestDecimal_Float64(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
		want float64
	}{
		{"zero", decimal.Decimal{}, 0},
		{"integer", decimal.Decimal{Integer: 123}, 123},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, 0.123},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, 123.123},
		{"negative", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, -123.123},
		{"indivisible", decimal.Decimal{Integer: 3, Fraction: 3333333333, Digits: 10}, 3.3333333333},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Float64(); got != tt.want {
				t.Errorf("Decimal.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_Float64(b *testing.B) {
	d := decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
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
	for i := 0; i < b.N; i++ {
		_ = d.String()
	}
}

func BenchmarkDecimal_StringPad(b *testing.B) {
	d := decimal.Decimal{Integer: 123, Fraction: 123, Digits: 6}
	for i := 0; i < b.N; i++ {
		_ = d.String()
	}
}

func BenchmarkDecimal_StringPadLong(b *testing.B) {
	d := decimal.Decimal{Integer: 12345, Fraction: 67890, Digits: 9}
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < b.N; i++ {
		_ = d.FromString("123.123")
	}
}

func BenchmarkDecimal_FromStringInt(b *testing.B) {
	d := decimal.Decimal{}
	for i := 0; i < b.N; i++ {
		_ = d.FromString("12345")
	}
}

func BenchmarkDecimal_FromStringLong(b *testing.B) {
	d := decimal.Decimal{}
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < b.N; i++ {
		_ = d.FromStringFuzzy("123.123")
	}
}

func BenchmarkDecimal_FromStringPrefix(b *testing.B) {
	d := decimal.Decimal{}
	for i := 0; i < b.N; i++ {
		_ = d.FromStringFuzzy("Angle 12.3")
	}
}

func BenchmarkDecimal_FromStringSuffix(b *testing.B) {
	d := decimal.Decimal{}
	for i := 0; i < b.N; i++ {
		_ = d.FromStringFuzzy("1234.5 lm")
	}
}

func TestDecimal_FromFloat64(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
		f    float64
	}{
		{"zero", decimal.Decimal{}, 0},
		{"integer", decimal.Decimal{}, 123},
		{"fraction", decimal.Decimal{}, 0.123},
		{"digits", decimal.Decimal{}, 123.123},
		{"negative", decimal.Decimal{}, -123.123},
		{"large", decimal.Decimal{}, 1234567890123456789.12345678901234567890},
		{"small", decimal.Decimal{}, 0.000000000000000001},
		{"negative_small", decimal.Decimal{}, -0.000000000000000001},
		{"indivisible", decimal.Decimal{}, 3.3333333333},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.FromFloat64(tt.f)
			if got := tt.d.Float64(); got != tt.f {
				t.Errorf("Decimal.FromFloat64() = %v, want %v", got, tt.f)
			}
		})
	}
}

func BenchmarkDecimal_FromFloat64(b *testing.B) {
	d := decimal.Decimal{}
	for i := 0; i < b.N; i++ {
		d.FromFloat64(123.123)
	}
}

func BenchmarkDecimal_FromFloat64Long(b *testing.B) {
	d := decimal.Decimal{}
	for i := 0; i < b.N; i++ {
		d.FromFloat64(12345.67890)
	}
}

func TestDecimal_Scan(t *testing.T) {
	tests := []struct {
		name    string
		d       decimal.Decimal
		value   interface{}
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
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < b.N; i++ {
		_, _ = d.Value()
	}
}

func TestDecimal_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		d       decimal.Decimal
		want    []byte
		wantErr bool
	}{
		{"zero", decimal.Decimal{}, []byte("0"), false},
		{"integer", decimal.Decimal{Integer: 123}, []byte("123"), false},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, []byte("0.123"), false},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, []byte("123.123"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Decimal.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("Decimal.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

// json.Marshal() does tend to parse the result of MarshalJSON() to reformat if necessary
// This test is to ensure that the value is not changed e.g. by being stored as a float and then converted back to a string, having trailing zeros removed, quotes added, etc.
func TestDecimal_MarshalJSON_JSON(t *testing.T) {
	d := decimal.Decimal{Integer: 123, Fraction: 12345, Digits: 5}
	res, err := json.Marshal(d)
	if err != nil {
		t.Error(err)
	}
	if string(res) != "123.12345" {
		t.Errorf("Decimal.MarshalJSON() = %v, want %v", string(res), "123.12345")
	}
	d.Fraction = 1234500
	d.Digits = 7
	res, err = json.Marshal(d)
	if err != nil {
		t.Error(err)
	}
	if string(res) != "123.1234500" {
		t.Errorf("Decimal.MarshalJSON() = %v, want %v", string(res), "123.1234500")
	}
}

func BenchmarkDecimal_MarshalJSON(b *testing.B) {
	d := decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}
	for i := 0; i < b.N; i++ {
		_, _ = d.MarshalJSON()
	}
}

func TestDecimal_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		d       decimal.Decimal
		data    []byte
		wantErr bool
	}{
		{"zero", decimal.Decimal{}, []byte("0.0"), false},
		{"integer", decimal.Decimal{}, []byte("123.0"), false},
		{"fraction", decimal.Decimal{}, []byte("0.123"), false},
		{"digits", decimal.Decimal{}, []byte("123.123"), false},
		{"invalid", decimal.Decimal{}, []byte("123.123.123"), true},
		{"null", decimal.Decimal{}, []byte("null"), true},
		{"true", decimal.Decimal{}, []byte("true"), true},
		{"false", decimal.Decimal{}, []byte("false"), true},
		{"empty_array", decimal.Decimal{}, []byte("[]"), true},
		{"empty_object", decimal.Decimal{}, []byte("{}"), true},
		{"array", decimal.Decimal{}, []byte("[123.123]"), true},
		{"object", decimal.Decimal{}, []byte(`{"value":123.123}`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.UnmarshalJSON(tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Decimal.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecimal_UnmarshalJSON_struct(t *testing.T) {
	type S struct {
		D decimal.Decimal `json:"d"`
	}
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		value   float64
	}{
		{"zero", []byte(`{"d":0.0}`), false, 0.0},
		{"integer", []byte(`{"d":123.0}`), false, 123.0},
		{"fraction", []byte(`{"d":0.123}`), false, 0.123},
		{"digits", []byte(`{"d":123.123}`), false, 123.123},
		{"invalid", []byte(`{"d":123.123.123}`), true, 0.0},
		{"null", []byte(`{"d":null}`), true, 0.0},
		{"true", []byte(`{"d":true}`), true, 0.0},
		{"false", []byte(`{"d":false}`), true, 0.0},
		{"empty_array", []byte(`{"d":[]}`), true, 0.0},
		{"empty_object", []byte(`{"d":{}}`), true, 0.0},
		{"array", []byte(`{"d":[123.123]}`), true, 0.0},
		{"object", []byte(`{"d":{"value":123.123}}`), true, 0.0},
		{"linebreaks", []byte("{\n\t\"d\": 123.123\n}\n"), false, 123.123},
		{"quote", []byte("{\"d\": \"123.123\"}"), false, 123.123},
		{"mismatched_quote", []byte("{\"d\": \"123.123}"), true, 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s S
			if err := json.Unmarshal(tt.data, &s); (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if s.D.Float64() != tt.value {
				t.Errorf("json.Unmarshal() = %v, want %v", s.D.Float64(), tt.value)
			}
		})
	}
}

func BenchmarkDecimal_UnmarshalJSON(b *testing.B) {
	d := decimal.Decimal{}
	for i := 0; i < b.N; i++ {
		_ = d.UnmarshalJSON([]byte("123.123"))
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
