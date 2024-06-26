// nolint
package errors_test

import (
	"fmt"
	"testing"

	"github.com/mrsoftware/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetChainFields(t *testing.T) {
	var (
		msg    = "some message"
		field1 = errors.Field{Key: "field1"}
		field2 = errors.Field{Key: "field2"}
		field3 = errors.Field{Key: "field3"}
	)

	err1 := errors.New(msg, field1)
	err2 := errors.Wrap(err1, msg, field2)
	err3 := errors.Wrap(err2, msg, field3)
	err4 := errors.Wrap(err3, msg)

	assert.Equal(t, []errors.Field{field3, field2, field1}, errors.GetChainFields(err4))
}

func TestFindFieldInChain(t *testing.T) {
	var (
		msg    = "some message"
		field1 = errors.Field{Key: "field1"}
		field2 = errors.Field{Key: "field2"}
		field3 = errors.Field{Key: "field3"}
	)

	err1 := errors.New(msg, field1)
	err2 := errors.Wrap(err1, msg, field2)
	err3 := errors.Wrap(err2, msg, field3)
	err4 := errors.Wrap(err3, msg)

	assert.Equal(t, "field4", errors.FindFieldInChain("field4", err4).Key)
	assert.Equal(t, errors.FieldTypeReflect, errors.FindFieldInChain("field4", err4).Type)
	assert.Equal(t, field2, errors.FindFieldInChain("field2", err4))
	assert.Equal(t, field1, errors.FindFieldInChain("field1", err1))
}

func TestField_Format(t *testing.T) {
	t.Parallel()

	t.Run("format q", func(t *testing.T) {
		field := errors.String("username", "mrsoftware")

		assert.Equal(t, "\"mrsoftware\"", fmt.Sprintf("%q", field))
	})

	t.Run("format s", func(t *testing.T) {
		field := errors.String("username", "mrsoftware")

		assert.Equal(t, "[username: mrsoftware]", fmt.Sprintf("%s", field))
	})

	t.Run("format v", func(t *testing.T) {
		field := errors.String("username", "mrsoftware")

		assert.Equal(t, "{Key: username, Value: mrsoftware}", fmt.Sprintf("%v", field))
	})

	t.Run("format +v", func(t *testing.T) {
		field := errors.String("username", "mrsoftware")

		assert.Equal(t, "{Key: username, Type: String, Value: mrsoftware}", fmt.Sprintf("%+v", field))
	})

	t.Run("format #v", func(t *testing.T) {
		field := errors.String("username", "mrsoftware")

		assert.Equal(t, "{username: \"mrsoftware\"}", fmt.Sprintf("%#v", field))
	})
}
