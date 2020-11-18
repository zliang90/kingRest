package slice

// HasInt judge v is in int slice
func HasInt(inits []int, v int) bool {
	for _, i := range inits {
		if i == v {
			return true
		}
	}
	return false
}

// HasString judge v is in string slice
func HasString(sl []string, v string) bool {
	for _, i := range sl {
		if i == v {
			return true
		}
	}
	return false
}
