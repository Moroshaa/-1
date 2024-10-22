package main

import (
	"fmt"
	"math/rand"
	"time"
)

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

func generateRandomNumbers(n int) {
	for i := 0; i < n; i++ {
		fmt.Println("Рандомное число:", rand.Intn(100))
		time.Sleep(time.Millisecond * 500)
	}
}

func sumOfSeries(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
		time.Sleep(time.Millisecond * 200)
	}
	return sum
}

func main() {
	// Start a goroutine for each function
	go func() {
		fmt.Println("Факториал 5:", factorial(5))
	}()

	go func() {
		generateRandomNumbers(5)
	}()

	go func() {
		fmt.Println("Sum of series up to 5:", sumOfSeries(5))
	}()

	time.Sleep(time.Second * 3)
}
