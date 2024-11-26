package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func SHA256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashByte := hash.Sum(nil)
	return hex.EncodeToString(hashByte)
}

func MD5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	hashByte := hash.Sum(nil)
	return hex.EncodeToString(hashByte)
}

func SHA512(data string) string {
	hash := sha512.New()
	hash.Write([]byte(data))
	hashByte := hash.Sum(nil)
	return hex.EncodeToString(hashByte)
}

func verifyHash(input string, providedHash string, hashFunc func(string) string) bool {
	// Генерируем хеш для введенной строки
	computedHash := hashFunc(input)

	// Сравниваем с предоставленным хешем
	return computedHash == providedHash
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Меню выбора
		fmt.Print("Выберите тип шифрования:\n1 - SHA-256\n2 - MD5\n3 - SHA-512\n4 - Проверка целостности\n0 - для выхода\n")
		var choice int
		fmt.Scanln(&choice)

		// Выход из программы
		if choice == 0 {
			fmt.Println("Выход из программы...")
			return // Завершаем программу
		}

		// Обработка проверки целостности данных
		if choice == 4 {
			fmt.Println("Проверка целостности данных:")

			// Ввод строки и её хеша
			fmt.Print("Введите строку: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			fmt.Print("Введите хеш строки: ")
			providedHash, _ := reader.ReadString('\n')
			providedHash = strings.TrimSpace(providedHash)

			// Запрос выбора алгоритма хеширования
			fmt.Print("Выберите тип шифрования для проверки:\n1 - SHA-256\n2 - MD5\n3 - SHA-512\n")
			var hashChoice int
			fmt.Scanln(&hashChoice)

			var result bool
			switch hashChoice {
			case 1:
				result = verifyHash(input, providedHash, SHA256)
			case 2:
				result = verifyHash(input, providedHash, MD5)
			case 3:
				result = verifyHash(input, providedHash, SHA512)
			default:
				fmt.Println("Не верный выбор алгоритма!")
				continue
			}

			// Вывод результата проверки
			if result {
				fmt.Println("Хеши совпадают. Целостность данных подтверждена.")
			} else {
				fmt.Println("Хеши не совпадают. Целостность данных нарушена.")
			}
		} else {
			// Ввод строки для шифрования
			fmt.Print("Введите строку: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			var result string
			switch choice {
			case 1:
				result = SHA256(input)
				fmt.Printf("Результат SHA-256 шифрования строки '%s': %s\n", input, result)
			case 2:
				result = MD5(input)
				fmt.Printf("Результат MD5 шифрования строки '%s': %s\n", input, result)
			case 3:
				result = SHA512(input)
				fmt.Printf("Результат SHA-512 шифрования строки '%s': %s\n", input, result)
			default:
				fmt.Println("Не верный выбор! Выберите от 0 до 3!")
			}
		}
	}
}
