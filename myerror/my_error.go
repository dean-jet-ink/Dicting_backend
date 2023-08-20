package myerror

import "errors"

var ErrRecordNotFound = errors.New("record not found")
var ErrMismatchedPassword = errors.New("mismatched password")
