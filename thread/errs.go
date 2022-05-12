package thread

import "github.com/pkg/errors"

var (
	ErrClosed = errors.New("thread pool already closed")
)
