package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

// Функция для генерации ключа на основе введенной строки (с использованием SHA-256)
func generateKey(secret string) ([]byte, error) {
	// Хешируем секретный ключ с помощью SHA-256
	hash := sha256.New()
	hash.Write([]byte(secret))
	return hash.Sum(nil), nil
}

// Функция для шифрования данных с использованием AES-GCM
func encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Генерация случайного nonce для GCM
	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	// Создание нового GCM шифратора
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Шифруем данные
	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

	// Возвращаем зашифрованные данные в формате base64
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

// Функция для расшифровки данных с использованием AES-GCM
func decrypt(ciphertext string, key []byte) (string, error) {
	// Декодируем данные из base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Разделяем nonce и ciphertext
	nonce, ciphertextData := data[:12], data[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Создание нового GCM шифратора
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Расшифровываем данные
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextData, nil)
	if err != nil {
		return "", err
	}

	// Возвращаем расшифрованную строку
	return string(plaintext), nil
}

func main() {
	// Чтение строки для шифрования
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите строку для шифрования: ")
	plaintext, _ := reader.ReadString('\n')
	plaintext = strings.TrimSpace(plaintext)

	// Чтение секретного ключа
	fmt.Print("Введите секретный ключ: ")
	secretKey, _ := reader.ReadString('\n')
	secretKey = strings.TrimSpace(secretKey)

	// Генерация ключа на основе секретной строки
	key, err := generateKey(secretKey)
	if err != nil {
		fmt.Println("Ошибка при генерации ключа:", err)
		return
	}

	// Шифрование строки
	encryptedText, err := encrypt(plaintext, key)
	if err != nil {
		fmt.Println("Ошибка при шифровании:", err)
		return
	}

	fmt.Printf("Зашифрованный текст: %s\n", encryptedText)

	// Чтение строки для расшифровки
	fmt.Print("Введите зашифрованный текст для расшифровки: ")
	encryptedTextInput, _ := reader.ReadString('\n')
	encryptedTextInput = strings.TrimSpace(encryptedTextInput)

	// Расшифровка строки
	decryptedText, err := decrypt(encryptedTextInput, key)
	if err != nil {
		fmt.Println("Ошибка при расшифровке:", err)
		return
	}

	fmt.Printf("Расшифрованный текст: %s\n", decryptedText)
}
