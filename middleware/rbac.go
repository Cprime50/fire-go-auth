package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"errors"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserNotFound = errors.New("User not found")
)

func RoleAuth(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userValue, exists := ctx.Get("user")
		if !exists {
			log.Println("User not found in context")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, ok := userValue.(*User)
		if !ok || user == nil {
			log.Println("Invalid user data in context")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if user.Role == "" {
			log.Println("User role not set")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if user.Role != requiredRole {
			log.Printf("User with email %s and role %s tried to access a route that was for the %s role only",
				user.Email, user.Role, requiredRole)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		log.Printf("User with email %s and role %s authorized", user.Email, user.Role)
		ctx.Next()
	}
}

func AssignRole(ctx context.Context, client *auth.Client, email string, role string) error {
	startTime := time.Now()
	defer func() {
		log.Println("Role assigned in:", time.Since(startTime))
	}()
	user, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		log.Printf("Assign Error: User with email %s not found", email)
		return ErrUserNotFound
	}
	currentCustomClaims := user.CustomClaims
	if currentCustomClaims == nil {
		currentCustomClaims = map[string]interface{}{}
	}
	currentCustomClaims["role"] = role
	if err := client.SetCustomUserClaims(ctx, user.UID, currentCustomClaims); err != nil {
		return fmt.Errorf("AssignRole Error: Error setting custom claims: %w", err)
	}
	return nil
}
