# Golang Concurrency Manager [![Build Status](https://travis-ci.com/zenthangplus/goccm.svg?branch=master)](https://travis-ci.com/zenthangplus/goccm)

Golang Concurrency Manager package limits the number of goroutines that are allowed to run concurrently.

### Installation

Run the following command to install this package:

```
$ go get -u github.com/zenthangplus/goccm
```

### Example

```go
package main

import (
    "fmt"
    "github.com/zenthangplus/goccm"
    "time"
)

func main()  {
    // Limit 3 goroutines to run concurrently.
    cm := goccm.New(3)
    for i := 0; i < 10; i++ {
        // This function have to call before any goroutine
        cm.Wait()
        go func(i int) {
            fmt.Printf("Job %d is running\n", i)
            time.Sleep(2 * time.Second)
            
            // This function have to when a goroutine has finished
            // Or you can use `defer c.Done()` at the top of goroutine.
            cm.Done()
        }(i)
    }
    
    // This function have to call to ensure all goroutines have finished 
    // after close the main program.
    cm.WaitAllDone()
}
```

### List of supported functions

```go
package main

import "github.com/zenthangplus/goccm"

func main() {
    // Create the concurrency manager
    // The first argument is the maximum number of goroutines to run concurrently.
    cm := goccm.New(10)
    
    // Wait until a slot is available for the new goroutine.
    cm.Wait()
    
    // Mark a goroutine as finished
    cm.Done()
    
    // Wait for all goroutines are done
    cm.WaitAllDone()
    
    // Returns the number of goroutines which are running
    cm.RunningCount()
}
```
