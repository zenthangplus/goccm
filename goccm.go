package goccm

type (
	// ConcurrencyManager Interface
	ConcurrencyManager interface {
		// Wait until a slot is available for the new goroutine.
		Wait()

		// Mark a goroutine as finished
		Done()

		// Close the manager manually
		Close()

		// Wait for all goroutines are done
		WaitAllDone()

		// Returns the number of goroutines which are running
		RunningCount() int
	}

	concurrencyManager struct {
		// The number of goroutines that are allowed to run concurrently
		max int

		// The manager channel to coordinate the number of concurrent goroutines.
		managerCh chan struct{}

		// This channel indicates when all goroutines have finished their job.
		allDoneCh chan struct{}

		// The close flag allows we know when we can close the manager
		closed bool
	}
)

// New concurrencyManager
func New(maxGoRoutines int) *concurrencyManager {
	// Create manager channel with maxGoRoutines size
	managerCh := make(chan struct{}, maxGoRoutines)

	// Fill the manager channel 
	for i := 0; i < maxGoRoutines; i++ {
		managerCh <- struct{}{}
	}

	// Initiate the manager
	return &concurrencyManager{
		max: maxGoRoutines,
		managerCh: managerCh,
		allDoneCh: make(chan struct{}, 1),
	}
}

// Wait until a slot is available for the new goroutine.
// A goroutine have to start after this function.
func (c *concurrencyManager) Wait() {
	// Try to receive from the manager channel. When we have something,
	// it means a slot is available and we can start a new goroutine.
	// Otherwise, it will block until a slot is available.
	<-c.managerCh
}

// Mark a goroutine as finished
func (c *concurrencyManager) Done() {
	// Say that another goroutine can now start
	c.managerCh <- struct{}{}

	// When the closed flag is set,
	// we need to close the manager if it doesn't have any running goroutine
	if c.closed && len(c.managerCh) == c.max {
		// Say that all goroutines are finished, we can close the manager
		c.allDoneCh <- struct{}{}
	}
}

// Close the manager manually
func (c *concurrencyManager) Close() {
	c.closed = true
}

// Wait for all goroutines are done
func (c *concurrencyManager) WaitAllDone() {
	// Close the manager automatic
	c.Close()

	// This will block until allDoneCh was marked
	<-c.allDoneCh
}

// Returns the number of goroutines which are running
func (c *concurrencyManager) RunningCount() int {
	return  c.max - len(c.managerCh)
}
