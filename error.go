package errors

import (
	"errors"
	"fmt"
)

// Error is an internal error with fields capabilities.
type Error struct {
	cause  error
	msg    string
	fields []Field
}

// New create a new error.
func New(msg string, fields ...Field) error {
	return &Error{msg: msg, fields: fields}
}

// Wrap creates a new error with given cause.
func Wrap(cause error, msg string, fields ...Field) error {
	return &Error{cause: cause, msg: msg, fields: fields}
}

// Wrapf is like Wrap, but it does format.
func Wrapf(cause error, format string, args ...interface{}) error {
	return Wrap(cause, fmt.Sprintf(format, args...))
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return Wrap(nil, fmt.Sprintf(format, args...))
}

// ErrorfWithFields is like Errorf and also support Field.
func ErrorfWithFields(format string, args []interface{}, fields ...Field) error {
	return Wrap(nil, fmt.Sprintf(format, args...), fields...)
}

// Error return error string.
func (e *Error) Error() string {
	if e.cause == nil {
		return e.msg
	}

	return e.msg + ": " + e.cause.Error()
}

// Cause return the cause if error.
func (e *Error) Cause() error { return e.cause }

// Unwrap return the cause if error.
func (e *Error) Unwrap() error { return e.cause }

// GetError check if the error is Error, create an empty Error if not.
func GetError(err error) (Err *Error) {
	var custom *Error
	ok := errors.As(err, &custom)
	if ok {
		return custom
	}

	return &Error{msg: err.Error(), cause: err}
}

// Cause return main error.
func Cause(err error) error {
	type causer interface {
		Unwrap() error
		Error() string
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}

		myCause := cause.Unwrap()
		if myCause == nil {
			err = cause

			break
		}

		err = myCause
	}

	return err
}

// Format is implement the fmt.Formatter for Error.
func (e *Error) Format(state fmt.State, verb rune) {
	if len(e.fields) == 0 {
		fmt.Fprint(state, e.Error())

		return
	}

	switch verb {
	case 'v':
		if state.Flag('+') {
			fmt.Fprintf(state, "%+v: %+v", e.Error(), e.fields)

			return
		}

		if state.Flag('#') {
			fmt.Fprintf(state, "%v: %#v", e.Error(), e.fields)

			return
		}

		fmt.Fprintf(state, "%v: %v", e.Error(), e.fields)
	case 's':
		fmt.Fprintf(state, "%s: %s", e.Error(), e.fields)
	case 'q':
		fmt.Fprintf(state, "%q: %q", e.Error(), e.fields)
	}
}

// AddFields to the passed error.
// passed error must be Error, if not, a new Error will create.
func AddFields(err error, fields ...Field) error {
	customError := GetError(err)
	customError.fields = append(customError.fields, fields...)

	return customError
}
