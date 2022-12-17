package compare

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"
)

func labeledError(err error, label string) error {
	if err == nil {
		return nil
	}
	return &labeledErr{
		err:   err,
		label: label,
	}
}

type labeledErr struct {
	err   error
	label string
}

func (e *labeledErr) Unwrap() error {
	return e.err
}

func (e *labeledErr) Error() string {
	return e.err.Error()
}

func formatDiffWithLabel(err error, label string) string {
	errs := multierr.Errors(err)
	if len(errs) == 0 {
		return ""
	}
	if len(errs) > 1 {
		var b []byte
		for i, err := range errs {
			if i > 0 {
				b = append(b, '\n')
			}
			b = append(b, formatDiffWithLabel(err, label)...)
		}
		return string(b)
	}

	err = errs[0]

	var lerr *labeledErr
	if errors.As(err, &lerr) {
		label := label
		if label == "" {
			label = lerr.label
		} else if lerr.label != "" {
			label += "." + lerr.label
		}
		return formatDiffWithLabel(lerr.err, label)
	}

	if label == "" {
		return err.Error()
	}
	return fmt.Sprintf("%s: %s", label, err.Error())
}

func formatDiff(err error) string {
	return formatDiffWithLabel(err, "")
}
