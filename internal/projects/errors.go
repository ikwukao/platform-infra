package projects

import "errors"

// ErrNotFound indicates that the requested project does not exist.
var ErrNotFound = errors.New("project not found")
