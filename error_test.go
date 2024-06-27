package errors_test

import (
	stdErrors "errors"
	"fmt"
	"testing"

	"github.com/mrsoftware/errors"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	msg := "some message"

	err := errors.New(msg)
	assert.Equal(t, msg, err.Error())
}

func TestGetError(t *testing.T) {
	t.Parallel()

	t.Run("compatible error", func(t *testing.T) {
		t.Parallel()

		var (
			msg    = "some message"
			fields = []errors.Field{{}}
		)

		err := errors.New(msg, fields...)
		assert.Equal(t, fields, errors.GetFields(err))
	})

	t.Run("not compatible wrapper error", func(t *testing.T) {
		t.Parallel()

		var (
			msg = "some message"
		)
		mainErr := errors.New(msg)
		err := fmt.Errorf("msg: %w", mainErr)
		assert.Equal(t, mainErr, errors.GetError(err))
	})

	t.Run("not compatible error", func(t *testing.T) {
		t.Parallel()

		var (
			msg = "some message"
		)

		err := stdErrors.New(msg)
		willCreateErr := errors.Wrap(err, msg)
		assert.Equal(t, willCreateErr, errors.GetError(err))
	})
}

func TestWrap(t *testing.T) {
	t.Parallel()

	msg := "some message"
	cErr := stdErrors.New("cause")

	err := errors.Wrap(cErr, msg)
	assert.Equal(t, "some message: cause", err.Error())
}

func TestWrapf(t *testing.T) {
	t.Parallel()

	format := "some message id: %d"
	cErr := stdErrors.New("cause")

	err := errors.Wrapf(cErr, format, 10)
	assert.Equal(t, "some message id: 10: cause", err.Error())
}

func TestErrorf(t *testing.T) {
	t.Parallel()

	format := "some message id: %d"

	err := errors.Errorf(format, 10)
	assert.Equal(t, "some message id: 10", err.Error())
}

func TestErrorfWithFields(t *testing.T) {
	t.Parallel()

	format := "some message id: %d"
	fields := []errors.Field{{}}

	err := errors.ErrorfWithFields(format, []interface{}{10}, fields...)
	assert.Equal(t, "some message id: 10", err.Error())
	assert.Equal(t, fields, errors.GetFields(err))
}

func TestCause(t *testing.T) {
	t.Parallel()

	t.Run("error with no cause", func(t *testing.T) {
		err := stdErrors.New("cause")

		assert.Equal(t, err, errors.Cause(err))
	})

	t.Run("error support causer", func(t *testing.T) {
		cause := stdErrors.New("cause")
		err := errors.Wrap(cause, "wrapper")

		assert.Equal(t, cause, errors.Cause(err))
	})

	t.Run("error support causer", func(t *testing.T) {
		cause := errors.New("cause")
		err := errors.Wrap(cause, "wrapper")

		assert.Equal(t, cause, errors.Cause(err))
	})
}

func TestError_Format(t *testing.T) {
	t.Parallel()

	t.Run("format any this with no field", func(t *testing.T) {
		err := errors.New("some error")

		assert.Equal(t, "some error", fmt.Sprint(err))
	})

	t.Run("format q", func(t *testing.T) {
		err := errors.New("some error", errors.String("username", "mrsoftware"))

		assert.Equal(t, "\"some error\": [\"mrsoftware\"]", fmt.Sprintf("%q", err))
	})

	t.Run("format s", func(t *testing.T) {
		err := errors.New("some error", errors.String("username", "mrsoftware"))

		assert.Equal(t, "some error: [[username: mrsoftware]]", fmt.Sprintf("%s", err))
	})

	t.Run("format v", func(t *testing.T) {
		err := errors.New("some error", errors.String("username", "mrsoftware"))

		assert.Equal(t, "some error: [{Key: username, Value: mrsoftware}]", fmt.Sprintf("%v", err))
	})

	t.Run("format +v", func(t *testing.T) {
		err := errors.New("some error", errors.String("username", "mrsoftware"))

		assert.Equal(t, "some error: [{Key: username, Type: String, Value: mrsoftware}]", fmt.Sprintf("%+v", err))
	})

	t.Run("format #v", func(t *testing.T) {
		err := errors.New("some error", errors.String("username", "mrsoftware"))

		assert.Equal(t, "some error: []errors.Field{{username: \"mrsoftware\"}}", fmt.Sprintf("%#v", err))
	})
}

func TestAddFields(t *testing.T) {
	t.Parallel()

	t.Run("err is not custom, expect wrap it with custom and add field", func(t *testing.T) {
		err := stdErrors.New("standard error")

		field := errors.String("code", "value")
		expect := errors.Wrap(err, err.Error(), field)

		withField := errors.AddFields(err, field)

		assert.Equal(t, expect, withField)
	})
}
