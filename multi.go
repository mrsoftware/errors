package errors

import (
	"bytes"
	"errors"
	"sync"
)

const (
	defaultErrorGroupSeparator = " | "
)

// MultiError is list of errors.
type MultiError struct {
	errors []error
	mx     sync.Mutex
}

// NewMultiError create new MultiError error.
func NewMultiError(errors ...error) *MultiError {
	multi := &MultiError{}
	for _, err := range errors {
		multi.Add(err)
	}

	return multi
}

// Error sum of all errors.
func (m *MultiError) Error() string {
	if m.Len() == 0 {
		return ""
	}

	buffer := &bytes.Buffer{}
	for index, err := range m.errors {
		buffer.WriteString(err.Error())

		if index < len(m.errors)-1 {
			buffer.WriteString(defaultErrorGroupSeparator)
		}
	}

	return buffer.String()
}

// Err of the multi error, return nil if no error is set.
func (m *MultiError) Err() error {
	if m.Len() == 0 {
		return nil
	}

	return m
}

// Errors return all list of errors.
func (m *MultiError) Errors() []error {
	return m.errors
}

// Add new error to list.
func (m *MultiError) Add(err error) {
	if err == nil {
		return
	}

	m.errors = append(m.errors, err)
}

// SafeAdd is like add but concurrent safe.
func (m *MultiError) SafeAdd(err error) {
	m.mx.Lock()
	m.Add(err)
	m.mx.Unlock()
}

// Len of errors.
func (m *MultiError) Len() int {
	return len(m.errors)
}

// SafeLen is like len but concurrent safe.
func (m *MultiError) SafeLen() int {
	m.mx.Lock()
	defer m.mx.Unlock()

	return len(m.errors)
}

// Unwrap return cause error (first error in list).
func (m *MultiError) Unwrap() error {
	if len(m.errors) == 0 {
		return nil
	}

	m.mx.Lock()
	defer m.mx.Unlock()

	return m.errors[0]
}

// Is check errors for match.
func (m *MultiError) Is(err error) bool {
	m.mx.Lock()
	defer m.mx.Unlock()

	for i := range m.errors {
		if errors.Is(err, m.errors[i]) {
			return true
		}
	}

	return false
}
