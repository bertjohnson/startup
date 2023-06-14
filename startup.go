// Package startup discovers configuration values when the service initializes.
package startup

import (
	"context"
	"sync"
)

const (
	// PackageType is the package type.
	PackageType = "startup"
)

var (
	// Whether the initializer has run.
	initialized bool

	// Mutex to manage access configuration objects and waiters maps.
	mutex sync.Mutex

	// Map of objects that are available.
	objects = map[string]bool{}

	// Waiters track methods waiting for a configuration object.
	waiters = map[string][]Waiter{}
)

// Ready declares that environment variables are ready and begins initialization.
func Ready() error {
	// Prevent initializing multiple times.
	if initialized {
		return nil
	}
	initialized = true

	// Read local environment variable overrides.
	err := readEnvFiles()
	if err != nil {
		return err
	}

	// Start loop to check environment variables.
	go checkEnv()

	// Publish that the environment configuration is available.
	Publish(context.Background(), PackageType)

	return nil
}

// EnvWait waits for one or more environment variables to be set.
func EnvWait(ctx context.Context, name ...string) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(name))
	for _, n := range name {
		go func(n string) {
			GetEnvWaiter(ctx, n).Wait(ctx)
			waitGroup.Done()
		}(n)
	}
	waitGroup.Wait()
}

// GetWaiter returns a waiter for a configuration object.
func GetWaiter(ctx context.Context, objectID string) (waiter Waiter) {
	waiter = make(Waiter)

	mutex.Lock()
	if _, ok := objects[objectID]; ok {
		// If the object is already available, resolve it immediately.
		go func() {
			waiter <- true
		}()
	} else {
		// If the object is not available yet, wait for it.
		objectWaiters, ok := waiters[objectID]
		if ok {
			waiters[objectID] = append(objectWaiters, waiter)
		} else {
			waiters[objectID] = []Waiter{waiter}
		}
	}
	mutex.Unlock()

	return waiter
}

// GetEnvWaiter returns a waiter for an environment variable.
func GetEnvWaiter(ctx context.Context, name string) (waitChannel Waiter) {
	return GetWaiter(ctx, "env."+name)
}

// Publish publishes an object ID as available.
func Publish(ctx context.Context, objectID string) {
	mutex.Lock()
	// Resolve channels waiting for this key.
	if keyWaiters, ok := waiters[objectID]; ok {
		for _, waiter := range keyWaiters {
			go func(waiter Waiter) {
				waiter <- true
			}(waiter)
		}
		delete(waiters, objectID)
	}
	objects[objectID] = true
	mutex.Unlock()
}

// Wait waits for one or more configuration objects to be initialized.
func Wait(ctx context.Context, objectID ...string) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(objectID))
	for _, o := range objectID {
		go func(o string) {
			GetWaiter(ctx, o).Wait(ctx)
			waitGroup.Done()
		}(o)
	}
	waitGroup.Wait()
}
