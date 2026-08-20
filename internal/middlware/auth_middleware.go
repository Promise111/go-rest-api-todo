package middlware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Promise111/go-rest-api-todo/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authHeader string = c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		// splitHeaderArr := strings.Split(authHeader, " ")
		// tokenString := splitHeaderArr[1]
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" || tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Invalid token claims",
			})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Invalid token claims",
			})
			c.Abort()
			return
		}

		if exp, ok := claims["exp"].(float64); !ok {
			var expiratnTime = time.Unix(int64(exp), 0)
			if time.Now().After(expiratnTime) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  false,
					"message": "Token expired",
				})
				c.Abort()
				return
			}
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
