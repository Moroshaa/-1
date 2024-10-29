package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware для логирования запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("%s %s %s\n", r.Method, r.URL.Path, duration)
	})
}

// Обработчик для GET /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Привет, это ваш HTTP сервер!")
}

// Обработчик для POST /data
func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}

	// Выводим содержимое JSON в консоль
	fmt.Printf("Получены данные: %v\n", data)
	fmt.Fprintf(w, "Данные получены и обработаны")
}

// Функция для инициализации маршрутов и запуска сервера
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/data", dataHandler)

	// Применение middleware для логирования
	loggedMux := loggingMiddleware(mux)

	// Запуск HTTP сервера на порту 8080
	fmt.Println("HTTP сервер запущен на порту 8080")
	err := http.ListenAndServe(":8080", loggedMux)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
