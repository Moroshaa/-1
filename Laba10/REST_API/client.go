package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var adminURL = "http://localhost:9090/admin/users"
var baseURL = "http://localhost:9090/users"
var authToken = "" // Авторизационный токен
var userRole = ""  // Роль пользователя

// User структура для пользователя
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Функция для отправки POST запроса на сервер для авторизации
func login(username string, password string) {
	url := "http://localhost:9090/login"
	client := &http.Client{}

	// Тело запроса
	data := map[string]string{
		"username": username,
		"password": password,
	}
	jsonData, _ := json.Marshal(data)

	// Отправка запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Ошибка при создании запроса: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса: %v", err)
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
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Ошибка при разборе ответа: %v", err)
		return
	}

	// Сохраняем токен и роль
	if token, exists := result["token"]; exists {
		authToken = token
		fmt.Println("Авторизация успешна, токен сохранен:", authToken)
	} else {
		fmt.Println("Не удалось получить токен.")
	}

	if role, exists := result["role"]; exists {
		userRole = role
		fmt.Println("Роль пользователя:", userRole)
	} else {
		fmt.Println("Не удалось получить роль пользователя.")
	}
}

// Проверка роли пользователя
func isAdmin() bool {
	return userRole == "admin"
}

// Отправка GET запроса
func sendGetRequest(url string) (*http.Response, error) {
	if authToken == "" {
		fmt.Println("Вы не авторизованы. Выполните вход, чтобы продолжить.")
		return nil, fmt.Errorf("не авторизован")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	client := &http.Client{}
	return client.Do(req)
}

// Функция для отображения всех пользователей
func getUsers() {
	resp, err := sendGetRequest(baseURL)
	if err != nil {
		log.Println("Ошибка при отправке запроса:", err)
		return
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

// Получение пользователя по ID
func getUserByID() {
	var id int
	fmt.Print("Введите ID пользователя: ")
	fmt.Scanln(&id)

	url := fmt.Sprintf("http://localhost:9090/users/%d", id)
	resp, err := sendGetRequest(url)
	if err != nil {
		log.Println("Ошибка при отправке запроса:", err)
		return
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

// Функции для добавления, обновления и удаления пользователей доступны только администратору

// Функция для добавления нового пользователя
func createUser() {
	if !isAdmin() {
		fmt.Println("Доступ запрещен. Только администратор может добавлять пользователей.")
		return
	}

	// Запрашиваем у пользователя имя и возраст
	var name string
	var age int
	fmt.Print("Введите имя нового пользователя: ")
	fmt.Scanln(&name)
	fmt.Print("Введите возраст нового пользователя: ")
	fmt.Scanln(&age)

	// Формируем данные для отправки
	newUser := User{Name: name, Age: age}

	// Преобразуем данные в формат JSON
	jsonData, err := json.Marshal(newUser)
	if err != nil {
		log.Fatal("Ошибка при преобразовании данных пользователя в JSON:", err)
		return
	}

	// Отправляем POST запрос для создания пользователя
	req, err := http.NewRequest("POST", adminURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Ошибка при создании запроса:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken) // Добавляем токен в заголовок
	req.Header.Set("Content-Type", "application/json")   // Указываем, что отправляем JSON

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		// Читаем тело ответа
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Ошибка при чтении ответа:", err)
			return
		}

		// Извлекаем данные о новом пользователе из ответа
		var createdUser User
		if err := json.Unmarshal(body, &createdUser); err == nil {
			fmt.Printf("Пользователь успешно добавлен: ID = %d, Name = %s, Age = %d\n", createdUser.ID, createdUser.Name, createdUser.Age)
		} else {
			fmt.Println("Ошибка при обработке ответа:", err)
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Ошибка при добавлении пользователя, код ошибки: %d, детали: %s\n", resp.StatusCode, string(body))
	}
}

// Обновление пользователя
func updateUser() {
	if !isAdmin() {
		fmt.Println("Доступ запрещен. Только администратор может обновлять пользователей.")
		return
	}

	var id int
	fmt.Print("Введите ID пользователя для обновления: ")
	fmt.Scanln(&id)

	var user User
	user.ID = id
	fmt.Print("Введите новое имя пользователя: ")
	fmt.Scanln(&user.Name)
	fmt.Print("Введите новый возраст пользователя: ")
	fmt.Scanln(&user.Age)

	data, _ := json.Marshal(user)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%d", adminURL, id), bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Пользователь успешно обновлен.")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Ошибка при обновлении пользователя, код ошибки: %d, детали: %s\n", resp.StatusCode, string(body))
	}
}

// Удаление пользователя
func deleteUser() {
	if !isAdmin() {
		fmt.Println("Доступ запрещен. Только администратор может удалять пользователей.")
		return
	}

	var id int
	fmt.Print("Введите ID пользователя для удаления: ")
	fmt.Scanln(&id)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%d", adminURL, id), nil)
	if err != nil {
		log.Fatal("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Пользователь успешно удален.")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Ошибка при удалении пользователя, код ошибки: %d, детали: %s\n", resp.StatusCode, string(body))
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

		// Отображаем дополнительные опции только для администратора
		if isAdmin() {
			fmt.Println("4. Добавить пользователя")
			fmt.Println("5. Обновить пользователя")
			fmt.Println("6. Удалить пользователя")
		}

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
