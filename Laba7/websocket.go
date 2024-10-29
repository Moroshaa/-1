package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Объект для обновления WebSocket соединений
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем все источники для упрощения
	},
}

// Структура, представляющая сервер чата
type ChatServer struct {
	clients   map[*websocket.Conn]bool // Хранит все активные подключения
	broadcast chan []byte              // Канал для рассылки сообщений
	mutex     sync.Mutex               // Мьютекс для синхронизации доступа к клиентам
}

// Создание нового ChatServer
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

// Метод для запуска сервера, прослушивания и рассылки сообщений
func (server *ChatServer) Start() {
	for {
		msg := <-server.broadcast
		server.mutex.Lock()
		for client := range server.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println("Ошибка отправки сообщения клиенту:", err)
				client.Close()
				delete(server.clients, client)
			}
		}
		server.mutex.Unlock()
	}
}

// Метод для добавления нового клиента
func (server *ChatServer) AddClient(conn *websocket.Conn) {
	server.mutex.Lock()
	server.clients[conn] = true
	server.mutex.Unlock()
}

// Обработчик для WebSocket соединений
func (server *ChatServer) handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при установке WebSocket соединения:", err)
		return
	}
	defer conn.Close()

	server.AddClient(conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Клиент отключился:", err)
			server.mutex.Lock()
			delete(server.clients, conn)
			server.mutex.Unlock()
			break
		}
		server.broadcast <- msg // Отправляем сообщение в канал
	}
}

func main() {
	server := NewChatServer()
	go server.Start()

	http.HandleFunc("/ws", server.handleConnections)

	fmt.Println("Сервер WebSocket запущен на :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка при запуске HTTP-сервера:", err)
	}
}
