package goccm

// Manager coordinates the number of concurrent goroutines.
type Manager chan struct{}

// New Manager that handles `maxGoRoutines` number of concurrency.
func New(maxGoRoutines int) Manager {
	// Create manager channel with maxGoRoutines size
	c := make(chan struct{}, maxGoRoutines)
	// Fill the manager channel
	for i := 0; i < maxGoRoutines; i++ {
		c <- struct{}{}
	}
	return c
}

// Wait blocks until a slot is allocated from the manager.
func (m Manager) Wait() {
	<-m
}

// Done returns the routine to the manager.
func (m Manager) Done() {
	m <- struct{}{}
}

// Wait for all goroutines to finish.
func (m Manager) WaitAllDone() {
	for len(m) != cap(m) {
	}
}

// RunningCount returns the number of running goroutines.
func (m Manager) RunningCount() int {
	return cap(m) - len(m)
}
