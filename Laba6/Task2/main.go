package main

import (
	"fmt"
	"sync"
)

func fibonacci(n int, c chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	a, b := 0, 1
	for i := 0; i < n; i++ {
		c <- a
		a, b = b, a+b
	}
	close(c)
}
func printFib(c chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range c {
		fmt.Println(num)
	}

}
func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	c := make(chan int)
	go fibonacci(10, c, &wg)
	go printFib(c, &wg)

	wg.Wait()
}
