package decimal_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/fossoreslp/decimal"
)

type cborMarshalFixture struct {
	name string
	d    decimal.Decimal
	hex  string
}

type cborUnmarshalFixture struct {
	name string
	hex  string
	want decimal.Decimal
}

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

func assertDecimalExact(t *testing.T, got, want decimal.Decimal) {
	t.Helper()
	if got != want {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestDecimal_MarshalCBOR_Exact(t *testing.T) {
	tests := []cborMarshalFixture{
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

func TestDecimal_UnmarshalCBOR_Exact(t *testing.T) {
	tests := []cborUnmarshalFixture{
		{"zero", "00", decimal.Decimal{}},
		{"one", "01", decimal.Decimal{Integer: 1}},
		{"negative_one", "20", decimal.Decimal{Integer: 1, Negative: true}},
		{"integer", "187b", decimal.Decimal{Integer: 123}},
		{"negative_integer", "387a", decimal.Decimal{Integer: 123, Negative: true}},
		{"uint16_range", "1903e8", decimal.Decimal{Integer: 1000}},
		{"large_integer", "1b112210f47de98115", decimal.Decimal{Integer: 1234567890123456789}},
		{"max_uint64", "1bffffffffffffffff", decimal.Decimal{Integer: 18446744073709551615}},
		{"fraction", "c48222187b", decimal.Decimal{Fraction: 123, Digits: 3}},
		{"digits", "c482221a0001e0f3", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3}},
		{"negative_fraction", "c482223a0001e0f2", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 3, Negative: true}},
		{"zero_with_digits", "c4822400", decimal.Decimal{Digits: 5}},
		{"trailing_zeros", "c482241a00bc5ea8", decimal.Decimal{Integer: 123, Fraction: 45000, Digits: 5}},
		{"leading_zeros", "c482241a00bbaf5b", decimal.Decimal{Integer: 123, Fraction: 123, Digits: 5}},
		{"small_fraction", "c4823101", decimal.Decimal{Fraction: 1, Digits: 18}},
		{"max_digits", "c482301baade3d294eb94b87", decimal.Decimal{Integer: 123, Fraction: 12345678901234567, Digits: 17}},
		{"point_five", "c4822005", decimal.Decimal{Fraction: 5, Digits: 1}},
		{"negative_point_five", "c4822024", decimal.Decimal{Fraction: 5, Digits: 1, Negative: true}},
		{"complex", "c482281b0db4da5f4b717715", decimal.Decimal{Integer: 987654321, Fraction: 123456789, Digits: 9}},
		{"bignum_positive", "c48232c24942becfe422c0618115", decimal.Decimal{Integer: 123, Fraction: 1234567890123456789, Digits: 19}},
		{"bignum_negative", "c48232c34942becfe422c0618114", decimal.Decimal{Negative: true, Integer: 123, Fraction: 1234567890123456789, Digits: 19}},
		{"bignum_large", "c48232c2500785ee10d5da46d900f4369fffffffff", decimal.Decimal{Integer: 999999999999999999, Fraction: 9999999999999999999, Digits: 19}},
		{"positive_exponent", "c4820205", decimal.Decimal{Integer: 500}},
		{"positive_exponent_bignum", "c48202c24105", decimal.Decimal{Integer: 500}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got decimal.Decimal
			if err := got.UnmarshalCBOR(mustCBORHex(t, tt.hex)); err != nil {
				t.Fatalf("UnmarshalCBOR() error = %v", err)
			}
			assertDecimalExact(t, got, tt.want)
		})
	}
}

func TestDecimal_UnmarshalCBOR_Floats(t *testing.T) {
	tests := []cborUnmarshalFixture{
		{"float16_1_5", "f93e00", decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}},
		{"float32_1_5", "fa3fc00000", decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}},
		{"float64_1_5", "fb3ff8000000000000", decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1}},
		{"float16_negative_1_5", "f9be00", decimal.Decimal{Integer: 1, Fraction: 5, Digits: 1, Negative: true}},
		{"float32_123_5", "fa42f70000", decimal.Decimal{Integer: 123, Fraction: 5, Digits: 1}},
		{"float64_123_5", "fb405ee00000000000", decimal.Decimal{Integer: 123, Fraction: 5, Digits: 1}},
		{"float64_negative_subnormal", "fb8000000000000001", decimal.Decimal{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got decimal.Decimal
			if err := got.UnmarshalCBOR(mustCBORHex(t, tt.hex)); err != nil {
				t.Fatalf("UnmarshalCBOR() error = %v", err)
			}
			assertDecimalExact(t, got, tt.want)
		})
	}
}

func TestDecimal_UnmarshalCBOR_TrailingGarbage(t *testing.T) {
	data := mustCBORHex(t, "c48232c24942becfe422c0618115")
	data = append(data, 0x00)

	var out decimal.Decimal
	if err := out.UnmarshalCBOR(data); err == nil {
		t.Fatal("UnmarshalCBOR should reject trailing garbage on bignum mantissa path")
	}
}

func TestDecimal_UnmarshalCBOR_Errors(t *testing.T) {
	tests := []struct {
		name string
		hex  string
	}{
		{"empty", ""},
		{"truncated_int_24", "18"},
		{"truncated_int_25", "1901"},
		{"truncated_int_26", "1a0102"},
		{"truncated_int_27", "1b01"},
		{"invalid_additional", "1c"},
		{"neg_overflow", "3bffffffffffffffff"},
		{"int_trailing", "0100"},
		{"unsupported_major", "60"},
		{"unknown_tag", "c1"},
		{"bignum_tag", "c24101"},
		{"decfrac_too_short", "c482"},
		{"decfrac_bad_array", "c48100"},
		{"decfrac_exp_overflow", "c482181500"},
		{"float16_wrong_len", "f900"},
		{"float32_wrong_len", "fa0000"},
		{"float64_wrong_len", "fb000000"},
		{"unknown_simple", "f4"},
		{"decfrac_mantissa_truncated", "c48200"},
		{"decfrac_mantissa_bad_major", "c4820060"},
		{"decfrac_mantissa_bad_tag", "c48200c1"},
		{"decfrac_mantissa_trailing", "c482000100"},
		{"decfrac_bignum_truncated", "c48200c2"},
		{"decfrac_bignum_not_bytes", "c48200c260"},
		{"decfrac_bignum_too_large", "c48200c2510102030405060708090a0b0c0d0e0f1011"},
		{"decfrac_bad_array_len", "c4830000"},
		{"decfrac_exp_bad_type", "c4824000"},
		{"decfrac_mantissa_after_long_exp", "c4821801"},
		{"decfrac_mantissa_parse_err", "c4820018"},
		{"decfrac_pos_exp_overflow", "c482121bffffffffffffffff"},
		{"decfrac_bignum_pos_exp_hi_overflow", "c48201c249010000000000000001"},
		{"decfrac_bignum_pos_exp_mul_overflow", "c48212c248ffffffffffffffff"},
		{"decfrac_bignum_neg_exp_hi_overflow", "c48220c249ff0000000000000001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d decimal.Decimal
			if err := d.UnmarshalCBOR(mustCBORHex(t, tt.hex)); err == nil {
				t.Fatalf("UnmarshalCBOR(%s) should have returned an error", tt.hex)
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
