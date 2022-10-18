package utils

func Zero[T any]() T {
	return *new(T)
}

func CastOrDefault[T any](a any) T {
	b, ok := a.(T)
	if !ok {
		return *new(T)
	}
	return b
}

func Ptr[T any](value T) *T {
	return &value
}
