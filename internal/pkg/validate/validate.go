package validate

func Len(l int, data ...string) bool {
	for _, d := range data {
		if len(d) < l {
			return false
		}
	}

	return true
}
