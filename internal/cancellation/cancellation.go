// Package cancellation owns Authoring's shared cancellation identity without
// exposing a second public sentinel.
package cancellation

import "errors"

// ErrCanceled is re-exported by packages that expose cancellation publicly.
var ErrCanceled = errors.New("authoring canceled")
