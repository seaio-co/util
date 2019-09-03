package conv

// SliceContains 判断任意类型s1中是否包含v
func SliceContains(sl []interface{}, v interface{}) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

// SliceContainsInt 判断Int类型s1中是否包含v
func SliceContainsInt(sl []int, v int) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
