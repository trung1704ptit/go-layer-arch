package services

import "errors"

var (
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
	ErrUserAlreadyExists   = errors.New("user with that email already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrRefreshTokenInvalid = errors.New("could not refresh access token")
	ErrRefreshUserNotFound = errors.New("the user belonging to this token no logger exists")
	ErrPostAlreadyExists   = errors.New("post with that title already exists")
	ErrPostNotFound        = errors.New("no post with that title exists")
)
