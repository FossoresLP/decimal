//go:build goexperiment.jsonv2

package decimal_test

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"testing"

	"github.com/fossoreslp/decimal"
)

func TestDecimal_MarshalJSONTo(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
		want string
	}{
		{"zero", decimal.Decimal{}, "0"},
		{"integer", decimal.Decimal{Integer: 123}, "123"},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, "0.123"},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, "123.123"},
		{"trailing_zeros", decimal.Decimal{Integer: 123, Fraction: 1234500, Digits: 7}, "123.1234500"},
		{"negative_integer", decimal.Decimal{Integer: 123, Negative: true}, "-123"},
		{"negative_fraction", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, "-123.456"},
		{"negative_fraction_only", decimal.Decimal{Fraction: 5, Digits: 1, Negative: true}, "-0.5"},
		{"max_uint64_integer", decimal.Decimal{Integer: ^uint64(0)}, "18446744073709551615"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			enc := jsontext.NewEncoder(&buf)
			if err := tt.d.MarshalJSONTo(enc); err != nil {
				t.Fatalf("MarshalJSONTo() error = %v", err)
			}
			if got := buf.String(); got != tt.want+"\n" {
				t.Errorf("MarshalJSONTo() = %q, want %q", got, tt.want+"\n")
			}
		})
	}
}

func BenchmarkDecimal_MarshalJSONTo(b *testing.B) {
	d := decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}
	var buf bytes.Buffer
	enc := jsontext.NewEncoder(&buf)
	for b.Loop() {
		buf.Reset()
		_ = d.MarshalJSONTo(enc)
	}
}

func TestDecimal_UnmarshalJSONFrom(t *testing.T) {
	sentinel := decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	tests := []struct {
		name    string
		data    string
		initial decimal.Decimal
		want    decimal.Decimal
		wantErr bool
	}{
		{"zero", "0.0", sentinel, decimal.Decimal{Digits: 1}, false},
		{"bare_integer", "123", sentinel, decimal.Decimal{Integer: 123}, false},
		{"integer", "123.0", sentinel, decimal.Decimal{Integer: 123, Digits: 1}, false},
		{"fraction", "0.123", sentinel, decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"digits", "123.123", sentinel, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"negative_integer", "-123.0", sentinel, decimal.Decimal{Integer: 123, Digits: 1, Negative: true}, false},
		{"negative_digits", "-123.456", sentinel, decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, false},
		{"negative_fraction_only", "-0.5", sentinel, decimal.Decimal{Fraction: 5, Digits: 1, Negative: true}, false},
		{"null", "null", sentinel, decimal.Decimal{}, false},
		{"quote", `"123.123"`, sentinel, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"quote_negative", `"-0.5"`, sentinel, decimal.Decimal{Fraction: 5, Digits: 1, Negative: true}, false},
		{"quote_maxuint64", `"18446744073709551615"`, sentinel, decimal.Decimal{Integer: ^uint64(0)}, false},
		{"quote_overflow", `"18446744073709551616"`, sentinel, sentinel, true},
		{"bad", "bad", sentinel, sentinel, true},
		{"true", "true", sentinel, sentinel, true},
		{"false", "false", sentinel, sentinel, true},
		{"empty_array", "[]", sentinel, sentinel, true},
		{"empty_object", "{}", sentinel, sentinel, true},
		{"array", "[123.123]", sentinel, sentinel, true},
		{"object", `{"value":123.123}`, sentinel, sentinel, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.initial
			dec := jsontext.NewDecoder(bytes.NewReader([]byte(tt.data)))
			if err := d.UnmarshalJSONFrom(dec); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSONFrom() error = %v, wantErr %v", err, tt.wantErr)
			}
			if d != tt.want {
				t.Errorf("UnmarshalJSONFrom() = %#v, want %#v", d, tt.want)
			}
		})
	}
}

