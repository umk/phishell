package slicesx

func Unique[V comparable](s []V) []V {
	m := make(map[V]bool)

	j := 0
	for i, v := range s {
		if _, ok := m[v]; ok {
			i++
		} else {
			m[v] = true
			s[i-j] = v
		}
	}

	return s[:len(s)-j]
}
