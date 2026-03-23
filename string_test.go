package decimal_test

import (
	"encoding"
	"encoding/xml"
	"testing"

	"github.com/fossoreslp/decimal"
)

// Compile-time interface checks.
var (
	_ encoding.TextMarshaler   = decimal.Decimal{}
	_ encoding.TextUnmarshaler = (*decimal.Decimal)(nil)
)

func TestDecimal_NewFromString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    decimal.Decimal
		wantErr bool
	}{
		{"zero", "0.0", decimal.Decimal{Digits: 1}, false},
		{"integer", "123.0", decimal.Decimal{Integer: 123, Digits: 1}, false},
		{"fraction", "0.123", decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"digits", "123.123", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"negative", "-123.123", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, false},
		{"invalid", "123.123.123", decimal.Decimal{}, true},
		{"leading_dot", ".123", decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"trailing_dot", "123.", decimal.Decimal{Integer: 123}, false},
		{"negative_integer", "-5", decimal.Decimal{Integer: 5, Negative: true}, false},
		{"plain_integer", "42", decimal.Decimal{Integer: 42}, false},
		{"empty", "", decimal.Decimal{}, true},
		{"bare_minus", "-", decimal.Decimal{}, true},
		{"bare_dot", ".", decimal.Decimal{}, true},
		{"minus_dot", "-.", decimal.Decimal{}, true},
		{"19_digit_with_decimal", "1234567890123456789.0", decimal.Decimal{Integer: 1234567890123456789, Digits: 1}, false},
		{"maxuint64_pos", "18446744073709551615", decimal.Decimal{Integer: 18446744073709551615}, false},
		{"maxuint64_neg", "-18446744073709551615", decimal.Decimal{Integer: 18446744073709551615, Negative: true}, false},
		{"maxuint64_plus_one", "18446744073709551616", decimal.Decimal{}, true},
		{"maxuint64_with_fraction", "18446744073709551615.0", decimal.Decimal{Integer: 18446744073709551615, Digits: 1}, false},
		{"maxuint64_plus_one_with_fraction", "18446744073709551616.0", decimal.Decimal{}, true},
		{"leading_zeros_fit", "0000000000000000000000000000000000000001", decimal.Decimal{Integer: 1}, false},
		{"leading_zeros_fit_with_fraction", "0000000000000000000000000000000000000001.25", decimal.Decimal{Integer: 1, Fraction: 25, Digits: 2}, false},
		{"leading_zeros_maxuint64", "00018446744073709551615", decimal.Decimal{Integer: 18446744073709551615}, false},
		{"leading_zeros_overflow", "00018446744073709551616", decimal.Decimal{}, true},
		{"negative_maxuint64_plus_one", "-18446744073709551616", decimal.Decimal{}, true},
		{"all_zeros_long", "0000000000000000000000000000000000000000", decimal.Decimal{}, false},
		{"negative_all_zeros_long", "-0000000000000000000000000000000000000000", decimal.Decimal{}, false},
		{"fraction_19_digits", "0.1234567890123456789", decimal.Decimal{Fraction: 1234567890123456789, Digits: 19}, false},
		{"fraction_overflow", "0.123456789012345678901", decimal.Decimal{}, true},
		{"integer_overflow", "123456789012345678901234567890", decimal.Decimal{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decimal.NewFromString(tt.s)
			if err != nil && !tt.wantErr {
				t.Errorf("Decimal.NewString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.want.Equal(got) {
				t.Errorf("Decimal.NewString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_NewFromString(b *testing.B) {
	benchmarks := []struct {
		name string
		s    string
	}{
		{"integer", "123"},
		{"fraction", "0.123"},
		{"digits", "123.123"},
		{"long", "12345.67890"},
	}
	for _, bb := range benchmarks {
		b.Run(bb.name, func(b *testing.B) {
			for b.Loop() {
				_, _ = decimal.NewFromString(bb.s)
			}
		})
	}
}

func TestDecimal_NewFromStringFuzzy(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    decimal.Decimal
		wantErr bool
	}{
		{"zero", "0.0", decimal.Decimal{Digits: 1}, false},
		{"integer_no_dot", "123", decimal.Decimal{Integer: 123}, false},
		{"integer_suffix_no_dot", "123abc", decimal.Decimal{Integer: 123}, false},
		{"integer", "123.0", decimal.Decimal{Integer: 123, Digits: 1}, false},
		{"fraction", "0.123", decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"digits", "123.123", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"negative", "-123.123", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, false},
		{"prefix", "prefix123.123", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"suffix", "123.123suffix", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"both", "prefix123.123suffix", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"utf8", "ä 123.123 🤦‍♀️", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"prefix_neg", "prefix-123.123", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, false},
		{"suffix_neg", "-123.123suffix", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, false},
		{"both_neg", "prefix-123.123suffix", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, false},
		{"utf8_neg", "ä -123.123 🤦‍♀️", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, false},
		{"leading_dot", ".123", decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"trailing_dot", "123.", decimal.Decimal{Integer: 123}, false},
		{"maxuint64_pos", "18446744073709551615", decimal.Decimal{Integer: 18446744073709551615}, false},
		{"maxuint64_neg", "-18446744073709551615", decimal.Decimal{Integer: 18446744073709551615, Negative: true}, false},
		{"prefix_maxuint64", "value=18446744073709551615;", decimal.Decimal{Integer: 18446744073709551615}, false},
		{"prefix_maxuint64_plus_one", "value=18446744073709551616;", decimal.Decimal{}, true},
		{"prefix_leading_zeros_fit", "id=00000000000000000000000042 end", decimal.Decimal{Integer: 42}, false},
		{"prefix_negative_maxuint64", "x=-18446744073709551615!", decimal.Decimal{Integer: 18446744073709551615, Negative: true}, false},
		{"prefix_negative_maxuint64_plus_one", "x=-18446744073709551616!", decimal.Decimal{}, true},
		{"skip_malformed_minus_prefix", "abc-xyz123", decimal.Decimal{Integer: 123}, false},
		{"skip_malformed_minus_dot_prefix", "abc-.x1.2", decimal.Decimal{Integer: 1, Fraction: 2, Digits: 1}, false},
		{"skip_malformed_dot_prefix", "abc..5", decimal.Decimal{Fraction: 5, Digits: 1}, false},
		{"skip_malformed_negative_prefix", "a-xyz7", decimal.Decimal{Integer: 7}, false},
		{"empty", "", decimal.Decimal{}, true},
		{"minus_no_digits", "-abc", decimal.Decimal{}, true},
		{"dot_no_digits", ".abc", decimal.Decimal{}, true},
		{"no_digits", "abc", decimal.Decimal{}, true},
		{"integer_overflow", "$123456789012345678901.00", decimal.Decimal{}, true},
		{"fraction_overflow", "$0.123456789012345678901", decimal.Decimal{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decimal.NewFromStringFuzzy(tt.s)
			if err != nil && !tt.wantErr {
				t.Errorf("Decimal.NewString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.want.Equal(got) {
				t.Errorf("Decimal.NewString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_NewFromStringFuzzy(b *testing.B) {
	benchmarks := []struct {
		name string
		s    string
	}{
		{"integer", "123"},
		{"fraction", "0.123"},
		{"digits", "123.123"},
		{"long", "12345.67890"},
		{"prefix", "$123.456"},
		{"suffix", "123.456 lm"},
		{"both", "Angle 123.456°"},
	}
	for _, bb := range benchmarks {
		b.Run(bb.name, func(b *testing.B) {
			for b.Loop() {
				_, _ = decimal.NewFromStringFuzzy(bb.s)
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
		{"max_uint64_integer", decimal.Decimal{Integer: ^uint64(0)}, "18446744073709551615"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("Decimal.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_StringRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
	}{
		{"zero", decimal.Decimal{}},
		{"integer", decimal.Decimal{Integer: 42}},
		{"negative_integer", decimal.Decimal{Integer: 42, Negative: true}},
		{"fraction", decimal.Decimal{Fraction: 125, Digits: 3}},
		{"leading_zero_fraction", decimal.Decimal{Integer: 12, Fraction: 34, Digits: 3}},
		{"trailing_zero_fraction", decimal.Decimal{Integer: 12, Fraction: 340, Digits: 3}},
		{"max_fraction_digits", decimal.Decimal{Fraction: 1234567890123456789, Digits: 19}},
		{"max_uint64_integer", decimal.Decimal{Integer: ^uint64(0)}},
		{"negative_max_uint64_integer", decimal.Decimal{Integer: ^uint64(0), Negative: true}},
		{"combined", decimal.Decimal{Integer: 12345678901234567890, Fraction: 1234567890123456789, Digits: 19}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.d.String()
			got, err := decimal.NewFromString(s)
			if err != nil {
				t.Fatalf("NewFromString(%q) error = %v", s, err)
			}
			if got != tt.d {
				t.Fatalf("round-trip mismatch: got %#v, want %#v", got, tt.d)
			}
		})
	}
}

func BenchmarkDecimal_String(b *testing.B) {
	benchmarks := []struct {
		name string
		d    decimal.Decimal
	}{
		{"integer", decimal.Decimal{Integer: 123, Fraction: 0, Digits: 0}},
		{"fraction", decimal.Decimal{Integer: 0, Fraction: 123, Digits: 3}},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}},
		{"long", decimal.Decimal{Integer: 12345, Fraction: 67890, Digits: 5}},
		{"pad", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 6}},
		{"pad_long", decimal.Decimal{Integer: 12345, Fraction: 67890, Digits: 9}},
	}
	for _, bb := range benchmarks {
		b.Run(bb.name, func(b *testing.B) {
			for b.Loop() {
				_ = bb.d.String()
			}
		})
	}
}

func TestDecimal_MarshalText(t *testing.T) {
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
		{"trailing_zeros", decimal.Decimal{Integer: 123, Fraction: 1234500, Digits: 7}, "123.1234500"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.MarshalText()
			if err != nil {
				t.Fatalf("MarshalText() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("MarshalText() = %q, want %q", got, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_MarshalText(b *testing.B) {
	d := decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}
	for b.Loop() {
		_, _ = d.MarshalText()
	}
}

func TestDecimal_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    decimal.Decimal
		wantErr bool
	}{
		{"zero", "0", decimal.Decimal{}, false},
		{"zero_frac", "0.0", decimal.Decimal{Digits: 1}, false},
		{"integer", "123", decimal.Decimal{Integer: 123}, false},
		{"fraction", "0.123", decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"digits", "123.123", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"negative", "-42.5", decimal.Decimal{Integer: 42, Fraction: 5, Digits: 1, Negative: true}, false},
		{"max_uint64", "18446744073709551615", decimal.Decimal{Integer: 18446744073709551615}, false},
		{"overflow", "18446744073709551616", decimal.Decimal{}, true},
		{"invalid", "abc", decimal.Decimal{}, true},
		{"empty", "", decimal.Decimal{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d decimal.Decimal
			err := d.UnmarshalText([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !d.Equal(tt.want) {
				t.Errorf("UnmarshalText() = %v, want %v", d, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_UnmarshalText(b *testing.B) {
	var d decimal.Decimal
	data := []byte("123.123")
	for b.Loop() {
		_ = d.UnmarshalText(data)
	}
}

func TestDecimal_UnmarshalText_ErrorKeepsReceiver(t *testing.T) {
	d := decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}
	if err := d.UnmarshalText([]byte("bad")); err == nil {
		t.Fatal("UnmarshalText(bad) error = nil, want non-nil")
	}
	if d != (decimal.Decimal{Integer: 7, Fraction: 5, Digits: 1}) {
		t.Errorf("UnmarshalText(bad) changed receiver: got %#v", d)
	}
}

func TestDecimal_TextRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
	}{
		{"zero", decimal.Decimal{}},
		{"integer", decimal.Decimal{Integer: 42}},
		{"negative_integer", decimal.Decimal{Integer: 42, Negative: true}},
		{"fraction", decimal.Decimal{Fraction: 125, Digits: 3}},
		{"combined", decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}},
		{"max_uint64", decimal.Decimal{Integer: ^uint64(0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.d.MarshalText()
			if err != nil {
				t.Fatalf("MarshalText() error = %v", err)
			}
			var got decimal.Decimal
			if err := got.UnmarshalText(data); err != nil {
				t.Fatalf("UnmarshalText(%q) error = %v", data, err)
			}
			if got != tt.d {
				t.Errorf("round-trip mismatch: got %#v, want %#v", got, tt.d)
			}
		})
	}
}

func TestDecimal_XML(t *testing.T) {
	type S struct {
		XMLName xml.Name        `xml:"root"`
		D       decimal.Decimal `xml:"d"`
	}

	// Marshal
	s := S{D: decimal.Decimal{Integer: 123, Fraction: 456, Digits: 3}}
	data, err := xml.Marshal(s)
	if err != nil {
		t.Fatalf("xml.Marshal() error = %v", err)
	}
	want := "<root><d>123.456</d></root>"
	if string(data) != want {
		t.Errorf("xml.Marshal() = %q, want %q", data, want)
	}

	// Unmarshal
	var got S
	if err := xml.Unmarshal(data, &got); err != nil {
		t.Fatalf("xml.Unmarshal() error = %v", err)
	}
	if got.D != s.D {
		t.Errorf("xml.Unmarshal() = %#v, want %#v", got.D, s.D)
	}
}
