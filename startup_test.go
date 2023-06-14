package startup

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMain runs tests.
func TestMain(m *testing.M) {
	Ready()

	// Run tests.
	retCode := m.Run()

	// Exit.
	os.Exit(retCode)
}

// TestEnvWait tests EnvWait().
func TestEnvWait(t *testing.T) {
	go func() {
		time.Sleep(time.Second)
		os.Setenv("wait1", "1")
		os.Setenv("wait2", "2")
		os.Setenv("wait3", "3")
	}()
	EnvWait(context.Background(), "wait1", "wait2", "wait3")
}

// TestGetWaiter tests TestGetWaiter() and Publish().
func TestGetWaiter(t *testing.T) {
	callbackMutex := sync.Mutex{}
	ctx := context.Background()
	firstCallbackCount := 0
	go func() {
		waiter := GetWaiter(ctx, "valA")
		waiter.Wait(ctx)
		callbackMutex.Lock()
		firstCallbackCount++
		callbackMutex.Unlock()
	}()

	secondCallbackCount := 0
	go func() {
		waiter := GetWaiter(ctx, "valA")
		waiter.Wait(ctx)
		callbackMutex.Lock()
		secondCallbackCount++
		callbackMutex.Unlock()
	}()

	thirdCallbackCount := 0
	go func() {
		waiter := GetWaiter(ctx, "valB")
		waiter.Wait(ctx)
		callbackMutex.Lock()
		thirdCallbackCount++
		callbackMutex.Unlock()
	}()

	callbackMutex.Lock()
	assert.Equal(t, 0, firstCallbackCount)
	assert.Equal(t, 0, secondCallbackCount)
	assert.Equal(t, 0, thirdCallbackCount)
	callbackMutex.Unlock()

	Publish(ctx, "valA")
	time.Sleep(time.Second)
	callbackMutex.Lock()
	assert.Equal(t, 1, firstCallbackCount)
	assert.Equal(t, 1, secondCallbackCount)
	assert.Equal(t, 0, thirdCallbackCount)
	callbackMutex.Unlock()

	Publish(ctx, "valB")
	time.Sleep(time.Second)
	callbackMutex.Lock()
	assert.Equal(t, 1, firstCallbackCount)
	assert.Equal(t, 1, secondCallbackCount)
	assert.Equal(t, 1, thirdCallbackCount)
	callbackMutex.Unlock()
}

// TestGetEnvWaiter tests GetEnvWaiter().
func TestGetEnvWaiter(t *testing.T) {
	callbackMutex := sync.Mutex{}
	firstCallback := false
	ctx := context.Background()
	go func() {
		waiter := GetEnvWaiter(ctx, "VAL1")
		waiter.Wait(ctx)
		callbackMutex.Lock()
		firstCallback = true
		callbackMutex.Unlock()
	}()

	secondCallback := false
	go func() {
		waiter := GetEnvWaiter(ctx, "VAL1")
		waiter.Wait(ctx)
		callbackMutex.Lock()
		secondCallback = true
		callbackMutex.Unlock()
	}()

	thirdCallback := false
	go func() {
		waiter := GetEnvWaiter(ctx, "VAL2")
		waiter.Wait(ctx)
		callbackMutex.Lock()
		thirdCallback = true
		callbackMutex.Unlock()
	}()

	callbackMutex.Lock()
	assert.Equal(t, false, firstCallback)
	assert.Equal(t, false, secondCallback)
	assert.Equal(t, false, thirdCallback)
	callbackMutex.Unlock()

	os.Setenv("VAL1", "ABC")
	time.Sleep(time.Second)
	callbackMutex.Lock()
	assert.Equal(t, true, firstCallback)
	assert.Equal(t, true, secondCallback)
	assert.Equal(t, false, thirdCallback)
	callbackMutex.Unlock()

	os.Setenv("VAL2", "ABC")
	time.Sleep(time.Second)
	callbackMutex.Lock()
	assert.Equal(t, true, firstCallback)
	assert.Equal(t, true, secondCallback)
	assert.Equal(t, true, thirdCallback)
	callbackMutex.Unlock()
}

// TestWait tests Wait().
func TestWait(t *testing.T) {
	ctx := context.Background()
	go func() {
		time.Sleep(time.Second)
		Publish(ctx, "wait1")
		Publish(ctx, "wait2")
		Publish(ctx, "wait3")
	}()
	Wait(ctx, "wait1", "wait2", "wait3")
}

// TestReadEnvFiles tests readEnvFiles().
func TestReadEnvFiles(t *testing.T) {
	assert.Equal(t, "hello", os.Getenv("ENV1"))
	assert.Equal(t, "world", os.Getenv("ENV2"))
}
