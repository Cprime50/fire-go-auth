package role

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidJson  = errors.New("invalid JSON format")
)

type EmailInput struct {
	Email string `json:"email"`
}

func MakeAdminHandler(ctx *gin.Context, service AdminService) {
	var input EmailInput
	if err := ctx.BindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJson})
		return
	}

	emailOk := validateEmail(input.Email)
	if !emailOk {
		log.Printf("Error validating email invalid email format: %v", input.Email)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrInvalidEmail})
		return
	}

	if err := service.MakeAdmin(input.Email); err != nil {
		log.Printf("Error assigning admin role: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s is now an admin", input.Email)})
}

func RemoveAdminHandler(ctx *gin.Context, service AdminService) {
	var input EmailInput
	if err := ctx.BindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJson})
		return
	}

	emailOk := validateEmail(input.Email)
	if !emailOk {
		log.Printf("Error validating email invalid email format: %v", input.Email)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrInvalidEmail})
		return
	}

	if err := service.RemoveAdmin(input.Email); err != nil {
		log.Printf("Error assigning user role: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s admin rights have been revoked", input.Email)})
}
