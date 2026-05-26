package decimal_test

import (
	"encoding/json"
	"testing"

	"github.com/fossoreslp/decimal"
)

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
		{"negative_integer", decimal.Decimal{Integer: 123, Negative: true}, []byte("-123"), false},
		{"negative_fraction", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, []byte("-123.456"), false},
		{"max_uint64_integer", decimal.Decimal{Integer: ^uint64(0)}, []byte("18446744073709551615"), false},
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
	for b.Loop() {
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
		{"null", decimal.Decimal{}, []byte("null"), false},
		{"invalid", decimal.Decimal{}, []byte("123.123.123"), true},
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
		{"null", []byte(`{"d":null}`), false, 0.0},
		{"maxuint64_string", []byte(`{"d":"18446744073709551615"}`), false, 18446744073709551615.0},
		{"maxuint64_plus_one_string", []byte(`{"d":"18446744073709551616"}`), true, 0.0},
		{"invalid", []byte(`{"d":123.123.123}`), true, 0.0},
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

func TestDecimal_UnmarshalJSON_ErrorKeepsReceiver(t *testing.T) {
	d := decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	if err := d.UnmarshalJSON([]byte("bad")); err == nil {
		t.Fatal("UnmarshalJSON(bad) error = nil, want non-nil")
	}
	if d != (decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}) {
		t.Errorf("UnmarshalJSON(bad) changed receiver: got %#v", d)
	}

	d = decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	if err := d.UnmarshalJSON([]byte(`"18446744073709551616"`)); err == nil {
		t.Fatal("UnmarshalJSON(overflow string) error = nil, want non-nil")
	}
	if d != (decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}) {
		t.Errorf("UnmarshalJSON(overflow string) changed receiver: got %#v", d)
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
