package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/BagasDhitya/owasp-secure-todo/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthConfig struct {
	AccessSecret string
}

func AuthRequired(cfg AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing access token"})
			return
		}
		tokStr := cookie.Value
		claims := &security.Claims{}
		tok, err := jwt.ParseWithClaims(tokStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.AccessSecret), nil
		}, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil || !tok.Valid || claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid/expired token"})
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// Simple CSRF check for unsafe methods when using cookie auth
func CSRFMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if m := strings.ToUpper(c.Request.Method); m == "POST" || m == "PUT" || m == "DELETE" || m == "PATCH" {
			csrfCookie, _ := c.Cookie("csrf")
			header := c.GetHeader("X-CSRF-Token")
			if csrfCookie == "" || header == "" || csrfCookie != header {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "csrf validation failed"})
				return
			}
		}
		c.Next()
	}
}
