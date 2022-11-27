package arrays

func Contains[T comparable](elements []T, target T) bool {
	for _, element := range elements {
		if target == element {
			return true
		}
	}
	return false
}
