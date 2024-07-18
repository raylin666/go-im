package utils

func InSliceByInt(value int, values []int) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}

func InSliceByString(value string, values []string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}
