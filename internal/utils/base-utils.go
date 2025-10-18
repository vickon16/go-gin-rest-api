package utils

import (
	"context"
	"time"
)

// CreateContext returns a new cancellable context with a default timeout.
func CreateContext(timeOut ...time.Duration) (context.Context, context.CancelFunc) {
	t := 5 * time.Second
	if len(timeOut) > 0 {
		t = timeOut[0]
	}
	return context.WithTimeout(context.Background(), t)
}
