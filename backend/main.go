package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/daitonium/pdf_alchemist/backend/internal/pdf"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
		var form pdf.Body
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var cmds []pdf.Command
		if err := json.Unmarshal([]byte(form.Commands), &cmds); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		file, err := form.FileHeader.Open()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		buf := bytes.NewBuffer(nil)

		if _, err := io.Copy(buf, file); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		errors := pdf.Validate(cmds)
		if len(errors) > 0 {
			c.JSON(400, gin.H{
				"message": "validation Failed",
				"errors":  errors,
			})
			return
		}
		document := pdf.Document{
			File:     bytes.NewReader(buf.Bytes()),
			Commands: cmds,
		}
		pdf.Transform(document)
	})
	r.GET("/documents/:id/download")

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
