package main

import (
	"bytes"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"image/png"
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

	// Barcode generation API
	r.POST("/api/barcode", func(c *gin.Context) {
		var request struct {
			Text string `json:"text" binding:"required"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Generate barcode
		barIntCS, err := code128.Encode(request.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate barcode"})
			return
		}

		// Convert to barcode.Barcode (Fix for CheckSum issue)
		bar, ok := barIntCS.(barcode.Barcode)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Type assertion to barcode.Barcode failed"})
			return
		}

		// Scale barcode
		bar, err = barcode.Scale(bar, 400, 100)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scale barcode"})
			return
		}

		// Encode barcode as PNG
		var buf bytes.Buffer
		if err := png.Encode(&buf, bar); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode barcode"})
			return
		}

		// Return barcode as image
		c.Data(http.StatusOK, "image/png", buf.Bytes())
	})

	r.Run(":8080") // Runs on port 8080
}
