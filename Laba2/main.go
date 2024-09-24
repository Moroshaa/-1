package main

import (
	"fmt"
	"unicode/utf8"
)

// Функция для 2 задания
func checkNumber(num1 int) string {
	if num1 > 0 {
		return "positive"
	} else if num1 < 0 {
		return "negative"
	} else {
		return "zero"
	}
}

// Функция для 4 задания
func stringLength(str string) int {
	return utf8.RuneCountInString(str)
}

// Функция для 6 задания
func averageValue(a, b float64) float64 {
	return (a + b) / 2
}

func main() {
	// 1 проверка на четность
	fmt.Println("Задани 1")
	num := 9
	if num%2 == 0 {
		fmt.Println("четное")
	} else {
		fmt.Println("нечетное")
	}
	fmt.Println()

	// 2 задание
	fmt.Println("Задани 2")
	num1 := -20
	result := checkNumber(num1)
	fmt.Println(num1, "is", result)
	fmt.Println()

	// 3 вывод чисел от 1 до 10 через for
	fmt.Println("Задани 3")
	for i := 0; i < 10; i++ {
		fmt.Println(i + 1)
	}
	fmt.Println()

	// 4 длина строки
	fmt.Println("Задани 4")
	str := "Привет мир"
	length := stringLength(str)
	fmt.Println("Длина строки", str, " - ", length, " символов")
	fmt.Println()

	// 6 задание
	fmt.Println("Задани 6")
	n1 := 10.0
	n2 := 25.0
	avg := averageValue(n1, n2)
	fmt.Println("Среднее значение 2х чисел =", avg)

}
