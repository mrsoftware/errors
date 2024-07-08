package errors

import (
	stdErr "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiError_Error(t *testing.T) {
	t.Run("no error in list, expect to return empty string", func(t *testing.T) {
		err := NewMultiError()

		expected := ""

		assert.Equal(t, expected, err.Error())
	})

	t.Run("some error is set in list, expect to merge them", func(t *testing.T) {
		err := NewMultiError(stdErr.New("error 1"), stdErr.New("error 2"))

		expected := "error 1 | error 2"

		assert.Equal(t, expected, err.Error())
	})
}

func TestMultiError_Err(t *testing.T) {
	t.Run("have no error, expect to get nil", func(t *testing.T) {
		err := NewMultiError(nil, nil)

		assert.Nil(t, err.Err())
	})

	t.Run("have errors, expect to get it", func(t *testing.T) {
		err := NewMultiError(stdErr.New("some error"))

		assert.Equal(t, err, err.Err())
	})
}

func TestMultiError_Errors(t *testing.T) {
	errList := []error{stdErr.New("some error")}

	err := NewMultiError(errList...)

	assert.Equal(t, errList, err.Errors())
}

func TestMultiError_SafeAdd(t *testing.T) {
	t.Run("add error to errors list", func(t *testing.T) {
		error1 := stdErr.New("error 1")
		err := NewMultiError()

		err.SafeAdd(error1)

		expected := &MultiError{errors: []error{error1}}

		assert.Equal(t, expected, err)
	})

	t.Run("the added error is nil, expect to ignore it", func(t *testing.T) {
		err := NewMultiError()

		err.SafeAdd(nil)

		expected := &MultiError{}

		assert.Equal(t, expected, err)
	})
}

func TestMultiError_LenSafe(t *testing.T) {
	errList := []error{stdErr.New("some error")}

	err := NewMultiError(errList...)

	assert.Equal(t, len(errList), err.SafeLen())
}

func TestMultiError_Unwrap(t *testing.T) {
	t.Run("no error in list, expect to return nil", func(t *testing.T) {
		err := NewMultiError()

		assert.Nil(t, err.Unwrap())
	})

	t.Run("errors is set, expect to return first error as root error", func(t *testing.T) {
		error1 := stdErr.New("error 1")
		error2 := stdErr.New("error 2")

		err := NewMultiError(error1, error2)

		assert.Equal(t, error1, err.Unwrap())
	})
}

func TestMultiError_Is(t *testing.T) {
	t.Run("requested err found in list, expect to get true", func(t *testing.T) {
		error1 := stdErr.New("error 1")
		error2 := stdErr.New("error 2")

		err := NewMultiError(error1, error2)

		assert.True(t, stdErr.Is(err, error1))
	})

	t.Run("requested err is not in list, expect to get false", func(t *testing.T) {
		error1 := stdErr.New("error 1")
		error2 := stdErr.New("error 2")
		error3 := stdErr.New("error 3")

		err := NewMultiError(error1, error2)

		assert.False(t, stdErr.Is(err, error3))
	})
}

func TestMultiError_As(t *testing.T) {
	t.Run("requested err found in list, expect to get true", func(t *testing.T) {
		error1 := stdErr.New("error 1")
		error2 := New("my error")

		err := NewMultiError(error1, error2)

		myErr := &Error{}

		assert.True(t, stdErr.As(err, &myErr))
		assert.Equal(t, error2, myErr)
	})

	t.Run("requested err is not in list, expect to get false", func(t *testing.T) {
		error1 := stdErr.New("error 1")
		error2 := stdErr.New("error 2")

		err := NewMultiError(error1, error2)

		myErr := &Error{}

		assert.False(t, stdErr.As(err, &myErr))
	})
}
