package swlib

import (
	"errors"
)

var (
	// ErrInvalidAttributeType is returned when an unknown attribute type is encountered.
	ErrInvalidAttributeType = errors.New("invalid attribute type")

	// ErrIncompatibleNetlink is returned when an incompatible netlink response is detected.
	ErrIncompatibleNetlink = errors.New("incompatible netlink response")
)
