package compare

import (
	"fmt"

	"go.uber.org/multierr"
)

func Comparable[T comparable](label string, v T) Comparer {
	return NewComparer(v, func(v1, v2 T) error {
		if v1 != v2 {
			return labeledError(fmt.Errorf("%#v != %#v", v1, v2), label)
		}
		return nil
	})
}

func ComparableSlice[T comparable](label string, v []T) Comparer {
	return Slice(label, v, func(v T) C {
		return C{
			Comparable("", v),
		}
	})
}

func ComparablePointer[T comparable](label string, v *T) Comparer {
	return NewComparer(v, func(v1, v2 *T) error {
		var err error
		if v1 == nil || v2 == nil {
			if v1 != v2 {
				err = fmt.Errorf("%#v != %#v", v1, v2)
			}
		} else if *v1 != *v2 {
			err = fmt.Errorf("%#v != %#v", *v1, *v2)
		}
		return labeledError(err, label)
	})
}

func Func[T any](label string, v T, diffFunc func(v1, v2 T) error) Comparer {
	return NewComparer(v, func(v1, v2 T) error {
		return labeledError(diffFunc(v1, v2), label)
	})
}

func Nest[T any](label string, v T, f func(T) C) Comparer {
	return NewComparer(v, func(v1, v2 T) error {
		return labeledError(diff(v1, v2, f), label)
	})
}

func Slice[T any](label string, v []T, f func(T) C) Comparer {
	return NewComparer(v, func(vs1, vs2 []T) error {
		if len(vs1) != len(vs2) {
			return fmt.Errorf("len: %d != %d", len(vs1), len(vs2))
		}
		var errs error
		for i, v1 := range vs1 {
			v2 := vs2[i]
			label := fmt.Sprintf("%s[%d]", label, i)
			errs = multierr.Append(errs, labeledError(diff(v1, v2, f), label))
		}
		return errs
	})
}

func NewComparer[T any](value T, diffFunc func(v1, v2 T) error) Comparer {
	return &comparer[T]{
		value:    value,
		diffFunc: diffFunc,
	}
}

type comparer[T any] struct {
	value    T
	diffFunc func(v1, v2 T) error
}

func (c *comparer[T]) Value() any {
	return c.value
}

func (c *comparer[T]) Diff(a1, a2 any) error {
	v1, ok := a1.(T)
	if !ok {
		return fmt.Errorf("invalid type: %T", a1)
	}
	v2, ok := a2.(T)
	if !ok {
		return fmt.Errorf("invalid type: %T", a2)
	}
	return c.diffFunc(v1, v2)
}
