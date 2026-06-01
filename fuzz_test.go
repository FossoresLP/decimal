package decimal_test

import (
	"math"
	"testing"

	"github.com/fossoreslp/decimal"
)

// outsideUint64 reports whether a float64 falls outside the uint64 magnitude range.
// New() relies on Go's float-to-uint64 conversion, which is platform-defined on
// overflow — the resulting Decimal is meaningless in that region so we skip it.
func outsideUint64(f float64) bool {
	return f >= float64(math.MaxUint64) || f <= -float64(math.MaxUint64)
}

// New() documents that Inf and NaN round-trip to zero. The round-trip property
// would assert |0 - Inf| < eps, which is vacuously true but uninformative, so skip.
func nonFinite(f float64) bool {
	return math.IsNaN(f) || math.IsInf(f, 0)
}

// FuzzNew_Float64 asserts that New(f64).Float64() recovers the input within
// the precision Decimal documents: 18 fractional digits, plus float reconstruction noise.
func FuzzNew_Float64(f *testing.F) {
	for _, tc := range []float64{123.456, 0.123, 1e-18, 1234567890123456789.12345678901234567890, 3.3333333333, -1.5, 0} {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, v float64) {
		if nonFinite(v) || outsideUint64(v) {
			return
		}
		back := decimal.New(v).Float64()
		// Decimal truncates at 1e-18 absolute; reconstruction adds |v| * 2^-52 ≈ 2.3e-16 relative noise.
		limit := 1e-15 + math.Abs(v)*1e-14
		if math.Abs(back-v) > limit {
			t.Errorf("New(%g).Float64() = %.18g, drift %.18g exceeds %.18g", v, back, back-v, limit)
		}
	})
}

// FuzzNew_Float32 asserts that New(f32).Float64() recovers float64(f32) within
// Decimal's documented 10-digit precision for float32 inputs.
func FuzzNew_Float32(f *testing.F) {
	for _, tc := range []float32{123.456, 0.123, 1e-10, 1234567.5, 3.3333333, -1.5, 0} {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, v float32) {
		vf := float64(v)
		if nonFinite(vf) || outsideUint64(vf) {
			return
		}
		back := decimal.New(v).Float64()
		// Decimal stores 10 fractional digits for float32 input → 1e-10 absolute floor;
		// float32 mantissa is ~7 decimal digits → |v| * 1e-7 relative noise.
		limit := 1e-9 + math.Abs(vf)*1e-7
		if math.Abs(back-vf) > limit {
			t.Errorf("New(float32(%g)).Float64() = %.18g, drift %.18g exceeds %.18g", vf, back, back-vf, limit)
		}
	})
}

// inFixedRange reports whether a float64 falls strictly inside the Fixed range.
// NewFixed performs `int32(v * 100)`, which is platform-defined on int32 overflow.
// We pad away from the exact boundaries because float multiplication can push
// values like 21474836.47 across the edge.
func inFixedRange(f float64) bool {
	return f > -21474836.0 && f < 21474836.0
}

// FuzzNewFixed_Float64 asserts that NewFixed(f64).Float64() reconstructs the input
// to within one cent — the precision Fixed itself documents — plus float noise.
func FuzzNewFixed_Float64(f *testing.F) {
	for _, tc := range []float64{0, 0.29, -0.29, 1.5, -1.5, 0.01, 0.99, 123.45, 21474836.47, -21474836.48} {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, v float64) {
		if nonFinite(v) || !inFixedRange(v) {
			return
		}
		back := decimal.NewFixed(v).Float64()
		// Fixed has 0.01 precision; |v|*1e-9 covers float multiplication noise.
		limit := 0.01 + math.Abs(v)*1e-9
		if math.Abs(back-v) > limit {
			t.Errorf("NewFixed(%g).Float64() = %.6f, drift %.18g exceeds %.18g", v, back, back-v, limit)
		}
	})
}

// FuzzNewFixed_Float32 asserts the same round-trip property against float32 inputs,
// with a looser noise term reflecting float32's 7-digit mantissa.
func FuzzNewFixed_Float32(f *testing.F) {
	for _, tc := range []float32{0, 0.29, -0.29, 1.5, -1.5, 0.01, 0.99, 123.45, 1e6} {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, v float32) {
		vf := float64(v)
		if nonFinite(vf) || !inFixedRange(vf) {
			return
		}
		back := decimal.NewFixed(v).Float64()
		limit := 0.01 + math.Abs(vf)*1e-6
		if math.Abs(back-vf) > limit {
			t.Errorf("NewFixed(float32(%g)).Float64() = %.6f, drift %.18g exceeds %.18g", vf, back, back-vf, limit)
		}
	})
}
