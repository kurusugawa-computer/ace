package cli

import "errors"

var (
	ErrUsage    = errors.New("Invalid usage")
	ErrInternal = errors.New("Internal error")
)
