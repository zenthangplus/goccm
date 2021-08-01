package goccm

// ConcurrencyManager Interface
type ConcurrencyManager interface {
	// Wait until a slot is available for the new goroutine.
	Wait()

	// Mark a goroutine as finished
	Done()

	// Wait for all goroutines are done
	WaitAllDone()

	// Returns the number of goroutines which are running
	RunningCount() int32
}

// Manager coordinates the maximum number of concurrent goroutines.
type concurrencyManager chan struct{}

// New Manager that handles `maxGoRoutines` number of concurrency.
func New(maxGoRoutines int) *concurrencyManager {
	// Create manager channel with maxGoRoutines size
	cm := concurrencyManager(make(chan struct{}, maxGoRoutines))
	return &cm
}

// Wait blocks until a slot is allocated from the manager.
func (cm *concurrencyManager) Wait() {
	*cm <- struct{}{}
}

// Done returns the routine to the manager.
func (cm *concurrencyManager) Done() {
	<-(*cm)
}

// Wait for all goroutines to finish.
func (cm *concurrencyManager) WaitAllDone() {
	for len(*cm) > 0 {
	}
}

// RunningCount returns the number of running goroutines.
func (cm *concurrencyManager) RunningCount() int {
	return len(*cm)
}
