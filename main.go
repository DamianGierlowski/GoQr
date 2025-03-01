package main

import (
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	// QR Code generation endpoint
	r.POST("/api/qrcode", func(c *gin.Context) {

		apiKey := c.GetHeader("X-API-Key")
		if apiKey != os.Getenv("API-KEY") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Get text from request JSON
		var request struct {
			Text string `json:"text" binding:"required"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		// Generate QR code
		var png []byte
		png, err := qrcode.Encode(request.Text, qrcode.High, 256)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
			return
		}

		// Return QR code as PNG
		c.Data(http.StatusOK, "image/png", png)
	})

	r.Run(":8080") // Runs on port 8080
}
