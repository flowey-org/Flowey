package utils

func PopSlice[T any](s []T) []T {
	if s == nil || len(s) < 2 {
		return []T{}
	}
	return s[1:]
}
