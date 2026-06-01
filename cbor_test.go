package decimal_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/fossoreslp/decimal"
)

func mustCBORHex(t *testing.T, s string) []byte {
	t.Helper()
	return cborHex(s)
}

func cborHex(s string) []byte {
	s = strings.ReplaceAll(s, " ", "")
	buf, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return buf
}

func TestDecimal_MarshalCBOR_Exact(t *testing.T) {
	tests := []struct {
		name string
		d    decimal.Decimal
		hex  string
	}{
		{"zero", decimal.Decimal{}, "00"},
		{"one", decimal.Decimal{Integer: 1}, "01"},
		{"negative_one", decimal.Decimal{Integer: 1, Negative: true}, "20"},
		{"integer", decimal.Decimal{Integer: 123}, "187b"},
		{"negative_integer", decimal.Decimal{Integer: 123, Negative: true}, "387a"},
		{"uint16_range", decimal.Decimal{Integer: 1000}, "1903e8"},
		{"large_integer", decimal.Decimal{Integer: 1234567890123456789}, "1b112210f47de98115"},
		{"max_uint64", decimal.Decimal{Integer: 18446744073709551615}, "1bffffffffffffffff"},
		{"fraction", decimal.Decimal{Fraction: 123, Digits: 3}, "c48222187b"},
		{"digits", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, "c482221a0001e0f3"},
		{"negative_fraction", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, "c482223a0001e0f2"},
		{"zero_with_digits", decimal.Decimal{Digits: 5}, "c4822400"},
		{"trailing_zeros", decimal.Decimal{Integer: 123, Fraction: 45000, Digits: 5}, "c482241a00bc5ea8"},
		{"leading_zeros", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 5}, "c482241a00bbaf5b"},
		{"small_fraction", decimal.Decimal{Fraction: 1, Digits: 18}, "c4823101"},
		{"max_digits", decimal.Decimal{Integer: 123, Fraction: 12345678901234567, Digits: 17}, "c482301baade3d294eb94b87"},
		{"point_five", decimal.Decimal{Fraction: 5, Digits: 1}, "c4822005"},
		{"negative_point_five", decimal.Decimal{Fraction: 5, Digits: 1, Negative: true}, "c4822024"},
		{"complex", decimal.Decimal{Integer: 987654321, Fraction: 123456789, Digits: 9}, "c482281b0db4da5f4b717715"},
		{"bignum_positive", decimal.Decimal{Integer: 123, Fraction: 1234567890123456789, Digits: 19}, "c48232c24942becfe422c0618115"},
		{"bignum_negative", decimal.Decimal{Negative: true, Integer: 123, Fraction: 1234567890123456789, Digits: 19}, "c48232c34942becfe422c0618114"},
		{"bignum_large", decimal.Decimal{Integer: 999999999999999999, Fraction: 9999999999999999999, Digits: 19}, "c48232c2500785ee10d5da46d900f4369fffffffff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.MarshalCBOR()
			if err != nil {
				t.Fatalf("MarshalCBOR() error = %v", err)
			}
			want := mustCBORHex(t, tt.hex)
			if string(got) != string(want) {
				t.Fatalf("MarshalCBOR() = %x, want %x", got, want)
			}
		})
	}
}

