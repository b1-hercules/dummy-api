package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Secret key untuk JWT
var secretKey = []byte("your-secret-key")

// Struktur untuk login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Struktur untuk buku
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// Data buku statis
var books = []Book{
	{ID: 1, Title: "Go Programming", Author: "John Doe"},
	{ID: 2, Title: "Microservices with Go", Author: "Jane Doe"},
	{ID: 3, Title: "Building Web Apps with Gin", Author: "Alice"},
}

// Handler untuk login
func login(c *gin.Context) {
    var request LoginRequest
    if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Logika autentikasi (contoh sederhana)
    if request.Username != "admin" || request.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Buat token JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": request.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

    c.JSON(http.StatusOK, gin.H{
		"response_code":        "00",
		"response_description": "success",
		"response_data": gin.H{
			"token": tokenString,
		},
	})
}

// Middleware untuk otentikasi JWT
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Cek apakah header Authorization kosong
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Parse token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})

		// Cek apakah token valid
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Handler untuk mendapatkan daftar buku
func getBooks(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "response_code":        "00",
        "response_description": "success",
        "response_data":                 books,
    })
}

func main() {
	r := gin.Default()

	// Endpoint untuk login
	r.POST("/login", login)

	// Endpoint untuk mendapatkan buku, dilindungi dengan middleware
	protected := r.Group("/")
	protected.Use(authMiddleware())
	protected.GET("/books", getBooks)

	// Jalankan server
	r.Run(":8080")
}
