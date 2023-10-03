package myerror

import "errors"

var ErrRecordNotFound = errors.New("record is not found")
var ErrMismatchedPassword = errors.New("mismatched password")
var ErrDuplicatedKey = errors.New("duplicate entry")
var ErrBindingFailure = errors.New("failed to bind")
var ErrValidation = errors.New("validation error")
var ErrCookieExpired = errors.New("cookie expired")
var ErrParsingFailure = errors.New("failed to parse")
var ErrUnverifiedEmail = errors.New("email is unverified")
var ErrFormFileNotFound = errors.New("form file is not found")
var ErrMissingJWT = errors.New("missing jwt token")
var ErrInvalidToken = errors.New("invalid token")
var ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
