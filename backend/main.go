package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize database
	dsn := "host=db user=myuser password=mypass dbname=mydb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.Exec("SELECT 1") // Simple query to verify connection
	// Initialize Gin router
	r := gin.Default()

	// GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// POST endpoint
	r.POST("/data", func(c *gin.Context) {
		var payload map[string]interface{}
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"status": "received", "data": payload})
	})

	// Start server
	r.Run(":8080")
}
