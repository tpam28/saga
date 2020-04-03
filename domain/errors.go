package domain

import (
	"errors"
)

var (
	SequenceErr          = errors.New("invalid sequence")
	KeysInvalidFormatErr = errors.New("keys has invalid format")
	ErrNoRows            = errors.New("step contain to many args")
	InvalidType          = errors.New("invalid type")
)
