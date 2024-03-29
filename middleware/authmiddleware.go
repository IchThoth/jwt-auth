package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ichthoth/jwt-auth/helpers"
	"net/http"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header")})
			c.Abort()
		}
		claims, err := helpers.ValidateTokens(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("last_name", claims.Last_name)
		c.Set("first_name", claims.First_name)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.User_Type)
		c.Next()
	}
}
