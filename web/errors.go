package web

import "errors"

var (
	ErrConfigurationFileNotExist = errors.New("Configuration file not exist")
	ErrNoSuchView                = errors.New("No such view")
	ErrInvalidConfiguration      = errors.New("Invalid configuration")
)
