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
	tests := []struct {
		name    string
		data    string
		want    decimal.Decimal
		wantErr bool
	}{
		{"zero", "0.0", decimal.Decimal{Digits: 1}, false},
		{"bare_integer", "123", decimal.Decimal{Integer: 123}, false},
		{"integer", "123.0", decimal.Decimal{Integer: 123, Digits: 1}, false},
		{"fraction", "0.123", decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"digits", "123.123", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"negative_integer", "-123.0", decimal.Decimal{Integer: 123, Digits: 1, Negative: true}, false},
		{"negative_digits", "-123.456", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, false},
		{"negative_fraction_only", "-0.5", decimal.Decimal{Fraction: 5, Digits: 1, Negative: true}, false},
		{"null", "null", decimal.Decimal{}, false},
		{"quote", `"123.123"`, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"quote_negative", `"-0.5"`, decimal.Decimal{Fraction: 5, Digits: 1, Negative: true}, false},
		{"quote_maxuint64", `"18446744073709551615"`, decimal.Decimal{Integer: ^uint64(0)}, false},
		{"quote_overflow", `"18446744073709551616"`, decimal.Decimal{}, true},
		{"true", "true", decimal.Decimal{}, true},
		{"false", "false", decimal.Decimal{}, true},
		{"empty_array", "[]", decimal.Decimal{}, true},
		{"empty_object", "{}", decimal.Decimal{}, true},
		{"array", "[123.123]", decimal.Decimal{}, true},
		{"object", `{"value":123.123}`, decimal.Decimal{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d decimal.Decimal
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
	tests := []struct {
		name    string
		data    string
		want    decimal.Decimal
		wantErr bool
	}{
		{"zero", `{"d":0.0}`, decimal.Decimal{Digits: 1}, false},
		{"bare_integer", `{"d":123}`, decimal.Decimal{Integer: 123}, false},
		{"integer", `{"d":123.0}`, decimal.Decimal{Integer: 123, Digits: 1}, false},
		{"fraction", `{"d":0.123}`, decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"digits", `{"d":123.123}`, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"negative", `{"d":-123.456}`, decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3, Negative: true}, false},
		{"null", `{"d":null}`, decimal.Decimal{}, false},
		{"maxuint64_string", `{"d":"18446744073709551615"}`, decimal.Decimal{Integer: ^uint64(0)}, false},
		{"maxuint64_plus_one_string", `{"d":"18446744073709551616"}`, decimal.Decimal{}, true},
		{"invalid", `{"d":123.123.123}`, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, true}, // v2 streaming decoder reads 123.123 as valid number before detecting the syntax error
		{"true", `{"d":true}`, decimal.Decimal{}, true},
		{"false", `{"d":false}`, decimal.Decimal{}, true},
		{"empty_array", `{"d":[]}`, decimal.Decimal{}, true},
		{"empty_object", `{"d":{}}`, decimal.Decimal{}, true},
		{"array", `{"d":[123.123]}`, decimal.Decimal{}, true},
		{"object", `{"d":{"value":123.123}}`, decimal.Decimal{}, true},
		{"linebreaks", "{\n\t\"d\": 123.123\n}\n", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"quote", `{"d": "123.123"}`, decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s S
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

func TestDecimal_UnmarshalJSONFrom_ErrorKeepsReceiver(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"bad", "bad"},
		{"overflow_string", `"18446744073709551616"`},
		{"true", "true"},
		{"false", "false"},
		{"array", "[1,2,3]"},
		{"object", `{"k":"v"}`},
	}
	original := decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := original
			dec := jsontext.NewDecoder(bytes.NewReader([]byte(tt.data)))
			if err := d.UnmarshalJSONFrom(dec); err == nil {
				t.Fatalf("UnmarshalJSONFrom(%s) error = nil, want non-nil", tt.data)
			}
			if d != original {
				t.Errorf("UnmarshalJSONFrom(%s) changed receiver: got %#v", tt.data, d)
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
