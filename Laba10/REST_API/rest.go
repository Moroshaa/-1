package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/lib/pq"
)

var (
	db     *sql.DB
	jwtKey = []byte("secret_key")
)

// Структура пользователя
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Подключение к базе данных
func connectDatabase() {
	var err error
	dsn := "postgres://morosha:123123123@localhost:8080/userdb?sslmode=disable" // Подключение к Postgres
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}
	log.Println("Подключение к базе данных успешно")

	// Настройки пула соединений
	db.SetMaxOpenConns(10)                  // Максимум 10 открытых соединений
	db.SetMaxIdleConns(5)                   // Максимум 5 неактивных соединений
	db.SetConnMaxLifetime(30 * time.Minute) // Время жизни соединения
}

// Функция для авторизации (login) и создания JWT токена
func login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Преобразуем JSON тело запроса в структуру
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Проверка логина и пароля
	var role string
	if credentials.Username == "admin" && credentials.Password == "12345" {
		role = "admin"
	} else if credentials.Username == "user" {
		role = "user"
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Создание JWT токена
	claims := jwt.MapClaims{
		"username": credentials.Username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// Отправка токена в ответ
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "role": role})
}

// Middleware для проверки JWT
func authMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем наличие заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["role"] == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		role := claims["role"].(string)
		if requiredRole == "admin" && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Set("role", role)
		c.Next()
	}
}
func GetUsers(c *gin.Context) {
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

	c.JSON(http.StatusOK, users)
}
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user User
	err := db.QueryRow("SELECT id, name, age FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Age)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	err := db.QueryRow("INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id", user.Name, user.Age).Scan(&user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	result, err := db.Exec("UPDATE users SET name = $1, age = $2 WHERE id = $3", user.Name, user.Age, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Основная функция
func main() {
	connectDatabase() // Подключаемся к базе данных

	router := gin.Default()

	// Маршруты для авторизации
	router.POST("/login", login)

	// Маршруты для админов
	admin := router.Group("/admin", authMiddleware("admin"))
	{
		admin.POST("/users", CreateUser)
		admin.PUT("/users/:id", UpdateUser)
		admin.DELETE("/users/:id", DeleteUser)
		admin.GET("/", GetUsers)
		admin.GET("/:id", GetUserByID)
	}

	// Маршруты для пользователей
	user := router.Group("/users", authMiddleware("user"))
	{
		user.GET("/", GetUsers)
		user.GET("/:id", GetUserByID)
	}

	router.Run(":9090") // Запуск сервера
}
