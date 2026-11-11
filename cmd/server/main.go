package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// init the default Gin engine with Logger and shutdown middleware
	r := gin.Default()

	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	log.Println("PaperCut is running on :8080...")
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
