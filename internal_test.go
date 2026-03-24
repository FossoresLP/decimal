package decimal

import "testing"

func BenchmarkFull(b *testing.B) {
	d := Decimal{}
	for b.Loop() {
		_ = d.full()
	}
}

func TestDiv128(t *testing.T) {
	tests := []struct {
		name                    string
		numHi, numLo            uint64
		denHi, denLo            uint64
		wantQ, wantRHi, wantRLo uint64
	}{
		{"den_hi_zero_num_hi_lt_den_lo", 1, 0, 0, 3, 6148914691236517205, 0, 1},
		{"den_hi_zero_num_hi_gte_den_lo", 5, 0, 0, 3, 12297829382473034410, 0, 2},
		{"top_bit_denominator_quotient_zero", 1 << 63, 0, 1 << 63, 1, 0, 1 << 63, 0},
		{"top_bit_denominator_quotient_one", 1 << 63, 1, 1 << 63, 1, 1, 0, 0},
		{"normalized_refinement", 3897563301579898974, 11706709894295411728, 1, 1122952782507670167, 3673912440421135436, 0, 16644754140121451324},
		{"normalized_refinement_overflow", 14138796634445629645, 12692267291059739667, 3, 8686375096939001715, 4073537142616101137, 3, 7080570136288529008},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQ, gotRHi, gotRLo := div128(tt.numHi, tt.numLo, tt.denHi, tt.denLo)
			if gotQ != tt.wantQ || gotRHi != tt.wantRHi || gotRLo != tt.wantRLo {
				t.Fatalf(
					"div128(%d,%d,%d,%d) = (%d,%d,%d), want (%d,%d,%d)",
					tt.numHi, tt.numLo, tt.denHi, tt.denLo,
					gotQ, gotRHi, gotRLo,
					tt.wantQ, tt.wantRHi, tt.wantRLo,
				)
			}
		})
	}
}
