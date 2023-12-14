package errors_test

import (
	"testing"

	"github.com/mrsoftware/errors"
)

func TestMy(t *testing.T) {
	stValue := struct {
		Username string
		Password string
	}{Username: "mrsoftware", Password: "kanman"}
	err := errors.New("my error", errors.String("myString", "string value"), errors.Reflect("user", stValue))

	// t.Logf("%q", err)
	// t.Logf("%s", err)
	// t.Logf("%v", err)
	t.Logf("%+s", err)

	t.Fail()
}

//
// func TestNew(t *testing.T) {
// 	msg := "some message"
//
// 	err := errors.New(msg)
// 	assert.Equal(t, msg, err.Error())
// }
//
// func TestGetError(t *testing.T) {
// 	var (
// 		msg    = "some message"
// 		fields = []errors.Field{{}}
// 	)
//
// 	err := errors.New(msg, fields...)
// 	assert.Equal(t, fields, errors.GetFields(err))
// }
//
// func TestGetErrorNotCompatible(t *testing.T) {
// 	var (
// 		msg    = "some message"
// 		fields = []errors.Field{{}}
// 	)
// 	mainErr := errors.New(msg, fields...)
// 	err := fmt.Errorf("msg", mainErr)
// 	assert.Equal(t, mainErr, errors.GetError(err))
// }
//
// func TestWrap(t *testing.T) {
// 	msg := "some message"
// 	cErr := stdErr.New("cause")
//
// 	err := errors.Wrap(cErr, msg)
// 	assert.Equal(t, "some message: cause", err.Error())
// }
//
// func TestWrapf(t *testing.T) {
// 	format := "some message id: %d"
// 	cErr := stdErr.New("cause")
//
// 	err := errors.Wrapf(cErr, format, 10)
// 	assert.Equal(t, "some message id: 10: cause", err.Error())
// }
//
// func TestErrorf(t *testing.T) {
// 	format := "some message id: %d"
//
// 	err := errors.Errorf(format, 10)
// 	assert.Equal(t, "some message id: 10", err.Error())
// }
//
// func TestErrorfWithFields(t *testing.T) {
// 	format := "some message id: %d"
// 	fields := []errors.Field{{}}
//
// 	err := errors.ErrorfWithFields(format, []interface{}{10}, fields...)
// 	assert.Equal(t, "some message id: 10", err.Error())
// 	assert.Equal(t, fields, errors.GetFields(err))
// }
