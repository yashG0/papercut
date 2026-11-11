package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// init the default Gin engine with Logger and shutdown middleware
	r := gin.Default()

	// this will create a health route in default Router at /healthz endpoint
	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// we will listen our server at 8080 Port
	log.Println("PaperCut is running on :8080...")
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
