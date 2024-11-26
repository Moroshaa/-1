// sender.go
package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex" // Импорт для работы с hex
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
)

func signMessage(privateKeyFile string, message string) ([]byte, error) {
	// Чтение приватного ключа
	privFile, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privFile)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Хэширование сообщения с помощью SHA-256
	hash := sha256.New()
	hash.Write([]byte(message))
	hashed := hash.Sum(nil)

	// Подписание сообщения
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func sendMessage(message string) {
	// Подписываем сообщение
	signature, err := signMessage("private_key.pem", message)
	if err != nil {
		fmt.Println("Error signing message:", err)
		return
	}

	// Преобразуем подпись в hex строку
	signatureHex := hex.EncodeToString(signature)

	// Создаем JSON-запрос с подписанным сообщением и подписью
	payload := fmt.Sprintf(`{"message": "%s", "signature": "%s"}`, message, signatureHex)

	resp, err := http.Post("http://localhost:9090/verify", "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response from server:", string(body))
}

func main() {
	message := "Это хайп братишка!"
	sendMessage(message)
}
