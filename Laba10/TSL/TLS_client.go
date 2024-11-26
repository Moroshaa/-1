package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
)

// TLS-клиент
func main() {
	// Загрузка клиентского сертификата и ключа
	clientCert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Fatalf("Ошибка загрузки сертификата клиента: %v", err)
	}

	// Загрузка корневого сертификата CA для проверки сервера
	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("Ошибка загрузки корневого сертификата: %v", err)
	}

	// Создание пула сертификатов для проверки серверного сертификата
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Настройка TLS
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert}, // Сертификат клиента
		RootCAs:            caCertPool,                    // Проверка сертификата сервера
		InsecureSkipVerify: false,                         // Обязательная проверка сертификата сервера
	}

	// Подключение к серверу
	conn, err := tls.Dial("tcp", "localhost:5454", tlsConfig)
	if err != nil {
		log.Fatalf("Ошибка соединения с сервером: %v", err)
	}
	defer conn.Close()

	// Отправка сообщения на сервер
	message := "Хайп братан!"
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	// Чтение ответа от сервера
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Ошибка чтения ответа от сервера: %v", err)
	}

	fmt.Printf("Ответ от сервера: %s\n", string(buf[:n]))
}
