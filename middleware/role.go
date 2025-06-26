package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ticketing-system-backend/auth-service/repository"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)

		user, err := repository.FindUserById(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "user tidak ditemukan",
			})
			c.Abort()
			return
		}

		hasAccess := false
		for _, ur := range user.Roles {
			for _, allowed := range allowedRoles {
				if ur.Level == allowed {
					hasAccess = true
					break
				}
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "akses ditolak: role tidak diizinkan",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
