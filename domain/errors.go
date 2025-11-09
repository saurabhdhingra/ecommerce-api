package domain

import (
	"errors"
)

var (
	ErrNotFound				= errors.New("record not found")
	ErrInsufficientInv 		= errors.New("insufficient product inventory")
	ErrCartEmpty			= errors.New("cannot checkout empty cart")
	ErrInvalidCredentials	= errors.New("invalid username or password")
)