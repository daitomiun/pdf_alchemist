package main

import (
	"encoding/json"
	"net/http"

	"github.com/daitonium/pdf_alchemist/backend/internal"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Define a simple GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/documents/transform", func(c *gin.Context) {
		var form internal.Body
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var cmds []internal.Command
		if err := json.Unmarshal([]byte(form.Commands), &cmds); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		errors := internal.Validate(cmds)
		if len(errors) > 0 {
			c.JSON(400, gin.H{
				"message": "validation Failed",
				"errors":  errors,
			})
			return
		}
	})
	r.GET("/documents/:id/download")

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
