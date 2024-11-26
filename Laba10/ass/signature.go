// sign_and_verify.go
package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func signMessage(privateKeyFile string, message string) ([]byte, error) {
	// Чтение приватного ключа
	privFile, err := os.ReadFile(privateKeyFile)
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

func verifySignature(publicKeyFile string, message string, signature []byte) error {
	// Чтение открытого ключа
	pubFile, err := os.ReadFile(publicKeyFile)
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

func main() {
	message := "Hello, this is a signed message!"

	// Подписываем сообщение
	signature, err := signMessage("private_key.pem", message)
	if err != nil {
		fmt.Println("Error signing message:", err)
		return
	}
	fmt.Printf("Signed message: %x\n", signature)

	// Проверяем подпись
	err = verifySignature("public_key.pem", message, signature)
	if err != nil {
		fmt.Println("Verification failed:", err)
	} else {
		fmt.Println("Signature verified successfully!")
	}
}
