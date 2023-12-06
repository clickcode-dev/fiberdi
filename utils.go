package fiberdi

func array[T any](t T) *T {
	return &t
}
