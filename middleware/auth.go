package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ticketing-system-backend/auth-service/utils"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token tidak ditemukan atau tidak valid",
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token tidak valid atau kadaluarsa",
			})
			return
		}

		// Konversi user_id ke uint
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "user_id tidak valid dalam token",
			})
			return
		}

		// Simpan ke context
		c.Set("user_id", uint(userID))
		c.Set("email", claims["email"])
		c.Next()
	}
}
