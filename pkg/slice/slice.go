package slice

func Concat[T any](slices ...[]T) []T {
	length := 0
	for _, values := range slices {
		length += len(values)
	}

	output := make([]T, 0, length)
	for _, values := range slices {
		output = append(output, values...)
	}

	return output
}

func Filter[T any](in []T, fn func(T) bool) []T {
	out := make([]T, 0, len(in))

	for i := 0; i < len(in); i++ {
		if fn(in[i]) {
			out = append(out, in[i])
		}
	}

	return out
}

func Map[In, Out any](in []In, fn func(In) Out) []Out {
	out := make([]Out, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = fn(in[i])
	}

	return out
}

func All[T any](in []T, fn func(T) bool) bool {
	for i := 0; i < len(in); i++ {
		if !fn(in[i]) {
			return false
		}
	}

	return true
}

func Any[T any](in []T, fn func(T) bool) bool {
	for i := 0; i < len(in); i++ {
		if fn(in[i]) {
			return true
		}
	}

	return false
}
