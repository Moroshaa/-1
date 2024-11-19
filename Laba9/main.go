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

// User - структура для пользователя
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Credentials - структура для авторизации
type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var db *sql.DB
var jwtSecretKey = []byte("my_secret_key") // Секретный ключ для токенов

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

// Генерация токена
func generateToken(username string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

// Авторизация пользователя и генерация токена
func login(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	// Проверка логина и пароля
	if creds.Username != "morosha" || creds.Password != "123123123" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := generateToken(creds.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Мидлвар для проверки токена
func authenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || len(tokenString) < 7 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Получить всех пользователей
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

// Получить пользователя по ID
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

// Добавить пользователя
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

// Обновить пользователя
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

// Удалить пользователя
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

// Главная функция
func main() {
	connectDatabase()
	defer db.Close()

	r := gin.Default()

	// Маршруты для логина
	r.POST("/login", login)

	// Защищённые маршруты
	protected := r.Group("/", authenticateMiddleware())
	protected.GET("/users", GetUsers)
	protected.GET("/users/:id", GetUserByID)
	protected.POST("/users", CreateUser)
	protected.PUT("/users/:id", UpdateUser)
	protected.DELETE("/users/:id", DeleteUser)

	r.Run(":9090")
}
