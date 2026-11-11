package middleware_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yashG0/papercut/internal/middleware"
)

func TestMaxBodySize_RejectsOversizedPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Limit request payload to 50 bytes for fast testing
	r.Use(middleware.MaxBodySize(50))
	r.POST("/upload", func(c *gin.Context) {
		// Read the full body to trigger MaxBytesReader enforcement
		_, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "payload too large"})
			return
		}
		c.Status(http.StatusOK)
	})

	// Send 100 bytes (exceeds the 50-byte threshold)
	oversizedData := bytes.Repeat([]byte("A"), 100)
	req, _ := http.NewRequest(http.MethodPost, "/upload", bytes.NewBuffer(oversizedData))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413 Payload Too Large, got %d", w.Code)
	}
}
