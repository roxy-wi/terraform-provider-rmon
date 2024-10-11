package rmon

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func intToBool(i float64) bool {
	return i == 1
}
