package main

import (
	"log"
	"manipulator-go/internal/handlers"
	"manipulator-go/internal/services"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create temp directory if it doesn't exist
	if err := os.MkdirAll("./data/tmp", 0755); err != nil {
		log.Fatal("Failed to create temp directory:", err)
	}

	// Initialize services
	pdfService := services.NewPDFService()

	// Initialize handlers
	convertHandler := handlers.NewConvertHandler(pdfService)
	extractHandler := handlers.NewExtractHandler(pdfService)
	mergeHandler := handlers.NewMergeHandler(pdfService)

	// Setup router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Routes
	r.POST("/api/convert", convertHandler.Convert)
	r.POST("/api/extract", extractHandler.ExtractPages)
	r.POST("/api/merge", mergeHandler.Merge)

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
