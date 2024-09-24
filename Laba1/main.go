package main

import (
	"fmt"
	"time"
)

func sumAndDifference(a, b float64) (float64, float64) {
	return a + b, a - b
}

func main() {
	// 1
	currentTime := time.Now()
	fmt.Println("Задание 1, вывод дата+время")
	fmt.Println("Текущая дата и время:", currentTime)

	// 2
	var intVar int = 10
	var floatVar float64 = 3.14
	var stringVar string = "Hello, World!"
	var boolVar bool = true

	fmt.Println("Задание 2, вывод int, float64(double), string, bool")

	fmt.Println("intVar:", intVar)
	fmt.Println("floatVar:", floatVar)
	fmt.Println("stringVar:", stringVar)
	fmt.Println("boolVar:", boolVar)
	// 3
	shortIntVar := 20
	shortFloatVar := 6.28
	shortStringVar := "Привет дружбан"
	shortBoolVar := false

	fmt.Println("Задание 3, краткая форма объявления переменных ")

	fmt.Println("shortIntVar:", shortIntVar)
	fmt.Println("shortFloatVar:", shortFloatVar)
	fmt.Println("shortStringVar:", shortStringVar)
	fmt.Println("shortBoolVar:", shortBoolVar)

	//  4
	a := 10
	b := 5

	fmt.Println("Задание 4, арифметические операции с двумя целыми числами ")

	fmt.Println("Сумма:", a+b)
	fmt.Println("Разность:", a-b)
	fmt.Println("Произведение:", a*b)
	fmt.Println("Частное:", a/b)

	// 5
	x := 170.14
	y := 2.71

	sum, diff := sumAndDifference(x, y)

	fmt.Println("Задание 5, вычисление суммы и разности двух чисел с плавающей запятой ")

	fmt.Println("Сумма:", sum)
	fmt.Println("Разность:", diff)

	// 6
	num1 := 5.0
	num2 := 7.0
	num3 := 9.0

	average := (num1 + num2 + num3) / 3

	fmt.Println("Задание 6, вычисление среднего значения трех чисел ")

	fmt.Println("Среднее значение:", average)
}
