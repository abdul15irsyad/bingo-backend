package util

func MapSlice[T any, K any](slice []T, mapper func(item T) K) []K {
	result := make([]K, len(slice))
	for i, value := range slice {
		result[i] = mapper(value)
	}
	return result
}

func FindSlice[T any](slice *[]T, predicate func(*T) bool) *T {
	for _, item := range *slice {
		if predicate(&item) {
			return &item
		}
	}
	return nil
}

func FilterSlice[T any](slice *[]T, predicate func(*T) bool) []T {
	var result []T
	for _, item := range *slice {
		if predicate(&item) {
			result = append(result, item)
		}
	}
	return result
}
