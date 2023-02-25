package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
)

func (h *Handler) getTokenHandler(c *gin.Context) {
	// Create a new JWT token with an expiration time of 5 minutes from now
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	})

	// Sign the token with a secret key and get the resulting signed string
	ss, err := token.SignedString([]byte("MySignature"))
	if err != nil {
		// If there was an error signing the token, return an internal server error with the error message
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the signed JWT token in the response body
	c.JSON(http.StatusOK, gin.H{
		"token": ss,
	})
}

func (h *Handler) validateToken(token string) error {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte("MySignature"), nil
	})

	return err
}

func (h *Handler) authorizationMiddleware(c *gin.Context) {
	s := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(s, "Bearer ")

	if err := h.validateToken(token); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
}
