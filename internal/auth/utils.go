package auth

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUserIdFromContext(c *gin.Context) (uint, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, errors.New("authorization header is missing")
	}

	token, err := ParseBearerToken(authHeader)
	if err != nil {
		return 0, fmt.Errorf("failed to parse bearer token: %w", err)
	}

	claims, err := ParseJWTToken(token)
	if err != nil {
		return 0, fmt.Errorf("failed to parse JWT token: %w", err)
	}

	floatUserId, ok := claims["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID in token claims: %v", claims["id"])
	}

	return uint(floatUserId), nil
}
