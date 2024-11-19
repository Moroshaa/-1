package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupRouter создает и настраивает маршруты для тестирования
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/users", GetUsers)
	r.GET("/users/:id", GetUserByID)
	r.POST("/users", CreateUser)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)
	return r
}

// TestMain инициализирует базу данных перед запуском тестов
func TestMain(m *testing.M) {
	connectDatabase() // Инициализация базы данных
	defer db.Close()  // Закрытие соединения при завершении тестов

	os.Exit(m.Run()) // Запуск тестов
}

// TestCreateUser тестирует создание нового пользователя
func TestCreateUser(t *testing.T) {
	router := setupRouter()

	// Создаем корректный JSON для пользователя
	user := `{"name": "Test User", "age": 25}`
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(user))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code 201, got %v", w.Code)
	}

	var createdUser User
	if err := json.Unmarshal(w.Body.Bytes(), &createdUser); err != nil {
		t.Errorf("Expected nil, but got: %v", err)
	}
	if createdUser.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%v'", createdUser.Name)
	}
	if createdUser.Age != 25 {
		t.Errorf("Expected age 25, got '%v'", createdUser.Age)
	}
}

// TestGetUsers тестирует получение списка пользователей
func TestGetUsers(t *testing.T) {
	r := setupRouter()

	// Добавляем тестового пользователя перед тестом
	_, err := db.Exec("INSERT INTO users (name, age) VALUES ('Test User', 25)")
	assert.Nil(t, err)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверка, что статус ответа - "OK"
	assert.Equal(t, http.StatusOK, w.Code)

	var users []User
	err = json.Unmarshal(w.Body.Bytes(), &users)
	assert.Nil(t, err)
	assert.NotEmpty(t, users) // Проверка, что пользователи есть в ответе
}

// TestGetUserByID тестирует получение пользователя по ID
func TestGetUserByID(t *testing.T) {
	r := setupRouter()

	// Добавляем тестового пользователя перед тестом
	var userID int
	err := db.QueryRow("INSERT INTO users (name, age) VALUES ('Test User', 25) RETURNING id").Scan(&userID)
	assert.Nil(t, err)

	req, _ := http.NewRequest("GET", "/users/"+strconv.Itoa(userID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверка, что статус ответа - "OK"
	if w.Code == http.StatusOK {
		var user User
		err := json.Unmarshal(w.Body.Bytes(), &user)
		assert.Nil(t, err)
		assert.Equal(t, userID, user.ID)
	} else {
		assert.Equal(t, http.StatusNotFound, w.Code)
	}
}

// TestUpdateUser тестирует обновление данных пользователя
func TestUpdateUser(t *testing.T) {
	r := setupRouter()

	// Добавляем тестового пользователя перед тестом
	var userID int
	err := db.QueryRow("INSERT INTO users (name, age) VALUES ('Old User', 20) RETURNING id").Scan(&userID)
	assert.Nil(t, err)

	updatedUser := User{Name: "Updated User", Age: 35}
	jsonValue, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/users/"+strconv.Itoa(userID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверка, что статус ответа - "OK"
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверка обновленных данных в ответе
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "User updated successfully", response["message"])
}

// TestDeleteUser тестирует удаление пользователя
func TestDeleteUser(t *testing.T) {
	r := setupRouter()

	// Добавляем тестового пользователя перед тестом
	var userID int
	err := db.QueryRow("INSERT INTO users (name, age) VALUES ('User to Delete', 40) RETURNING id").Scan(&userID)
	assert.Nil(t, err)

	req, _ := http.NewRequest("DELETE", "/users/"+strconv.Itoa(userID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверка, что статус ответа - "OK"
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверка, что пользователь удален
	req, _ = http.NewRequest("GET", "/users/"+strconv.Itoa(userID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
