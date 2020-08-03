package main

func AbsI64(a int64) int64 {
	if a < 0 {
		a = a * -1
	}
	return a
}

func minint(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
