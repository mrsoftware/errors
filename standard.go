package errors

import stdErr "errors"

// Is reports whether any error in err's tree matches target. (calling go standard errors.Is).
func Is(err, target error) bool {
	return stdErr.Is(err, target)
}

// As finds the first error in err's tree that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false. (calling go standard errors.As).
func As(err error, target interface{}) bool {
	return stdErr.As(err, target)
}
