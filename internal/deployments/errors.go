package deployments

import "errors"

var (
	// ErrNotFound indicates that a deployment does not exist.
	ErrNotFound = errors.New("deployment not found")
)

// ErrInvalidStatus indicates that a deployment status is not supported.
var ErrInvalidStatus = errors.New("invalid deployment status")

// ValidStatus reports whether status is a supported deployment lifecycle state.
func ValidStatus(status string) bool {
	switch status {
	case StatusPending,
		StatusRunning,
		StatusSucceeded,
		StatusFailed:
		return true
	default:
		return false
	}
}
