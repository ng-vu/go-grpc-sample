package base

import "errors"

// ErrorCollector ...
type ErrorCollector struct {
	Count  int
	Errors []error
}

// Collect collects errors
func (e *ErrorCollector) Collect(err error) {
	e.Count++
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

// All returns error if all results collected are error
func (e *ErrorCollector) All() error {
	if len(e.Errors) == 0 || e.Count > len(e.Errors) {
		return nil
	}
	return e.concat()
}

// Any returns error if any result collected is error
func (e *ErrorCollector) Any() error {
	return e.concat()
}

func (e *ErrorCollector) concat() error {
	if len(e.Errors) == 0 {
		return nil
	}
	if len(e.Errors) == 1 {
		return e.Errors[0]
	}
	s := ""
	for i, err := range e.Errors {
		if i > 0 {
			s += "; "
		}
		s += err.Error()
	}
	return errors.New(s)
}
