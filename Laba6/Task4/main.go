package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	counter int
	mu      sync.Mutex
)

func increment(wg *sync.WaitGroup) {
	mu.Lock()
	defer mu.Unlock()
	counter++
	wg.Done()
}
func main() {
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			increment(&wg)

		}()
		fmt.Println("Счетчик: ", counter)
		time.Sleep(time.Millisecond * 500)
	}
	wg.Wait()

}