func TestDecimal_UnmarshalCBOR(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		want    decimal.Decimal
		wantErr bool
	}{
		// integer encodings (CBOR major type 0/1)
		{"zero", "00", decimal.Decimal{}, false},
		{"one", "01", decimal.Decimal{Integer: 1}, false},
		{"negative_one", "20", decimal.Decimal{Integer: 1, Negative: true}, false},
		{"integer", "187b", decimal.Decimal{Integer: 123}, false},
		{"negative_integer", "387a", decimal.Decimal{Integer: 123, Negative: true}, false},
		{"uint16_range", "1903e8", decimal.Decimal{Integer: 1000}, false},
		{"large_integer", "1b112210f47de98115", decimal.Decimal{Integer: 1234567890123456789}, false},
		{"max_uint64", "1bffffffffffffffff", decimal.Decimal{Integer: 18446744073709551615}, false},
		// decimal-fraction tag
		{"fraction", "c48222187b", decimal.Decimal{Fraction: 123, Digits: 3}, false},
		{"digits", "c482221a0001e0f3", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}, false},
		{"negative_fraction", "c482223a0001e0f2", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}, false},
		{"zero_with_digits", "c4822400", decimal.Decimal{Digits: 5}, false},
		{"trailing_zeros", "c482241a00bc5ea8", decimal.Decimal{Integer: 123, Fraction: 45000, Digits: 5}, false},
		{"leading_zeros", "c482241a00bbaf5b", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 5}, false},
		{"small_fraction", "c4823101", decimal.Decimal{Fraction: 1, Digits: 18}, false},
		{"max_digits", "c482301baade3d294eb94b87", decimal.Decimal{Integer: 123, Fraction: 12345678901234567, Digits: 17}, false},
		{"point_five", "c4822005", decimal.Decimal{Fraction: 5, Digits: 1}, false},
		{"negative_point_five", "c4822024", decimal.Decimal{Fraction: 5, Digits: 1, Negative: true}, false},
		{"complex", "c482281b0db4da5f4b717715", decimal.Decimal{Integer: 987654321, Fraction: 123456789, Digits: 9}, false},
		{"bignum_positive", "c48232c24942becfe422c0618115", decimal.Decimal{Integer: 123, Fraction: 1234567890123456789, Digits: 19}, false},
		{"bignum_negative", "c48232c34942becfe422c0618114", decimal.Decimal{Negative: true, Integer: 123, Fraction: 1234567890123456789, Digits: 19}, false},
		{"bignum_large", "c48232c2500785ee10d5da46d900f4369fffffffff", decimal.Decimal{Integer: 999999999999999999, Fraction: 9999999999999999999, Digits: 19}, false},
		{"positive_exponent", "c4820205", decimal.Decimal{Integer: 500}, false},
		{"positive_exponent_bignum", "c48202c24105", decimal.Decimal{Integer: 500}, false},
		// floating-point encodings
		{"float16_1_5", "f93e00", decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}, false},
		{"float32_1_5", "fa3fc00000", decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}, false},
		{"float64_1_5", "fb3ff8000000000000", decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}, false},
		{"float16_negative_1_5", "f9be00", decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1, Negative: true}, false},
		{"float32_123_5", "fa42f70000", decimal.Decimal{Integer: 123, Fraction: 5, Digits: 1}, false},
		{"float64_123_5", "fb405ee00000000000", decimal.Decimal{Integer: 123, Fraction: 5, Digits: 1}, false},
		{"float64_negative_subnormal", "fb8000000000000001", decimal.Decimal{}, false},
		// errors — receiver state is implementation-defined; we only assert err is non-nil.
		{"empty", "", decimal.Decimal{}, true},
		{"truncated_int_24", "18", decimal.Decimal{}, true},
		{"truncated_int_25", "1901", decimal.Decimal{}, true},
		{"truncated_int_26", "1a0102", decimal.Decimal{}, true},
		{"truncated_int_27", "1b01", decimal.Decimal{}, true},
		{"invalid_additional", "1c", decimal.Decimal{}, true},
		{"neg_overflow", "3bffffffffffffffff", decimal.Decimal{}, true},
		{"int_trailing", "0100", decimal.Decimal{}, true},
		{"unsupported_major", "60", decimal.Decimal{}, true},
		{"unknown_tag", "c1", decimal.Decimal{}, true},
		{"bignum_tag", "c24101", decimal.Decimal{}, true},
		{"decfrac_too_short", "c482", decimal.Decimal{}, true},
		{"decfrac_bad_array", "c48100", decimal.Decimal{}, true},
		{"decfrac_exp_overflow", "c482181500", decimal.Decimal{}, true},
		{"float16_wrong_len", "f900", decimal.Decimal{}, true},
		{"float32_wrong_len", "fa0000", decimal.Decimal{}, true},
		{"float64_wrong_len", "fb000000", decimal.Decimal{}, true},
		{"unknown_simple", "f4", decimal.Decimal{}, true},
		{"decfrac_mantissa_truncated", "c48200", decimal.Decimal{}, true},
		{"decfrac_mantissa_bad_major", "c4820060", decimal.Decimal{}, true},
		{"decfrac_mantissa_bad_tag", "c48200c1", decimal.Decimal{}, true},
		{"decfrac_mantissa_trailing", "c482000100", decimal.Decimal{}, true},
		{"decfrac_bignum_truncated", "c48200c2", decimal.Decimal{}, true},
		{"decfrac_bignum_not_bytes", "c48200c260", decimal.Decimal{}, true},
		{"decfrac_bignum_too_large", "c48200c2510102030405060708090a0b0c0d0e0f1011", decimal.Decimal{}, true},
		{"decfrac_bad_array_len", "c4830000", decimal.Decimal{}, true},
		{"decfrac_exp_bad_type", "c4824000", decimal.Decimal{}, true},
		{"decfrac_mantissa_after_long_exp", "c4821801", decimal.Decimal{}, true},
		{"decfrac_mantissa_parse_err", "c4820018", decimal.Decimal{}, true},
		{"decfrac_pos_exp_overflow", "c482121bffffffffffffffff", decimal.Decimal{}, true},
		{"decfrac_bignum_pos_exp_hi_overflow", "c48201c249010000000000000001", decimal.Decimal{}, true},
		{"decfrac_bignum_pos_exp_mul_overflow", "c48212c248ffffffffffffffff", decimal.Decimal{}, true},
		{"decfrac_bignum_neg_exp_hi_overflow", "c48220c249ff0000000000000001", decimal.Decimal{}, true},
		{"bignum_trailing_garbage", "c48232c24942becfe422c061811500", decimal.Decimal{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got decimal.Decimal
			err := got.UnmarshalCBOR(mustCBORHex(t, tt.hex))
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalCBOR() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("UnmarshalCBOR() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func BenchmarkDecimal_MarshalCBOR(b *testing.B) {
	d := decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}
	for b.Loop() {
		_, _ = d.MarshalCBOR()
	}
}

func BenchmarkDecimal_UnmarshalCBOR(b *testing.B) {
	data := cborHex("c482221a0001e0f3")
	var result decimal.Decimal
	for b.Loop() {
		_ = result.UnmarshalCBOR(data)
	}
}

func TestFixed_MarshalCBOR(t *testing.T) {
	tests := []struct {
		name string
		f    decimal.Fixed
		hex  string
	}{
		{"zero", 0, "c4822100"},
		{"hundredth", 1, "c4822101"},
		{"fraction", 5, "c4822105"},
		{"fraction_full", 99, "c482211863"},
		{"one", 100, "c482211864"},
		{"integer", 12300, "c4822119300c"},
		{"digits", 12345, "c48221193039"},
		{"max_int32", 2147483647, "c482211a7fffffff"},
		{"negative_hundredth", -1, "c4822120"},
		{"negative_fraction", -5, "c4822124"},
		{"negative_integer", -500, "c482213901f3"},
		{"negative_digits", -12345, "c48221393038"},
		{"min_int32", -2147483648, "c482213a7fffffff"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.MarshalCBOR()
			if err != nil {
				t.Fatalf("MarshalCBOR() error = %v", err)
			}
			want := mustCBORHex(t, tt.hex)
			if string(got) != string(want) {
				t.Fatalf("MarshalCBOR() = %x, want %x", got, want)
			}
		})
	}
}

