package decimal_test

import (
	"encoding/json"
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestDecimal_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
		want string
	}{
		{"zero", decimal.Decimal{}, "0"},
		{"integer", decimal.Decimal{Integer: 123}, "123"},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, "0.123"},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, "123.123"},
		{"negative_integer", decimal.Decimal{Integer: 123, Negative: true}, "-123"},
		{"negative_fraction", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, "-123.456"},
		{"max_uint64_integer", decimal.Decimal{Integer: ^uint64(0)}, "18446744073709551615"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.MarshalJSON()
			if err != nil {
				t.Fatalf("Decimal.MarshalJSON() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("Decimal.MarshalJSON() = %v, want %v", string(got), tt.want)
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
	for b.Loop() {
		_, _ = d.MarshalJSON()
	}
}

func TestDecimal_UnmarshalJSON(t *testing.T) {
	sentinel := decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	tests := []struct {
		name    string
		data    []byte
		initial decimal.Decimal
		want    decimal.Decimal
		wantErr bool
	}{
		{"zero", []byte("0.0"), sentinel, decimal.Decimal{Digits: 1}, false},
		{"integer", []byte("123.0"), sentinel, decimal.Decimal{Integer: 123, Digits: 1}, false},
		{"fraction", []byte("0.123"), sentinel, decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"digits", []byte("123.123"), sentinel, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"quote", []byte(`"123.123"`), sentinel, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"null", []byte("null"), sentinel, decimal.Decimal{}, false},
		{"overflow_string", []byte(`"18446744073709551616"`), sentinel, sentinel, true},
		{"bad", []byte("bad"), sentinel, sentinel, true},
		{"invalid", []byte("123.123.123"), sentinel, sentinel, true},
		{"true", []byte("true"), sentinel, sentinel, true},
		{"false", []byte("false"), sentinel, sentinel, true},
		{"empty_array", []byte("[]"), sentinel, sentinel, true},
		{"empty_object", []byte("{}"), sentinel, sentinel, true},
		{"array", []byte("[123.123]"), sentinel, sentinel, true},
		{"object", []byte(`{"value":123.123}`), sentinel, sentinel, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.initial
			if err := d.UnmarshalJSON(tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Decimal.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if d != tt.want {
				t.Errorf("Decimal.UnmarshalJSON() = %#v, want %#v", d, tt.want)
			}
		})
	}
}

func TestDecimal_UnmarshalJSON_struct(t *testing.T) {
	type S struct {
		D decimal.Decimal `json:"d"`
	}
	sentinel := decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	tests := []struct {
		name    string
		data    []byte
		initial decimal.Decimal
		want    decimal.Decimal
		wantErr bool
	}{
		{"zero", []byte(`{"d":0.0}`), sentinel, decimal.Decimal{Digits: 1}, false},
		{"integer", []byte(`{"d":123.0}`), sentinel, decimal.Decimal{Integer: 123, Digits: 1}, false},
		{"fraction", []byte(`{"d":0.123}`), sentinel, decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"digits", []byte(`{"d":123.123}`), sentinel, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"null", []byte(`{"d":null}`), sentinel, decimal.Decimal{}, false},
		{"maxuint64_string", []byte(`{"d":"18446744073709551615"}`), sentinel, decimal.Decimal{Integer: 18446744073709551615}, false},
		{"maxuint64_plus_one_string", []byte(`{"d":"18446744073709551616"}`), sentinel, sentinel, true},
		{"invalid", []byte(`{"d":123.123.123}`), sentinel, sentinel, true},
		{"true", []byte(`{"d":true}`), sentinel, sentinel, true},
		{"false", []byte(`{"d":false}`), sentinel, sentinel, true},
		{"empty_array", []byte(`{"d":[]}`), sentinel, sentinel, true},
		{"empty_object", []byte(`{"d":{}}`), sentinel, sentinel, true},
		{"array", []byte(`{"d":[123.123]}`), sentinel, sentinel, true},
		{"object", []byte(`{"d":{"value":123.123}}`), sentinel, sentinel, true},
		{"linebreaks", []byte("{\n\t\"d\": 123.123\n}\n"), sentinel, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"quote", []byte("{\"d\": \"123.123\"}"), sentinel, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"mismatched_quote", []byte("{\"d\": \"123.123}"), sentinel, sentinel, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := S{D: tt.initial}
			if err := json.Unmarshal(tt.data, &s); (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if s.D != tt.want {
				t.Errorf("json.Unmarshal() = %#v, want %#v", s.D, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_UnmarshalJSON(b *testing.B) {
	d := decimal.Decimal{}
	data := []byte("123.123")
	for b.Loop() {
		_ = d.UnmarshalJSON(data)
	}
}

func TestFixed_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		f    decimal.Fixed
		want string
	}{
		{"zero", 0, "0.00"},
		{"fraction", 5, "0.05"},
		{"fraction_full", 99, "0.99"},
		{"one", 100, "1.00"},
		{"integer", 12300, "123.00"},
		{"digits", 12345, "123.45"},
		{"max_int32", 2147483647, "21474836.47"},
		{"negative_fraction", -5, "-0.05"},
		{"negative_integer", -500, "-5.00"},
		{"negative_digits", -12345, "-123.45"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.MarshalJSON()
			if err != nil {
				t.Fatalf("Fixed.MarshalJSON() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("Fixed.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

// json.Marshal() reformats the result of MarshalJSON() if necessary.
// This test ensures the value passes through unchanged, keeping its trailing zeros and not gaining quotes.
func TestFixed_MarshalJSON_JSON(t *testing.T) {
	res, err := json.Marshal(decimal.Fixed(12300))
	if err != nil {
		t.Error(err)
	}
	if string(res) != "123.00" {
		t.Errorf("json.Marshal() = %v, want %v", string(res), "123.00")
	}
}

func BenchmarkFixed_MarshalJSON(b *testing.B) {
	f := decimal.Fixed(12345)
	for b.Loop() {
		_, _ = f.MarshalJSON()
	}
}

func TestFixed_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		initial decimal.Fixed
		want    decimal.Fixed
		wantErr bool
	}{
		{"zero", []byte("0.0"), 12345, 0, false},
		{"integer", []byte("123"), 0, 12300, false},
		{"fraction", []byte("0.12"), 0, 12, false},
		{"digits", []byte("123.45"), 0, 12345, false},
		{"negative", []byte("-123.45"), 0, -12345, false},
		{"quote", []byte(`"123.45"`), 0, 12345, false},
		{"null", []byte("null"), 12345, 0, false},
		{"truncated", []byte("123.456"), 12345, 12345, true},
		{"invalid", []byte("123.123.123"), 12345, 12345, true},
		{"true", []byte("true"), 12345, 12345, true},
		{"empty_array", []byte("[]"), 12345, 12345, true},
		{"empty_object", []byte("{}"), 12345, 12345, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.initial
			err := f.UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fixed.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if f != tt.want {
				t.Errorf("Fixed.UnmarshalJSON() = %d, want %d", int32(f), int32(tt.want))
			}
		})
	}
}

func TestFixed_UnmarshalJSON_struct(t *testing.T) {
	type S struct {
		F decimal.Fixed `json:"f"`
	}
	tests := []struct {
		name    string
		data    []byte
		initial decimal.Fixed
		want    decimal.Fixed
		wantErr bool
	}{
		{"zero", []byte(`{"f":0.0}`), 12345, 0, false},
		{"integer", []byte(`{"f":123}`), 0, 12300, false},
		{"digits", []byte(`{"f":123.45}`), 0, 12345, false},
		{"negative", []byte(`{"f":-123.45}`), 0, -12345, false},
		{"null", []byte(`{"f":null}`), 12345, 0, false},
		{"quote", []byte(`{"f":"123.45"}`), 0, 12345, false},
		{"truncated", []byte(`{"f":123.456}`), 12345, 12345, true},
		{"invalid", []byte(`{"f":123.123.123}`), 12345, 12345, true},
		{"linebreaks", []byte("{\n\t\"f\": 123.45\n}\n"), 0, 12345, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := S{F: tt.initial}
			err := json.Unmarshal(tt.data, &s)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if s.F != tt.want {
				t.Errorf("json.Unmarshal() = %d, want %d", int32(s.F), int32(tt.want))
			}
		})
	}
}

func BenchmarkFixed_UnmarshalJSON(b *testing.B) {
	var f decimal.Fixed
	data := []byte("123.45")
	for b.Loop() {
		_ = f.UnmarshalJSON(data)
	}
}
