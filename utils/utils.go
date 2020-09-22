package utils

// AbsI64 returns the absolute value of a 64 bit integer
func AbsI64(a int64) int64 {
	if a < 0 {
		a = a * -1
	}
	return a
}

// MinInt returns the minimum integer between two
func MinInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