func TestFixed_UnmarshalCBOR(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		initial decimal.Fixed
		want    decimal.Fixed
		wantErr bool
	}{
		// values produced by Fixed.MarshalCBOR
		{"zero", "c4822100", 12345, 0, false},
		{"digits", "c48221193039", 0, 12345, false},
		{"max_int32", "c482211a7fffffff", 0, 2147483647, false},
		{"negative_digits", "c48221393038", 0, -12345, false},
		{"min_int32", "c482213a7fffffff", 0, -2147483648, false},
		// plain CBOR encodings are accepted when representable as Fixed
		{"plain_integer", "187b", 0, 12300, false},
		{"plain_negative_integer", "387a", 0, -12300, false},
		{"float_1_5", "fb3ff8000000000000", 0, 150, false},
		{"float_negative_1_5", "f9be00", 0, -150, false},
		{"more_digits", "c482221a0001e0f3", 12345, 12345, true},
		{"invalid", "1c", 12345, 12345, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.initial
			err := got.UnmarshalCBOR(mustCBORHex(t, tt.hex))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalCBOR() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("UnmarshalCBOR() = %d, want %d", int32(got), int32(tt.want))
			}
		})
	}
}

func TestFixed_CBORRoundTrip(t *testing.T) {
	values := []decimal.Fixed{0, 1, 99, 100, 12345, 2147483647, -1, -500, -12345, -2147483648}
	for _, v := range values {
		data, err := v.MarshalCBOR()
		if err != nil {
			t.Fatalf("MarshalCBOR(%d) error = %v", int32(v), err)
		}
		var got decimal.Fixed
		if err := got.UnmarshalCBOR(data); err != nil {
			t.Fatalf("UnmarshalCBOR(%x) error = %v", data, err)
		}
		if got != v {
			t.Errorf("round-trip mismatch: got %d, want %d", int32(got), int32(v))
		}
	}
}

func BenchmarkFixed_MarshalCBOR(b *testing.B) {
	f := decimal.Fixed(12345)
	for b.Loop() {
		_, _ = f.MarshalCBOR()
	}
}

func BenchmarkFixed_UnmarshalCBOR(b *testing.B) {
	data := cborHex("c48221193039")
	var f decimal.Fixed
	for b.Loop() {
		_ = f.UnmarshalCBOR(data)
	}
}
