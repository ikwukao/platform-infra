package services

import "errors"

// ErrNotFound indicates that the requested service does not exist.
var ErrNotFound = errors.New("service not found")
