package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка подключения к серверу:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите сообщение для отправки серверу: ")
	message, _ := reader.ReadString('\n')

	// Отправка сообщения серверу
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Ошибка отправки сообщения:", err)
		return
	}

	// Получение ответа от сервера
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения ответа от сервера:", err)
		return
	}
	fmt.Println("Ответ от сервера:", response)
}
