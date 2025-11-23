package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"library-go/config"
	"library-go/database"
	"library-go/handlers"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func main() {
	// Initialize configuration
	config.LoadConfig()

	// Initialize database connection
	database.InitDB()
	defer database.CloseDB()

	// Set up Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	bookHandler := handlers.NewBookHandler()
	readerHandler := handlers.NewReaderHandler()
	borrowHandler := handlers.NewBorrowHandler()

	// API routes
	api := r.Group("/api")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", authHandler.AuthMiddleware(), authHandler.GetMe)
			auth.POST("/logout", authHandler.AuthMiddleware(), authHandler.Logout)
		}

		// Books routes (protected)
		books := api.Group("/books").Use(authHandler.AuthMiddleware())
		{
			books.GET("/", bookHandler.GetBooks)
			books.GET("/:id", bookHandler.GetBook)
			books.POST("/", bookHandler.CreateBook)
			books.PUT("/:id", bookHandler.UpdateBook)
			books.DELETE("/:id", bookHandler.DeleteBook)
		}

		// Readers routes (protected)
		readers := api.Group("/readers").Use(authHandler.AuthMiddleware())
		{
			readers.GET("/", readerHandler.GetReaders)
			readers.GET("/:id", readerHandler.GetReader)
			readers.POST("/", readerHandler.CreateReader)
			readers.PUT("/:id", readerHandler.UpdateReader)
			readers.DELETE("/:id", readerHandler.DeleteReader)
		}

		// Borrows routes (protected)
		borrows := api.Group("/borrows").Use(authHandler.AuthMiddleware())
		{
			borrows.GET("/", borrowHandler.GetBorrows)
			borrows.GET("/:id", borrowHandler.GetBorrow)
			borrows.POST("/", borrowHandler.CreateBorrow)
			borrows.PUT("/:id", borrowHandler.UpdateBorrow)
			borrows.DELETE("/:id", borrowHandler.DeleteBorrow)
		}
	}

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Library Management API",
		})
	})

	// Dashboard endpoint (serving static files)
	r.Static("/static", "./templates")
	r.GET("/dashboard", func(c *gin.Context) {
		c.File("./templates/index.html")
	})

	// Determine port to run on
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(r.Run(":" + port))
}