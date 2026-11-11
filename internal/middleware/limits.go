package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MaxBodySize(maxBytes int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, int64(maxBytes))
		ctx.Next()
	}
}
