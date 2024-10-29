package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

// handleConnection обрабатывает клиентские соединения
func handleConnection(conn net.Conn) {
	defer wg.Done() // Сигнал завершения горутины
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Клиент отключился:", conn.RemoteAddr())
			return
		}
		fmt.Print("Получено сообщение от клиента:", message)
		// Отправка подтверждения клиенту
		_, err = conn.Write([]byte("Сообщение получено\n"))
		if err != nil {
			fmt.Println("Ошибка отправки данных клиенту:", err)
			return
		}
	}
}

// startServer запускает сервер, который слушает на указанном порту и принимает входящие соединения
func startServer(port string, shutdownCh <-chan struct{}) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
	defer listener.Close()

	fmt.Println("Сервер запущен на порту", port)

	// Цикл для прослушивания входящих подключений
	for {
		select {
		case <-shutdownCh:
			fmt.Println("Сервер завершает работу...")
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Ошибка при подключении:", err)
				continue
			}
			wg.Add(1)
			go handleConnection(conn) // Обработка соединения в отдельной горутине
		}
	}
}

func main() {
	shutdownCh := make(chan struct{})
	port := "8080"

	// Запуск сервера в отдельной горутине
	go startServer(port, shutdownCh)

	// Отслеживание сигнала завершения работы сервера
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Закрытие сервера
	close(shutdownCh)
	wg.Wait()
	fmt.Println("Сервер остановлен")
}
