package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yashG0/papercut/internal/handler"
	"github.com/yashG0/papercut/internal/middleware"
	"github.com/yashG0/papercut/internal/model"
	"github.com/yashG0/papercut/internal/service"
)

func main() {
	// 1. Initialize Gin engine with default Logger and Recovery
	r := gin.Default()

	// 2. Attach our memory protection middleware (limit request payloads to 10MB)
	r.Use(middleware.MaxBodySize(10 << 20))

	// 3. Initialize layers (Dependency Injection)
	docService := service.NewDocumentService()
	docHandler := handler.NewDocumentHandler(docService)

	// 4. Liveness / Readiness health probe
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.HealthResponse{
			Service: "papercut",
			Status:  "healthy",
		})
	})

	// 5. Versioned API routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/render/pdf", docHandler.RenderPDF)
	}

	log.Println("PaperCut microservice listening on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
