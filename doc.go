// Package errors provide simple error handling primitives.
//
// The traditional error handling idiom in Go is roughly akin to
//
//	if err != nil {
//	        return err
//	}
//
// which when applied recursively up the call stack results in error reports
// without context or debugging information. The errors package allows
// programmers to add context to the failure path in their code in a way
// that does not destroy the original value of the error.
//
// # Adding context to an error
//
// The errors.Wrap function returns a new error that adds context to the
// original error by recording a stack trace at the point Wrap is called,
// together with the supplied message. For example
//
//	_, err := ioutil.ReadAll(r)
//	if err != nil {
//	        return errors.Wrap(err, "read failed")
//	}
//
// If additional control is required, the errors.WithStack and
// errors.WithMessage functions destructure errors.Wrap into its component
// operations: annotating an error with a stack trace and with a message,
// respectively.
//
// # Retrieving the cause of an error
//
// Using errors.Wrap constructs a stack of errors, adding context to the
// preceding error. Depending on the nature of the error it may be necessary
// to reverse the operation of errors.Wrap to retrieve the original error
// for inspection. Any error value which implements this interface
//
//	type causer interface {
//	        Cause() error
//	}
//
// can be inspected by errors.Cause. Errors.Cause will recursively retrieve
// the topmost error that does not implement causer, which is assumed to be
// the original cause. For example:
//
//	switch err := errors.Cause(err).(type) {
//	case *MyError:
//	        // handle specifically
//	default:
//	        // unknown error
//	}
//
// Although the causer interface is not exported by this package, it is
// considered a part of its stable public interface.
//
// # Formatted printing of errors
//
// All error values returned from this package implement fmt.Formatter and can
// be formatted by the fmt package. The following verbs are supported:
//
//  consider this error:
//  user := struct { Username string }{Username: "mrsoftware"}
//	errors.New("some error", errors.String("name", "mohammad"), errors.Reflect("user", user))
//
//	%q    print the error. If the error has a Cause it will be
//	      printed recursively. the fields value print as a list.
//
//  sample:
//
//  "some error": ["mohammad" {"mrsoftware"}]
//
//	%s    print the error. If the error has a Cause it will be
//	      printed recursively. the fields key/value print as a list
//
//  sample:
//
//  some error: [[name: mohammad] [user: {mrsoftware}]]

//	%v    extended format. If the error has a Cause, it will be
//	      printed recursively. the fields key/value print as a list like struct.
//
// sample:
//
// some error: [{Key: name, Value: mohammad} {Key: user, Value: {Username:mrsoftware}}]
//
//	%+v   extended format. If the error has a Cause, it will be
//	  printed recursively. the fields key/type/value print as a list like struct.
//
// sample:
//
// some error: [{Key: name, Type: String, Value: mohammad} {Key: user, Type: Reflect, Value: {Username:mrsoftware}}]
//
//	%#v   extended format. If the error has a Cause, it will be
//	  printed recursively. the fields key/type/value print as a list like struct.
//
// sample:
//
// some error: []errors.Field{{name: "mohammad"}, {user: struct { Username string }{Username:"mrsoftware"}}}
package errors