func TestDecimal_UnmarshalJSONFrom_struct(t *testing.T) {
	type S struct {
		D decimal.Decimal `json:"d"`
	}
	sentinel := decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	// "invalid": v2 streaming decoder reads 123.123 as a valid number, mutates the field,
	// then detects the trailing ".123" syntax error. So the receiver ends up as 123.123 rather than sentinel.
	parsedBeforeSyntaxError := decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}
	tests := []struct {
		name    string
		data    string
		initial decimal.Decimal
		want    decimal.Decimal
		wantErr bool
	}{
		{"zero", `{"d":0.0}`, sentinel, decimal.Decimal{Digits: 1}, false},
		{"bare_integer", `{"d":123}`, sentinel, decimal.Decimal{Integer: 123}, false},
		{"integer", `{"d":123.0}`, sentinel, decimal.Decimal{Integer: 123, Digits: 1}, false},
		{"fraction", `{"d":0.123}`, sentinel, decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"digits", `{"d":123.123}`, sentinel, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"negative", `{"d":-123.456}`, sentinel, decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, false},
		{"null", `{"d":null}`, sentinel, decimal.Decimal{}, false},
		{"maxuint64_string", `{"d":"18446744073709551615"}`, sentinel, decimal.Decimal{Integer: ^uint64(0)}, false},
		{"maxuint64_plus_one_string", `{"d":"18446744073709551616"}`, sentinel, sentinel, true},
		{"invalid", `{"d":123.123.123}`, sentinel, parsedBeforeSyntaxError, true},
		{"true", `{"d":true}`, sentinel, sentinel, true},
		{"false", `{"d":false}`, sentinel, sentinel, true},
		{"empty_array", `{"d":[]}`, sentinel, sentinel, true},
		{"empty_object", `{"d":{}}`, sentinel, sentinel, true},
		{"array", `{"d":[123.123]}`, sentinel, sentinel, true},
		{"object", `{"d":{"value":123.123}}`, sentinel, sentinel, true},
		{"linebreaks", "{\n\t\"d\": 123.123\n}\n", sentinel, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"quote", `{"d": "123.123"}`, sentinel, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := S{D: tt.initial}
			dec := jsontext.NewDecoder(bytes.NewReader([]byte(tt.data)))
			if err := json.UnmarshalDecode(dec, &s); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalDecode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if s.D != tt.want {
				t.Errorf("UnmarshalDecode() = %#v, want %#v", s.D, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_UnmarshalJSONFrom(b *testing.B) {
	data := []byte("123.123")
	r := bytes.NewReader(data)
	dec := jsontext.NewDecoder(r)
	for b.Loop() {
		var d decimal.Decimal
		_ = d.UnmarshalJSONFrom(dec)
		r.Reset(data)
	}
}

func TestFixed_MarshalJSONTo(t *testing.T) {
	tests := []struct {
		name string
		f    decimal.Fixed
		want string
	}{
		{"zero", 0, "0.00"},
		{"fraction", 5, "0.05"},
		{"one", 100, "1.00"},
		{"integer", 12300, "123.00"},
		{"digits", 12345, "123.45"},
		{"max_int32", 2147483647, "21474836.47"},
		{"negative_digits", -12345, "-123.45"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			enc := jsontext.NewEncoder(&buf)
			if err := tt.f.MarshalJSONTo(enc); err != nil {
				t.Fatalf("MarshalJSONTo() error = %v", err)
			}
			if got := buf.String(); got != tt.want+"\n" {
				t.Errorf("MarshalJSONTo() = %q, want %q", got, tt.want+"\n")
			}
		})
	}
}

func BenchmarkFixed_MarshalJSONTo(b *testing.B) {
	f := decimal.Fixed(12345)
	var buf bytes.Buffer
	enc := jsontext.NewEncoder(&buf)
	for b.Loop() {
		buf.Reset()
		_ = f.MarshalJSONTo(enc)
	}
}

func TestFixed_UnmarshalJSONFrom(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		initial decimal.Fixed
		want    decimal.Fixed
		wantErr bool
	}{
		{"zero", "0.0", 12345, 0, false},
		{"bare_integer", "123", 0, 12300, false},
		{"fraction", "0.12", 0, 12, false},
		{"digits", "123.45", 0, 12345, false},
		{"negative_digits", "-123.45", 0, -12345, false},
		{"null", "null", 12345, 0, false},
		{"quote", `"123.45"`, 0, 12345, false},
		{"bad", "bad", 12345, 12345, true},
		{"truncated", "123.456", 12345, 12345, true},
		{"true", "true", 12345, 12345, true},
		{"false", "false", 12345, 12345, true},
		{"empty_array", "[]", 12345, 12345, true},
		{"empty_object", "{}", 12345, 12345, true},
		{"array", "[123.45]", 12345, 12345, true},
		{"object", `{"value":123.45}`, 12345, 12345, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.initial
			dec := jsontext.NewDecoder(bytes.NewReader([]byte(tt.data)))
			err := f.UnmarshalJSONFrom(dec)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSONFrom() error = %v, wantErr %v", err, tt.wantErr)
			}
			if f != tt.want {
				t.Errorf("UnmarshalJSONFrom() = %d, want %d", int32(f), int32(tt.want))
			}
		})
	}
}

func TestFixed_UnmarshalJSONFrom_struct(t *testing.T) {
	type S struct {
		F decimal.Fixed `json:"f"`
	}
	tests := []struct {
		name    string
		data    string
		initial decimal.Fixed
		want    decimal.Fixed
		wantErr bool
	}{
		{"zero", `{"f":0.0}`, 12345, 0, false},
		{"bare_integer", `{"f":123}`, 0, 12300, false},
		{"digits", `{"f":123.45}`, 0, 12345, false},
		{"negative", `{"f":-123.45}`, 0, -12345, false},
		{"null", `{"f":null}`, 12345, 0, false},
		{"quote", `{"f":"123.45"}`, 0, 12345, false},
		{"linebreaks", "{\n\t\"f\": 123.45\n}\n", 0, 12345, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := S{F: tt.initial}
			dec := jsontext.NewDecoder(bytes.NewReader([]byte(tt.data)))
			err := json.UnmarshalDecode(dec, &s)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalDecode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if s.F != tt.want {
				t.Errorf("UnmarshalDecode() = %d, want %d", int32(s.F), int32(tt.want))
			}
		})
	}
}

func BenchmarkFixed_UnmarshalJSONFrom(b *testing.B) {
	data := []byte("123.45")
	r := bytes.NewReader(data)
	dec := jsontext.NewDecoder(r)
	for b.Loop() {
		var f decimal.Fixed
		_ = f.UnmarshalJSONFrom(dec)
		r.Reset(data)
	}
}
