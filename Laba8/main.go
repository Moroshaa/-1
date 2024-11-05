package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq" // Подключаем драйвер PostgreSQL
)

// Структура пользователя
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required,gte=0,lte=130"`
}

// Объявляем глобальные переменные для работы с базой и валидацией
var (
	db       *sql.DB
	validate *validator.Validate
)

// Подключение к базе данных
func connectDatabase() {
	var err error
	dsn := "postgres://morosha:123123123@localhost:8080/userdb?sslmode=disable" // добавьте sslmode=disable, если SSL не используется
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}
	log.Println("Подключение к базе данных успешно")
}

// Получение всех пользователей с пагинацией и фильтрацией
func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Получение пользователя по ID
func getUserByID(c *gin.Context) {
	id := c.Param("id")
	var user User
	err := db.QueryRow("SELECT id, name, age FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Age)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Создание нового пользователя
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация данных
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вставка в базу данных
	err := db.QueryRow("INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id", user.Name, user.Age).Scan(&user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Обновление информации о пользователе
func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	// Прочитайте JSON-данные из запроса
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	// Выполните SQL-запрос для обновления данных
	result, err := db.Exec("UPDATE users SET name = $1, age = $2 WHERE id = $3", user.Name, user.Age, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверьте, что обновилась хотя бы одна строка
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// Удаление пользователя
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь удален"})
}

func main() {
	// Подключение к базе данных
	connectDatabase()
	defer db.Close()

	// Инициализация валидатора
	validate = validator.New()

	// Создание маршрутов
	r := gin.Default()
	r.GET("/users", getUsers)
	r.GET("/users/:id", getUserByID)
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	// Запуск сервера
	if err := r.Run(":9090"); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v\n", err)
	}

}
