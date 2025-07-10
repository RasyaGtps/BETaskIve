package models

import "errors"

var (
	ErrPasswordTooShort = errors.New("password must be at least 6 characters long")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrEmailExists      = errors.New("email already exists")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUserNotFound     = errors.New("user not found")
	ErrUnauthorized     = errors.New("unauthorized access")
	ErrForbidden        = errors.New("forbidden access")
) 