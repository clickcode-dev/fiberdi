package fiberdi

func ternary[T any](condition bool, ifTrue T, ifFalse T) T {
	if condition {
		return ifTrue
	}

	return ifFalse
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) && len(v) > 7 {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
