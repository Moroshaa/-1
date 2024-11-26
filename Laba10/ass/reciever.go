// receiver.go
package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Функция для подписания сообщения
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
	signature, err := rsa.SignPKCS1v15(nil, privKey, crypto.SHA256, hashed)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// Функция для проверки подписи
func verifySignature(publicKeyFile string, message string, signature []byte) error {
	// Чтение открытого ключа
	pubFile, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(pubFile)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an RSA public key")
	}

	// Хэширование сообщения с помощью SHA-256
	hash := sha256.New()
	hash.Write([]byte(message))
	hashed := hash.Sum(nil)

	// Проверка подписи
	err = rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA256, hashed, signature)
	if err != nil {
		return fmt.Errorf("signature verification failed: %v", err)
	}

	return nil
}

// Обработчик запроса для проверки подписи и отправки ответа
func verifyHandler(w http.ResponseWriter, r *http.Request) {
	// Логирование запроса
	fmt.Println("Received request:", r.Method, r.URL)

	var data struct {
		Message   string `json:"message"`
		Signature string `json:"signature"`
	}

	// Чтение тела запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	// Разбор JSON
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Печатаем информацию о сообщении и подписи
	fmt.Printf("Received message: %s\n", data.Message)
	fmt.Printf("Received signature: %s\n", data.Signature)

	// Преобразуем подпись из hex в байты
	signature, err := hex.DecodeString(data.Signature)
	if err != nil {
		http.Error(w, "Invalid signature format", http.StatusBadRequest)
		return
	}

	// Проверка подписи
	err = verifySignature("public_key.pem", data.Message, signature)
	if err != nil {
		http.Error(w, fmt.Sprintf("Verification failed: %v", err), http.StatusUnauthorized)
		return
	}

	// Подготовка нового сообщения для отправки обратно
	newMessage := "Да это хайп братишка! Я получил твое сообщение!"

	// Подписываем новое сообщение
	newSignature, err := signMessage("private_key.pem", newMessage)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error signing new message: %v", err), http.StatusInternalServerError)
		return
	}

	// Преобразуем подпись нового сообщения в hex строку
	newSignatureHex := hex.EncodeToString(newSignature)

	// Создаем JSON-ответ с новым сообщением и подписью
	response := struct {
		Message   string `json:"message"`
		Signature string `json:"signature"`
	}{
		Message:   newMessage,
		Signature: newSignatureHex,
	}

	// Отправляем ответ обратно отправителю
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Логирование старта сервера
	fmt.Println("Server is listening on port 9090...")

	// Убедитесь, что обработчик маршрута /verify настроен
	http.HandleFunc("/verify", verifyHandler)

	// Запуск HTTP-сервера
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
