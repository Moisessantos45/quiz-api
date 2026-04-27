package middleware

import (
	"log"
	"quiz/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(maker *utils.PasetoMaker, rd *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		log.Printf("Starting AuthMiddleware for path: %s", c.FullPath())
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 7 || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"msg": "missing or invalid Authorization header"})
			return
		}

		token := strings.TrimSpace(authHeader[7:])
		if token == "" || len(token) < 20 {
			log.Printf("Invalid token length: %d", len(token))
			c.AbortWithStatusJSON(401, gin.H{"msg": "invalid token"})
			return
		}

		log.Printf("Received token: %s", token[:8]+"...")

		isRevoked, err := rd.Exists(c, token).Result()
		if err != nil {
			log.Printf("Redis EXISTS error: %v", err)
			c.AbortWithStatusJSON(500, gin.H{"msg": "internal error"})
			return
		}

		log.Printf("Token revocation check for token: %s, isRevoked: %d", token[:8]+"...", isRevoked)

		if isRevoked == 0 {
			log.Printf("Token expired or logged out: %s", token[:8]+"...")
			c.AbortWithStatusJSON(401, gin.H{"msg": "token revoked or expired"})
			return
		}

		payload, err := maker.VerifyToken(token)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			c.AbortWithStatusJSON(401, gin.H{"msg": "invalid token"})
			return
		}

		c.Set("email", payload.Email)
		c.Set("payload", payload)
		c.Set("token", token)
		c.Next()
	}

}
