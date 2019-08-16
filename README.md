# Golang concurrency manager

Golang Concurrency Manager package manages the number of goroutines that are allowed to run concurrently.

### Example

```go
package main

import (
    "fmt"
    "goccm"
    "time"
)

func main()  {
    c := goccm.New(3)// Limit 3 goroutines to run concurrently.
    for i := 1; i <= 10; i++ {
        c.Wait()
        go func(i int) {
            fmt.Printf("Job %d is running\n", i)
            time.Sleep(2 * time.Second)
            c.Done()
        }(i)
    }
    c.WaitAllDone()
}
```
