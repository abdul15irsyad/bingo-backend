package util

func RandomSlice[T any](array []T) T {
	var element T
	if len(array) == 0 {
		return element
	}
	return array[RandomInt(0, len(array)-1)]
}

func MapSlice[T any, K any](slice *[]T, mapper func(item T) K) []K {
	result := make([]K, len(*slice))
	for i, value := range *slice {
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

func UniqueSlice[T any, R comparable](slice *[]T, predicate func(*T) R) []T {
	exists := []R{}
	result := []T{}
	for _, item := range *slice {
		value := predicate(&item)
		if exist := FindSlice(&exists, func(r *R) bool {
			return *r == value
		}); exist == nil {
			result = append(result, item)
		}
	}
	return result
}

func ReduceSlice[T any, U any](
	slice *[]T,
	reducer func(prev U, curr T) U,
	initial U,
) U {
	accumulator := initial
	for _, v := range *slice {
		accumulator = reducer(accumulator, v)
	}
	return accumulator
}
