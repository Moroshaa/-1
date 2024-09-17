

# Задания для лабораторной работы 1. 

1. Написать программу, которая выводит текущее время и дату.
2. Создать переменные различных типов (int, float64, string, bool) и вывести их на экран. 
3. Использовать краткую форму объявления переменных для создания и вывода переменных. 4. Написать программу для выполнения арифметических операций с двумя целыми числами и выводом результатов.
5. Реализовать функцию для вычисления суммы и разности двух чисел с плавающей запятой. 
6. Написать программу, которая вычисляет среднее значение трех чисел.

## Установка и запуск

Для запуска программы, скачайте zip-файл , откройте main.go в удобной для вас среде разработки и пропишиште в терминале go run main.go , предварительно указав в терминали путь до папки. 

Если файл по какой-то причине не запускается, проверьте установлен у вас язык go на пк. Если это не помогло, укажите path до go.exe в системных переменных.




## Объявление пакета, к которому пренадлжит файл:

	package main

## Импорты пакетов fmt(ввод-вывод) и time(для работы с датами и временем) 

	import (
		"fmt"
		"time"
	)

## 1. Данная часть кода обращается к функции Now() из пакета time и выводит информацию в консоль приложения 

    	currentTime := time.Now()
	fmt.Println("Задание 1, вывод дата+время")
	fmt.Println("Текущая дата и время:", currentTime)


## 2. В этом блоке кода приведены некоторые типы данных(целые числа, числа с плавающей точкой, строка и булево), которые выводятся в консоль 
    	var intVar int = 10
	var floatVar float64 = 3.14
	var stringVar string = "Hello, World!"
	var boolVar bool = true

	fmt.Println("Задание 2, вывод int, float64(double), string, bool")

	fmt.Println("intVar:", intVar)
	fmt.Println("floatVar:", floatVar)
	fmt.Println("stringVar:", stringVar)
	fmt.Println("boolVar:", boolVar)

## 3. Здесь мы можем наблюдать краткую запись переменных через оператор :=, а так же их вывод в консоль.


    shortIntVar := 20
	shortFloatVar := 6.28
	shortStringVar := "Привет дружбан"
	shortBoolVar := false

	fmt.Println("Задание 3, краткая форма объявления переменных ")

	fmt.Println("shortIntVar:", shortIntVar)
	fmt.Println("shortFloatVar:", shortFloatVar)
	fmt.Println("shortStringVar:", shortStringVar)
	fmt.Println("shortBoolVar:", shortBoolVar)

## 4. Операции с двумя заданными целыми числами, и их вывод

   	a := 10
	b := 5

	fmt.Println("Задание 4, арифметические операции с двумя целыми числами ")

	fmt.Println("Сумма:", a+b)
	fmt.Println("Разность:", a-b)
	fmt.Println("Произведение:", a*b)
	fmt.Println("Частное:", a/b)

## 5. Сумма и разность двух чисел с плавающей точкой
    x := 170.14
	y := 2.71

	sum, diff := sumAndDifference(x, y)

	fmt.Println("Задание 5, вычисление суммы и разности двух чисел с плавающей запятой ")

	fmt.Println("Сумма:", sum)
	fmt.Println("Разность:", diff)

## 6. Вычисление среднего значения трех чисел с плавающей точкой

    num1 := 5.0
	num2 := 7.0
	num3 := 9.0

	average := (num1 + num2 + num3) / 3

	fmt.Println("Задание 6, вычисление среднего значения трех чисел ")

	fmt.Println("Среднее значение:", average)


## Общий листинг программы:

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

## Вывод:

PS C:\Лабы\1 семестр\Прикладное программирование\Лаба 1> go run main.go <br>
Задание 1, вывод дата+время<br>
Текущая дата и время: 2024-09-17 23:43:24.3879788 +0400 +04 m=+0.000000001 <br>
Задание 2, вывод int, float64(double), string, bool<br>
intVar: 10<br>
floatVar: 3.14<br>
stringVar: Hello, World!<br>
boolVar: true<br>
Задание 3, краткая форма объявления переменных <br>
shortIntVar: 20<br>
shortFloatVar: 6.28<br>
shortStringVar: Привет дружбан<br>
shortBoolVar: false<br>
Задание 4, арифметические операции с двумя целыми числами <br>
Сумма: 15<br>
Разность: 5<br>
Произведение: 50<br>
Частное: 2<br>
Задание 5, вычисление суммы и разности двух чисел с плавающей запятой<br>
Сумма: 172.85<br>
Разность: 167.42999999999998<br>
Задание 6, вычисление среднего значения трех чисел<br>
Среднее значение: 7


