package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	for {

		var TaskNumber int
		fmt.Println("Выберите номер задания:\n1 - Работа с картами\n2 - Работа с регистром\n3 - Разворот массива\nДля выхода выберите 0")
		fmt.Scan(&TaskNumber)

		switch TaskNumber {

		case 1:
			people()
		case 2:
			upperCase()
		case 3:
			reverseArr()
		case 0:
			fmt.Println("Пока!")
			return

		default:
			fmt.Println("Не верный номер задание, выберите от 0 до 4")
			os.Exit(1)
		}
	}
}

func people() {
	humanMap := make(map[string]int)

	humanMap["Дмитрий"] = 45
	humanMap["Александр"] = 20
	humanMap["Егор"] = 12
	humanMap["Василий"] = 30
	humanMap["Роман"] = 20

	for {

		var choice int
		fmt.Print("Что вы хотите сделать?\n1 - Добавть\n2 - Удалить\n3 - Показать карту\n4 - Посчитать средний возраст\n0 - Выход\n")
		fmt.Scan(&choice)
		switch choice {
		case 1:
			var value int
			var key string
			fmt.Println("Введите имя: ")

			fmt.Scan(&key)

			fmt.Println("Введите возраст: ")

			fmt.Scan(&value)

			humanMap[key] = value
			fmt.Println("Добавлен эелемент: ", key, "=", value)

		case 2:
			var key string
			fmt.Println("Введите ключ для удаления:")
			fmt.Scan(&key)

			if _, exists := humanMap[key]; exists {
				delete(humanMap, key)
				fmt.Println("Элемент с ключом", key, "удалён.")
			} else {
				fmt.Println("Ключ не найден:", key)
			}

		case 3:
			fmt.Println("Текущая карта: ", humanMap)

		case 4:
			if len(humanMap) == 0 {
				fmt.Println("Карта пуста.")
			} else {
				sum := 0
				for _, value := range humanMap {
					sum += value
				}
				average := float64(sum) / float64(len(humanMap))
				fmt.Printf("Среднее значение: %.2f\n", average)
			}

		case 0:
			fmt.Println("Выход из задания!")
			return

		default:
			fmt.Println("Выберите от 1 до 3! ")
		}
	}
}
func upperCase() {
	var str string
	fmt.Println("Введите строку в нижнем регисте:")
	fmt.Scan(&str)

	upper := strings.ToUpper(str)
	fmt.Println("Строка в верхнем регистре: ", upper)
}
func reverseArr() {
	arr := []int{12, 23, 34, 5, 157}
	fmt.Println("Исходный масств: ", arr)
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[len(arr)-1-i] = arr[len(arr)-1-i], arr[i]

	}
	fmt.Println("Перевернутый массив:", arr)
}
