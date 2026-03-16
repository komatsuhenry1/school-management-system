package middleware

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu         sync.Mutex
	clients    = make(map[string]*client)
	blockedIPs = make(map[string]time.Time)
)

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("AuthUser")
		ip := c.ClientIP()
		limiter := getClient(ip)
		if limiter == nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "IP temporarily blocked"})
			return
		}

		if !limiter.Allow() {
			mu.Lock()
			blockedIPs[ip] = time.Now().Add(1 * time.Minute)
			mu.Unlock()

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests - IP blocked for 1 minute"})
			return
		}

		const BearerSchema = "Bearer "
		header := c.GetHeader("Authorization")
		fmt.Println("header", header)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
			return
		}

		tokenString := header[len(BearerSchema):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signature method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		fmt.Println(token)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "ADMIN" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": "false", "data": []interface{}{}, "message": "Rota com acesso apenas para usuários administradores"})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func getClient(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if unblockTime, blocked := blockedIPs[ip]; blocked {
		if time.Now().Before(unblockTime) {
			return nil
		}
		delete(blockedIPs, ip)
	}

	c, exists := clients[ip]
	if !exists {
		limiter := rate.NewLimiter(1, 5)
		clients[ip] = &client{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	c.lastSeen = time.Now()
	return c.limiter
}

func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}
