package middlewares

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"neilz.space/web/controllers"
)

// RequireAuth - verify with JWT
func RequireAuth(c *gin.Context) {
	accessToken, err := c.Cookie("access-token")
	if err != nil {
		controllers.ErrorRediect(c, "Authorization Cookie is required")
		c.Abort()
	}
	if accessToken == "" {
		controllers.ErrorRediect(c, "Authorization Cookie is required")
		c.Abort()
		return
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_KEY")), nil
	})

	if err != nil || !token.Valid {
		if err == jwt.ErrTokenExpired {
			controllers.ErrorRediect(c, "Token Expired")
			c.Abort()
			return
		}
		controllers.ErrorRediect(c, "Invalid token")
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Set("claims", claims)
	}

	c.Next()
}
