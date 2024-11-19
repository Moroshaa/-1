package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var baseURL = "http://localhost:9090/users"
var authToken = "" // Авторизационный токен

// User структура для пользователя
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Функция для отправки POST запроса на сервер для авторизации
func login(username, password string) {
	url := "http://localhost:9090/login"
	client := &http.Client{}

	// Тело запроса
	data := map[string]string{
		"username": username,
		"password": password,
	}
	jsonData, _ := json.Marshal(data)

	// Отправка запроса
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка при авторизации, код ошибки: %d\n", resp.StatusCode)
		return
	}

	// Читаем ответ
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]string
	_ = json.Unmarshal(body, &result)

	// Сохраняем токен
	if token, exists := result["token"]; exists {
		authToken = token
		fmt.Println("Авторизация успешна, токен сохранен:", authToken)
	} else {
		fmt.Println("Не удалось получить токен.")
	}
}

// Функция для отображения всех пользователей
func getUsers() {
	if authToken == "" {
		fmt.Println("Вы не авторизованы. Выполните вход, чтобы продолжить.")
		return
	}

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		log.Fatal("Ошибка при создании запроса:", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ошибка при отправке запроса:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		var users []User
		if err := json.Unmarshal(body, &users); err == nil {
			for _, user := range users {
				fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
			}
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Ошибка при получении данных, код ошибки: %d, детали: %s\n", resp.StatusCode, string(body))
	}
}

// Функция для получения пользователя по ID
func getUserByID() {
	if authToken == "" {
		fmt.Println("Вы не авторизованы. Выполните вход, чтобы продолжить.")
		return
	}

	var id int
	fmt.Print("Введите ID пользователя: ")
	fmt.Scanln(&id)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", baseURL, id), nil)
	if err != nil {
		log.Fatal("Ошибка при создании запроса:", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ошибка при отправке запроса:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		var user User
		if err := json.Unmarshal(body, &user); err == nil {
			fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Ошибка при получении данных, код ошибки: %d, детали: %s\n", resp.StatusCode, string(body))
	}
}

// Функция для добавления нового пользователя
func createUser() {
	if authToken == "" {
		fmt.Println("Вы не авторизованы. Выполните вход, чтобы продолжить.")
		return
	}

	var name string
	var age int
	fmt.Print("Введите имя пользователя: ")
	fmt.Scanln(&name)
	fmt.Print("Введите возраст пользователя: ")
	fmt.Scanln(&age)

	newUser := User{Name: name, Age: age}
	body, err := json.Marshal(newUser)
	if err != nil {
		log.Fatal("Ошибка при сериализации данных:", err)
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("Ошибка при создании запроса:", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ошибка при отправке запроса:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Пользователь успешно добавлен.")
	} else {
		fmt.Println("Ошибка при добавлении пользователя, код ошибки:", resp.StatusCode)
	}
}

// Функция для обновления пользователя
func updateUser() {
	if authToken == "" {
		fmt.Println("Вы не авторизованы. Выполните вход, чтобы продолжить.")
		return
	}

	var id, age int
	var name string
	fmt.Print("Введите ID пользователя для обновления: ")
	fmt.Scanln(&id)
	fmt.Print("Введите новое имя пользователя: ")
	fmt.Scanln(&name)
	fmt.Print("Введите новый возраст пользователя: ")
	fmt.Scanln(&age)

	updatedUser := User{Name: name, Age: age} // Убираем ID из тела запроса
	body, err := json.Marshal(updatedUser)
	if err != nil {
		log.Fatal("Ошибка при сериализации данных:", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%d", baseURL, id), bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("Ошибка при создании запроса:", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ошибка при отправке запроса:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Пользователь успешно обновлен.")
	} else {
		fmt.Println("Ошибка при обновлении пользователя, код ошибки:", resp.StatusCode)
	}
}

// Функция для удаления пользователя
func deleteUser() {
	if authToken == "" {
		fmt.Println("Вы не авторизованы. Выполните вход, чтобы продолжить.")
		return
	}

	var id int
	fmt.Print("Введите ID пользователя для удаления: ")
	fmt.Scanln(&id)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%d", baseURL, id), nil)
	if err != nil {
		log.Fatal("Ошибка при создании запроса:", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ошибка при отправке запроса:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Пользователь успешно удален.")
	} else {
		fmt.Println("Ошибка при удалении пользователя, код ошибки:", resp.StatusCode)
	}
}

// Главное меню
func main() {
	var choice int
	for {
		fmt.Println("\nМеню:")
		fmt.Println("1. Авторизоваться")
		fmt.Println("2. Получить всех пользователей")
		fmt.Println("3. Получить пользователя по ID")
		fmt.Println("4. Добавить пользователя")
		fmt.Println("5. Обновить пользователя")
		fmt.Println("6. Удалить пользователя")
		fmt.Println("7. Выйти")
		fmt.Print("Выберите опцию: ")

		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Ошибка ввода. Попробуйте еще раз.")
			continue
		}

		switch choice {
		case 1:
			var username, password string
			fmt.Print("Введите имя пользователя: ")
			fmt.Scanln(&username)
			fmt.Print("Введите пароль: ")
			fmt.Scanln(&password)
			login(username, password)
		case 2:
			getUsers()
		case 3:
			getUserByID()
		case 4:
			createUser()
		case 5:
			updateUser()
		case 6:
			deleteUser()
		case 7:
			fmt.Println("Выход...")
			return
		default:
			fmt.Println("Неверная опция.")
		}
	}
}
