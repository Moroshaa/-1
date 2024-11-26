package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

// Запуск TLS-сервера
func main() {
	// Загрузка сертификата сервера и его ключа
	serverCert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Ошибка загрузки сертификата сервера: %v", err)
	}

	// Загрузка корневого сертификата CA для проверки клиентов
	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("Ошибка загрузки корневого сертификата: %v", err)
	}

	// Создание пула сертификатов для проверки клиентских сертификатов
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Настройка TLS-сертификатов и проверки клиента
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert, // Требование клиентского сертификата
	}

	// Настройка TCP-листенера для TLS-соединений
	listener, err := tls.Listen("tcp", "localhost:5454", tlsConfig)
	if err != nil {
		log.Fatalf("Ошибка создания TLS-листенера: %v", err)
	}
	defer listener.Close()

	fmt.Println("TLS сервер запущен на порту 5454...")

	for {
		// Ожидание подключения клиента
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Ошибка при принятии соединения:", err)
			continue
		}
		go handleConnection(conn)
	}
}

// Обработка клиентского соединения
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Чтение данных от клиента
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Ошибка чтения от клиента:", err)
		return
	}

	fmt.Printf("Получено сообщение от клиента: %s\n", string(buf[:n]))

	// Ответ клиенту
	_, err = conn.Write([]byte("И тебе хайпик!"))
	if err != nil {
		log.Println("Ошибка отправки данных клиенту:", err)
		return
	}
}
