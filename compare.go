package compare

import "go.uber.org/multierr"

type C []Comparer

type Comparer interface {
	Value() any
	Diff(v1, v2 any) error
}

func Diff[T any](v1, v2 T, f func(T) C) string {
	return formatDiff(diff(v1, v2, f))
}

func diff[T any](v1, v2 T, f func(T) C) error {
	cs1 := f(v1)
	cs2 := f(v2)
	var errs error
	for i, c1 := range cs1 {
		c2 := cs2[i]
		if err := c1.Diff(c1.Value(), c2.Value()); err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	return errs
}
