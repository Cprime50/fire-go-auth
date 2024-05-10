package main

import (
	"log"
	"net/http"
	"os"

	"firebase.google.com/go/v4/auth"

	"github.com/cprime50/fire-go-auth/role"

	"github.com/cprime50/fire-go-auth/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	// Initialize Firebase authentication middleware
	client, err := middleware.InitAuth()
	if err != nil {
		log.Fatalf("Error initializing Firebase auth: %v", err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	// Register routes
	RegisterRoutes(r, client)
	RegisterAdminRoutes(r, client)

	// Set port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Gin server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start Gin server: %v", err)
	}
}

func loadEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	log.Println(".env file loaded successfully")
}

func RegisterRoutes(r *gin.Engine, client *auth.Client) {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "fire-go-auth success")
	})
	r.GET("/user", middleware.Auth(client), func(c *gin.Context) {
		user, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User details not found"})
			return
		}
		userDetails := user.(*middleware.User)
		c.JSON(http.StatusOK, userDetails)
	})
}

// Admin routes
func RegisterAdminRoutes(r *gin.Engine, client *auth.Client) {
	adminService := role.NewAdminService(client)

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.Auth(client), middleware.RoleAuth("admin"))
	{
		adminRoutes.GET("/", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "admin")
		})
		adminRoutes.POST("/make", func(ctx *gin.Context) {
			role.MakeAdminHandler(ctx, adminService)
		})
		adminRoutes.DELETE("/remove", func(ctx *gin.Context) {
			role.RemoveAdminHandler(ctx, adminService)
		})
	}
}
