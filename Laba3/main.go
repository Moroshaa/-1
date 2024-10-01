package main

import (
	"Laba3/mathutils"
	"Laba3/stringutils"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	// факториал числа
	fmt.Print("Введите число: ")
	input, _ := reader.ReadString('\n')

	input = strings.TrimSpace(input)

	number, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Ошибка: введите целое число.")
		return
	}

	result := mathutils.Factorial(number)
	if result == 0 {
		fmt.Println("Ошибка: факториал не определён для отрицательных чисел.")
	} else {
		fmt.Printf("Факториал числа %d: %d\n", number, result)
	}
	// переворот строки
	fmt.Print("Введите строку: ")
	inputStr, _ := reader.ReadString('\n')
	inputStr = strings.TrimSpace(inputStr)
	reversed := stringutils.StrRevers(inputStr)

	fmt.Println("Перевернутая строка: ", reversed)

	// массив из 5 чисел
	fmt.Println("Запуск генератора массива")
	ArrGenerator()
	fmt.Println("Создание строк")
	StrLen()
}

func ArrGenerator() {
	rand.Seed(time.Now().UnixMicro())
	var numbers [5]int
	for i := 0; i < 5; i++ {
		numbers[i] = rand.Intn(100)
	}
	fmt.Println("Сгенерированный массив: ", numbers)
	slice := numbers[:]

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\nВыберите действие:\n1 - Добавить элемент\n2 - Удалить элемент")
	scanner.Scan()
	var action int
	fmt.Sscan(scanner.Text(), &action)

	if action == 1 {
		var value, index int
		fmt.Println("Введите значение для добавления:")
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &value)
		fmt.Println("Введите индекс для добавления (0 -", len(slice), "):")
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &index)

		if index < 0 || index > len(slice) {
			fmt.Println("Неверный индекс")
		} else {
			slice = append(slice[:index], append([]int{value}, slice[index:]...)...)
		}
	} else if action == 2 {
		var index int
		fmt.Println("Введите индекс для удаления (0 -", len(slice)-1, "):")
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &index)

		if index < 0 || index >= len(slice) {
			fmt.Println("Неверный индекс")
		} else {
			slice = append(slice[:index], slice[index+1:]...)
		}
	} else {
		fmt.Println("Неверное действие")
	}

	fmt.Println("\nЗначения среза после выполнения действия:")
	for i, value := range slice {
		fmt.Printf("slice[%d] = %d\n", i, value)
	}
}

func StrLen() {
	var strings []string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Введите строки (пустая строка для завершения):")
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		strings = append(strings, line)
	}

	fmt.Println("\nЗначения среза:")
	for i, value := range strings {
		fmt.Printf("strings[%d] = %s\n", i, value)
	}

	var longestString string
	for _, value := range strings {
		if len(value) > len(longestString) {
			longestString = value
		}
	}
	fmt.Printf("\nСамая длинная строка: %s\n", longestString)
}
