package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func generateKeys() error {
	// Генерация приватного ключа
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Сохранение приватного ключа в файл
	privateFile, err := os.Create("private_key.pem")
	if err != nil {
		return err
	}
	defer privateFile.Close()

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	privateFile.Write(privateKeyPEM)

	// Сохранение открытого ключа в файл
	publicKey := &privateKey.PublicKey
	publicKeyPEM, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	publicFile, err := os.Create("public_key.pem")
	if err != nil {
		return err
	}
	defer publicFile.Close()

	publicFile.Write(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyPEM,
	}))

	fmt.Println("Keys generated and saved!")
	return nil
}

func main() {
	if err := generateKeys(); err != nil {
		fmt.Println("Error generating keys:", err)
	}
}
