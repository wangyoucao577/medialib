// Package fixedpoint provides utils for fixed-point value calculation.
package fixedpoint

// From16x16 converts fixed-point 16.16 to normal decimal.
func From16x16(v float64) float64 {
	return v / 65536.0
}
