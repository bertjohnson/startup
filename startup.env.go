package startup

import (
	"math/rand"
	"os"
	"strings"
	"time"
)

// increment defines sleep increments.
type increment struct {
	FixedAmount    time.Duration
	VariableAmount time.Duration
}

var (
	// Standard jitter increments.
	increments = []increment{
		{FixedAmount: 250 * time.Millisecond, VariableAmount: 250 * time.Millisecond},
		{FixedAmount: 500 * time.Millisecond, VariableAmount: 500 * time.Millisecond},
		{FixedAmount: 1 * time.Second, VariableAmount: 1 * time.Second},
		{FixedAmount: 2 * time.Second, VariableAmount: 1 * time.Second},
		{FixedAmount: 4 * time.Second, VariableAmount: 1 * time.Second},
		{FixedAmount: 8 * time.Second, VariableAmount: 1 * time.Second},
		{FixedAmount: 16 * time.Second, VariableAmount: 1 * time.Second},
		{FixedAmount: 32 * time.Second, VariableAmount: 2 * time.Second},
		{FixedAmount: 64 * time.Second, VariableAmount: 2 * time.Second},
	}

	// Length of jitter increments.
	incrementsMax = len(increments) - 1

	// Current jitter increment.
	currentIncrement = 0
)

// checkEnv checks environment variables with a backoff schedule.
func checkEnv() {
	for {
		mutex.Lock()
		vars := os.Environ()
		resolved := false
		for _, envVar := range vars {
			envVarParts := strings.SplitN(envVar, "=", 2)
			if len(envVarParts) < 2 {
				continue
			}

			// Resolve channels waiting for this key.
			key := "env." + envVarParts[0]
			if keyWaiters, ok := waiters[key]; ok {
				for _, waiter := range keyWaiters {
					go func(waiter Waiter) {
						waiter <- true
					}(waiter)
				}
				delete(waiters, key)
				resolved = true
			}
			objects[key] = true
		}
		mutex.Unlock()

		// Sleep a jittered amount of time.
		if resolved {
			currentIncrement = 0
		} else {
			if currentIncrement < incrementsMax {
				currentIncrement++
			}
		}

		delay := increments[currentIncrement].FixedAmount + time.Duration(rand.Int31n(int32(increments[currentIncrement].VariableAmount)))
		time.Sleep(delay)
	}
}
