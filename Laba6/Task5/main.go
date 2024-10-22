package main

import (
	"fmt"
)

type Request struct {
	op   string
	a, b float64
	res  chan float64
}

func addition(a float64, b float64, res chan float64) {
	res <- a + b
	close(res)
}

func subtraction(a float64, b float64, res chan float64) {
	res <- a - b
	close(res)
}

func multiplication(a float64, b float64, res chan float64) {
	res <- a * b
	close(res)
}

func division(a float64, b float64, res chan float64) {
	res <- a / b
	close(res)
}

func main() {
	var ans string
	var a, b float64

	for {
		fmt.Println("Выберите действие: + - / * (или 0 для выхода)")
		fmt.Scanln(&ans)

		if ans == "0" {
			break
		}

		fmt.Println("Введите a:")
		fmt.Scanln(&a)
		fmt.Println("Введите b:")
		fmt.Scanln(&b)
		res := make(chan float64)

		switch ans {
		case "+":
			go addition(a, b, res)
		case "-":
			go subtraction(a, b, res)
		case "/":
			go division(a, b, res)
		case "*":
			go multiplication(a, b, res)
		default:
			fmt.Println("Выберите допустимые действия!")
			continue
		}

		fmt.Println("Результат:", <-res)
	}

	fmt.Println("Выход из калькулятора.")
}
