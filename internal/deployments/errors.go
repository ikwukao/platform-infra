package deployments

import "errors"

var (
	// ErrNotFound indicates that a deployment does not exist.
	ErrNotFound = errors.New("deployment not found")
)
