package goccm

import (
	"sync/atomic"
)

type (
	// ConcurrencyManager Interface
	ConcurrencyManager interface {
		// Wait until a slot is available for the new goroutine.
		Wait()

		// Done Mark a goroutine as finished
		Done()

		// Close the manager manually
		Close()

		// WaitAllDone Wait for all goroutines are done
		WaitAllDone()

		// RunningCount Returns the number of goroutines which are running
		RunningCount() int32

		// Release Closes all the channels associated with goccm
		Release()
	}

	concurrencyManager struct {
		// The number of goroutines that are allowed to run concurrently
		max int

		// The manager channel to coordinate the number of concurrent goroutines.
		managerCh chan interface{}

		// The done channel indicates when a single goroutine has finished its job.
		doneCh chan bool

		// This channel indicates when all goroutines have finished their job.
		allDoneCh chan bool

		// The close flag allows we know when we can close the manager
		closed atomic.Bool

		// The running count allows we know the number of goroutines are running
		runningCount int32
	}
)

// New concurrencyManager
func New(maxGoRoutines int) *concurrencyManager {
	// Initiate the manager object
	c := concurrencyManager{
		max:       maxGoRoutines,
		managerCh: make(chan interface{}, maxGoRoutines),
		doneCh:    make(chan bool),
		allDoneCh: make(chan bool),
	}

	// Fill the manager channel by placeholder values
	for i := 0; i < c.max; i++ {
		c.managerCh <- nil
	}

	// Start the controller to collect all the jobs
	go c.controller()

	return &c
}

// Create the controller to collect all the jobs.
// When a goroutine is finished, we can release a slot for another goroutine.
func (c *concurrencyManager) controller() {
	for {
		// This will block until a goroutine is finished
		<-c.doneCh

		// Say that another goroutine can now start
		c.managerCh <- nil

		// When the closed flag is set,
		// we need to close the manager if it doesn't have any running goroutine
		if c.IsClosed() && c.RunningCount() == 0 {
			break
		}
	}

	// Say that all goroutines are finished, we can close the manager
	c.allDoneCh <- true
}

func (c *concurrencyManager) IsClosed() bool {
	return c.closed.Load()
}

// Wait until a slot is available for the new goroutine.
// A goroutine have to start after this function.
func (c *concurrencyManager) Wait() {

	// Try to receive from the manager channel. When we have something,
	// it means a slot is available and we can start a new goroutine.
	// Otherwise, it will block until a slot is available.
	<-c.managerCh

	// Increase the running count to help we know how many goroutines are running.
	atomic.AddInt32(&c.runningCount, 1)
}

// Done Mark a goroutine as finished
func (c *concurrencyManager) Done() {
	// Decrease the number of running count
	atomic.AddInt32(&c.runningCount, -1)
	c.doneCh <- true
}

// Close the manager manually
func (c *concurrencyManager) Close() {
	c.closed.Store(true)
}

// WaitAllDone Wait for all goroutines are done
func (c *concurrencyManager) WaitAllDone() {
	// Close the manager automatic
	c.Close()

	// This will block until allDoneCh was marked
	<-c.allDoneCh
}

// RunningCount Returns the number of goroutines which are running
func (c *concurrencyManager) RunningCount() int32 {
	return atomic.LoadInt32(&c.runningCount)
}

// Release Closes all the channels associated with goccm
func (c *concurrencyManager) Release() {
	close(c.managerCh)
	close(c.doneCh)
	close(c.allDoneCh)
}
