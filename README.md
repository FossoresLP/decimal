# decimal

`decimal` is a small Go package for fixed-layout decimal values with value semantics and no heap allocation in the common paths.

It is designed for high performance, compact storage, straightforward formatting, and predictable behavior that stays close to primitive integer operations in Go.

It is **not** an arbitrary-precision decimal library.

## Representation

A `Decimal` stores:

- a sign bit
- an unsigned 64-bit integer part
- an unsigned 64-bit fractional part
- a decimal scale in `Digits`
- as a 24 byte struct

The supported range is:

- integer part: `0` to `math.MaxUint64`
- fractional precision: `0` to `19` decimal digits

That means the type can represent values such as:

- `123`
- `-123.45`
- `0.000000000000000001`

but not arbitrary precision or more than 19 digits after the decimal point.

## Installation

```bash
go get github.com/fossoreslp/decimal
```

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/fossoreslp/decimal"
)

func main() {
	a, _ := decimal.NewFromString("123.4500")
	b := decimal.New(2)

	sum := decimal.Add(a, b)
	expected, _ := decimal.NewFromString("125.45")

	fmt.Println(a.String())   // 123.4500
	fmt.Println(sum.String()) // 125.4500
	fmt.Println(decimal.Equal(sum, expected)) // true
}
```

## Constructors And Parsing

### `New`

`New` is a generic constructor that accepts Go builtin numeric types:

- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `int`, `int8`, `int16`, `int32`, `int64`
- `float32`, `float64`

Behavior:

- integer inputs are represented exactly
- `float32` values are converted with 10 fractional digits
- `float64` values are converted with 18 fractional digits
- trailing fractional zeros are removed after float conversion
- `NaN`, `+Inf`, and `-Inf` are converted to zero
- finite floats larger than the `uint64` range inherit Go's `float`-to-`uint64` overflow behavior

### `NewFromString`

`NewFromString` parses a strict decimal string:

```go
d, err := decimal.NewFromString("-123.4500")
```

Notes:

- accepts an optional leading `-`
- accepts a leading decimal point such as `.25`
- accepts a trailing decimal point such as `123.`
- rejects extra surrounding characters
- allows at most 19 digits after the decimal point
- accepts any integer length whose numeric value still fits in `uint64`
- leading zeros are allowed and do not by themselves cause integer overflow

### `NewFromStringFuzzy`

`NewFromStringFuzzy` scans the first decimal-looking token from a larger string:

```go
d, err := decimal.NewFromStringFuzzy("price: 123.45 EUR")
```

This is intended for loose extraction, not strict validation.

### `Zero`

`Zero` creates a zero-valued decimal value.

## Arithmetic

The package currently provides:

- `Add(a, b) Decimal`
- `Subtract(a, b) Decimal`
- `Multiply(a, b) Decimal`
- `Divide(a, b) Decimal`
- `Negate(n) Decimal`
- `Absolute(n) Decimal`

### Arithmetic Semantics

Operations are designed to behave like primitive integer operations in Go with respect to overflow and division-by-zero semantics. They will wrap on overflow while divide-by-zero panics.

Arithmetic operations will expand the number of decimal digits as far as necessary (up to the limit of 19 digits) to represent the result accurately.

Operations that require more than 19 digits after the decimal point may lose precision.

All arithmetic operations accept `Decimal` values as well as any of the primitive Go numeric types, converting to `Decimal` as described by `New`.

## Equality / Comparison

`Equal` compares values after canonical trailing-zero truncation.

These compare equal:

```go
a, _ := decimal.NewFromString("0.1")
b, _ := decimal.NewFromString("0.10")
decimal.Equal(a, b) // true
```

`Compare` compares to values and returns `-1` if `a < b`, `0` if `a = b` and 1 if `a > b`.

```go
a := decimal.New(1)
b := decimal.New(-2)
decimal.Compare(a, b) // 1
```

`IsZero` checks if a decimal is exactly zero.

```go
decimal.New(0).IsZero() // true
```

## Precision

Creating decimal values from strings adheres to the precision on the input.
Conversion of numeric types uses the least amount of fractional digits possible to represent the exact converted value.

Arithmetic operations extend precision as necessary to represent the resulting value exactly unless it overflows the limits.

To adjust precision manually, there are three options:

- `ToDigits(uint8) Decimal`
- `Round(uint8) Decimal`
- `Truncate() Decimal`

`ToDigits` extends precision by adding trailing zeros, or reduces precision by truncation. It truncates toward zero and does not round.
`Round` extends precision by adding trailing zeros, or reduces precision by rounding to nearest, ties away from zero.
`Truncate` removes unnecessary trailing zeros while ensuring to never change the value.

Example:

```go
d, _ := decimal.NewFromString("1.999")

d.ToDigits(2) // 1.99
d.Round(2)    // 2.00
```

## Formatting And Conversion

- `String()` returns a decimal string without scientific notation
- `Float64()` converts to `float64`

`Float64()` is a convenience conversion and can lose precision, just like any decimal-to-float conversion.

### JSON

The type implements `json.Marshaler` and `json.Unmarshaler`.

Behavior:

- values are marshaled as JSON numbers, not JSON strings
- `null` unmarshals to zero
- quoted JSON strings are accepted only when the content is a plain decimal string
- escaped or otherwise non-plain JSON strings are not supported

Examples:

```json
123.45
```

```json
"123.45"
```

### SQL

The type implements `sql.Scanner` and `driver.Valuer`.

Supported `Scan` inputs:

- `nil`
- `[]byte`
- `string`
- `float64`
- `int64`
- `uint64`

Behavior:

- `nil` resets the receiver to zero
- `Value()` returns the decimal string form

### CBOR

The type implements CBOR marshaling and unmarshaling.

Behavior:

- integers are encoded as CBOR integers when `Digits == 0`
- values with a fractional part are encoded as RFC 8949 decimal fractions (tag 4)
- large mantissas are encoded using CBOR bignum tags when needed
- unmarshaling also accepts CBOR integers and CBOR float16/32/64 values

Interoperability note:

- the encoder uses preferred encoding
- the decoder may not support every valid alternate CBOR encoding of the same semantic value

## Designed Limits And Invariants

This package is intentionally low-level. The `Decimal` fields are exported, so callers can construct values directly. Public methods assume the following invariants are respected:

- `Digits` must be between `0` and `19`
- `Fraction` must match the declared scale in `Digits`
- negative zero is unsupported

If you construct values manually and break those invariants, behavior is undefined for this package and may include incorrect formatting, incorrect arithmetic, or panics.

Examples of supported manual values:

```go
decimal.Decimal{Integer: 12, Fraction: 34, Digits: 2}   // 12.34
decimal.Decimal{Integer: 12, Fraction: 340, Digits: 3}  // 12.340
decimal.Decimal{Integer: 12, Fraction: 34, Digits: 3}   // 12.034
```

Examples of unsupported manual values:

```go
decimal.Decimal{Fraction: 123, Digits: 2}  // Fraction value exceeds defined digits
decimal.Decimal{Fraction: 123, Digits: 20} // Digits exceeds 19
decimal.Decimal{Negative: true}            // Negative zero
```
