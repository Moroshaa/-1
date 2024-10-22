package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func randNum(n int, c chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		num := rand.Intn(101)
		time.Sleep(time.Millisecond * 500)
		c <- num
	}
	close(c)
}

func evenNum(c chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case num, ok := <-c:
			if !ok {
				return
			}
			if num%2 == 0 {
				fmt.Printf("Число %d четное\n", num)

			} else {
				fmt.Printf("Число %d нечетное\n", num)
			}
		}
	}

}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	c := make(chan int)
	go randNum(10, c, &wg)
	go evenNum(c, &wg)
	wg.Wait()
}
