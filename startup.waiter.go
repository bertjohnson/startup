package startup

import (
	"context"
)

// Waiter is used to wait for asynchronous values.
type Waiter chan bool

// Wait waits until a value is available.
func (w Waiter) Wait(ctx context.Context) {
	<-w
}
